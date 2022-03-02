#!/bin/bash -eux

pushd dp-interactives-importer
  make build
  cp build/dp-interactives-importer docker/mime.types Dockerfile.concourse ../build
popd
