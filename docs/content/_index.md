---
title: wp-opentofu
---

[![Build Status](https://ci.thegeeklab.de/api/badges/thegeeklab/wp-opentofu/status.svg)](https://ci.thegeeklab.de/repos/thegeeklab/wp-opentofu)
[![Docker Hub](https://img.shields.io/badge/dockerhub-latest-blue.svg?logo=docker&logoColor=white)](https://hub.docker.com/r/thegeeklab/wp-opentofu)
[![Quay.io](https://img.shields.io/badge/quay-latest-blue.svg?logo=docker&logoColor=white)](https://quay.io/repository/thegeeklab/wp-opentofu)
[![Go Report Card](https://goreportcard.com/badge/github.com/thegeeklab/wp-opentofu)](https://goreportcard.com/report/github.com/thegeeklab/wp-opentofu)
[![GitHub contributors](https://img.shields.io/github/contributors/thegeeklab/wp-opentofu)](https://github.com/thegeeklab/wp-opentofu/graphs/contributors)
[![Source: GitHub](https://img.shields.io/badge/source-github-blue.svg?logo=github&logoColor=white)](https://github.com/thegeeklab/wp-opentofu)
[![License: Apache-2.0](https://img.shields.io/github/license/thegeeklab/wp-opentofu)](https://github.com/thegeeklab/wp-opentofu/blob/main/LICENSE)

Woodpecker CI plugin to manage infrastructure with [OpenTofu](https://github.com/opentofu/opentofu).

<!-- prettier-ignore-start -->
<!-- spellchecker-disable -->
{{< toc >}}
<!-- spellchecker-enable -->
<!-- prettier-ignore-end -->

## Usage

```YAML
steps:
  - name: tofu
    image: quay.io/thegeeklab/wp-opentofu
    settings:
      actions:
        - validate
        - plan
```

### Parameters

<!-- prettier-ignore-start -->
<!-- spellchecker-disable -->
{{< propertylist name=wp-opentofu.data sort=name >}}
<!-- spellchecker-enable -->
<!-- prettier-ignore-end -->

## Build

Build the binary with the following command:

```Shell
make build
```

Build the Container image with the following command:

```Shell
docker build --file Containerfile.multiarch --tag thegeeklab/wp-opentofu .
```

## Test

```Shell
docker run -it thegeeklab/wp-opentofu --action validate
```
