#!/bin/bash -eux

cwd=$(pwd)

pushd $cwd/dp-elasticsearch
  make lint
popd