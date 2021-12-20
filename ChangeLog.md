ChangeLog
=========

All noticeable changes in the project  are documented in this file.

Format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

This project uses [semantic versions](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
