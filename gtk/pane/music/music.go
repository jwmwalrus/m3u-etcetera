package musicpane

import (
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/playlists"
	log "github.com/sirupsen/logrus"
)

var (
	musicCollectionsSignals *onMusicCollections
	musicQuerySignals       *onMusicQuery
	musicPlaylistSignals    *onMusicPlaylist

	musicQueueSignals *playlists.OnQueue
)

// Setup sets the music pane.
func Setup(signals *map[string]interface{}) (err error) {
	log.Info("Setting up music")

	musicCollectionsSignals, err = createMusicCollections()
	if err != nil {
		return
	}
	(*signals)["on_collections_sel_changed"] = musicCollectionsSignals.selChanged
	(*signals)["on_collections_view_button_press_event"] = musicCollectionsSignals.context
	(*signals)["on_collections_view_context_append_activate"] = musicCollectionsSignals.contextAppend
	(*signals)["on_collections_view_context_prepend_activate"] = musicCollectionsSignals.contextAppend
	(*signals)["on_collections_view_context_play_now_activate"] = musicCollectionsSignals.contextPlayNow
	(*signals)["on_collections_view_row_activated"] = musicCollectionsSignals.dblClicked
	(*signals)["on_collections_filter_search_changed"] = musicCollectionsSignals.filtered
	(*signals)["on_collections_hierarchy_changed"] = musicCollectionsSignals.hierarchyChanged
	(*signals)["on_collections_hierarchy_grouped_toggled"] = musicCollectionsSignals.hierarchyGroupToggled

	musicQueueSignals, err = playlists.CreateQueue(m3uetcpb.Perspective_MUSIC, "music_queue_view", "music_queue_view_context")
	if err != nil {
		return
	}
	(*signals)["on_music_queue_sel_changed"] = musicQueueSignals.SelChanged
	(*signals)["on_music_queue_view_row_activated"] = musicQueueSignals.DblClicked
	(*signals)["on_music_queue_view_button_press_event"] = musicQueueSignals.Context
	(*signals)["on_music_queue_view_key_press_event"] = musicQueueSignals.Key
	(*signals)["on_music_queue_view_context_play_now_activate"] = musicQueueSignals.ContextPlayNow
	(*signals)["on_music_queue_view_context_enqueue_activate"] = musicQueueSignals.ContextEnqueue
	(*signals)["on_music_queue_view_context_top_activate"] = musicQueueSignals.ContextMove
	(*signals)["on_music_queue_view_context_up_activate"] = musicQueueSignals.ContextMove
	(*signals)["on_music_queue_view_context_down_activate"] = musicQueueSignals.ContextMove
	(*signals)["on_music_queue_view_context_bottom_activate"] = musicQueueSignals.ContextMove
	(*signals)["on_music_queue_view_context_delete_activate"] = musicQueueSignals.ContextDelete
	(*signals)["on_music_queue_view_context_clear_activate"] = musicQueueSignals.ContextClear
	(*signals)["on_music_queue_view_context_popped_up"] = musicQueueSignals.ContextPoppedUp
	(*signals)["on_music_queue_view_context_hide"] = musicQueueSignals.ContextHide

	musicQuerySignals, err = createMusicQueries()
	if err != nil {
		return
	}
	(*signals)["on_queries_sel_changed"] = musicQuerySignals.selChanged
	(*signals)["on_queries_view_button_press_event"] = musicQuerySignals.context
	(*signals)["on_queries_view_context_edit_activate"] = musicQuerySignals.contextEdit
	(*signals)["on_queries_view_context_to_queue_activate"] = musicQuerySignals.contextAppend
	(*signals)["on_queries_view_context_to_playlist_activate"] = musicQuerySignals.contextAppend
	(*signals)["on_queries_view_context_new_playlist_activate"] = musicQuerySignals.contextNewPlaylist
	(*signals)["on_queries_view_context_delete_activate"] = musicQuerySignals.contextDelete
	(*signals)["on_queries_view_row_activated"] = musicQuerySignals.dblClicked
	(*signals)["on_queries_filter_search_changed"] = musicQuerySignals.filtered

	if err = musicQuerySignals.createDialog(); err != nil {
		return
	}
	(*signals)["on_queries_add_clicked"] = musicQuerySignals.defineQuery
	(*signals)["on_query_dialog_search_clicked"] = musicQuerySignals.doSearch
	(*signals)["on_query_dialog_toggle_selection_clicked"] = musicQuerySignals.toggleSelection

	musicPlaylistSignals, err = createMusicPlaylists()
	if err != nil {
		return
	}
	(*signals)["on_music_playlists_filter_search_changed"] = musicPlaylistSignals.filtered
	(*signals)["on_music_playlists_view_row_activated"] = musicPlaylistSignals.dblClicked
	(*signals)["on_music_playlists_view_button_press_event"] = musicPlaylistSignals.context
	(*signals)["on_music_playlists_sel_changed"] = musicPlaylistSignals.selChanged
	(*signals)["on_music_playlists_view_context_delete_activate"] = musicPlaylistSignals.contextDelete
	(*signals)["on_music_playlists_view_context_edit_activate"] = musicPlaylistSignals.contextEdit
	(*signals)["on_music_playlists_view_context_open_activate"] = musicPlaylistSignals.contextOpen
	(*signals)["on_music_playlists_view_context_export_activate"] = musicPlaylistSignals.contextExport
	return
}
