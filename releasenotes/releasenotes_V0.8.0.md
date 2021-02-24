# Release Notes develop

## New Features

- Moved Auto-PR from Travis CI to GitHub actions [#2750](https://github.com/keptn/keptn/2750)
- Moved unit tests from Travis CI to GitHub actions [#2796](https://github.com/keptn/keptn/2796)
- Removed WebSocket functionality [#2727](https://github.com/keptn/keptn/2727)
- Added events for configure-monitoring task [#2727](https://github.com/keptn/keptn/2727)
- Added convenience method for retrieving open `.triggered` events [#2533](https://github.com/keptn/keptn/2533)
- Added events and helper functions for the new project/service creation mechanism [#2266](https://github.com/keptn/keptn/2266)
- Added shkeptnspecversion property to CloudEvent context [#2982](https://github.com/keptn/keptn/issues/2982)
- Added `.invalidated` event type [#spec-55](https://github.com/keptn/spec/issues/55)
- Adapted to changes in Keptn API in API client helpers [#3001](https://github.com/keptn/keptn/issues/3001) [#2999](https://github.com/keptn/keptn/issues/2999)
- Added convenience methods for sending `.started`, `.status-changed` and `.finished` events [#3035](https://github.com/keptn/keptn/issues/3035)
- Deprecated `EventBrokerURL` property used for the `NewKeptn()` function. Make `EventSender` injectable to `KeptnHandler` to allow easier unit testing [#2919](https://github.com/keptn/keptn/issues/2919)
- Removed obsolete CloudEvent structures [#2830](https://github.com/keptn/keptn/issues/2830) [#2922](https://github.com/keptn/keptn/issues/2922)

## Fixed Issues

## Known Limitations

- This version introduces a couple of changes within the structure of the module. When upgrading from previous versions, please make sure to follow the instructions
in the [README.md](https://github.com/keptn/go-utils/tree/release-0.7.2#upgrade-to-072-from-previous-versions)

