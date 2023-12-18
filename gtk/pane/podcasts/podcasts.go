package podcastspane

import (
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/gtk/playlists"
)

var (
	podcastsQueueSignals *playlists.OnQueue
)

// Setup sets the podcasts pane.
func Setup(signals *builder.Signals) (err error) {

	podcastsQueueSignals, err = playlists.CreateQueue(
		m3uetcpb.Perspective_PODCASTS,
		"podcasts_queue_view",
		"podcasts_queue_view_context",
	)
	if err != nil {
		return
	}

	(*signals).AddDetail(
		"podcasts_queue_sel",
		"changed",
		podcastsQueueSignals.SelChanged,
	)
	(*signals).AddDetail(
		"podcasts_queue_view",
		"row-activated",
		podcastsQueueSignals.DblClicked,
	)
	(*signals).AddDetail(
		"podcasts_queue_view",
		"button-press-event",
		podcastsQueueSignals.Context,
	)
	(*signals).AddDetail(
		"podcasts_queue_view_context_play_now",
		"activate",
		podcastsQueueSignals.ContextPlayNow,
	)
	(*signals).AddDetail(
		"podcasts_queue_view_context_enqueue",
		"activate",
		podcastsQueueSignals.ContextEnqueue,
	)
	(*signals).AddDetail(
		"podcasts_queue_view_context_top",
		"activate",
		podcastsQueueSignals.ContextMove,
	)
	(*signals).AddDetail(
		"podcasts_queue_view_context_up",
		"activate",
		podcastsQueueSignals.ContextMove,
	)
	(*signals).AddDetail(
		"podcasts_queue_view_context_down",
		"activate",
		podcastsQueueSignals.ContextMove,
	)
	(*signals).AddDetail(
		"podcasts_queue_view_context_bottom",
		"activate",
		podcastsQueueSignals.ContextMove,
	)
	(*signals).AddDetail(
		"podcasts_queue_view_context_delete",
		"activate",
		podcastsQueueSignals.ContextDelete,
	)
	(*signals).AddDetail(
		"podcasts_queue_view_context_clear",
		"activate",
		podcastsQueueSignals.ContextClear,
	)
	return
}
