package cmd

import "strings"

const (
	macro_start string = "__"
	macro_end   string = "__"
)

func macroSubstition(inputString, searchValue, replaceValue string) (modifiedString string) {
	//	macroString := macro_start + searchValue + macro_end
	return strings.ReplaceAll(inputString, searchValue, replaceValue)
}
