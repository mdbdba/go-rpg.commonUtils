package commonUtils

import "fmt"

// MapStringIntToString converts a map of integers keyed by strings to a comma
// delimited string.
func MapStringIntToString(src map[string]int) (tgt string) {
	tgt = "["
	firstLoop := true
	for k,v := range src {
		joinChr := ", "
		if firstLoop {
			joinChr = ""
			firstLoop = false
	    }
		tgt = fmt.Sprintf("%s%s\"%s\": %d",tgt,joinChr,k,v)
	 }
	tgt = fmt.Sprintf("%s]", tgt)
	return
 }

