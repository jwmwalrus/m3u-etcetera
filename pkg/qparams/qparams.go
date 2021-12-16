package qparams

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jwmwalrus/bnp/slice"
)

// QParam defines a query parameter
// A query will be something like
// ```
// title=some*thing and artist=some*other genre=[Rr]ock or genre=ska not genre=pop
// ```
// Which is sort of an SQL where condition, with proper wildcards
type QParam struct {
	Or  bool
	Not bool
	Key string
	Val string
}

func (qp *QParam) ToFuzzy() *QParam {
	out := *qp
	if strings.ContainsAny(out.Val, "*?[]") {
		return &out
	}

	out.Val = "*" + out.Val + "*"
	if strings.Count(out.Val, " ") == 0 {
		return &out
	}

	out.Val = strings.Join(strings.Split(out.Val, " "), "*")
	return &out
}

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

func ParseParams(params string) (qp []*QParam, err error) {
	conditions := []string{"and", "or", "not"}

	if len(params) == 0 {
		err = errors.New("Cannot parse empty string")
		return
	}

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

	list := []string{words[0]}
	for i := 1; i < len(words); i++ {
		if slice.Contains(conditions, words[i-1]) &&
			slice.Contains(conditions, words[i]) {
			list[len(list)-1] = words[i]
			continue
		}
		list = append(list, words[i])
	}
	words = list

	getKeyVal := func(s string) (k, v string) {
		idx := strings.Index(s, "=")
		if idx < 0 {
			return
		}
		k = strings.TrimSpace(strings.ToLower(s[:idx]))
		v = strings.TrimSpace(s[idx+1:])
		return
	}

	or := false
	not := false
	from = 0
	for i, w := range words {
		cond := strings.ToLower(w)
		if slice.Contains(conditions, cond) {
			kv := strings.Join(words[from:i], " ")
			k, v := getKeyVal(kv)
			if k == "" {
				err = fmt.Errorf("No key found in: %v", kv)
				return
			}
			newq := QParam{Or: or, Not: not, Key: k, Val: v}
			qp = append(qp, &newq)

			or = false
			not = false
			from = i + 1
			switch cond {
			case "or":
				or = true
			case "not":
				not = true
			default:
			}
		}
	}
	kv := strings.Join(words[from:], " ")
	k, v := getKeyVal(kv)
	if k == "" {
		err = fmt.Errorf("No key found in: %v", kv)
		return
	}
	newq := QParam{Or: or, Not: not, Key: k, Val: v}
	qp = append(qp, &newq)
	return
}
