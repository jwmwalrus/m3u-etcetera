package models

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/jwmwalrus/bnp/pointers"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/m3u-etcetera/internal/config"
	"github.com/jwmwalrus/m3u-etcetera/internal/subscription"
	"github.com/jwmwalrus/m3u-etcetera/pkg/qparams"
	"github.com/jwmwalrus/onerror"
	rtc "github.com/jwmwalrus/rtcycler"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

// QueryIndex defines indexes for collections.
type QueryIndex int

const (
	// HistoryQuery for the playback history.
	HistoryQuery QueryIndex = iota + 1

	// TopTracksQuery for the top 100 tracks in playback history.
	TopTracksQuery
)

func (idx QueryIndex) String() string {
	return [...]string{"", "\t", "\t\t"}[idx]
}

func (idx QueryIndex) Description() string {
	return [...]string{"", "Playback History", "Playback Top Tracks"}[idx]
}

// Get returns the query associated to the index.
func (idx QueryIndex) Get() (qy *Query, err error) {
	qy = &Query{}
	err = db.Where("idx = ?", int(idx)).First(qy).Error
	return
}

// QueryEvent defines a query event.
type QueryEvent int

// QueryEvent enum.
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

// QueryBoundaryTx defines the transactional query boundary interface.
type QueryBoundaryTx interface {
	FindTracksTx(*gorm.DB) []*Track
	DataDeleterTx
	DataUpdaterTx
}

// QueryBoundaryID defines the query boundary ID interface.
type QueryBoundaryID interface {
	GetQueryID() int64
}

var supportedParams = []string{
	"id",
	"title",
	"artist",
	"album",
	"albumartist",
	"composer",
	"genre",
	"year",
	"date",
	"rating",
}

// CountSupportedParams returns the count of supported parameters in a slice.
func CountSupportedParams(qp []qparams.QParam) (n int) {
	for _, x := range qp {
		if !slices.Contains(supportedParams, x.Key) {
			continue
		}
		n++
	}
	return
}

// Query Defines a query.
type Query struct {
	ID          int64  `json:"id" gorm:"primaryKey"`
	Idx         int    `json:"idx" gorm:"not null,default:0"`
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

// Read implements DataReader interface.
func (qy *Query) Read(id int64) error {
	return db.First(qy, id).Error
}

// Save implements DataSaver interface.
func (qy *Query) Save() error {
	qy.ProvideName()
	return db.Save(qy).Error
}

// FromProtobuf implements ProtoIn interface.
func (qy *Query) FromProtobuf(in proto.Message) {
	protobufToQuery(in.(*m3uetcpb.Query), qy)
}

// ToProtobuf implements ProtoOut interface.
func (qy *Query) ToProtobuf() proto.Message {
	bv, err := json.Marshal(qy)
	if err != nil {
		log.Error(err)
		return &m3uetcpb.Query{}
	}

	out := &m3uetcpb.Query{}
	err = jsonUnmarshaler.Unmarshal(bv, out)
	onerror.Log(err)

	// Unmatched
	out.ReadOnly = qy.IsReadOnly()

	cqs := qy.GetCollections()
	for _, x := range cqs {
		out.CollectionIds = append(out.CollectionIds, x.CollectionID)
	}

	return out
}

// AfterCreate is a GORM hook.
func (qy *Query) AfterCreate(tx *gorm.DB) error {
	go func() {
		if rtc.FlagTestMode() {
			return
		}
		subscription.Broadcast(
			subscription.ToQueryStoreEvent,
			subscription.Event{
				Idx:  int(QueryEventItemAdded),
				Data: qy,
			},
		)
	}()
	return nil
}

// AfterUpdate is a GORM hook.
func (qy *Query) AfterUpdate(tx *gorm.DB) error {
	go func() {
		if rtc.FlagTestMode() {
			return
		}
		subscription.Broadcast(
			subscription.ToQueryStoreEvent,
			subscription.Event{
				Idx:  int(QueryEventItemChanged),
				Data: qy,
			},
		)
	}()
	return nil
}

// AfterDelete is a GORM hook.
func (qy *Query) AfterDelete(tx *gorm.DB) error {
	go func() {
		if rtc.FlagTestMode() {
			return
		}
		subscription.Broadcast(
			subscription.ToQueryStoreEvent,
			subscription.Event{
				Idx:  int(QueryEventItemRemoved),
				Data: qy,
			},
		)
	}()
	return nil
}

// Delete deletes a query from the DB.
func (qy *Query) Delete(qybs ...QueryBoundaryTx) error {
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

// FindTracks return the list of tracks that match the query.
func (qy *Query) FindTracks(qybs []QueryBoundaryTx) (ts []*Track) {
	entry := log.WithFields(log.Fields{
		"qy":        qy,
		"len(qybs)": len(qybs),
	})
	entry.Info("Finding tracks")

	switch QueryIndex(qy.Idx) {
	case HistoryQuery:
		return findUniqueHistoryTracks()
	case TopTracksQuery:
		return findTopTracks()
	default:
	}

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
				if !slices.Contains(supportedParams, strings.ToLower(x.Key)) {
					entry.Warnf("Ignored query paranmeter: %v", x.Key)
					continue
				}
				comp := " LIKE ?"
				if x.Key == "id" {
					if _, err := strconv.ParseInt(x.Val, 10, 64); err != nil {
						entry.Warnf("Ignoring `id` value due to parsing error:%v", x.Val)
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
		ts = append(ts, qybs[i].FindTracksTx(tx)...)
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

// GetCollections returns all the collections associated to the given query.
// This adds forward support for CollectionQuery, required for ToProtobuf.
func (qy *Query) GetCollections() []*CollectionQuery {
	cqs := []CollectionQuery{}

	err := db.Where("query_id = ?", qy.ID).Find(&cqs).Error
	if err != nil {
		log.Error(err)
		return []*CollectionQuery{}
	}
	return pointers.FromSlice(cqs)
}

// IsReadOnly returns true if the query is read-only.
func (qy *Query) IsReadOnly() bool {
	return qy.Idx > 0
}

// ProvideName provides a name before saving a query.
func (qy *Query) ProvideName() bool {
	if qy.Name == "" {
		qy.Name = "Query from " + time.Now().Format(time.RFC3339)
		return true
	}
	return false
}

// SaveBound persists a query in the DB and bounds it to a list of collections.
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

// FromProtobuf returns a Query type populated from the given m3uetcpb.Query.
func FromProtobuf(in *m3uetcpb.Query) (qy *Query) {
	qy = &Query{}
	protobufToQuery(in, qy)
	return
}

// GetAllQueries returns all queries, constrained by a limit and collections.
func GetAllQueries(limit int, qybs ...QueryBoundaryID) (s []*Query) {
	entry := log.WithFields(log.Fields{
		"limit":     limit,
		"len(qybs)": len(qybs),
	})
	entry.Info("Getting all queries")

	s = []*Query{}

	qys := []Query{}
	if err := db.Find(&qys).Error; err != nil {
		entry.Error(err)
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

// SupportedParams returns the list of supported string parameters.
func SupportedParams() []string {
	return supportedParams
}

func findHistoryTracks() (ts []*Track, lpf []int64) {
	log.Info("Finding history tracks")

	limit := config.DefaultQueryMaxLimit
	if base.Conf.Server.Query.Limit > 0 {
		limit = base.Conf.Server.Query.Limit
	}

	tx := db.Session(&gorm.Session{SkipHooks: true})
	tx = tx.Limit(limit)

	ts = []*Track{}
	list := []PlaybackHistory{}

	err := tx.Order("created_at DESC").
		Find(&list).
		Error
	onerror.Log(err)

	for _, h := range list {
		var err error
		t := &Track{}
		if h.TrackID > 0 {
			if err = t.Read(h.TrackID); err != nil {
				log.Error(err)
				continue
			}
		} else {
			tx := db.Session(&gorm.Session{SkipHooks: true})
			err = tx.Where("location = ?", h.Location).First(t).Error
			if err != nil {
				t, err = ReadTagsForLocation(h.Location)
				if err != nil {
					log.Error(err)
					continue
				}
				t.Lastplayed = h.CreatedAt
				err = t.createTransient(tx, nil)
				if err != nil {
					log.Warn(err)
					continue
				}
			}
		}
		lpf = append(lpf, h.Duration)
		ts = append(ts, t)
	}

	return
}

func findTopTracks() (ts []*Track) {
	log.Info("Finding top tracks")

	tx := db.Session(&gorm.Session{SkipHooks: true})
	tx = tx.Limit(100)

	list := []Track{}

	err := tx.Where("playcount > 0").
		Order("playcount DESC, updated_at DESC").
		Find(&list).
		Error
	onerror.Log(err)

	ts = pointers.FromSlice(list)
	return

}

func findUniqueHistoryTracks() (ts []*Track) {
	list, _ := findHistoryTracks()

	unique := make(map[int64]*Track)
	for i := range list {
		_, ok := unique[list[i].ID]
		if ok {
			continue
		}
		unique[list[i].ID] = list[i]
	}

	ts = []*Track{}
	for i := range unique {
		ts = append(ts, unique[i])
	}
	return
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
