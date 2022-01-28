package podcastspane

import (
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/playlists"
)

var (
	podcastsQueueSignals *playlists.OnQueue
)

func Setup(signals *map[string]interface{}) (err error) {

	podcastsQueueSignals, err = playlists.CreateQueue(m3uetcpb.Perspective_PODCASTS, "podcasts_queue_view", "podcasts_queue_view_context")
	if err != nil {
		return
	}
	(*signals)["on_podcasts_queue_sel_changed"] = podcastsQueueSignals.SelChanged
	(*signals)["on_podcasts_queue_view_row_activated"] = podcastsQueueSignals.DblClicked
	(*signals)["on_podcasts_queue_view_button_press_event"] = podcastsQueueSignals.Context
	(*signals)["on_podcasts_queue_view_context_play_now_activate"] = podcastsQueueSignals.ContextPlayNow
	(*signals)["on_podcasts_queue_view_context_enqueue_activate"] = podcastsQueueSignals.ContextEnqueue
	(*signals)["on_podcasts_queue_view_context_top_activate"] = podcastsQueueSignals.ContextMove
	(*signals)["on_podcasts_queue_view_context_up_activate"] = podcastsQueueSignals.ContextMove
	(*signals)["on_podcasts_queue_view_context_down_activate"] = podcastsQueueSignals.ContextMove
	(*signals)["on_podcasts_queue_view_context_bottom_activate"] = podcastsQueueSignals.ContextMove
	(*signals)["on_podcasts_queue_view_context_delete_activate"] = podcastsQueueSignals.ContextDelete
	(*signals)["on_podcasts_queue_view_context_clear_activate"] = podcastsQueueSignals.ContextClear
	return
}
