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
	"github.com/dmage/lproc/pkg/rewriter"
	"github.com/dmage/lproc/pkg/transaction"
)

func ImportCSV(filename string, format config.Format, classifiers []classifier.Config) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

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

	var transactions []*transaction.Transaction
	r := bufio.NewReader(rd)
	lineno := 0
	for {
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
				return fmt.Errorf("%s:%d: got %d columns, want %d", filename, lineno, len(fields), len(format.Columns))
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
					return fmt.Errorf("%s:%d: rewrite rule %d: %s", filename, lineno, ruleno+1, err)
				}
			}

			transactions = append(transactions, &transaction.Transaction{
				ID:       state.Get("TransactionID"),
				CSV:      state.Get("CSV"),
				Date:     state.Get("Date"),
				Payee:    state.Get("Payee"),
				Amount:   state.Get("Amount"),
				Currency: state.Get("Currency"),
				To:       state.Get("To"),
				From:     state.Get("From"),
			})
		}
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("%s:%d: %s", filename, lineno, err)
		}
	}
	if format.Reverse {
		l := len(transactions)
		for i := 0; i < l/2; i++ {
			transactions[i], transactions[l-1-i] = transactions[l-1-i], transactions[i]
		}
	}
	for _, tr := range transactions {
		tr.Print()
	}
	return nil
}
