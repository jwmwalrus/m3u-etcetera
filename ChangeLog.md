ChangeLog
=========

All noticeable changes in the project  are documented in this file.

Format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

This project uses [semantic versions](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.21.0] 2023-10-27

Adapt to dependencies' changes

### Updated

* Adapt to rtcycler changes
* Rework middleware code
* Update dependencies, cleanup comments

## [0.20.0] 2023-05-22

Twenty percent

### Fixed

* Fix segfaults

### Added

* Opening supported playlist files
* Read-only queries
* Task to generate deb package
* Handling images generation

### Modified

* Display proper duration from cli
* Collecion name's index must be unique
* Use rtcycler
* Inject playback events into services
* Refactoring, misc. tweaks
    * Rearrange some packages and fix tests
    * Update .golanci.yml and fix new issues

## [0.19.0] 2023-03-03

Misc. changes

## Fixed

* Panic on subscription-related enum
* Race conditions

## Added

* Resizable playlist columns

## Modified

* Updated protoc-gen-go version
* Updated Go version and dependencies
* Increased isolation of concerns
* Improved unloader performance
* Improved thread safety and performance
* Cleaned up subscriptions and PlaybackChanged channel
* Simplified Taskfile tasks
* Improved logrus.Entry usage
* Made init() the topmost function

## [0.18.0] 2022-07-18

New and improved features

### Added

* Exporting playlist from context menu
* Multiple selection in playlists
* Toggle selection in query results
* Removing a playlist group
* Collection hierarchy switch
* The playlist merge task
* Perspective activation
* Playback seek
* Keyboard events for the delete key.

### Modified

* Rework collection settings implementation
* Progress bar seek cleanup
* Improve/complete collection management features
* Refactor store.values in terms of gtk.TreeModel.
* Implement a better column renderer generator
* Lock selected tree-view values until consumed or reset.
* Further separate `store` from calls (a.k.a. `dialer`)
    * Create a gtk/utils package

## [0.17.0] 2022-06-27

Features and fixes

### Fixed

* Handling of deleted item returned by `poser.DeleteAt`
* Ensure MPRIS updates are exclusive
* Ensure queue and playlists are ordered before calling poser

### Added

* Status bar digest context
* Track duration discovery
* Config file for golangci-lint
* MPPRIS' Player.CanGoNext

### Modified

* Emit proper MPRIS signal when playback changes
* Ensure proper playback unload
* Update protoc-gen-go-grpc version
* Handle deprecations
* Ignore playlist.setFocused error
* Pick a better icon for dynamic mode

## [0.16.0] 2022-06-21

Fixes and dependencies

## Fixed

* `sanityCheck` on GTK's side
* Ensure treview selection values are valid

### Modified

* Upgrade Go version
* Update dependencies
* Refactor/complete MPRIS implementation
* Simplify handling of pointer slices
* Abstract list handler out
* Make opening query as playlist the default
* Refactor Taskfile
* Rename web directory to mobile
* Defer deletion of transient playlists
* Embed GTK resources

## [0.15.0] 2022-02-14

Papercuts

## Fixed

* Latency issue when stopping playback

### Added

* Focus request for GUI playlists
* GtkApplication, to ensure there's only one GTK instance running

### Modified

* Played threshold when previousEvent is issued
    * If track position is below threshold, restart current playback
* Disable video and text in playback
* Reduce the number of IdleAdd calls for playbar-related updates
* Don't rebuild playlists from scratch after every update
* Use a separate db session for scanning tracks
* Update tests

## [0.14.0] 2022-02-05

Features, fixes and improvements

### Fixed

* Segfault when finishing playback
* Crash related to query_parse_seeking
* Query dialog limit and rating

### Added

* Implement PLS import/export
* GUI validations on unique names
* Playlist export task
* Adding a ollection from the GUI
* Creating playlist groups

### Modified

* Adjustments to icon
* Switch GStreamer bindings
* Disable items not yet implemented
* Add default playlist groups
    * There should be one default playlist group per perspective
    * `Transient` should be a property of the playlist
* Truncate duration for display purposes
* Improve database cleanup
* Split large source files
    * Add some consistency to naming

## [0.13.0] 2022-01-31

Playlists and context menus

### Fixed

* Playlist closing issue, by not relying upon GtkNotebook's page number
* Trim playlist name when importing M3U file

### Added

* Appending to playlist/queue from query
* Creating playlist from query
* Unimplemented MPRIS properties
* Playback list task
* New icon

### Modified

* Unify context menus common code
* Conditionally disable menu items
* Reduced the number of IdleAdd calls
* Allow editing closed playlists
* Do not send tracks for closed open-playlists

## [0.12.0] 2022-01-29

M3U playlists and MPRIS support

### Added

* Importing M3U playlists
* Partial MPRIS support

### Modified

* Move application interfaces out of internal
* Conditionally update playbar model from playback

## [0.11.0] 2022-01-27

Implement playbars

## [0.10.0] 2022-01-14

UI improvements and updated dependencies

### Fixed

* Query tree not refreshing properly after delete
* Play previous
* Issues with TreeView selections

### Added

* Display logo as cover
* Set application's subtitle to current playback

### Modified

* Updated dependencies
* Issue initial event after any collection change
* Reimplement collection tree

## [0.9.0] 2022-01-12

Features, fixes and improvements, again

### Fixed

* The task `serve off --force` not woring while the GTK client is active

### Added

* Own Interceptors
* Progress bar
* Context menu to collections
* Context menu to music queue
* Settings menu
    * Collections manager
    * Quit all
* Queries tab
* Show cover if available
* The `noWait` option to serve off
* A hack to allow exiting the server while stream is paused

### Modified

* Improved Idler encapsulation
* The `Off` method shold be exempted from going to the idle stack
* Simplify client connections
* Refactorings
* Client connections should be handled from gtk/store
* Obtaining values from GtkTreeModel should be handled in just one place
* Anticipate gtk/pane becoming way too big to handle
* Unloaders should be handled in a specific order. In particular:
    * Engine must not be unloaded while server requests are possible;
    * Database should still be available when unloading engine
* Deduplicate covers during scanning
* Minimize db interaction during playback
* Allow non-printable characters in qparams
* Use Label for title, artist, source in gtk, instead of TextView

### Removed

* CollectionTrack (since it was one-to-one)
* Debug channel

## [0.8.0] 2022-01-03

Even more features, fixes and improvements

### Fixed

* Compilation warnings

### Added

* Testing
* Playback, Queue and Collection subscriptions
* GTK playback, music queue and collections treeview
* Interrupt signal

### Modified
* Misc. changes to queue and collections treeviews
* Database-related tweaks
* Allow explicitly forcing application's exit
* Get rid of oneof in playback proto message
* Use incerceptors for logging and idle handling

## [0.7.0] 2021-12-20

More features and improvements

### Fixed

* Issue when stopping engine with pending queue

### Added

* Some testing
* Random column to query lists
* Showing track info from queue when available
* Some interfaces for the future

### Modified
* Query boundaries should always apply
* Refactor base
* Move location and id sanity checks to the API level

## [0.6.0] 2021-12-17

Fixes, features, improvements

## Fixed

* Multiple idle requests and playback stop issues
* Id-as-location issue
* Pointer-in-stack issue

## Added

* Implemented queue move
* Added configurable query limit

## Modified

* Simplified query by
* Reduced log pollution by debug

## [0.5.0] 2021-12-16

### Fixed

* Added missing seed

### Added

* Query task for searching tracks in collections

### Modified

* Display related enhancements to the playback and query tasks
* Reimplemented base.Idle by cancellable context

## [0.4.0] 2021-12-14

Implement collection

## [0.3.0] 2021-12-13

Complete queue implementation

## [0.2.0] 2021-12-13

Implement basic playback

## Added

* Unloader and idle stuff
* Database layer
* Logger middleware
* Perspectives
* Basic playback, including queue capabilities

## Modified

* Renamed protobuf generated package
    * This is to avoid confusion with the playback's short name and prefix

## [0.1.0] 2021-12-12
 
Initial release
