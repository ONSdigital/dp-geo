#!/bin/bash -eux

cwd=$(pwd)

pushd $cwd/dp-geo
  make audit
popd