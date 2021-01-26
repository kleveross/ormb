# ORMB

[English](../README.md) | 中文

`ORMB` 是一个使用镜像仓库来分发机器学习模型的项目。用户可以通过和 `docker pull`， `docker push` 类似的 `ormb` 命令 `ormb pull`，`ormb push` 来上传和下载模型。

查阅我们的 [介绍文档]() 来了解为什么要使用 `ORMB`。

## 什么是 ormbfile？

`ormbfile.yaml` 类似于 `Dockerfile`，是 `ORMB` 定义的配置文件。用户通过填写 `ormbfile.yaml` 中的 `format`，`framework`，`signature` 等字段来更好的描述模型信息。详细字段请参考 [spec-v1alpha1.md](../spec-v1alpha1.md)

## `ORMB` 命令

### ormb login

`ORMB` 使用镜像仓库来存储模型，像分发 Docker 镜像一样分发机器学习模型。`ormb login` 命令用于登录镜像仓库。

```bash
$ ormb login  --insecure harbor.ormb.xyz -u root -p password
Using config file: /Users/Fog/.ormb/config.json
Using /Users/Fog/.ormb as the root path
WARNING! Using --password via the CLI is insecure. Use --password-stdin.
Login insecurely
Login succeeded
```

### ormb save

`ormb save` 会将当前目录下的模型文件存储到本地缓存中，并同时通过解析 `ormbfile.yaml` 得到模型的 `metadata`。模型会先转换成如下格式：

```go
type Model struct {
	// Metadata is the contents of the Chartfile.
	Metadata *Metadata `json:"metadata,omitempty"`
	Path     string    `json:"path,omitempty"`
	Content  []byte    `json:"content,omitempty"`
	Config   []byte    `json:"config,omitempty"`
}
```

并以如下格式存储在缓存中：

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

命令格式：

```bash
ormb save <model directory> projectName/modelName:modelVersion
```

### ormb tag

`ormb tag` 与 `docker tag` 类似，可以更改模型的 tag。

命令格式：

```bash
ormb tag ormbtest/fashion_model:v1 ormbtest/fashion_model:v2 
```

### ormb push

`ormb push` 会将存储在本地缓存中的模型推送到远端仓库中。

命令格式：

```bash
ormb push projectName/modelName:modelVersion
```

### ormb pull

`ormb pull` 会将存储在远端仓库的模型拉取到本地缓存中。

命令格式：

```bash
ormb pull projectName/modelName:modelVersion
```

### ormb export

`ormb export` 会将存储在缓存中的模型导出到当前目录。

命令格式：

```bash
ormb export projectName/modelName:modelVersion
```

## 成为贡献者之一

`ORMB` 是一个遵循 [Apache 2.0 开源协议](https://www.apache.org/licenses/LICENSE-2.0) 的开源项目。如果您愿意为 `ORMB` 项目做出贡献，请参阅我们的 [贡献者指南](/CONTRIBUTING.md)
