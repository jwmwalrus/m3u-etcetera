package audiobookspane

import (
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/gtk/playlists"
)

var (
	audiobooksQueueSignals *playlists.OnQueue
)

func Setup(signals *map[string]interface{}) (err error) {
	audiobooksQueueSignals, err = playlists.CreateQueue(m3uetcpb.Perspective_AUDIOBOOKS, "audiobooks_queue_view", "audiobooks_queue_view_context")
	if err != nil {
		return
	}
	(*signals)["on_audiobooks_queue_sel_changed"] = audiobooksQueueSignals.SelChanged
	(*signals)["on_audiobooks_queue_view_row_activated"] = audiobooksQueueSignals.DblClicked
	(*signals)["on_audiobooks_queue_view_button_press_event"] = audiobooksQueueSignals.Context
	(*signals)["on_audiobooks_queue_view_context_play_now_activate"] = audiobooksQueueSignals.ContextPlayNow
	(*signals)["on_audiobooks_queue_view_context_enqueue_activate"] = audiobooksQueueSignals.ContextEnqueue
	(*signals)["on_audiobooks_queue_view_context_top_activate"] = audiobooksQueueSignals.ContextMove
	(*signals)["on_audiobooks_queue_view_context_up_activate"] = audiobooksQueueSignals.ContextMove
	(*signals)["on_audiobooks_queue_view_context_down_activate"] = audiobooksQueueSignals.ContextMove
	(*signals)["on_audiobooks_queue_view_context_bottom_activate"] = audiobooksQueueSignals.ContextMove
	(*signals)["on_audiobooks_queue_view_context_delete_activate"] = audiobooksQueueSignals.ContextDelete
	(*signals)["on_audiobooks_queue_view_context_clear_activate"] = audiobooksQueueSignals.ContextClear
	return
}
