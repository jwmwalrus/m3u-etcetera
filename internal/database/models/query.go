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
	"github.com/jwmwalrus/m3u-etcetera/internal/subscription"
	"github.com/jwmwalrus/m3u-etcetera/pkg/config"
	"github.com/jwmwalrus/m3u-etcetera/pkg/qparams"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

// QueryEvent defines a query event
type QueryEvent int

// QueryEvent enum
const (
	QueryEventNone QueryEvent = iota
	QueryEventInitial
	_
	_
	QueryEventItemAdded
	QueryEventItemChanged
	QueryEventItemRemoved
)

func (qye QueryEvent) String() string {
	return []string{
		"none",
		"initial",
		"initial-item",
		"initial-done",
		"query-added",
		"query-changed",
		"query-removed",
	}[qye]
}

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
	CreatedAt   int64  `json:"createdAt" gorm:"autoCreateTime:nano"`
	UpdatedAt   int64  `json:"updatedAt" gorm:"autoUpdateTime:nano"`
}

// Read implements DataReader interface
func (qy *Query) Read(id int64) error {
	return db.First(qy, id).Error
}

// Save implements DataSaver interface
func (qy *Query) Save() error {
	qy.ProvideName()
	return db.Save(qy).Error
}

// FromProtobuf implements ProtoIn interface
func (qy *Query) FromProtobuf(in proto.Message) {
	protobufToQuery(in.(*m3uetcpb.Query), qy)
	return
}

// ToProtobuf implements ProtoOut interface
func (qy *Query) ToProtobuf() proto.Message {
	bv, err := json.Marshal(qy)
	if err != nil {
		log.Error(err)
		return &m3uetcpb.Query{}
	}

	out := &m3uetcpb.Query{}
	err = json.Unmarshal(bv, out)
	onerror.Log(err)

	// Unmatched
	out.CreatedAt = qy.CreatedAt
	out.UpdatedAt = qy.UpdatedAt

	cqs := qy.GetCollections()
	for _, x := range cqs {
		out.CollectionIds = append(out.CollectionIds, x.CollectionID)
	}
	return out
}

// AfterCreate is a GORM hook
func (qy *Query) AfterCreate(tx *gorm.DB) error {
	go func() {
		if !base.FlagTestingMode && globalCollectionEvent != CollectionEventInitial {
			subscription.Broadcast(
				subscription.ToQueryStoreEvent,
				subscription.Event{
					Idx:  int(QueryEventItemAdded),
					Data: qy,
				},
			)
		}
	}()
	return nil
}

// AfterSave is a GORM hook
func (qy *Query) AfterSave(tx *gorm.DB) error {
	go func() {
		if !base.FlagTestingMode && globalCollectionEvent != CollectionEventInitial {
			subscription.Broadcast(
				subscription.ToQueryStoreEvent,
				subscription.Event{
					Idx:  int(QueryEventItemChanged),
					Data: qy,
				},
			)
		}
	}()
	return nil
}

// AfterDelete is a GORM hook
func (qy *Query) AfterDelete(tx *gorm.DB) error {
	go func() {
		if !base.FlagTestingMode && globalCollectionEvent != CollectionEventInitial {
			subscription.Broadcast(
				subscription.ToQueryStoreEvent,
				subscription.Event{
					Idx:  int(QueryEventItemRemoved),
					Data: qy,
				},
			)
		}
	}()
	return nil
}

// Delete deletes a query from the DB
func (qy *Query) Delete(qybs ...QueryBoundaryTx) (err error) {
	return db.Transaction(func(tx *gorm.DB) error {
		for _, b := range qybs {
			if err := b.DeleteTx(tx); err != nil {
				return err
			}
		}
		if err := tx.Delete(qy).Error; err != nil {
			return err
		}
		return nil
	})
}

// FindTracks return the list of tracks that match the query
func (qy *Query) FindTracks(qybs []QueryBoundaryTx) (ts []*Track) {
	log.WithFields(log.Fields{
		"qy":        qy,
		"len(qybs)": len(qybs),
	}).
		Info("Finding tracks")

	limit := config.DefaultQueryMaxLimit
	if base.Conf.Server.Query.Limit > 0 {
		limit = base.Conf.Server.Query.Limit
	}
	if qy.Limit > 0 {
		limit = qy.Limit
	}

	buildStmt := func() *gorm.DB {
		tx := db.Limit(limit)

		if qy.Params != "" {
			parsed, _ := qparams.ParseParams(qy.Params)
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

		if qy.Rating > 0 {
			tx.Where("rating >= ?", qy.Rating)
		}

		if qy.From > 0 {
			fromYear := time.Unix(qy.From, 0).Year()
			tx.Where("year >= ?", fromYear)
		}
		if qy.To > 0 {
			toYear := time.Unix(qy.To, 0).Year()
			tx.Where("year <= ?", toYear)
		}

		if qy.Random {
			tx.Order("random()")
		}
		return tx
	}

	ts = []*Track{}
	for i := range qybs {
		tx := buildStmt()
		list := qybs[i].FindTracksTx(tx)
		ts = appendToTrackList(ts, list)
	}

	if qy.Random {
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
func (qy *Query) GetCollections() (cqs []*CollectionQuery) {
	cqs = []*CollectionQuery{}

	list := []CollectionQuery{}
	err := db.Where("query_id = ?", qy.ID).Find(&list).Error
	if err != nil {
		log.Error(err)
		return
	}

	for i := range list {
		cqs = append(cqs, &list[i])
	}
	return
}

// ProvideName provides a name before saving a query
func (qy *Query) ProvideName() bool {
	if qy.Name == "" {
		qy.Name = "Query from " + time.Now().Format(time.RFC3339)
		return true
	}
	return false
}

// SaveBound persists a query in the DB and bounds it to a list of collections
func (qy *Query) SaveBound(qybs []QueryBoundaryTx) error {
	return db.Transaction(func(tx *gorm.DB) error {
		wasSet := qy.ProvideName()
		if err := tx.Save(qy).Error; err != nil {
			if wasSet {
				qy.Name = ""
			}
			return err
		}

		for _, b := range qybs {
			if err := b.SaveTx(tx); err != nil {
				return err
			}
		}
		return nil
	})
}

// FromProtobuf returns a Query type populated from the given m3uetcpb.Query
func FromProtobuf(in *m3uetcpb.Query) (qy *Query) {
	qy = &Query{}
	protobufToQuery(in, qy)
	return
}

// GetAllQueries returns all queries, constrained by a limit and collections
func GetAllQueries(limit int, qybs ...QueryBoundaryID) (s []*Query) {
	log.WithFields(log.Fields{
		"limit":     limit,
		"len(qybs)": len(qybs),
	}).
		Info("Getting all queries")

	s = []*Query{}

	qys := []Query{}
	if err := db.Find(&qys).Error; err != nil {
		log.Error(err)
		return
	}

	for k, v := range qys {
		if len(qybs) > 0 {
			match := false
			for _, i := range qybs {
				if v.ID == i.GetQueryID() {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}
		s = append(s, &qys[k])
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
