<p align="center">
<img src="docs/images/logo.png" height="150">
</p>


[![Build Status](https://travis-ci.com/kleveross/ormb.svg?branch=master)](https://travis-ci.com/kleveross/ormb)
[![Coverage Status](https://coveralls.io/repos/github/kleveross/ormb/badge.svg?branch=master)](https://coveralls.io/github/kleveross/ormb?branch=master)

ormb is an open-source model registry to manage machine learning model. 

ormb helps you manage your Machine Learning/Deep Learning models. It makes your models easy to create, version, share and publish.

[![asciicast](https://asciinema.org/a/345812.svg)](https://asciinema.org/a/345812)

## Installation

You can install the pre-compiled binary, or compile from source.

### Install the pre-compiled binary

Download the pre-compiled binaries from [the releases](https://github.com/kleveross/ormb/releases) page and copy to the desired location.

### Compile from source

Clone:

```
$ git clone https://github.com/kleveross/ormb
$ cd ormb
```

Get the dependencies:

```
$ go mod tidy
```

Build:

```
$ make build-local
```

Verify it works:

```
$ ./bin/ormb --help
```

## Introduction

[ormb 介绍](./docs/intro_zh.md)

## Tutorial

### Distribute models with ormb and Docker Registry

Please have a look at [docs/tutorial.md](docs/tutorial.md)

### Serving with Seldon Core

Please have a look at [docs/tutorial-serving-seldon.md](docs/tutorial-serving-seldon.md)

## OCI Model Configuration Specification

Please have a look at [docs/spec_v1alpha1.md](docs/spec-v1alpha1.md)

## Community

ormb project is part of Klever, a Cloud Native Machine Learning platform.

The Klever slack workspace is klever.slack.com. To join, click this [invitation to our Slack workspace](https://join.slack.com/t/kleveross/shared_invite/zt-g0eoiyq9-9OwiI7c__oV79bh_94MyTw).
