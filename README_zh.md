<p align="center">
<img src="docs/images/logo.png" height="150">
</p>


[![Build Status](https://travis-ci.com/kleveross/ormb.svg?branch=master)](https://travis-ci.com/kleveross/ormb)
[![Coverage Status](https://coveralls.io/repos/github/kleveross/ormb/badge.svg?branch=master)](https://coveralls.io/github/kleveross/ormb?branch=master)

[English](./README.md) | 中文

`ormb` 是一个用于管理机器学习模型的开源模型仓库。

`ormb` 可以帮助用户更好的管理他们的机器学习 / 深度学习模型。通过 `ormb`，模型能更易于创建、版本化、共享以及发布。

[![asciicast](https://asciinema.org/a/345812.svg)](https://asciinema.org/a/345812)

## 安装

您可以安装预编译的二进制文件，或是从源代码进行编译。

### 安装二进制文件

从 [releases](https://github.com/kleveross/ormb/releases) 页面中下载预编译的二进制文件并且将其复制到所需的位置。

### 从源代码编译

下载源码:

```
$ git clone https://github.com/kleveross/ormb
$ cd ormb
```

安装依赖:

```
$ go mod tidy
```

编译:

```
$ make build-local
```

验证运行:

```
$ ./bin/ormb --help
```

## 介绍

[ormb 介绍](./docs/intro_zh.md)

## 教程

### 使用 `ormb` 和 Docker 镜像仓库分发模型

请查阅 [tutorial.md](docs/tutorial.md)

### 使用 `Seldon Core` 启动模型服务

请查阅 [tutorial-serving-seldon.md](docs/tutorial-serving-seldon.md)

## OCI 模型配置规范

请查阅 [spec_v1alpha1.md](docs/spec-v1alpha1.md)

## 社区

`ormb` 是 Klever 云原生机器学习平台的子项目。

Klever 的 Slack 是 klever.slack.com. 请点击 [邀请链接](https://join.slack.com/t/kleveross/shared_invite/zt-g0eoiyq9-9OwiI7c__oV79bh_94MyTw) 加入 Slack 讨论。
