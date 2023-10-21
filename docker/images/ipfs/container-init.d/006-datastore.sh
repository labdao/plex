#!/bin/bash
set -e

if [ "${IPFS_DEBUG}" == "true" ]; then
  set -x
fi

if [ -n "${IPFS_DATASTORE_STORAGEMAX}" ]; then
  ipfs config Datastore.StorageMax "${IPFS_DATASTORE_STORAGEMAX}"
fi

if [ -N "${IPFS_DATASTORE_STORAGEGCWATERMARK}" ]; then
  ipfs config Datastore.StorageGCWatermark "${IPFS_DATASTORE_STORAGEGCWATERMARK}"
fi

if [ -n "${IPFS_DATASTORE_GCPERIOD}" ]; then
  ipfs config Datastore.GCPeriod "${IPFS_DATASTORE_GCPERIOD}"
fi
