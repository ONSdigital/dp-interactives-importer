#!/bin/bash -eux

pushd dp-interactives-importer
  make build
  cp build/dp-interactives-importer Dockerfile.concourse ../build
popd
