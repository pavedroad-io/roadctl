package cmd

import (
	"strings"
)

type macroSubstitutions struct {
	originalValue    string
	replacementValue string
}

func (s *macroSubstitutions) replaceAll(input string, def bpDef) (newString string) {
	macroList := s.createList(def)

	listLength := len(macroList)

	if listLength == 1 {
		return strings.ReplaceAll(input, macroList[0].originalValue, macroList[0].replacementValue)
	}

	if listLength > 1 {
		modifiedString := strings.ReplaceAll(input, macroList[0].originalValue, macroList[0].replacementValue)
		for i := 1; i < listLength; i++ {
			modifiedString = strings.ReplaceAll(modifiedString, macroList[i].originalValue, macroList[i].replacementValue)
		}
		return modifiedString
	}

	return newString
}

func (s *macroSubstitutions) createList(def bpDef) (list []macroSubstitutions) {
	list = []macroSubstitutions{
		{originalValue: "template",
			replacementValue: def.Info.Name},
		{originalValue: "organization",
			replacementValue: def.Info.Organization},
		{originalValue: "__template__",
			replacementValue: def.Info.Name},
		{originalValue: "__organization__",
			replacementValue: def.Info.Organization},
	}
	return list
}
