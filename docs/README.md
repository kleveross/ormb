# ORMB

English | [中文](docs_zh/README.md)

`ORMB` is a project that uses image registry to distribute machine learning models. Users can upload and download models through commands `ormb pull` and `ormb push` just like `docker pull` and `docker push`.

Check out the [ORMB Introduction Doc](introduction.md) to see why you want to use `ORMB`.

## What is ormbfile?

`ormbfile.yaml`, which is similar to `Dockerfile`, is a configuration file defined by `ORMB`. Users can describe the model information by filling in `format`, `framework`, `signature` and other fields in `ormbfile.yaml`. Refer to [spec-v1alpha1.md](../spec-v1alpha1.md) for detailed fields.

## `ORMB` commands

### ormb login

`ORMB` uses image registry to store models and distributes machine learning models just like Docker images. You can login to the image registry through  `ormb login`.

```bash
$ ormb login  --insecure harbor.ormb.xyz -u root -p password
Using config file: /Users/Fog/.ormb/config.json
Using /Users/Fog/.ormb as the root path
WARNING! Using --password via the CLI is insecure. Use --password-stdin.
Login insecurely
Login succeeded
```

### ormb save

`ormb save` will store the model files in the current directory into the local cache and simultaneously parses `ormbfile.yaml` to get the model metadata. The model will be converted to the following format:

```go
type Model struct {
	// Metadata is the contents of the Chartfile.
	Metadata *Metadata `json:"metadata,omitempty"`
	Path     string    `json:"path,omitempty"`
	Content  []byte    `json:"content,omitempty"`
	Config   []byte    `json:"config,omitempty"`
}
```

Then the model will be stored in the cache in the following format:

```go
// CacheRefSummary contains as much info as available describing a chart reference in cache
// Note: fields here are sorted by the order in which they are set in FetchReference method
type CacheRefSummary struct {
	Name         string
	Repo         string
	Tag          string
	Exists       bool
	Manifest     *ocispec.Descriptor
	Config       *ocispec.Descriptor
	ContentLayer *ocispec.Descriptor
	Size         int64
	Digest       digest.Digest
	CreatedAt    time.Time
	Model        *model.Model
}
```

Command format:

```bash
ormb save projectName/modelName:modelVersion
```

### ormb tag

`ormb tag` is similar to `docker tag`, the tag of the model can be changed by the command.

Command format:

```bash
ormb tag ormbtest/fashion_model:v1 ormbtest/fashion_model:v2 
```

### ormb push

`ormb push` pushes the models from local cache to the remote registry.

Command format:

```bash
ormb push projectName/modelName:modelVersion
```

### ormb pull

`ormb pull` pulls the models from the remote registry to local cache.

Command format:

```bash
ormb pull projectName/modelName:modelVersion
```

### ormb export

`ormb export` exports the models from local cache to the current directory.

Command format:

```bash
ormb export projectName/modelName:modelVersion
```
## Contributing to ORMB

`ORMB` is an open source project based on [Apache 2.0 open source licenses](https://www.apache.org/licenses/LICENSE-2.0). If you'd like to contribute to the `ORMB` project, please refer to our [Contributor's guide](/CONTRIBUTING.md).
