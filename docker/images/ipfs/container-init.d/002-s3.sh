#!/bin/sh
set -e

if [ "${IPFS_DEBUG}" == "true" ]; then
  set -x
fi

if [ "${IPFS_S3_ENABLED}" == "true" ]; then
  echo "IPFS PATH: ${IPFS_PATH}"

  # Check if running in ECS, set credentialsEndpoint
  if [ "${AWS_EXECUTION_ENV}" == "AWS_ECS_FARGATE" ]; then
    export CLUSTER_CREDENTIALS_ENDPOINT=http://169.254.170.2${AWS_CONTAINER_CREDENTIALS_RELATIVE_URI}
  fi

  # We backup old config file
  cp "${IPFS_PATH}"/config "${IPFS_PATH}"/config_bak

  # We inject the S3 plugin datastore
  # Important: Make sure your fill out the optionnal parameters $CLUSTER_S3_BUCKET, $CLUSTER_AWS_KEY, $CLUSTER_AWS_SECRET in the cloudformation parameters
  cat "${IPFS_PATH}"/config_bak | \
  jq ".Datastore.Spec = {
      mounts: [
          {
            child: {
              type: \"s3ds\",
              region: \"${AWS_REGION}\",
              bucket: \"${CLUSTER_S3_BUCKET}\",
              rootDirectory: \"${CLUSTER_PEERNAME}\",
              accessKey: \"${CLUSTER_AWS_KEY}\",
              secretKey: \"${CLUSTER_AWS_SECRET}\",
              keyTransform: \"${CLUSTER_KEY_TRANSFORM:-next-to-last/2}\",
              credentialsEndpoint: \"${CLUSTER_CREDENTIALS_ENDPOINT}\"
            },
            mountpoint: \"/blocks\",
            prefix: \"s3.datastore\",
            type: \"measure\"
          },
          {
            child: {
              compression: \"none\",
              path: \"datastore\",
              type: \"levelds\"
            },
            mountpoint: \"/\",
            prefix: \"leveldb.datastore\",
            type: \"measure\"
          }
      ],
      type: \"mount\"
  }" > "${IPFS_PATH}"/config

  # We override the ${IPFS_PATH}/datastore_spec file
  echo "{\"mounts\":[{\"bucket\":\"${CLUSTER_S3_BUCKET}\",\"mountpoint\":\"/blocks\",\"region\":\"${AWS_REGION}\",\"rootDirectory\":\"${CLUSTER_PEERNAME}\"},{\"mountpoint\":\"/\",\"path\":\"datastore\",\"type\":\"levelds\"}],\"type\":\"mount\"}" > "${IPFS_PATH}"/datastore_spec
fi
