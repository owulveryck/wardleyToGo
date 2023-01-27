package svg

import (
	"strings"
)

func splitString(s string, max int) []string {
	output := make([]string, 1)
	words := strings.Fields(s)
	for _, w := range words {
		switch {
		case len(output[len(output)-1]) >= max:
			output = append(output, w)
		case len(output[len(output)-1]) == 0:
			output[len(output)-1] = w
		default:
			output[len(output)-1] = output[len(output)-1] + " " + w
		}
	}
	return output
}
