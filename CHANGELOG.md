# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

### [0.16.1](https://github.com/keptn/go-utils/compare/v0.16.0...v0.16.1) (2022-06-03)


### Bug Fixes

* Set the path properly for calls to api-service ([#471](https://github.com/keptn/keptn/issues/471)) ([0650f77](https://github.com/keptn/go-utils/commit/0650f7701c35cfd296d85e6e224fab6770028b09))

## [0.16.0](https://github.com/keptn/go-utils/compare/v0.15.0...v0.16.0) (2022-06-02)


### Features

* Added option for configuring number of retries for http event sender, as well as additional logging ([#465](https://github.com/keptn/keptn/issues/465)) ([2052e14](https://github.com/keptn/go-utils/commit/2052e1404e4c238aca16b2f2ea510d042373df4e))
* Provide option to specify readiness condition ([#464](https://github.com/keptn/keptn/issues/464)) ([c5e1b75](https://github.com/keptn/go-utils/commit/c5e1b7519becaa8cd0f3501a4909df74f061843f))


### Bug Fixes

* Fixed wrong paths for apiutils ([#469](https://github.com/keptn/keptn/issues/469)) ([623e06c](https://github.com/keptn/go-utils/commit/623e06c0a2f683415c29d1a2ace85169a369fbee))

## [0.15.0](https://github.com/keptn/go-utils/compare/v0.14.0...v0.15.0) (2022-05-06)


### ⚠ BREAKING CHANGES

* Rename the create/update project parameter `GitProxyInsecure` to `InsecureSkipTLS` * removed unused GitCommit from finished events

### Features

* Introduce proxy parameters to ExpandedProject struct ([#433](https://github.com/keptn/keptn/issues/433)) ([6c53542](https://github.com/keptn/go-utils/commit/6c53542a31b0a4036e2dc792bab4d0ca7528d642))
* Rename GitProxyInsecure to InsecureSkipTLS since that property should not only be tied to the proxy feature ([#445](https://github.com/keptn/keptn/issues/445)) ([003ee3f](https://github.com/keptn/go-utils/commit/003ee3f85292e5ab6a049f8412bbb8fa28d4b6eb))


### Bug Fixes

* Avoid nil pointer exception in AddTemporaryData ([#460](https://github.com/keptn/keptn/issues/460)) ([5672c07](https://github.com/keptn/go-utils/commit/5672c074a6de2e3d6c868fc1abd9a87744ff66e3))
* correct error message in fake/EventSender.AssertSentEventTypes ([2466736](https://github.com/keptn/go-utils/commit/24667368a1594c17cfa725b88a97d19ddfef109e))
* **go-utils:** Add parameters to metadata model ([#434](https://github.com/keptn/keptn/issues/434)) ([297c1b2](https://github.com/keptn/go-utils/commit/297c1b2ddd7c40e518659dceddc20532ad99d321))
* **go-utils:** Make AutomaticProvisioning parameter required in metadata model ([#435](https://github.com/keptn/keptn/issues/435)) ([0b73d75](https://github.com/keptn/go-utils/commit/0b73d757bfd6589eb8dde8e802e27cc3cfea997d))
* **go-utils:** Make GitProxyInsecure parameter required ([#436](https://github.com/keptn/keptn/issues/436)) ([802847e](https://github.com/keptn/go-utils/commit/802847e8175045550ea877694ef1f6e71c33fa15))
* Gracefully handle missing event labels ([#446](https://github.com/keptn/keptn/issues/446)) ([2e23eb7](https://github.com/keptn/go-utils/commit/2e23eb712e3db97c5f1136e67293bbbbc4111e05))
* Restore opentelemetry dependency updates ([#456](https://github.com/keptn/keptn/issues/456)) ([a0381c5](https://github.com/keptn/go-utils/commit/a0381c53c6d819e63bcbdb1881a1fb5a03332158))


### Other

* Removed unneeded Git Commit from finished events ([#430](https://github.com/keptn/keptn/issues/430)) ([c6d4983](https://github.com/keptn/go-utils/commit/c6d49838bec8f86f6bbab373474d85734e738ad7))

## [0.13.0](https://github.com/keptn/go-utils/compare/v0.12.0...v0.13.0) (2022-02-17)


### Features

* Add SSH publicKey auth support ([#392](https://github.com/keptn/keptn/issues/392)) ([be3425c](https://github.com/keptn/go-utils/commit/be3425c548a783d8e571492763199fea4921f82b))
* added oauthutils package ([#395](https://github.com/keptn/keptn/issues/395)) ([f30183e](https://github.com/keptn/go-utils/commit/f30183e06eee9cb7a3e182f2ea8b8378e403e0d1))
* added query parameters to resource getter (keptn/keptn/[#6349](https://github.com/keptn/keptn/issues/6349)) ([#375](https://github.com/keptn/keptn/issues/375)) ([b7470c0](https://github.com/keptn/go-utils/commit/b7470c0a2a7c5a4321d15754c12fbddc1a9e2607))
* Introduced interfaces for different types of APIs ([#379](https://github.com/keptn/keptn/issues/379)) ([349cd94](https://github.com/keptn/go-utils/commit/349cd94a73287f7151281c7da5cea973217e2491))
* introducing `APISet` for more convenient access to keptn APIs ([#377](https://github.com/keptn/keptn/issues/377)) ([5c52509](https://github.com/keptn/go-utils/commit/5c525092bb7634df2edaa5f35e3bdce16c7dff2e))
* Propagate git commit ID for sequence in CloudEvent context ([#374](https://github.com/keptn/keptn/issues/374)) ([fa37290](https://github.com/keptn/go-utils/commit/fa37290ac704af4f14ff6fa5b865e26183e8891c))


### Bug Fixes

* Add missing Method to KeptnInterface/APISet ([#393](https://github.com/keptn/keptn/issues/393)) ([6b99172](https://github.com/keptn/go-utils/commit/6b991728fc642743b4a80f272f203bddc8dd18d7))
* Revert old getters and deprecate them ([#381](https://github.com/keptn/keptn/issues/381)) ([376fb7b](https://github.com/keptn/go-utils/commit/376fb7bad8f1570410988bfda9950ca64ee52199))


### Other

* enhance configuration of APISet ([#378](https://github.com/keptn/keptn/issues/378)) ([d68990e](https://github.com/keptn/go-utils/commit/d68990e6fabd8919e252586656d3d3fa6f7328e9))

## [0.12.0](https://github.com/keptn/go-utils/compare/v0.11.0...v0.12.0) (2022-02-16)


### Features

* added ComparedValues to SLIResult (keptn/[#5496](https://github.com/keptn/keptn/issues/5496)) ([#358](https://github.com/keptn/keptn/issues/358)) ([e95de56](https://github.com/keptn/go-utils/commit/e95de56c5c09a0b8bd24ff00b07495e1cc6b2c59))
* added get-action data to contain the action index (keptn/keptn/[#4206](https://github.com/keptn/keptn/issues/4206)) ([#361](https://github.com/keptn/keptn/issues/361)) ([08c82f0](https://github.com/keptn/go-utils/commit/08c82f03f59f79818817ae48e22ac4eac405f9de))


### Bug Fixes

* Add error check when creating requests ([#369](https://github.com/keptn/keptn/issues/369)) ([dcfdacb](https://github.com/keptn/go-utils/commit/dcfdacb805942a76d688dbf84758adce2125c18d))
* adding missing error checks ([#371](https://github.com/keptn/keptn/issues/371)) ([5626bf9](https://github.com/keptn/go-utils/commit/5626bf92b8c69571a4f7eb72564c29ce37e4fb00))
* if the integrationId is not set we should not ping (keptn/[#6309](https://github.com/keptn/keptn/issues/6309)) ([#370](https://github.com/keptn/keptn/issues/370)) ([de65cd4](https://github.com/keptn/go-utils/commit/de65cd48c936cc5b40ad8ac568a2e8575a1b1598))
* Make fake event sender thread safe by adding a lock ([#357](https://github.com/keptn/keptn/issues/357)) ([fe1fb0c](https://github.com/keptn/go-utils/commit/fe1fb0c473a48a094f9dc7b593f4138004934fe9))


### Other

*  forced grpc to use latest x/net library (snyc security treat) ([#362](https://github.com/keptn/keptn/issues/362)) ([8dcf434](https://github.com/keptn/go-utils/commit/8dcf434a97080bbd607fd6b6f32cf683e5e88d3e))


### Refactoring

* Move code for (de)serialization from/to JSON to model structs ([#376](https://github.com/keptn/keptn/issues/376)) ([544c270](https://github.com/keptn/go-utils/commit/544c27052949921a0a8aa676c8eeb19aeadf598c))
---

✔ Running lifecycle script "postchangelog"
ℹ - execute command: "./gh-actions-scripts/post-changelog-actions.sh"
✔ committing CHANGELOG.md
✔ tagging release v0.12.0
ℹ Run `git push --follow-tags origin HEAD` to publish
