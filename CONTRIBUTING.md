# Contributing to ORMB

Thank you for contributing your time and expertise to ORMB. This document describes the contribution guidelines for the project.

## Contributing to ORMB code

### Commit Conventions

Below we describe commit message conventions.

#### Goals

* improve repository maintainability
* provide better history information
* allow auto-generating CHANGELOG.md
* allow ignoring commits by git bisect (e.g. not important commits like formatting)

#### Format of commit message

```
<type>(<scope>): <subject>
<BLANK LINE>
<body>
<BLANK LINE>
```

##### `<type>` (required)

Type is required to better capture the area of the commit. Must be one of the following:

* **feat**: A new feature
* **fix**: A bug fix
* **docs**: Documentation only changes
* **style**: Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)
* **refactor**: A code change that neither fixes a bug nor adds a feature
* **perf**: A code change that improves performance
* **test**: Adding missing or correcting existing tests
* **chore**: Changes to the build process or auxiliary tools and libraries such as documentation generation

##### `<scope>` (optional)

Scope is optional, it is could be anything specifying place of the commit change. Github issue link is
also a valid scope. For example: fix(cli), feat(api), fix(#101), etc.

You can use `*` when the change affects more than a single scope.

##### `<subject>` (required)

Subject line contains succinct description of the change.

* be descriptive, e.g. "fix(controller): fix incorrect image name causing image pull error", not "fix small bug"
* use imperative, present tense: “change” not “changed” nor “changes”
* don't capitalize first letter
* no dot (.) at the end

##### `<body>` (optional)

Body messge is optional, it should provide detailed desciption for large change

* use imperative, present tense: “change” not “changed” nor “changes”
* includes motivation for the change and contrasts with previous behavior

## Examples

```
docs(golang): reword golang getting started guide
```

```
feat(#121): set failure status for worker job

Failure status for worker job is not set, possible statuses are: Accepted, Running, Failed.
```

```
fix(apidocs): show api docs url by default
```
