package cmdutils

import (
	"strings"
	"unicode"

	"github.com/galactixx/codescout"
)

func findLengthOfOutput(output string) int {
	maxLineLength := 0
	for _, line := range strings.Split(output, "\n") {
		maxLineLength = getMax(maxLineLength, len(line))
	}
	return maxLineLength
}

func getMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func capitalizeString(str string) string {
	words := strings.Split(str, "-")
	for idx, word := range words {
		runes := []rune(word)
		runes[0] = unicode.ToUpper(runes[0])
		words[idx] = string(runes)
	}
	return strings.Join(words, " ")
}

func getNameFromNodes(node interface{}) string {
	if v, ok := node.(codescout.NodeInfo); ok {
		return v.Name()
	} else if v, ok := node.(*codescout.NodeInfo); ok {
		return (*v).Name()
	} else {
		return ""
	}
}
