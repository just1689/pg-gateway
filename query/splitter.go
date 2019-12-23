package query

import (
	"strings"
)

func BuildQueryFromURL(u string) (result Query, err error) {
	split := splitFullURL(u)
	result = Query{
		Entity:      split[0],
		Comparisons: []Comparison{},
	}
	if len(split) > 1 {
		for id, row := range split {
			if id == 0 {
				continue
			}
			coArr := strings.Split(row, "=")
			co := Comparison{
				Field: coArr[0],
			}
			riArr := strings.Split(coArr[1], ".")
			if len(riArr) == 1 {
				err = result.processOtherQuery(coArr[0], riArr[0])
				if err != nil {
					return
				}
				continue
			}
			co.Comparator = riArr[0]
			co.Value = riArr[1]
			result.Comparisons = append(result.Comparisons, co)
		}
	}
	return
}
