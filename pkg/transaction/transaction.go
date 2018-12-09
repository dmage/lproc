package transaction

import "fmt"

type Transaction struct {
	ID       string
	CSV      string
	Date     string
	Payee    string
	Amount   string
	Currency string
	To       string
	From     string
}

func (t Transaction) Print() {
	fmt.Printf("%s * %s\n", t.Date, t.Payee)
	fmt.Printf("    ; TransactionID: %s\n", t.ID)
	fmt.Printf("    ; CSV: %s\n", t.CSV)
	fmt.Printf("    %-30s %13s %s\n", t.To, t.Amount, t.Currency)
	fmt.Printf("    %s\n", t.From)
	fmt.Printf("\n")
}
