### [0.8.5](https://github.com/keptn/go-utils/compare/v0.8.4...v0.8.5)


### Features
* Added `sh.keptn.log.error` event struct [#4306](https://github.com/keptn/keptn/issues/4306)
* Added support for the interaction with Keptn's Log API [#4030](https://github.com/keptn/keptn/issues/4030) [#4032](https://github.com/keptn/keptn/issues/4032)
* Extended the secretUtils library to support the retrieval of secret metadata from Keptn's secret API [#4061](https://github.com/keptn/keptn/issues/4061)
* Added support for the interaction with Keptn's Uniform Integration Registration API [#4031](https://github.com/keptn/keptn/issues/4031)
* Reduced the number of dependencies [#4063](https://github.com/keptn/keptn/issues/4063)

### [0.8.4](https://github.com/keptn/go-utils/compare/v0.8.3...v0.8.4)


### Features
* Go Version 1.16 support #288 [#2936](https://github.com/keptn/keptn/issues/2936)
* Added `triggeredAfter` property to `Task` structs #281 [#3681](https://github.com/keptn/keptn/issues/3681)
* Added structs for supporting the closed loop remediation use case introduced in Keptn v0.8.3: Added `GetAction` struct, and adapted `Problem` structs #287 [#3682](https://github.com/keptn/keptn/issues/3682)
* Removed obsolete remediation use case structs #293 [#4084](https://github.com/keptn/keptn/issues/4084)
* Added structs to support user-managed deployment URIs #289 [#3757](https://github.com/keptn/keptn/issues/3757)
* Added function ExecuteCommandWithEnv in cmdUtils #286
* Extended common utils package(s) #281

### [0.8.3](https://github.com/keptn/go-utils/compare/v0.8.2...v0.8.3)


### Features
* Introduced common utils package(s) [#279](https://github.com/keptn/go-utils/pull/279)


### [0.8.2](https://github.com/keptn/go-utils/compare/v0.8.1...v0.8.2)


### Features
* Added Client for new secrets API [#3465](https://github.com/keptn/keptn/issues/3465)
* Added `Deployment` field to `GetSLITriggeredEventData` [#72](https://github.com/keptn/spec/issues/72)
* Removed file I/O from `GetKeptnResource` [#3465](https://github.com/keptn/keptn/issues/3465)
* Added `displayName` to SLO and SLI result [#3345](https://github.com/keptn/keptn/issues/3345)


### [0.8.1](https://github.com/keptn/go-utils/compare/v0.8.0...v0.8.1)


### Features
* Added Client for new secrets API [#3465](https://github.com/keptn/keptn/issues/3465)
* Added `Deployment` field to `GetSLITriggeredEventData` [#72](https://github.com/keptn/spec/issues/72)
* Removed file I/O from `GetKeptnResource` [#3465](https://github.com/keptn/keptn/issues/3465)
* Added `displayName` to SLO and SLI result [#3345](https://github.com/keptn/keptn/issues/3345)


## [0.8.0](https://github.com/keptn/go-utils/compare/v0.7.2-alpha...v0.8.0)


### Features
* Moved Auto-PR from Travis CI to GitHub actions [#2750](https://github.com/keptn/keptn/2750)
* Moved unit tests from Travis CI to GitHub actions [#2796](https://github.com/keptn/keptn/2796)
* Removed WebSocket functionality [#2727](https://github.com/keptn/keptn/2727)
* Added events for configure-monitoring task [#2727](https://github.com/keptn/keptn/2727)
* Added convenience method for retrieving open `.triggered` events [#2533](https://github.com/keptn/keptn/2533)
* Added events and helper functions for the new project/service creation mechanism [#2266](https://github.com/keptn/keptn/2266)
* Added shkeptnspecversion property to CloudEvent context [#2982](https://github.com/keptn/keptn/issues/2982)
* Added `.invalidated` event type [#spec-55](https://github.com/keptn/spec/issues/55)
* Adapted to changes in Keptn API in API client helpers [#3001](https://github.com/keptn/keptn/issues/3001) [#2999](https://github.com/keptn/keptn/issues/2999)
* Added convenience methods for sending `.started`, `.status-changed` and `.finished` events [#3035](https://github.com/keptn/keptn/issues/3035)
* Deprecated `EventBrokerURL` property used for the `NewKeptn()` function. Make `EventSender` injectable to `KeptnHandler` to allow easier unit testing [#2919](https://github.com/keptn/keptn/issues/2919)
* Removed obsolete CloudEvent structures [#2830](https://github.com/keptn/keptn/issues/2830) [#2922](https://github.com/keptn/keptn/issues/2922)


### Known Limitations

- This version introduces a couple of changes within the structure of the module. When upgrading from previous versions, please make sure to follow the instructions
  in the [README.md](https://github.com/keptn/go-utils/tree/release-0.7.2#upgrade-to-072-from-previous-versions)
