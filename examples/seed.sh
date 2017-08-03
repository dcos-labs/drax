#!/bin/bash
#
# Launches all Marathon app specs in the directory.
# These Marathon apps serve as a test bed for DRAX.

set -o errexit
set -o pipefail
set -o nounset

command -v dcos >/dev/null 2>&1 || { echo >&2 "Need the DC/OS CLI to carry out my work, but it seems it's not installed. See https://dcos.io/docs/latest/cli/install/ for how to get it ..."; exit 1; }

for ma in *.json ; do
  echo "Trying to launching Marathon app defined in: $ma"
  dcos marathon app add $ma
  echo "$ma launched"
done
