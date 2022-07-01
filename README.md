# Keptn go-utils
![GitHub release (latest by date)](https://img.shields.io/github/v/release/keptn/go-utils)
![tests](https://github.com/keptn/go-utils/workflows/tests/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/keptn/go-utils)](https://goreportcard.com/report/github.com/keptn/go-utils)

<img src="./gopher.png" alt="go-utils-gopher" width="210"/>

This repository contains packages for common functionality around the [Keptn Project](https://github.com/keptn).
Please post any issues to [keptn/keptn repository](https://github.com/keptn/keptn/issues) and label them with `area:go-utils`.

## Installation

Get the latest version using
```console
go get github.com/keptn/go-utils
```
Also consider browsing our [Releases Page](https://github.com/keptn/go-utils/releases) to find out about all releases.

## Contributing

If you want to contribute, just create a PR on the **master** branch.
Please also see [CONTRIBUTING.md](CONTRIBUTING.md) instructions on how to contribute.


## Create a Keptn service using `cp-connector`

One way to create a Keptn integration (a.k.a. Keptn service) is to use the `cp-connector` library which abstracts away the
details of how to interact with the keptn api to register your implementation as an Keptn integration to the control plane.

[Example](./examples/cp-connector)

## Create a Keptn service using the `go-sdk` (experimental)

If you want to use more features besides what the `cp-connector` provides, you can use the Keptn `go-sdk` which
basically wraps around `cp-connector` and provides features like automatic sending of `.started/.finished` or error events.

[Example](./examples/go-sdk)
