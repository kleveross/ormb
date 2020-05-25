<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Tutorial](#tutorial)
  - [Prepare the model](#prepare-the-model)
  - [Login to a remote registry.](#login-to-a-remote-registry)
  - [Package and cache the model](#package-and-cache-the-model)
  - [Push the model to a remote Docker Registry](#push-the-model-to-a-remote-docker-registry)
  - [Remove the model from local cache and pull it from the remote registry](#remove-the-model-from-local-cache-and-pull-it-from-the-remote-registry)
  - [Export the model to the directory](#export-the-model-to-the-directory)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Tutorial

## Prepare the model

The model directory need to be prepared. In this tutorial, we use [examples/PMML-model](../examples/PMML-model) as the model.

```bash
$ tree ./examples/PMML-model
PMML-model
├── model
│   └── single_audit_mlp.pmml
└── ormbfile.yaml
```

`ormbfile.yaml` is required in the model directory:

```bash
$ cat ./examples/PMML-model/ormbfile.yaml
author: "Ce Gao <gaoce@caicloud.io>"
format: PMML
```

In this case, we set the author and format field of the model.

## Login to a remote registry.

We need to login to a remote registry. We can run the command to setup a local Docker registry:

```bash
$ docker run -it --rm -p 5000:5000 registry
WARN[0000] No HTTP secret provided - generated random secret. This may cause problems with uploads if multiple registries are behind a load-balancer. To provide a shared secret, fill in http.secret in the configuration file or set the REGISTRY_HTTP_SECRET environment variable.  go.version=go1.11.2 instance.id=3a703617-a91b-4735-a4f2-43bf7c80f027 service=registry version=v2.7.1
INFO[0000] redis not configured                          go.version=go1.11.2 instance.id=3a703617-a91b-4735-a4f2-43bf7c80f027 service=registry version=v2.7.1
INFO[0000] Starting upload purge in 33m0s                go.version=go1.11.2 instance.id=3a703617-a91b-4735-a4f2-43bf7c80f027 service=registry version=v2.7.1
INFO[0000] using inmemory blob descriptor cache          go.version=go1.11.2 instance.id=3a703617-a91b-4735-a4f2-43bf7c80f027 service=registry version=v2.7.1
INFO[0000] listening on [::]:5000                        go.version=go1.11.2 instance.id=3a703617-a91b-4735-a4f2-43bf7c80f027 service=registry version=v2.7.1
```

## Package and cache the model

Now we can save the model to cache in the local file system.

```bash
$ ormb save ./examples/PMML-model localhost:5000/ormb/pmml:v1
ref:     localhost:5000/ormb/pmml:v1
digest:  c41f49e5d49d301d58b9b37590c833f2a85c3d083e09b680735f190fbdb8158a
size:    6.1 KiB
v1: saved
```

## Push the model to a remote Docker Registry

After we store the model in the local cache, we can push it to a remote registry.

```bash
$ ormb push localhost:5000/ormb/pmml:v1
The push refers to repository [localhost:5000/ormb/pmml]
ref:     localhost:5000/ormb/pmml:v1
digest:  c41f49e5d49d301d58b9b37590c833f2a85c3d083e09b680735f190fbdb8158a
size:    6.1 KiB
v1: pushed to remote (1 layer, 6.1 KiB total)
```

## Remove the model from local cache and pull it from the remote registry

We can remove the model from the local cache, then try to pull it from the remote registry.

```bash
$ rm -rf cache/
$ ormb pull localhost:5000/ormb/pmml:v1
v1: Pulling from localhost:5000/ormb/pmml
ref:     localhost:5000/ormb/pmml:v1
digest:  c41f49e5d49d301d58b9b37590c833f2a85c3d083e09b680735f190fbdb8158a
size:    6.1 KiB
Status: Downloaded newer model for localhost:5000/ormb/pmml:v1
```

## Export the model to the directory

We can export the model from the local cache to the target directory.

```bash
$ ormb pull localhost:5000/ormb/pmml:v1
ref:     localhost:5000/ormb/pmml:v1
digest:  c41f49e5d49d301d58b9b37590c833f2a85c3d083e09b680735f190fbdb8158a
size:    6.1 KiB
```
