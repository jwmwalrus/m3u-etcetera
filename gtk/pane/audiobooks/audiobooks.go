package audiobookspane

import (
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/gtk/playlists"
)

var (
	audiobooksQueueSignals *playlists.OnQueue
)

// Setup sets the audiobooks pane.
func Setup(signals *builder.Signals) (err error) {
	audiobooksQueueSignals, err = playlists.CreateQueue(
		m3uetcpb.Perspective_AUDIOBOOKS,
		"audiobooks_queue_view",
		"audiobooks_queue_view_context",
	)
	if err != nil {
		return
	}
	(*signals).AddDetail(
		"audiobooks_queue_sel",
		"changed",
		audiobooksQueueSignals.SelChanged,
	)
	(*signals).AddDetail(
		"audiobooks_queue_view",
		"row_activated",
		audiobooksQueueSignals.DblClicked,
	)
	(*signals).AddDetail(
		"audiobooks_queue_view",
		"button_press_event",
		audiobooksQueueSignals.Context,
	)
	(*signals).AddDetail(
		"audiobooks_queue_view_context_play_now",
		"activate",
		audiobooksQueueSignals.ContextPlayNow,
	)
	(*signals).AddDetail(
		"audiobooks_queue_view_context_enqueue",
		"activate",
		audiobooksQueueSignals.ContextEnqueue,
	)
	(*signals).AddDetail(
		"audiobooks_queue_view_context_top",
		"activate",
		audiobooksQueueSignals.ContextMove,
	)
	(*signals).AddDetail(
		"audiobooks_queue_view_context_up",
		"activate",
		audiobooksQueueSignals.ContextMove,
	)
	(*signals).AddDetail(
		"audiobooks_queue_view_context_down",
		"activate",
		audiobooksQueueSignals.ContextMove,
	)
	(*signals).AddDetail(
		"audiobooks_queue_view_context_bottom",
		"activate",
		audiobooksQueueSignals.ContextMove,
	)
	(*signals).AddDetail(
		"audiobooks_queue_view_context_delete",
		"activate",
		audiobooksQueueSignals.ContextDelete,
	)
	(*signals).AddDetail(
		"audiobooks_queue_view_context_clear",
		"activate",
		audiobooksQueueSignals.ContextClear,
	)
	return
}
