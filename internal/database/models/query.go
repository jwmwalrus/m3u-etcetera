package models

import (
	"net/url"
	"strconv"
	"time"
)

// Query Defines a query
type Query struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	Random    bool      `json:"random"` // query allows random results
	Rating    int       `json:"rating"` // minimum rating to consider, from 1 to 10
	Limit     int       `json:"limit"`  // maximum number of tracks permitted
	Genre     string    `json:"genre"`
	Pattern   string    `json:"pattern"` // string to look for in track's indexed columns
	From      time.Time `json:"from"`    // from datetime in range
	To        time.Time `json:"to"`      // to datetime in range
	CreatedAt int64     `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt int64     `json:"updatedAt" gorm:"autoUpdateTime"`
}

func (q *Query) String() (str string) {
	params := url.Values{}

	if q.Random {
		params.Add("random", "1")
	}
	if q.Rating > 0 {
		params.Add("rating", strconv.Itoa(q.Rating))
	}
	if q.Limit > 0 {
		params.Add("limit", strconv.Itoa(q.Limit))
	}
	if q.Genre != "" {
		params.Add("genre", q.Genre)
	}
	if q.Pattern != "" {
		params.Add("pattern", q.Pattern)
	}
	if !q.From.IsZero() {
		params.Add("from", q.From.Format(time.RFC3339))
	}
	if !q.To.IsZero() {
		params.Add("to", q.To.Format(time.RFC3339))
	}

	str = params.Encode()

	return
}

// Populate parses a query string and populates the struct
func (q *Query) Populate(str string) (err error) {
	v, err := url.ParseQuery(str)
	if err != nil {
		return
	}

	fields := []string{"random", "rating", "limit", "genre", "pattern", "from", "to"}

	for _, f := range fields {
		fv := v.Get(f)
		if fv == "" {
			continue
		}
		var i int64
		switch f {
		case "random":
			q.Random = fv == "true"
		case "rating":
			i, err = strconv.ParseInt(fv, 10, 64)
			if err != nil {
				return
			}
			q.Rating = int(i)
		case "limit":
			i, err = strconv.ParseInt(fv, 10, 64)
			if err != nil {
				return
			}
			q.Limit = int(i)
		case "genre":
			q.Genre = fv
		case "pattern":
			q.Pattern = fv
		case "from":
			// TODO: parse datetime
		case "to":
			// TODO: parse datetime
		default:
		}
	}

	return
}

// ReplaceWith replaces the values of query q by those of query r
func (q *Query) ReplaceWith(r *Query, save bool) (err error) {

	q.Random = r.Random
	q.Rating = r.Rating
	q.Limit = r.Limit
	q.Genre = r.Genre
	q.Pattern = r.Pattern
	q.From = r.From
	q.To = r.To

	if save {
		err = db.Save(q).Error
	}

	return
}

// GetQueryStringByID returns the query string associated with the given id
func GetQueryStringByID(id int64) (str string, err error) {
	q := Query{}

	if err = db.Find(&q, id).Error; err != nil {
		return
	}

	str = q.String()
	return
}
