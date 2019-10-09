# Keptn go-utils
[![Build Status](https://travis-ci.org/keptn/go-utils.svg?branch=master)](https://travis-ci.org/keptn/go-utils)

This repo serves as a util package for common functionalities such as logging of the [Keptn Project](https://github.com/keptn).

Please post any issues with this package to the [keptn/keptn repository](https://github.com/keptn/keptn/issues) and label them with `area:go-utils`.

## Usage

```
import {
  keptnutils "github.com/keptn/go-utils/pkg/utils"
}

```

## Logging

```
keptnutils.Debug(keptncontext, message)
keptnutils.Info(keptncontext, message)
keptnutils.Error(keptncontext, message)
```

