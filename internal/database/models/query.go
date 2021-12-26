package models

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/bnp/slice"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/m3u-etcetera/pkg/config"
	"github.com/jwmwalrus/m3u-etcetera/pkg/qparams"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

// QueryBoundaryTx defines the transactional query boundary interface
type QueryBoundaryTx interface {
	FindTracksTx(*gorm.DB) []*Track
	DataDeleterTx
	DataUpdaterTx
}

// QueryBoundaryID defines the query boundary ID interface
type QueryBoundaryID interface {
	GetQueryID() int64
}

var supportedParams = []string{
	"id",
	"title",
	"artist",
	"album",
	"albumartist",
	"genre",
}

// CountSupportedParams returns the count of supported parameters in a slice
func CountSupportedParams(qp []qparams.QParam) (n int) {
	for _, x := range qp {
		if !slice.Contains(supportedParams, x.Key) {
			continue
		}
		n++
	}
	return
}

// Query Defines a query
type Query struct {
	ID          int64  `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`        // query name
	Description string `json:"description"` // query description
	Random      bool   `json:"random"`      // query allows random results
	Rating      int    `json:"rating"`      // minimum rating to consider, from 1 to 10
	Limit       int    `json:"limit"`       // maximum number of tracks permitted
	Params      string `json:"params"`      // patterns to look for in track's indexed columns
	From        int64  `json:"from"`        // from datetime in range
	To          int64  `json:"to"`          // to datetime in range
	CreatedAt   int64  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   int64  `json:"updatedAt" gorm:"autoUpdateTime"`
}

// Read implements DataReader interface
func (q *Query) Read(id int64) error {
	return db.First(q, id).Error
}

// Save implements DataSaver interface
func (q *Query) Save() error {
	return db.Save(q).Error
}

// FromProtobuf implements ProtoIn interface
func (q *Query) FromProtobuf(in proto.Message) {
	protobufToQuery(in.(*m3uetcpb.Query), q)
	return
}

// ToProtobuf implements ProtoOut interface
func (q *Query) ToProtobuf() proto.Message {
	bv, err := json.Marshal(q)
	if err != nil {
		log.Error(err)
		return &m3uetcpb.Query{}
	}

	out := &m3uetcpb.Query{}
	err = json.Unmarshal(bv, out)
	onerror.Log(err)

	// Unmatched
	out.CreatedAt = q.CreatedAt
	out.UpdatedAt = q.UpdatedAt

	cqs := q.GetCollections()
	for _, x := range cqs {
		out.CollectionIds = append(out.CollectionIds, x.CollectionID)
	}
	return out
}

// Delete deletes a query from the DB
func (q *Query) Delete(qbs ...QueryBoundaryTx) (err error) {
	return db.Transaction(func(tx *gorm.DB) error {
		for _, b := range qbs {
			if err := b.DeleteTx(tx); err != nil {
				return err
			}
		}
		if err := tx.Delete(q).Error; err != nil {
			return err
		}
		return nil
	})
}

// FindTracks return the list of tracks that match the query
func (q *Query) FindTracks(qbs []QueryBoundaryTx) (ts []*Track) {
	log.WithFields(log.Fields{
		"q":        q,
		"len(qbs)": len(qbs),
	}).
		Info("Finding tracks")

	limit := config.DefaultQueryMaxLimit
	if base.Conf.Server.Query.Limit > 0 {
		limit = base.Conf.Server.Query.Limit
	}
	if q.Limit > 0 {
		limit = q.Limit
	}

	buildStmt := func() *gorm.DB {
		tx := db.Limit(limit)

		if q.Params != "" {
			parsed, _ := qparams.ParseParams(q.Params)
			for _, x := range parsed {
				if !slice.Contains(supportedParams, strings.ToLower(x.Key)) {
					log.Warnf("Ignored query paranmeter: %v", x.Key)
					continue
				}
				comp := " LIKE ?"
				if x.Key == "id" {
					if _, err := strconv.ParseInt(x.Val, 10, 64); err != nil {
						log.Warnf("Ignoring `id` value due to parsing error:%v", x.Val)
						continue
					}
					comp = " = ?"
				}
				y := x.ToFuzzy().ToSQL()
				if y.Or {
					tx.Or("track."+y.Key+comp, y.Val)
				} else if y.Not {
					tx.Not("track."+y.Key+comp, y.Val)
				} else {
					tx.Where("track."+y.Key+comp, y.Val)
				}
			}
		}

		if q.Rating > 0 {
			tx.Where("rating >= ?", q.Rating)
		}

		if q.From > 0 {
			fromYear := time.Unix(q.From, 0).Year()
			tx.Where("year >= ?", fromYear)
		}
		if q.To > 0 {
			toYear := time.Unix(q.To, 0).Year()
			tx.Where("year <= ?", toYear)
		}

		if q.Random {
			tx.Order("random()")
		}
		return tx
	}

	ts = []*Track{}
	for i := range qbs {
		tx := buildStmt()
		list := qbs[i].FindTracksTx(tx)
		ts = appendToTrackList(ts, list)
	}

	if q.Random {
		shuff := getSuffler(len(ts))
		newts := make([]*Track, len(ts))
		for k, v := range shuff {
			newts[k] = ts[v]
		}
	}

	if len(ts) > limit {
		ts = ts[:limit]
	}

	return
}

// GetCollections returns all the collections associated to the given query
// This adds forward support for CollectionQuery, required for ToProtobuf
func (q *Query) GetCollections() (cqs []*CollectionQuery) {
	cqs = []*CollectionQuery{}

	list := []CollectionQuery{}
	err := db.Where("query_id = ?", q.ID).Find(&list).Error
	if err != nil {
		log.Error(err)
		return
	}

	for i := range list {
		cqs = append(cqs, &list[i])
	}
	return
}

// SaveBound persists a query in the DB and bounds it to a list of collections
func (q *Query) SaveBound(qbs []QueryBoundaryTx) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(q).Error; err != nil {
			return err
		}

		for _, b := range qbs {
			if err := b.SaveTx(tx); err != nil {
				return err
			}
		}
		return nil
	})
}

// FromProtobuf returns a Query type populated from the given m3uetcpb.Query
func FromProtobuf(in *m3uetcpb.Query) (q *Query) {
	q = &Query{}
	protobufToQuery(in, q)
	return
}

// GetAllQueries returns all queries, constrained by a limit and collections
func GetAllQueries(limit int, qbs ...QueryBoundaryID) (s []*Query) {
	log.WithFields(log.Fields{
		"limit":    limit,
		"len(qbs)": len(qbs),
	}).
		Info("Getting all queries")

	s = []*Query{}

	qs := []Query{}
	if err := db.Find(&qs).Error; err != nil {
		log.Error(err)
		return
	}

	for k, v := range qs {
		if len(qbs) > 0 {
			match := false
			for _, i := range qbs {
				if v.ID == i.GetQueryID() {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}
		s = append(s, &qs[k])
		if limit > 0 && len(s) == limit {
			break
		}
	}
	return
}

// SupportedParams returns the list of supported string parameters
func SupportedParams() []string {
	return supportedParams
}

func protobufToQuery(in *m3uetcpb.Query, out *Query) {
	out.Name = in.Name
	out.Description = in.Description
	out.Random = in.Random
	out.Rating = int(in.Rating)
	out.Limit = int(in.Limit)
	out.Params = in.Params
	out.From = in.From
	out.To = in.To
}
