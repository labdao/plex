#!/bin/bash

echo "Installing awscli"
apt-get update
apt-get -y --no-install-recommends install awscli git python3-pip

echo "Install amazon-ssm-agent"
snap install amazon-ssm-agent --classic
systemctl enable snap.amazon-ssm-agent.amazon-ssm-agent.service
systemctl start snap.amazon-ssm-agent.amazon-ssm-agent.service

echo "Install Ansible"
pip3 install ansible
