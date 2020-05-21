<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Golang Template Project](#golang-template-project)
  - [About the project](#about-the-project)
    - [API docs](#api-docs)
    - [Design](#design)
    - [Status](#status)
    - [See also](#see-also)
  - [Getting started](#getting-started)
    - [Layout](#layout)
  - [Notes](#notes)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Golang Template Project

## About the project

The template is used to create golang project. All golang projects must follow the conventions in the
template. Calling for exceptions must be brought up in the engineering team.

### API docs

The template doesn't have API docs. For web service, please include API docs here, whether it's
auto-generated or hand-written. For auto-generated API docs, you can also give instructions on the
build process.

### Design

The template follows project convention doc.

* [Repository Conventions](https://github.com/caicloud/engineering/blob/master/guidelines/repo_conventions.md)

### Status

The template project is in alpha status.

### See also

* [nirvana project template](https://github.com/caicloud/nirvana-template-project)
* [python project template](https://github.com/caicloud/python-template-project)
* [common project template](https://github.com/caicloud/common-template-project)

## Getting started

Below we describe the conventions or tools specific to golang project.

### Layout

```tree
├── .github
│   ├── ISSUE_TEMPLATE.md
│   └── PULL_REQUEST_TEMPLATE.md
├── .gitignore
├── .golangci.yml
├── CHANGELOG.md
├── Makefile
├── OWNERS
├── README.md
├── build
│   ├── admin
│   │   └── Dockerfile
│   └── controller
│       └── Dockerfile
├── cmd
│   ├── admin
│   │   └── admin.go
│   └── controller
│       └── controller.go
├── docs
│   └── README.md
├── hack
│   ├── README.md
│   ├── deployment.yaml
│   └── script.sh
├── pkg
│   ├── apis
│   │   └── v1
│   │       └── README.md
│   ├── utils
│   │   └── net
│   │       └── net.go
│   └── version
│       └── version.go
├── release
│   ├── template-admin.yaml
│   └── template-controller.yaml
├── test
│   ├── README.md
│   └── test_make.sh
├── third_party
│   └── README.md
└── vendor
    └── README.md
```

A brief description of the layout:

* `.github` has two template files for creating PR and issue. Please see the files for more details.
* `.gitignore` varies per project, but all projects need to ignore `bin` directory.
* `.golangci.yml` is the golangci-lint config file.
* `Makefile` is used to build the project. **You need to tweak the variables based on your project**.
* `CHANGELOG.md` contains auto-generated changelog information.
* `OWNERS` contains owners of the project.
* `README.md` is a detailed description of the project.
* `bin` is to hold build outputs.
* `cmd` contains main packages, each subdirecoty of `cmd` is a main package.
* `build` contains scripts, yaml files, dockerfiles, etc, to build and package the project.
* `docs` for project documentations.
* `hack` contains scripts used to manage this repository, e.g. codegen, installation, verification, etc.
* `pkg` places most of project business logic and locate `api` package.
* `release` [chart](https://github.com/caicloud/charts) for production deployment.
* `test` holds all tests (except unit tests), e.g. integration, e2e tests.
* `third_party` for all third party libraries and tools, e.g. swagger ui, protocol buf, etc.
* `vendor` contains all vendored code.

## Notes

* Makefile **MUST NOT** change well-defined command semantics, see Makefile for details.
* Every project **MUST** use `dep` for vendor management and **MUST** checkin `vendor` direcotry.
* `cmd` and `build` **MUST** have the same set of subdirectories for main targets
  * For example, `cmd/admin,cmd/controller` and `build/admin,build/controller`.
  * Dockerfile **MUST** be put under `build` directory even if you have only one Dockerfile.
