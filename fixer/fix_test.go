package fixer

import "testing"

var fixer = &Fixer{
	Filepath: testFile,
	Callback: func(sheet *FixedSheet) error {
		sheet.Amount = ""
		sheet.Description = ""
		sheet.Category = ""
		sheet.Date = ""

		return nil
	},
}

func TestAll(t *testing.T) {
	callbacks := []Test{testRead, testInterpret, testWrite, testNew}
	for _, cb := range callbacks {
		t.Logf("Running Test Function: %v\n", name(cb))
		err := cb()
		test(t, nil, err)
	}

	test(t, nil, fixer.Write(testOutFile))
}

func testNew() error {
	var err error
	fixer, err = NewFixer("test.csv")
	return err
}

func testRead() error {
	return fixer.read()
}

func testInterpret() error {
	return fixer.interpret()
}
func testWrite() error {
	return fixer.Write(testOutFile)
}
