#!/usr/bin/env bash

# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# The Azure provided machines typically have the following disk allocation:
# Total space: 85GB
# Allocated: 67 GB
# Free: 17 GB
# This script frees up 28 GB of disk space by deleting unneeded packages and 
# large directories.
# The Flink end to end tests download and generate more than 17 GB of files,
# causing unpredictable behavior and build failures.

echo "=============================================================================="
echo "Freeing up disk space on CI system"
echo "=============================================================================="

echo "Listing 100 largest packages"
dpkg-query -Wf '${Installed-Size}\t${Package}\n' | sort -n | tail -n 100
df -h

echo "Removing large packages"
sudo apt-get remove -y 'humanity-icon-theme'
sudo apt-get remove -y '^dotnet-.*'
sudo apt-get remove -y 'php.*'
sudo apt-get remove -y '^mongodb-.*'
sudo apt-get remove -y '^mysql-.*' '^libmysql.*'
sudo apt-get remove -y '^postgresql-.*'
sudo apt-get remove -y '^ruby-.*'

# looping over package to avoid passing in a whole list together with some packages
# missing causing 'Unable to locate package' error and stopping other packages for being removed
for package in azure-cli google-chrome-stable firefox microsoft-edge-stable powershell mono-devel libgl1-mesa-dri iso-codes azure-cli google-cloud-cli; do
  sudo apt-get remove -y ${package}
done

sudo apt-get autoremove -y
sudo apt-get clean

df -h

echo "Removing large directories"
sudo rm -rf /usr/share/dotnet/
sudo rm -rf /usr/local/graalvm/
sudo rm -rf /usr/local/.ghcup/
sudo rm -rf /usr/local/share/powershell
sudo rm -rf /usr/local/share/chromium
sudo rm -rf /usr/local/lib/android
sudo rm -rf /usr/local/lib/node_modules
sudo rm -rf /var/lib/gems/
sudo rm -rf /opt/hostedtoolcache/go /opt/hostedtoolcache/Ruby /opt/hostedtoolcache/node

df -h

echo "Removing docs and man pages"
sudo rm -rf /usr/share/doc/
sudo rm -rf /usr/share/man/
sudo rm -rf /usr/share/locale/

df -h
