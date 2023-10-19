# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

### [0.20.3](https://github.com/keptn/go-utils/compare/v0.20.2...v0.20.3) (2023-10-19)

### [0.20.2](https://github.com/keptn/go-utils/compare/v0.20.1...v0.20.2) (2023-08-01)


### Other

* fix dependencies ([#591](https://github.com/keptn/keptn/issues/591)) ([4085f33](https://github.com/keptn/go-utils/commit/4085f3329695274ed263f240c09c64b1fc1a43c0))

### [0.20.1](https://github.com/keptn/go-utils/compare/v0.20.0...v0.20.1) (2023-04-04)


### Features

* add possibility to read mongo secret from file ([#566](https://github.com/keptn/keptn/issues/566)) ([6b8d8f2](https://github.com/keptn/go-utils/commit/6b8d8f22108199a622236195f76816924354fc84))

## [0.20.0](https://github.com/keptn/go-utils/compare/v0.19.0...v0.20.0) (2022-12-07)


### Features

* **go-sdk:** Helper for retrieving SLI resources ([#541](https://github.com/keptn/keptn/issues/541)) ([cc6d13c](https://github.com/keptn/go-utils/commit/cc6d13cd1a40a4d02f1a2031b82c49451320f4ad))


### Bug Fixes

* Fix event handling for duplicated subscriptions in `go-sdk` remote execution plane use case ([#542](https://github.com/keptn/keptn/issues/542)) ([8655d73](https://github.com/keptn/go-utils/commit/8655d7329e7c5008d7acd95411a8dd00e2c456e4))


### Other

* **go-utils:** Changed dependencies ([#545](https://github.com/keptn/keptn/issues/545)) ([1f92e1e](https://github.com/keptn/go-utils/commit/1f92e1e08722c66e3d92b943f203ee95b8e5b42c))

## [0.19.0](https://github.com/keptn/go-utils/compare/v0.18.0...v0.19.0) (2022-09-23)


### Features

* Expose methods to automatically generate events from parent event ([#538](https://github.com/keptn/keptn/issues/538)) ([cf4bfd8](https://github.com/keptn/go-utils/commit/cf4bfd8b27ffe36fe330722b1fb826b5ced4f80e))
* Introduce IsUpstreamAutoProvisioned to ExpandedProject model ([#536](https://github.com/keptn/keptn/issues/536)) ([dc8c096](https://github.com/keptn/go-utils/commit/dc8c0968b1331f46db4d5e7846d7e09b1ca1c4a5))
* Introduce skipping of automatic event responses per task handler ([#537](https://github.com/keptn/keptn/issues/537)) ([278ad8b](https://github.com/keptn/go-utils/commit/278ad8bdb21774ecce3e5de80ec6bbd8e0e0d9c0))


### Bug Fixes

* Use resourceHandler from apiv2 package instead of newly creating an internal-only, unauthenticated client ([#535](https://github.com/keptn/keptn/issues/535)) ([f07eb2f](https://github.com/keptn/go-utils/commit/f07eb2f4a74be3e0b54972d1ff52f303708d4a24))

## [0.18.0](https://github.com/keptn/go-utils/compare/v0.17.0...v0.18.0) (2022-07-28)


### ⚠ BREAKING CHANGES

* **go-utils:** Since the configuration-service is deprecated, all references to “configuration-service” are now replaced by “resource-service”. This will make the go-utils library from version 0.18.0 INCOMPATIBLE with installations using configuration-service.

### Features

* Disable NATS connection on remote execution-plane configuration ([#524](https://github.com/keptn/keptn/issues/524)) ([866624f](https://github.com/keptn/go-utils/commit/866624f8ce4274c3a9ae9b597fba04a338e20dc3))
* **go-utils:** Added retry logic to cp-connector for contacting Keptn's control plane for registration and renewal of such ([#503](https://github.com/keptn/keptn/issues/503)) ([69c90ea](https://github.com/keptn/go-utils/commit/69c90ea3a7bcabe095d06144d8bb3b09a209c89b))
* Wait for all event handlers to complete before exiting `controlPlane.Register()` ([#496](https://github.com/keptn/keptn/issues/496)) ([d9a621b](https://github.com/keptn/go-utils/commit/d9a621b678384a8a1e0e7fcbe2ddcfef292aaa87))


### Bug Fixes

* Add mutex to protect `connection` in `nats.NatsConnector` ([#514](https://github.com/keptn/keptn/issues/514)) ([3a171cc](https://github.com/keptn/go-utils/commit/3a171cc6c166bec3c632d47d3ddfbdb5a006054d))
* Shut down control plane components before calling wg.Wait() ([#523](https://github.com/keptn/keptn/issues/523)) ([6b12679](https://github.com/keptn/go-utils/commit/6b126790c54c120282e8eb69ab882825dac8943c))
* Time property is not lost between Keptn and CloudEvent conversion ([#495](https://github.com/keptn/keptn/issues/495)) ([3ef0a10](https://github.com/keptn/go-utils/commit/3ef0a10730ad12b551b3ac5885495b108fac0fe7))


### Other

* **go-utils:** Changed configuration-service to resource-service ([#491](https://github.com/keptn/keptn/issues/491)) ([6550348](https://github.com/keptn/go-utils/commit/65503489d75aac53e74e089babc7a5258887d93c))
* increase test coverage of go-sdk ([#526](https://github.com/keptn/keptn/issues/526)) ([c15488f](https://github.com/keptn/go-utils/commit/c15488f2145b2c0235687a0e3884123f4ec47354))
* Remove unneeded code ([#490](https://github.com/keptn/keptn/issues/490)) ([d00898a](https://github.com/keptn/go-utils/commit/d00898a7fdcf3eebd427a6b6c93866c6cdaea1aa))

## [0.17.0](https://github.com/keptn/go-utils/compare/v0.16.0...v0.17.0) (2022-07-05)


### ⚠ BREAKING CHANGES

* Git credentials for git authentication were moved to a separate sub-structure and split to either ssh or http sub-structures depending on the used authentication method.
 
### Features

* Add `v2.InternalAPISet` that implements the `v2.KeptnInterface` ([#487](https://github.com/keptn/keptn/issues/487)) ([eb5fb9b](https://github.com/keptn/go-utils/commit/eb5fb9ba43e021fd8de4467ac086d1bda4df6da0))
* Add `v2.KeptnInterface` that adds `context.Context` support to `api.KeptnInterface` ([#449](https://github.com/keptn/keptn/issues/449)) ([0874051](https://github.com/keptn/go-utils/commit/0874051eb4dd2f20e6c4b19e750d37931b426330)), closes [#479](https://github.com/keptn/keptn/issues/479)
* Move commonly used modules from keptn/keptn into sub-packages of go-utils ([#483](https://github.com/keptn/keptn/issues/483)) ([3ed2fc6](https://github.com/keptn/go-utils/commit/3ed2fc6cf1edd6cafa55a2b0b8cab19edcd0fc7d))
* Refactor git remote repository credentials models ([#475](https://github.com/keptn/keptn/issues/475)) ([fc5b6f9](https://github.com/keptn/go-utils/commit/fc5b6f967e50bc7765311a706e052bc21c19ad07))


### Bug Fixes

* **go-utils:** Pass logger implementation from go-sdk to cp-connector ([#494](https://github.com/keptn/keptn/issues/494)) ([29e14a0](https://github.com/keptn/go-utils/commit/29e14a06fcb70ada7f3704dc80c02cb69b522aba))
* Make GetAllServiceResources compatible with Keptn 0.16.0+ ([#480](https://github.com/keptn/keptn/issues/480)) ([0d19a1b](https://github.com/keptn/go-utils/commit/0d19a1b4509e82c678bb77d3a216fe2597eef37e))
* Make unit tests work with `-race` flag ([#489](https://github.com/keptn/keptn/issues/489)) ([9b0c779](https://github.com/keptn/go-utils/commit/9b0c779950b66d699df08987dd26a40b14cccf8b))
* Set the path properly for calls to api-service ([#470](https://github.com/keptn/keptn/issues/470)) ([a3c50ce](https://github.com/keptn/go-utils/commit/a3c50ce6446d7bc5226ff71c0a31c193972e5aca))
* Use ExecuteCommand implementation from kubernetes-utils ([#482](https://github.com/keptn/keptn/issues/482)) ([8d145bc](https://github.com/keptn/go-utils/commit/8d145bc902942b486c6872d4cbcbfaf188c3c4ec))
* Use NetworkingV1 instead of deprecated ExtensionsV1beta1 ([#492](https://github.com/keptn/keptn/issues/492)) ([0fc8c36](https://github.com/keptn/go-utils/commit/0fc8c36d130663fbfe443872f4020a220c5cac24))


### Other

* **go-utils:** Removed deprecated subscription from uniform ([#474](https://github.com/keptn/keptn/issues/474)) ([647fbac](https://github.com/keptn/go-utils/commit/647fbacf27d0d60775142524d516fa529ca21066))
* Introduce needed methods before deprecating kubernetes-utils ([#477](https://github.com/keptn/keptn/issues/477)) ([4d49101](https://github.com/keptn/go-utils/commit/4d49101f88b408c90fc691d893e480f809634178))


### Docs

* **go-utils:** Update `README.md` documentation of `go-utils` ([#493](https://github.com/keptn/keptn/issues/493)) ([8369229](https://github.com/keptn/go-utils/commit/83692294c5793f5bcf7eda43aac1884c15c961f5))

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
