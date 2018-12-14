package ledger

import (
	"encoding/xml"
	"fmt"
	"os/exec"
)

type xmlAmount struct {
	Commodity string `xml:"commodity>symbol"`
	Quantity  string `xml:"quantity"`
}

type xmlPosting struct {
	Account string    `xml:"account>name"`
	Amount  xmlAmount `xml:"post-amount>amount"`
}

type xmlMetadata struct {
	Key   string `xml:"key,attr"`
	Value string `xml:"string"`
}

type xmlTransaction struct {
	Date     string        `xml:"date"`
	Payee    string        `xml:"payee"`
	Metadata []xmlMetadata `xml:"metadata>value"`
	Postings []xmlPosting  `xml:"postings>posting"`
}

type xmlLedger struct {
	Transactions []xmlTransaction `xml:"transactions>transaction"`
}

type Transaction struct {
	ID    string
	Date  string
	Payee string
}

func transactionIDFromMetadata(metadata []xmlMetadata) (string, error) {
	n := 0
	transactionID := ""
	for _, value := range metadata {
		if value.Key != "TransactionID" {
			continue
		}
		transactionID = value.Value
		n++
	}
	if n > 1 {
		return transactionID, fmt.Errorf("got %d values for TransactionID", n)
	}
	return transactionID, nil
}

func GetTransactions() ([]Transaction, error) {
	cmd := exec.Command("ledger", "xml")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var data xmlLedger
	if err = xml.Unmarshal(output, &data); err != nil {
		return nil, err
	}
	var transactions []Transaction
	for _, tr := range data.Transactions {
		id, err := transactionIDFromMetadata(tr.Metadata)
		if err != nil {
			return transactions, err
		}
		transactions = append(transactions, Transaction{
			ID:    id,
			Date:  tr.Date,
			Payee: tr.Payee,
		})
	}
	return transactions, err
}
