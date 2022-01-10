package commonUtils

import "fmt"

// StringSliceToString converts a slice of strings to a comma delimited string
func StringSliceToString(src []string) (tgt string) {
	tgt ="["
	for i:=0; i<len(src);i++ {
	   joinChr := ", "
		if i == 0 {
			joinChr = ""
		}
	   tgt = fmt.Sprintf("%s%s\"%s\"", tgt,joinChr,src[i])
	}
   tgt = fmt.Sprintf("%s]", tgt)
   return
 }

