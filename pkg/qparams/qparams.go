package qparams

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jwmwalrus/bnp/slice"
)

// QParam defines a query parameter
// A QParam is sort of an SQL where condition, with proper wildcards.
//
// A query will be something like
// ```sql
// title=some*thing and artist=some*other genre=[Rr]ock or genre=ska not genre=pop
// ```
//
// ## CSV-like conditions:
// The following two expressions are equivalent
// ```sql
// genre=pop,rock,punk
// genre=pop or genre=rock or genre=punk
// ```
type QParam struct {
	Or  bool
	Not bool
	Key string
	Val string
}

// ToFuzzy converts the given value into a fuzzy one
// * Numbers and proper wildcards are never converted
func (qp *QParam) ToFuzzy() *QParam {
	out := *qp
	if strings.ContainsAny(out.Val, "*?[]") {
		return &out
	}

	if _, err := strconv.ParseInt(out.Val, 10, 64); err == nil {
		return &out
	}

	out.Val = "*" + out.Val + "*"
	if strings.Count(out.Val, " ") == 0 {
		return &out
	}

	out.Val = strings.Join(strings.Split(out.Val, " "), "*")
	return &out
}

// ToSQL converts the given wildcards to SQL
func (qp *QParam) ToSQL() (out QParam) {
	out = *qp
	var sb strings.Builder
	literal := 0
	for _, r := range out.Val {
		switch r {
		case '[':
			literal++
		case ']':
			literal--
		case '*':
			if literal == 0 {
				sb.WriteRune('%')
				continue
			}
		case '?':
			if literal == 0 {
				sb.WriteRune('_')
				continue
			}
		default:
		}
		sb.WriteRune(r)
	}
	out.Val = sb.String()
	return
}

// ParseParams parse a params string and return an equivalent slice
func ParseParams(params string) (qp []*QParam, err error) {
	qp = []*QParam{}

	if len(params) == 0 {
		err = errors.New("Cannot parse empty string")
		return
	}

	out := []QParam{}

	words := []string{}
	literal := 0
	from := 0
	for i, r := range params {
		switch r {
		case '[':
			literal++
			continue
		case ']':
			literal--
			continue
		case ' ':
			words = append(words, params[from:i])
			from = i + 1
		}
	}
	words = append(words, params[from:])

	aux := []string{}
	for i := 0; i < len(words); i++ {
		if i > 0 &&
			isCondition(words[i-1]) &&
			isCondition(words[i]) {
			aux[len(aux)-1] = words[i]
			continue
		}
		aux = append(aux, words[i])
	}

	if !isCondition(words[0]) {
		words = []string{"and"}
		words = append(words, aux...)
	} else {
		words = aux
	}

	var kv string
	var newq QParam
	from = 1
	cond := words[0]
	for i := 1; i < len(words); i++ {
		if isCondition(words[i]) {
			kv = strings.Join(words[from:i], " ")
			if newq, err = createParam(cond, kv); err != nil {
				return
			}
			out = append(out, newq)
			cond = words[i]
			from = i + 1
		}
	}
	kv = strings.Join(words[from:], " ")
	if newq, err = createParam(cond, kv); err != nil {
		return
	}
	out = append(out, newq)

	for _, v := range out {
		csv := splitCSV(&v)
		qp = append(qp, csv...)
	}
	return
}

func createParam(cond, kv string) (newq QParam, err error) {
	var k, v string
	var or, not bool
	if k, v, err = getKeyVal(kv); err != nil {
		return
	}
	or, not = parseCondition(cond)
	newq = QParam{Or: or, Not: not, Key: k, Val: v}
	return
}

func getKeyVal(s string) (k, v string, err error) {
	idx := strings.Index(s, "=")
	if idx < 0 {
		err = fmt.Errorf("No key=value pair found in %s", s)
		return
	}
	k = strings.TrimSpace(strings.ToLower(s[:idx]))
	v = strings.TrimSpace(s[idx+1:])
	if k == "" || v == "" {
		err = fmt.Errorf("No key or value found in: %s", s)
		return
	}
	return
}

func isCondition(s string) bool {
	conditions := []string{"and", "or", "not"}
	cond := strings.ToLower(s)
	return slice.Contains(conditions, cond)
}

func parseCondition(cond string) (or, not bool) {
	switch strings.ToLower(cond) {
	case "or":
		or = true
	case "not":
		not = true
	default:
	}
	return
}

func splitCSV(in *QParam) (out []*QParam) {
	if strings.Index(in.Val, ",") < 0 {
		out = append(out, in)
		return
	}

	s := strings.Split(in.Val, ",")
	if s[0] != "" {
		nosp := strings.TrimSpace(s[0])
		out = append(out, &QParam{Or: in.Or, Not: in.Not, Key: in.Key, Val: nosp})
	}
	for i := 1; i < len(s); i++ {
		if s[i] != "" {
			not := false
			or := true
			if in.Not {
				not = true
				or = false
			}
			nosp := strings.TrimSpace(s[i])
			out = append(out,
				&QParam{Not: not, Or: or, Key: in.Key, Val: nosp})
		}
	}
	return
}
