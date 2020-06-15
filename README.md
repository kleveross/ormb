<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [ormb](#ormb)
  - [Installation](#installation)
  - [Usage](#usage)
    - [Save the model](#save-the-model)
    - [Push the model to a remote registry](#push-the-model-to-a-remote-registry)
    - [Pull the model from a remote registry](#pull-the-model-from-a-remote-registry)
    - [Export the model to the current directory](#export-the-model-to-the-current-directory)
  - [Tutorial](#tutorial)
  - [OCI Model Configuration Specification](#oci-model-configuration-specification)
  - [Community](#community)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# ormb

ormb is an open-source model registry to manage machine learning model. 

ormb helps you manage your Machine Learning/Deep Learning models. It makes your models easy to create, version, share and publish.

## Installation

You can install the pre-compiled binary, or compile from source.

### Install the pre-compiled binary

Download the pre-compiled binaries from [the releases](https://github.com/caicloud/ormb/releases) page and copy to the desired location.

### Compile from source

Clone:

```
$ git clone https://github.com/caicloud/ormb
$ cd goreleaser
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

## Usage

### Save the model

```bash
$ ormb save ./resnet_v2_fp32_savedmodel_NCHW caicloud/resnetv2:v1
ref:     caicloud/resnetv2:v1
digest:  f51973c855608ab06d8f5e4333925a635f87f01ff992ffc5f9988f26d1da24e9
size:    90.6 MiB
format:  SavedModel
v1: saved
```

### Push the model to a remote registry

```bash
$ ormb push caicloud/resnetv2:v1
The push refers to repository [caicloud/resnetv2]
ref:     caicloud/resnetv2:v1
digest:  f51973c855608ab06d8f5e4333925a635f87f01ff992ffc5f9988f26d1da24e9
size:    90.6 MiB
format:  SavedModel
v1: pushed to remote (1 layer, 90.6 MiB total)
```

### Pull the model from a remote registry

```bash
$ ormb pull caicloud/resnetv2:v1 
v1: Pulling from caicloud/resnetv2
ref:     caicloud/resnetv2:v1
digest:  f51973c855608ab06d8f5e4333925a635f87f01ff992ffc5f9988f26d1da24e9
size:    90.6 MiB
format:  SavedModel
Status: Downloaded newer model for caicloud/resnetv2:v1
```

### Export the model to the current directory

```bash
$ ormb export caicloud/resnetv2:v1
ref:     localhost:5000/caicloud/resnetv2:v1
digest:  f51973c855608ab06d8f5e4333925a635f87f01ff992ffc5f9988f26d1da24e9
size:    90.6 MiB

$ tree ./resnet_v2_fp32_savedmodel_NCHW
resnet_v2_fp32_savedmodel_NCHW
├── 1538687196
│   ├── saved_model.pb
│   └── variables
│       ├── variables.data-00000-of-00001
│       └── variables.index

2 directories, 3 files
```

## Tutorial

Please have a look at [docs/tutorial.md](docs/tutorial.md)

## OCI Model Configuration Specification

Please have a look at [docs/spec.md](docs/spec.md)

## Community

ormb project is part of Clever, a Cloud Native Machine Learning platform. We are going to open source a community edition soon.

The Clever slack workspace is caicloud-clever.slack.com. To join, click this [invitation to our Slack workspace](https://join.slack.com/t/caicloud-clever/shared_invite/zt-efz4rdrm-kcOg0Qvs_B8aIWGdZv9E6g).
