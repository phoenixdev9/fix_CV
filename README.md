# fixcsv
### Version 0.0.1

Takes CSV bank statements from PSCU and organizes the information to be appended to the Google Sheets Monthly Budget template

## Usage
You install and run this repo using
```sh
go install
fixcsv pscu_statement.csv output_for_sheets.csv
```

## Library Usage
You can also use the [contents of the fixer folder](fixer/) as a library


#### Callback
If you use fixcsv as a library you can pass a callback to your fixer which allows you to customize the output. fixcsv will use a default callback if none is given

#####The callback function signature
func(*fixer.FixedSheet) error

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/maxmansfield/fixcsv/fixer"
)

func main() {
	usage := "USAGE: " + os.Args[0] + " [IN] [OUT]\n"

	if len(os.Args) != 3 {
		fmt.Println(usage)
    os.Exit(0)
	}
	f, err := fixer.NewFixer(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	f.Write(os.Args[2])
	return
}

// This is an example of a callback. It simply checks if the description starts with 'deposit'.
// If it does then it removes it from the description and makes the category 'deposit'
func fixsheets (sheet *fixer.FixedSheet) error {
  prefix := "Deposit"
  if strings.HasPrefix(sheet.Description, prefix) {
    sheet.Description = string.TrimPrefix(sheet.Deposit, prefix)
    sheet.Category = prefix
  }
  
  return nil
}
```
