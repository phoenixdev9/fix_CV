package fixer

import (
	"reflect"
	"runtime"
	"strings"
	"testing"
)

//Test is the function schema that is used by all tests
type Test func() error
type empty struct{}

var testFile = "test.csv"
var testOutFile = "result.csv"

//TestFilepath the location of the test csv file to use
func TestFilepath() string {
	return testFile
}

func expected(t *testing.T, e error, g error) {
	t.Errorf("Expected(%T) '%v' got(%T) '%v'\n", e, e, g, g)
}

func die(t *testing.T, e error, g error) {
	expected(t, e, g)
	t.Fatal("*dead*")
}

func test(t *testing.T, e error, g error) {
	if e != g {
		die(t, e, g)
	}
}

func name(i interface{}) string {
	prefix := reflect.TypeOf(empty{}).PkgPath()
	return strings.TrimPrefix(runtime.FuncForPC(reflect.ValueOf(i).Pointer()).
		Name(), prefix+".")
}
