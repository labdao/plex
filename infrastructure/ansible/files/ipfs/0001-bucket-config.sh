#!/bin/sh
set -ex

echo "IPFS PATH: ${IPFS_PATH}"

# We backup old config file
cp ${IPFS_PATH}/config ${IPFS_PATH}/config_bak

# We inject the S3 plugin datastore
# Important: Make sure your fill out the optionnal parameters $CLUSTER_S3_BUCKET, $CLUSTER_AWS_KEY, $CLUSTER_AWS_SECRET in the cloudformation parameters
cat ${IPFS_PATH}/config_bak | \
jq ".Datastore.Spec = {
    mounts: [
        {
          child: {
            type: \"s3ds\",
            region: \"{{ ansible_ec2_placement_region }}\",
            bucket: \"labdao-ipfs-bucket-{{ ansible_ec2_tags_instance_Env | lower }}\",
            rootDirectory: \"${CLUSTER_PEERNAME}\",
            accessKey: \"${CLUSTER_AWS_KEY}\",
            secretKey: \"${CLUSTER_AWS_SECRET}\",
            keyTransform: \"next-to-last/2\"
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
}" > ${IPFS_PATH}/config

# We override the ${IPFS_PATH}/datastore_spec file
echo "{\"mounts\":[{\"bucket\":\"labdao-ipfs-bucket-{{ ansible_ec2_tags_instance_Env | lower }}\",\"mountpoint\":\"/blocks\",\"region\":\"{{ ansible_ec2_placement_region }}\",\"rootDirectory\":\"${CLUSTER_PEERNAME}\"},{\"mountpoint\":\"/\",\"path\":\"datastore\",\"type\":\"levelds\"}],\"type\":\"mount\"}" > ${IPFS_PATH}/datastore_spec
