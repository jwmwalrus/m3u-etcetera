package musicpane

import (
	log "github.com/sirupsen/logrus"
)

var (
// collsFilter *gtk.TreeModelFilter
)

var (
	musicCollectionsSignals = &onMusicCollections{}
	musicQueueSignals       = &onMusicQueue{}
	musicQuerySignals       = &onMusicQuery{}
)

// Setup sets the music pane
func Setup(signals *map[string]interface{}) (err error) {
	log.Info("Setting up music")

	if err = createMusicCollections(); err != nil {
		return
	}
	(*signals)["on_collections_sel_changed"] = musicCollectionsSignals.selChanged
	(*signals)["on_collections_view_button_press_event"] = musicCollectionsSignals.context
	(*signals)["on_collections_view_context_append_activate"] = musicCollectionsSignals.contextAppend
	(*signals)["on_collections_view_context_preppend_activate"] = musicCollectionsSignals.contextAppend
	(*signals)["on_collections_view_context_play_now_activate"] = musicCollectionsSignals.contextPlayNow
	(*signals)["on_collections_view_row_activated"] = musicCollectionsSignals.dblClicked
	(*signals)["on_collections_filter_search_changed"] = musicCollectionsSignals.filtered

	if err = createMusicQueue(); err != nil {
		return
	}
	(*signals)["on_music_queue_sel_changed"] = musicQueueSignals.selChanged
	(*signals)["on_music_queue_view_button_press_event"] = musicQueueSignals.context
	(*signals)["on_music_queue_view_context_play_now_activate"] = musicQueueSignals.contextPlayNow
	(*signals)["on_music_queue_view_context_enqueue_activate"] = musicQueueSignals.contextEnqueue
	(*signals)["on_music_queue_view_context_top_activate"] = musicQueueSignals.contextMove
	(*signals)["on_music_queue_view_context_up_activate"] = musicQueueSignals.contextMove
	(*signals)["on_music_queue_view_context_down_activate"] = musicQueueSignals.contextMove
	(*signals)["on_music_queue_view_context_bottom_activate"] = musicQueueSignals.contextMove
	(*signals)["on_music_queue_view_context_delete_activate"] = musicQueueSignals.contextDelete
	(*signals)["on_music_queue_view_context_clear_activate"] = musicQueueSignals.contextClear
	(*signals)["on_music_queue_view_row_activated"] = musicQueueSignals.dblClicked

	if err = createMusicQueries(); err != nil {
		return
	}
	(*signals)["on_queries_sel_changed"] = musicQuerySignals.selChanged
	(*signals)["on_queries_view_button_press_event"] = musicQuerySignals.context
	(*signals)["on_queries_view_context_append_activate"] = musicQuerySignals.contextAppend
	(*signals)["on_queries_view_context_delete_activate"] = musicQuerySignals.contextDelete
	(*signals)["on_queries_view_row_activated"] = musicQuerySignals.dblClicked
	(*signals)["on_queries_filter_search_changed"] = musicQuerySignals.filtered

	if err = musicQuerySignals.createDialog(); err != nil {
		return
	}
	(*signals)["on_queries_add_clicked"] = musicQuerySignals.defineQuery
	(*signals)["on_query_dialog_search_clicked"] = musicQuerySignals.doSearch

	return
}
