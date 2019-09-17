package importcsv

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/text/encoding/ianaindex"

	"github.com/dmage/lproc/pkg/classifier"
	"github.com/dmage/lproc/pkg/config"
	"github.com/dmage/lproc/pkg/ledger"
	"github.com/dmage/lproc/pkg/rewriter"
	"github.com/dmage/lproc/pkg/transaction"
)

func ImportCSV(filename string, format config.Format, existingTransactionsList []ledger.Transaction, classifiers []classifier.Config) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	existingTransactions := make(map[string]ledger.Transaction)
	for _, tr := range existingTransactionsList {
		if tr.ID == "" {
			continue
		}
		if _, ok := existingTransactions[tr.ID]; ok {
			return fmt.Errorf("found duplicates in the existing transactions (TransactionID: %s)", tr.ID)
		}
		existingTransactions[tr.ID] = tr
	}

	var rd io.Reader
	if format.Encoding != "" {
		enc, err := ianaindex.IANA.Encoding(format.Encoding)
		if err != nil {
			return err
		}

		rd = enc.NewDecoder().Reader(f)
	} else {
		rd = f
	}

	classifierFactory := classifier.NewFactory(classifiers)
	defer func() {
		if err := classifierFactory.Close(); err != nil {
			log.Printf("unable to close classifiers factory: %s", err)
		}
	}()

	var transactions []*transaction.Transaction
	r := bufio.NewReader(rd)
	lineno := 0
	errors := 0
	tooManyErrors := false
	for {
		if errors >= 10 {
			tooManyErrors = true
			break
		}

		lineno++
		line, err := r.ReadString('\n')
		if err == nil || line != "" {
			if lineno <= format.HeaderSkip {
				continue
			}
			line = strings.TrimRight(line, "\r\n")

			csvReader := csv.NewReader(strings.NewReader(line))
			if format.Comma != "" {
				var comma rune
				for _, comma = range format.Comma {
					// get the first character
					break
				}
				csvReader.Comma = comma
			}

			fields, err := csvReader.Read()
			if err != nil {
				log.Printf("unable to parse CSV line: %s", line)
				return fmt.Errorf("%s:%d: %s", filename, lineno, err)
			}

			if len(fields) != len(format.Columns) {
				log.Printf("unable to process CSV line: %s", line)
				log.Printf("%s:%d: got %d columns, want %d", filename, lineno, len(fields), len(format.Columns))
				errors += 1
				continue
			}

			state := rewriter.NewState(classifierFactory)
			state.Assign("CSV", line)
			for i := range fields {
				if format.Columns[i] == "" {
					continue
				}
				state.Assign(format.Columns[i], fields[i])
			}

			for ruleno, rule := range format.Rewrite {
				err := rule.Execute(state)
				if err != nil {
					log.Printf("%s:%d: rewrite rule %d: %s", filename, lineno, ruleno+1, err)
					errors += 1
					continue
				}
			}

			payee := strings.Trim(state.Get("Payee"), " \t")
			if payee == "" {
				payee = "<Unspecified payee>"
			}

			tr := &transaction.Transaction{
				ID:       state.Get("TransactionID"),
				CSV:      state.Get("CSV"),
				Date:     state.Get("Date"),
				Payee:    payee,
				Amount:   state.Get("Amount"),
				Currency: state.Get("Currency"),
				To:       state.Get("To"),
				From:     state.Get("From"),
			}
			if tr.ID == "" {
				log.Printf("%s:%d: invalid transaction: TransactionID should have a non-empty value", filename, lineno)
				errors += 1
				continue
			}
			transactions = append(transactions, tr)
		}
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("%s:%d: %s", filename, lineno, err)
		}
	}
	if errors > 0 {
		if tooManyErrors {
			return fmt.Errorf("too many errors, aborted")
		}
		return fmt.Errorf("got %d error(s)", errors)
	}

	if format.Reverse {
		l := len(transactions)
		for i := 0; i < l/2; i++ {
			transactions[i], transactions[l-1-i] = transactions[l-1-i], transactions[i]
		}
	}

	var filteredTransactions []*transaction.Transaction
	for _, tr := range transactions {
		if etr, ok := existingTransactions[tr.ID]; ok {
			if tr.Date == etr.Date && tr.Payee == etr.Payee {
				log.Printf("SKIP: %s %s", tr.Date, tr.Payee)
			} else {
				log.Printf("!!!!: %s %s", tr.Date, tr.Payee)
			}
			continue
		}
		log.Printf(" ADD: %s %s", tr.Date, tr.Payee)
		filteredTransactions = append(filteredTransactions, tr)
	}

	for _, tr := range filteredTransactions {
		tr.Print()
	}

	return nil
}
