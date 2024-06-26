package models

import (
	"context"
	"log/slog"
	"time"

	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/gear-pieces/idler"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	rtc "github.com/jwmwalrus/rtcycler"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// QueueTrack defines a track in the queue.
type QueueTrack struct { // too transient
	Model
	Position int    `json:"position"`
	Played   bool   `json:"played"`
	Location string `json:"location" gorm:"not null"`
	TrackID  int64  `json:"trackId" gorm:"index:idx_queue_track_track_id"`
	QueueID  int64  `json:"queueId" gorm:"index:idx_queue_track_queue_id,not null"`
	Queue    Queue  `json:"queue" gorm:"foreignKey:QueueID"`
}

func (qt *QueueTrack) Create() error {
	return qt.CreateTx(db)
}

func (qt *QueueTrack) CreateTx(tx *gorm.DB) error {
	return tx.Create(qt).Error
}

func (qt *QueueTrack) Read(id int64) error {
	return qt.ReadTx(db, id)
}

func (qt *QueueTrack) ReadTx(tx *gorm.DB, id int64) error {
	return tx.First(qt, id).Error
}

func (qt *QueueTrack) Save() error {
	return qt.SaveTx(db)
}

func (qt *QueueTrack) SaveTx(tx *gorm.DB) error {
	return tx.Save(qt).Error
}

func (qt *QueueTrack) ToProtobuf() proto.Message {
	q := Queue{}
	onerror.Log(q.Read(qt.QueueID))

	return &m3uetcpb.QueueTrack{
		Id:       qt.ID,
		Position: int32(qt.Position),
		Played:   qt.Played, Location: qt.Location,
		Perspective: m3uetcpb.Perspective(q.Perspective.Idx),
		TrackId:     qt.TrackID,
		CreatedAt:   timestamppb.New(time.Unix(0, qt.CreatedAt)),
		UpdatedAt:   timestamppb.New(time.Unix(0, qt.UpdatedAt)),
	}
}

// AfterCreate is a GORM hook.
func (qt *QueueTrack) AfterCreate(tx *gorm.DB) error {
	go func() {
		if !rtc.FlagTestMode() &&
			!idler.IsAppBusyBy(idler.StatusEngineLoop) {
			TriggerPlaybackChange()
		}
	}()
	return nil
}

// GetPosition implements the Poser interface.
func (qt *QueueTrack) GetPosition() int {
	return qt.Position
}

// SetPosition implements the Poser interface.
func (qt *QueueTrack) SetPosition(pos int) {
	qt.Position = pos
}

// GetIgnore implements the Poser interface.
func (qt *QueueTrack) GetIgnore() bool {
	return qt.Played
}

// SetIgnore implements the Poser interface.
func (qt *QueueTrack) SetIgnore(ignore bool) {
	qt.Played = ignore
}

// GetAllQueueTracks returns all queue tracks for the given perspective,
// constrained by a limit.
func GetAllQueueTracks(idx PerspectiveIndex, limit int) (qts []*QueueTrack, ts []*Track) {
	logw := slog.With(
		"idx", idx,
		"limit", limit,
	)
	logw.Info("Getting all queue tracks")

	q, err := idx.GetPerspectiveQueue()
	if err != nil {
		logw.Error("Failed to get perspective queue", "error", err)
		return
	}

	tx := db.
		Where("played = 0 AND queue_id = ?", q.ID).
		Order("position ASC")

	if limit > 0 {
		tx.Limit(limit)
	}

	qts = []*QueueTrack{}
	ts = []*Track{}
	s := []QueueTrack{}
	if err = tx.Find(&s).Error; err != nil {
		logw.Error("Failed to find queue tracks in database", "error", err)
		return
	}

	ids := []int64{}
	locations := []string{}
	for i := range s {
		if s[i].TrackID > 0 {
			ids = append(ids, s[i].TrackID)
		} else {
			locations = append(locations, s[i].Location)
		}
		qts = append(qts, &s[i])
	}

	ts, _ = FindTracksIn(ids)
	for _, l := range locations {
		var t *Track
		if t, err = ReadTagsForLocation(l); err != nil {
			continue
		}
		ts = append(ts, t)
	}
	return
}

// GetQueueStore returns all queue tracks for all perspectives.
func GetQueueStore() (qs []*QueueTrack, ts []*Track, dig []*PerspectiveDigest) {
	slog.Info("Getting queue store")

	dig = []*PerspectiveDigest{}

	for _, idx := range PerspectiveIndexList() {
		qsaux, tsaux := GetAllQueueTracks(idx, 0)
		qs = append(qs, qsaux...)
		pd := PerspectiveDigest{}
		pd.Idx = idx
		for i := range tsaux {
			ts = append(ts, tsaux[i])
			pd.Duration += tsaux[i].Duration
		}
		dig = append(dig, &pd)
	}
	return
}

func findQueueTrack(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case id := <-queueTrackNeeded:
			go func(id int64) {
				qt := QueueTrack{}
				err := qt.Read(id)
				if err != nil {
					slog.With(
						"id", id,
						"error", err,
					).Error("Failed to read queue track")
					return
				}
				if qt.TrackID > 0 {
					return
				}

				logw := slog.With("qt", qt)
				logw.Info("Finding track for queue entry")

				t := Track{}
				if err := db.Where("location = ?", qt.Location).First(&t).Error; err != nil {
					return
				}
				qt.TrackID = t.ID
				onerror.NewRecorder(logw).Log(qt.Save())
			}(id)
		}
	}
}
