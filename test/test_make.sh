#!/bin/bash

set -e

ROOT=$(dirname "${BASH_SOURCE}")/..

function test_make() {
  cd $ROOT
  make lint
  make test
  make build
  make build-linux
  make container
  make clean
  cd ..
}

test_make
