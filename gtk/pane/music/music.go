package musicpane

import (
	"log/slog"

	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/gtk/playlists"
)

var (
	musicCollectionsSignals *onMusicCollections
	musicQuerySignals       *onMusicQuery
	musicPlaylistSignals    *onMusicPlaylist

	musicQueueSignals *playlists.OnQueue
)

// Setup sets the music pane.
func Setup(signals *builder.Signals) (err error) {
	slog.Info("Setting up music")

	musicCollectionsSignals, err = createMusicCollections()
	if err != nil {
		return
	}
	(*signals).AddDetail(
		"collections_sel",
		"changed",
		musicCollectionsSignals.selChanged,
	)
	(*signals).AddDetail(
		"collections_view",
		"button-press-event",
		musicCollectionsSignals.context,
	)
	(*signals).AddDetail(
		"collections_view_context_append",
		"activate",
		musicCollectionsSignals.contextAppend,
	)
	(*signals).AddDetail(
		"collections_view_context_prepend",
		"activate",
		musicCollectionsSignals.contextAppend,
	)
	(*signals).AddDetail(
		"collections_view_context_play_now",
		"activate",
		musicCollectionsSignals.contextPlayNow,
	)
	(*signals).AddDetail(
		"collections_view",
		"row-activated",
		musicCollectionsSignals.dblClicked,
	)
	(*signals).AddDetail(
		"collections_filter",
		"search-changed",
		musicCollectionsSignals.filtered,
	)
	(*signals).AddDetail(
		"collections_hierarchy",
		"changed",
		musicCollectionsSignals.hierarchyChanged,
	)
	(*signals).AddDetail(
		"collections_hierarchy_grouped",
		"toggled",
		musicCollectionsSignals.hierarchyGroupToggled,
	)

	musicQueueSignals, err = playlists.CreateQueue(m3uetcpb.Perspective_MUSIC, "music_queue_view", "music_queue_view_context")
	if err != nil {
		return
	}
	(*signals).AddDetail(
		"music_queue_sel",
		"changed",
		musicQueueSignals.SelChanged,
	)
	(*signals).AddDetail(
		"music_queue_view",
		"row-activated",
		musicQueueSignals.DblClicked,
	)
	(*signals).AddDetail(
		"music_queue_view",
		"button-press-event",
		musicQueueSignals.Context,
	)
	(*signals).AddDetail(
		"music_queue_view",
		"key-press-event",
		musicQueueSignals.Key,
	)
	(*signals).AddDetail(
		"music_queue_view_context_play_now",
		"activate",
		musicQueueSignals.ContextPlayNow,
	)
	(*signals).AddDetail(
		"music_queue_view_context_enqueue",
		"activate",
		musicQueueSignals.ContextEnqueue,
	)
	(*signals).AddDetail(
		"music_queue_view_context_top",
		"activate",
		musicQueueSignals.ContextMove,
	)
	(*signals).AddDetail(
		"music_queue_view_context_up",
		"activate",
		musicQueueSignals.ContextMove,
	)
	(*signals).AddDetail(
		"music_queue_view_context_down",
		"activate",
		musicQueueSignals.ContextMove,
	)
	(*signals).AddDetail(
		"music_queue_view_context_bottom",
		"activate",
		musicQueueSignals.ContextMove,
	)
	(*signals).AddDetail(
		"music_queue_view_context_delete",
		"activate",
		musicQueueSignals.ContextDelete,
	)
	(*signals).AddDetail(
		"music_queue_view_context_clear",
		"activate",
		musicQueueSignals.ContextClear,
	)
	(*signals).AddDetail(
		"music_queue_view_context",
		"popped-up",
		musicQueueSignals.ContextPoppedUp,
	)
	(*signals).AddDetail(
		"music_queue_view_context",
		"hide",
		musicQueueSignals.ContextHide,
	)

	musicQuerySignals, err = createMusicQueries()
	if err != nil {
		return
	}
	(*signals).AddDetail(
		"queries_sel",
		"changed",
		musicQuerySignals.selChanged,
	)
	(*signals).AddDetail(
		"queries_view",
		"button-press-event",
		musicQuerySignals.context,
	)
	(*signals).AddDetail(
		"queries_view_context_edit",
		"activate",
		musicQuerySignals.contextEdit,
	)
	(*signals).AddDetail(
		"queries_view_context_to_queue",
		"activate",
		musicQuerySignals.contextAppend,
	)
	(*signals).AddDetail(
		"queries_view_context_to_playlist",
		"activate",
		musicQuerySignals.contextAppend,
	)
	(*signals).AddDetail(
		"queries_view_context_new_playlist",
		"activate",
		musicQuerySignals.contextNewPlaylist,
	)
	(*signals).AddDetail(
		"queries_view_context_delete",
		"activate",
		musicQuerySignals.contextDelete,
	)
	(*signals).AddDetail(
		"queries_view",
		"row-activated",
		musicQuerySignals.dblClicked,
	)
	(*signals).AddDetail(
		"queries_filter",
		"search-changed",
		musicQuerySignals.filtered,
	)

	if err = musicQuerySignals.createDialog(); err != nil {
		return
	}
	(*signals).AddDetail(
		"queries_add",
		"clicked",
		musicQuerySignals.defineQuery,
	)
	(*signals).AddDetail(
		"query_dialog_search",
		"clicked",
		musicQuerySignals.doSearch,
	)
	(*signals).AddDetail(
		"query_dialog_toggle_selection",
		"clicked",
		musicQuerySignals.toggleSelection,
	)

	musicPlaylistSignals, err = createMusicPlaylists()
	if err != nil {
		return
	}
	(*signals).AddDetail(
		"music_playlists_filter",
		"search-changed",
		musicPlaylistSignals.filtered,
	)
	(*signals).AddDetail(
		"music_playlists_view",
		"row-activated",
		musicPlaylistSignals.dblClicked,
	)
	(*signals).AddDetail(
		"music_playlists_view",
		"button-press-event",
		musicPlaylistSignals.context,
	)
	(*signals).AddDetail(
		"music_playlists_sel",
		"changed",
		musicPlaylistSignals.selChanged,
	)
	(*signals).AddDetail(
		"music_playlists_view_context_delete",
		"activate",
		musicPlaylistSignals.contextDelete,
	)
	(*signals).AddDetail(
		"music_playlists_view_context_edit",
		"activate",
		musicPlaylistSignals.contextEdit,
	)
	(*signals).AddDetail(
		"music_playlists_view_context_open",
		"activate",
		musicPlaylistSignals.contextOpen,
	)
	(*signals).AddDetail(
		"music_playlists_view_context_export",
		"activate",
		musicPlaylistSignals.contextExport,
	)
	return
}
