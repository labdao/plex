#!/bin/bash

set -e

export ANSIBLE_VERSION="8.2.0"
export ANSIBLE_CHECKOUT="main"

echo "Installing awscli"
apt-get update
apt-get -y --no-install-recommends install awscli git python3-pip

echo "Install amazon-ssm-agent"
snap install amazon-ssm-agent --classic
systemctl enable snap.amazon-ssm-agent.amazon-ssm-agent.service
systemctl start snap.amazon-ssm-agent.amazon-ssm-agent.service

echo "Determine instance type info from tag"
INSTANCE_TYPE_TAG=$(curl -s http://instance-data/latest/meta-data/tags/instance/Type)

# Check if vars file needs loading - by checking if file with environment name exists
LOAD_EXTRA_VARS_FILE="false"
# If environment tag is not empty
if [[ "x${environment}" != "x" ]] ;then
  VARS_FILE="infrastructure/ansible/vars/${environment}.yaml"
  VARS_FILE_URL="https://raw.githubusercontent.com/labdao/plex/$${ANSIBLE_CHECKOUT}/$${VARS_FILE}"
  # Check if extra-vars file exists
  if [ $(curl -LI "$${VARS_FILE_URL}" -o /dev/null -w '%%{http_code}\n' -s) == "200" ]; then
    LOAD_EXTRA_VARS_FILE="true"
    wget "$${VARS_FILE_URL}"
  fi
fi

echo "Install Ansible $${ANSIBLE_VERSION}"
pip3 install ansible==$${ANSIBLE_VERSION}

echo "Run provided playbooks"

for playbook in infrastructure/ansible/install_requirements.yaml infrastructure/ansible/provision_$${INSTANCE_TYPE_TAG}.yaml; do
  ANSIBLE_TMP_DIR="$(mktemp -d)"

  # Check if vars file needs loading
  if [[ "$${LOAD_EXTRA_VARS_FILE}" == "true" ]] ;then

    echo "Running \"$${playbook}\" under temp directory $${ANSIBLE_TMP_DIR} with extra-vars file \"$${VARS_FILE}\""
    /usr/local/bin/ansible-pull \
     --accept-host-key \
     --verbose \
     --extra-vars 'target_hosts=localhost' \
     --extra-vars "@${environment}.yaml" \
     --connection 'local' \
     --url https://github.com/labdao/plex.git \
     --checkout "$${ANSIBLE_CHECKOUT}" \
     --directory "$${ANSIBLE_TMP_DIR}" \
     --purge \
     --limit localhost \
     "$${playbook}"
  else

    echo "Running \"$${playbook}\" under temp directory $${ANSIBLE_TMP_DIR}"
    /usr/local/bin/ansible-pull \
     --accept-host-key \
     --verbose \
     --extra-vars 'target_hosts=localhost' \
     --connection 'local' \
     --url https://github.com/labdao/plex.git \
     --checkout "$${ANSIBLE_CHECKOUT}" \
     --directory "$${ANSIBLE_TMP_DIR}" \
     --purge \
     --limit localhost \
     "$${playbook}"
  fi

  # remove ansible temp directory
  rm -rf "$${ANSIBLE_TMP_DIR}"
done
