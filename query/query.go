package query

import (
	"errors"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

var comparisonMap = map[string]string{
	"eq":  "=",
	"gt":  ">",
	"lt":  "<",
	"gte": ">=",
	"lte": "<=",
	"neg": "!=",
	"is":  "is",
}

type Query struct {
	Entity      string
	Select      []string
	Comparisons []Comparison
	Limit       int
}

func (q *Query) ToSelectQuery() (queryString string, bindArr []interface{}) {
	queryString = generateSelect(q)
	queryString += " FROM " + q.Entity + " as tbl"
	queryString, bindArr = generateWhere(q, queryString, bindArr)

	if q.Limit != 0 {
		queryString += " LIMIT $" + strconv.Itoa(len(bindArr)+1)
		bindArr = append(bindArr, strconv.Itoa(q.Limit))
	}
	return
}

func (q *Query) ToUpdateStatement(values map[string]interface{}) (sql string, bindArr []interface{}) {
	sql = "UPDATE " + q.Entity + " SET "
	for field, value := range values {
		sql += field + "=$" + strconv.Itoa(len(bindArr)+1)
		bindArr = append(bindArr, value)
	}
	sql, bindArr = generateWhere(q, sql, bindArr)
	return
}

func (q *Query) ToURL(baseURL string) string {
	url := baseURL + "/" + q.Entity
	changed := 0
	for _, sel := range q.Select {
		changed++
		if changed == 1 {
			url += "?select="
		} else {
			url += ","
		}
		url += sel
	}

	for _, c := range q.Comparisons {
		changed++
		if changed == 1 {
			url += "?"
		} else {
			url += "&"
		}
		url += c.ComparatorToURL()
	}

	if q.Limit != 0 {
		changed++
		if changed == 1 {
			url += "?"
		} else {
			url += "&"
		}
		url += "limit=" + strconv.Itoa(q.Limit)
	}
	return url

}

func generateSelect(q *Query) (queryString string) {
	queryString = "SELECT "
	if len(q.Select) == 0 {
		queryString += "row_to_json(" + "tbl" + ") as row"
	} else {
		queryString += "(select row_to_json(_) as row from (select tbl."
		for id, row := range q.Select {
			if id != 0 {
				queryString += ", tbl."
			}
			queryString += row
		}
		queryString += ") as _) as schemaname"
	}
	return
}

func generateWhere(q *Query, queryString string, bindArrIn []interface{}) (string, []interface{}) {
	if len(q.Comparisons) == 0 {
		return queryString, bindArrIn
	}
	queryString += " WHERE "
	for id, row := range q.Comparisons {
		if id > 0 {
			queryString += " AND "
		}
		queryString += row.Field + row.ComparatorToSQL() + "$" + strconv.Itoa(len(bindArrIn)+1)
		bindArrIn = append(bindArrIn, row.Value)
	}
	return queryString, bindArrIn
}

func (q *Query) processOtherQuery(field, value string) (err error) {
	if field == "limit" {
		i, _ := strconv.Atoi(value)
		q.Limit = i
		return
	} else if field == "select" {
		q.Select = strings.Split(value, ",")
		return
	}
	err = errors.New("not sure how to handle " + field + "=" + value + " for processOtherQuery. Are you missing an operator?")
	logrus.Errorln(err)
	return
}

type Comparison struct {
	Field      string
	Comparator string
	Value      string
}

func (c *Comparison) ComparatorToURL() string {
	return c.Field + "=" + c.Comparator + "." + c.Value
}
func (c *Comparison) ComparatorToSQL() string {
	v, found := comparisonMap[c.Comparator]
	if !found {
		return c.Comparator
	}
	return v

}
