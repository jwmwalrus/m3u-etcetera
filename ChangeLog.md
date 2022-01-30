ChangeLog
=========

All noticeable changes in the project  are documented in this file.

Format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

This project uses [semantic versions](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
