# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

## [0.9.0](https://github.com/keptn/go-utils/compare/v0.9.0-next.0...v0.9.0) (2021-08-26)

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

* This version introduces a couple of changes within the structure of the module. When upgrading from previous versions, please make sure to follow the instructions
  in the [README.md](https://github.com/keptn/go-utils/tree/release-0.7.2#upgrade-to-072-from-previous-versions)

### [0.7.2](https://github.com/keptn/go-utils/compare/v0.7.1...v0.7.2-alpha)


### Features
* Added a new helper function for triggering evaluations via the Keptn
  API [#2387](https://github.com/keptn/keptn/issues/2387)
* Include a list of compared evaluation-done events in the details of an
  evaluation [#2388](https://github.com/keptn/keptn/issues/2388)
* Added metadata properties (git upstream URL, git commit ID, branch) to struct representing the responses from the
  resources API within Keptn [#2307](https://github.com/keptn/keptn/issues/2307)
* Added Next-gen Keptn events that will be used in Keptn 0.8.x [#2107](https://github.com/keptn/keptn/issues/2107)
* Added support for CloudEvents v1.0 [#2254](https://github.com/keptn/keptn/issues/2254)

### Known Limitations

* This version introduces a couple of changes within the structure of the module. When upgrading from previous versions,
  please make sure to follow the instructions in
  the [README.md](https://github.com/keptn/go-utils/tree/release-0.7.2#upgrade-to-072-from-previous-versions)

### [0.7.1](https://github.com/keptn/go-utils/compare/v0.7.0...v0.7.1)


### Features

* Added structs for next generation of Shipyard [#2016](https://github.com/keptn/keptn/issues/2016)
* Added page size parameter to EventFilter
* Added structs for next generation of Keptn CloudEvents [#2107](https://github.com/keptn/keptn/issues/2107)
* Added sh.keptn.internal.event.service.delete event [#2199](https://github.com/keptn/keptn/issues/2199)
* Added helper function to delete service [#2199](https://github.com/keptn/keptn/issues/2199)
* Added helper functions to send .started and .finished events

### Bug Fixes

* Set error code to 404 if no event could be found [#1655](https://github.com/keptn/keptn/issues/1655)

## [0.7.0](https://github.com/keptn/go-utils/compare/v0.6.2...v0.7.0)


### Features
* Added models for Keptn metadata [#181](https://github.com/keptn/go-utils/issues/181)
* Allow retrieval of multiple events from Keptn datastore [#1749](https://github.com/keptn/keptn/issues/1749)
* Allow fine-grained filtering of Keptn events [#161](https://github.com/keptn/go-utils/issues/161)
* Added models for delivery assistant use case [#1749](https://github.com/keptn/keptn/issues/1749)
* Added models for remediation workflow [#1816](https://github.com/keptn/keptn/issues/1816) [#1848](https://github.com/keptn/keptn/issues/1848)
* Simplify logging [#1607](https://github.com/keptn/keptn/issues/1607)
* Added `triggeredid` property to CloudEvents [#1815](https://github.com/keptn/keptn/issues/1815)

### Bug Fixes
* Allow distinguishing between not-available resource and internal configuration-service error [#1480](https://github.com/keptn/keptn/issues/1480)
* Removed fixed host header `api.keptn` for http requests to the api [#1797](https://github.com/keptn/keptn/issues/1797)

### [0.6.2](https://github.com/keptn/go-utils/compare/v0.6.1...v0.6.2)


### Features
* Added a helper function to list all projects [#1549](https://github.com/keptn/keptn/issues/1549)
* Retry to send CloudEvents in case of an error [#1279](https://github.com/keptn/keptn/issues/1279)
* Updated URLs of internal Keptn services [1589](https://github.com/keptn/keptn/issues/1589)
* Added helper functions for sending CloudEvents and retrieving service endpoints [#1079](https://github.com/keptn/keptn/issues/1079)
* Refactored and restructured the complete module [#1492](https://github.com/keptn/keptn/issues/1492)

### [0.6.1](https://github.com/keptn/go-utils/compare/v0.6.0...v0.6.1)


### Features
* Added `DeploymentURILocal` and `DeploymentURIPublic` properties to `DeploymentFinishedEventData` struct. [#1403](https://github.com/keptn/keptn/issues/1403)

## [0.6.0](https://github.com/keptn/go-utils/compare/v0.5.0...v0.6.0)


### Features
* Added result property to `TestsFinishedEventData` [#542](https://github.com/keptn/keptn/issues/542)
* Added method for validating Keptn entity name [#1261](https://github.com/keptn/keptn/issues/1261)
* Added utility to create namespaces [#1231](https://github.com/keptn/keptn/issues/1231)
* Added helper function to get SLI config for service considering stage and project configs [#1192](https://github.com/keptn/keptn/issues/1192)

## [0.5.0](https://github.com/keptn/go-utils/compare/v0.4.0...v0.5.0)


### Features
* Added deployment-type to get-sli events [#1161](https://github.com/keptn/keptn/issues/1161)
* Always set host to `api.keptn` when sending cluster-internal API requests [#1167](https://github.com/keptn/keptn/issues/1167)
* Added `label` property to all events involved in the CD workflow [#1147](https://github.com/keptn/keptn/issues/1147)

## [0.4.0](https://github.com/keptn/go-utils/compare/v0.3.0...v0.4.0)


### Features
* Evaluation Done events contain more info about SLIs and SLOs [#1058](https://github.com/keptn/keptn/issues/1058)
* Flattened events in MongoDB [#1061](https://github.com/keptn/keptn/issues/1061)
* Add testStrategy and deploymentStrategy in several events [#1098](https://github.com/keptn/keptn/issues/1098)

### Bug Fixes
* Fixed an endless loop when fetching resources [#1043](https://github.com/keptn/keptn/issues/1043)

## [0.3.0](https://github.com/keptn/go-utils/compare/v0.2.0...v0.3.0)


### Features
* Provide REST endpoint for project and service [#893](https://github.com/keptn/keptn/issues/893)
* Add automated testing via Travis CI [#944](https://github.com/keptn/keptn/issues/944)
* Allow NodePort for Istio ingressgateway [#462](https://github.com/keptn/keptn/issues/462)
* Provide utility functions to retrieve Keptn events via API [#949](https://github.com/keptn/keptn/issues/949)

### 0.2.1


### Features
* Add functions for deleting a project [#887](https://github.com/keptn/keptn/issues/887)

## 0.2.0


### Features
* Added functions for interacting with the new configuration service to upload and retrieve:
  * Projects
  * Stages
  * Services
  * Resources
* Collection of Keptn events

### Bug Fixes
* Correctly wait for deployment in a namespace to be complete

### 0.1.1


### Features
* Helper method for expanding tilde in filepath [#528](https://github.com/keptn/keptn/issues/528)

### Bug Fixes
* Wait for all deployments to be available when doing helm upgrade [#483](https://github.com/keptn/keptn/issues/483)

## 0.1.0


### Features
* Library for commonly used uitl functions [#418](https://github.com/keptn/keptn/issues/418)
