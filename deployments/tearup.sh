#!/usr/bin/env bash

set -e

cmd=$1

current_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
path_to_helm="$current_dir/$cmd"

if [ ! -d "$path_to_helm" ]; then
    echo "Error: command (directory) $1 does not exists."
    exit 9999 # die with error code 9999
fi

echo "tearup snappcloud-status-backend-$cmd ..."

helm upgrade --install "snappcloud-status-backend-$cmd" "$path_to_helm"
