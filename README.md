<p align="center">
<img src="docs/images/logo.png" height="150">
</p>


[![Build Status](https://travis-ci.com/kleveross/ormb.svg?branch=master)](https://travis-ci.com/kleveross/ormb)
[![Coverage Status](https://coveralls.io/repos/github/kleveross/ormb/badge.svg?branch=master)](https://coveralls.io/github/kleveross/ormb?branch=master)

English | [中文](./README_zh.md)

`ORMB` is an open-source model registry to manage machine learning model. 

`ORMB` helps you manage your Machine Learning/Deep Learning models with image registry. It makes your models easy to create, version, share and publish.

## Getting Started

You can watch our sample usage video or read the text version below.

[![asciicast](https://asciinema.org/a/345812.svg)](https://asciinema.org/a/345812)

First, open a model folder that meets the specification of `ORMB`. (See our documentation for more information about [ormbfile.yaml](/docs/README.md#What-is-ormbfile?))

```bash
# View the local file directory
$ tree examples/SavedModel-fashion
examples/SavedModel-fashion
├── model
│   ├── saved_model.pb
│   └── variables
│       ├── variables.data-00000-of-00001
│       └── variables.index
├── ormbfile.yaml
└── training-serving.ipynb

2 directories, 5 files
```

Next, we can push the trained model from local to remote image registry.

```bash
# Save the model in local cache first
$ ormb save gaocegege/fashion_model:v1
ref:       gaocegege/fashion_model:v1
digest:    6b08cd25d01f71a09c1eb852b3a696ee2806abc749628de28a71b507f9eab996
size:      162.1 KiB
format:    SavedModel
v1: saved

# Push the model from local cache to remote registry
$ ormb push gaocegege/fashion_model:v1
The push refers to repository [gaocegege/fashion_model]
ref:       gaocegege/fashion_model:v1
digest:    6b08cd25d01f71a09c1eb852b3a696ee2806abc749628de28a71b507f9eab996
size:      162.1 KiB
format:    SavedModel
v1: pushed to remote (1 layer, 162.1 KiB total)
```

Taking [Harbor](https://github.com/goharbor/harbor) as an example, we can see the model's metadata in Harbor registry.

<p align="center">
<img src="/docs/images/intro/harbor.png" height="350">
</p>

Then, we can download the model from the registry. The download process is similar to the push.

```bash
# Pull the model from remote registry to local cache
$ ormb pull gaocegege/fashion_model:v1
v1: Pulling from gaocegege/fashion_model
ref:     gaocegege/fashion_model:v1
digest:  6b08cd25d01f71a09c1eb852b3a696ee2806abc749628de28a71b507f9eab996
size:    162.1 KiB
Status: Downloaded newer model for gaocegege/fashion_model:v1

# Export the model from local cache to current directory
$ ormb export gaocegege/fashion_model:v1
ref:     localhost/gaocegege/fashion_model:v1
digest:  6b08cd25d01f71a09c1eb852b3a696ee2806abc749628de28a71b507f9eab996
size:    162.1 KiB

# View the local file directory
$ tree examples/SavedModel-fashion
examples/SavedModel-fashion
├── model
│   ├── saved_model.pb
│   └── variables
│       ├── variables.data-00000-of-00001
│       └── variables.index
├── ormbfile.yaml
└── training-serving.ipynb

2 directories, 5 files
```

## Installation

### Install the image registry

`ORMB` uses the image registry to store model, you can choose to [install Harbor](https://github.com/goharbor/harbor-helm) or [use Docker Registry](https://docs.docker.com/registry/deploying/). We recommended Harbor here.

### Install `ORMB`

You can install the pre-compiled binary, or compile from source.

#### Install the pre-compiled binary

Download the pre-compiled binaries from [the releases](https://github.com/kleveross/ormb/releases) page and copy to the desired location.

#### Compile from source

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

## Understanding ORMB

### Why choose ORMB?

See [`ORMB` introduction](/docs/introduction.md) for more information.

### Official Documentation

See [`ORMB` docs](/docs/README.md) for more information.

### Tutorials

* Distribute models with `ORMB` and Docker Registry: [tutorial.md](docs/tutorial.md)
* Serving model with Seldon Core: [tutorial-serving-seldon.md](docs/tutorial-serving-seldon.md)

### OCI Model Configuration Specification

Please have a look at [docs/spec_v1alpha1.md](docs/spec-v1alpha1.md)

## Community

`ORMB` project is part of Klever, a Cloud Native Machine Learning platform.

The Klever slack workspace is klever.slack.com. To join, click this [invitation to our Slack workspace](https://join.slack.com/t/kleveross/shared_invite/zt-g0eoiyq9-9OwiI7c__oV79bh_94MyTw).
