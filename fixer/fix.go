package fixer

import (
	"bufio"
	"encoding/csv"
	"os"
	"strings"
)

//Callback is the function type to use when defining a custom callback
type Callback func(*FixedSheet) error

type reader interface {
	read() ([][]string, error)
	interpret() ([][]string, error)
}

//Fixer structure to encapsulate CSV functions
type Fixer struct {
	Filepath string
	Input    [][]string
	Output   [][]string
	Sheets   []FixedSheet
	Callback Callback
}

//NewFixer returns an initialized instance of Fixer
func NewFixer(fp string) (*Fixer, error) {
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		return &Fixer{}, err
	}

	fixer := &Fixer{
		Filepath: fp,
		Callback: DefaultCallback,
	}

	err := fixer.read()
	if err != nil {
		return &Fixer{}, err
	}

	err = fixer.interpret()
	if err != nil {
		return &Fixer{}, nil
	}

	return fixer, nil
}

//Open open CSV file at filepath fp returning errors
func (f *Fixer) read() error {
	file, err := os.Open(f.Filepath)
	if err != nil {
		return err
	}

	sheets := &f.Sheets

	r := csv.NewReader(bufio.NewReader(file))
	record, err := r.ReadAll()
	// Stop at EOF.
	if err != nil {
		return err
	}

	for i, rec := range record {
		// Skip the headers
		if i == 0 {
			continue
		}
		s := Statement{
			Type:   rec[0],
			Date:   rec[1],
			Amount: rec[2],
			ID:     rec[3],
			Name:   rec[4],
			Memo:   rec[5],
		}

		sheet := FixedSheet{
			IsExpense:   strings.HasPrefix(s.Amount, "-"),
			Date:        s.Date,
			Amount:      "$" + strings.TrimPrefix(s.Amount, "-"),
			Description: s.Memo,
			Category:    strings.Fields(s.Name)[0],
		}

		*sheets = append(*sheets, sheet)
	}

	f.Input = record
	return nil
}

func (f *Fixer) interpret() error {

	output := make([][]string, len(f.Sheets))

	for i := range output {
		output[i] = make([]string, 10)
		out := output[i]
		sheet := &f.Sheets[i]

		f.Callback(sheet)

		dateCol := 1
		if !sheet.IsExpense {
			dateCol = (len(output[i]) / 2) + 1
		}

		out[dateCol] = sheet.Date
		out[dateCol+1] = sheet.Amount
		out[dateCol+2] = sheet.Description
		out[dateCol+3] = sheet.Category
	}

	f.Output = output
	return nil
}

//DefaultCallback will be used if no callback is set
func DefaultCallback(sheet *FixedSheet) error {
	prefixes := []string{"Withdrawal Credit/Debit Card Debit Card",
		"Withdrawal", "Depoisit"}

	for _, prefix := range prefixes {
		if strings.HasPrefix(sheet.Description, prefix) {
			sheet.Description = strings.TrimPrefix(sheet.Description, prefix)
		}
	}

	if strings.HasSuffix(sheet.Description, "PSPE") {
		sheet.IsExpense = false
		sheet.Category = "Savings"
		sheet.Description = "Transfer from Checking to Savings PSPE"
	} else if strings.HasSuffix(sheet.Description, "INMOTION HOSTING") {
		sheet.Category = "Paycheck"
	}

	return nil
}

//Write will write the fixed information to a CSV at the given filepath
func (f *Fixer) Write(fp string) error {
	outfile, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer outfile.Close()

	writer := csv.NewWriter(outfile)
	defer writer.Flush()

	for _, value := range f.Output {
		if writer.Write(value) != nil {
			return err
		}
	}

	return nil
}

//Statement encapsulates the CSV statement from my bank
type Statement struct {
	ID     string
	Date   string
	Type   string
	Name   string
	Memo   string
	Amount string
}

//FixedSheet encapsulates the data needed for the Google Sheet
type FixedSheet struct {
	Date        string
	Description string
	Amount      string
	Category    string
	IsExpense   bool
}
