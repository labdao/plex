#!/bin/bash
# This script is tested to run on a fresh Ubuntu 20.04 LTS Nvidia GPU compute instance
# Specifically tested on an AWS P3, G4, and G5 instance

# this exception allows exit 77 to exit the whole script within subshells
# https://unix.stackexchange.com/questions/48533/exit-shell-script-from-a-subshell
set -E
trap '[ "$?" -ne 77 ] || exit 77' ERR

installDocker() {
    # Docker install directions from https://docs.docker.com/engine/install/ubuntu/
    echo "Installing Docker"
    sudo apt-get update
    sudo apt-get install -y \
        ca-certificates \
        curl \
        gnupg \
        lsb-release
    sudo mkdir -p /etc/apt/keyrings
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --yes --dearmor -o /etc/apt/keyrings/docker.gpg
    echo \
        "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
        $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    sudo apt-get update
    sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin
    sudo groupadd -f docker
    sudo usermod -aG docker $USER
}

testDockerInstall() {
    echo "Testing Docker Install"
    /usr/bin/newgrp docker <<EONG
    if docker run hello-world ; then
        echo "Docker succesfully installed"
    else
        echo "Docker install failed"
        exit 77
    fi
EONG
}

installNvidiaDrivers() {
    # Nvidia driver instructions from https://docs.nvidia.com/datacenter/tesla/tesla-installation-notes/index.html
    echo "Installing Nvidia Drivers"
    sudo apt install -y build-essential
    sudo apt-get install linux-headers-$(uname -r)
    distribution=$(. /etc/os-release;echo $ID$VERSION_ID | sed -e 's/\.//g')
    wget https://developer.download.nvidia.com/compute/cuda/repos/$distribution/x86_64/cuda-keyring_1.0-1_all.deb
    sudo dpkg -i cuda-keyring_1.0-1_all.deb
    sudo apt-get update
    sudo apt-get -y install cuda-drivers
    export PATH=/usr/local/cuda-12.0/bin${PATH:+:${PATH}}
    sudo /usr/bin/nvidia-persistenced --verbose
}

testNvidiaInstall() {
    echo "Testing Nvidia Driver Install"
    if cat /proc/driver/nvidia/version; then
        echo "Nvidia succesfully installed"
    else
        echo "Nvidia install failed"
        exit
    fi
}

installNvidiaContainerToolKit() {
    # Nvidia Container Toolkit instructions from https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/install-guide.html
    echo "Installing Nvidia Container Toolkit"
    distribution=$(. /etc/os-release;echo $ID$VERSION_ID) \
        && curl -fsSL https://nvidia.github.io/libnvidia-container/gpgkey | sudo gpg --yes --dearmor -o /usr/share/keyrings/nvidia-container-toolkit-keyring.gpg \
        && curl -s -L https://nvidia.github.io/libnvidia-container/$distribution/libnvidia-container.list | \
            sed 's#deb https://#deb [signed-by=/usr/share/keyrings/nvidia-container-toolkit-keyring.gpg] https://#g' | \
            sudo tee /etc/apt/sources.list.d/nvidia-container-toolkit.list
    sudo apt-get update
    sudo apt-get install -y nvidia-docker2
    sudo systemctl restart docker
}

testNvidiaContainerToolkitInstall() {
    /usr/bin/newgrp docker <<EONG
    echo "Testing Nvidia Container Toolkit Install"
    if docker run --gpus all nvidia/cuda:11.6.2-base-ubuntu20.04 nvidia-smi ; then
        echo "Nvidia Container Toolkit succesfully installed"
    else
        echo "Nvidia Container Toolkit install failed"
        exit 77
    fi
EONG
}

installGolang() {
    echo "Installing GoLang"
    wget https://go.dev/dl/go1.19.6.linux-amd64.tar.gz
    sudo tar -C /usr/local -xvzf go1.19.6.linux-amd64.tar.gz
    rm go1.19.6.linux-amd64.tar.gz
    export PATH=$PATH:/usr/local/go/bin # for current shell
    echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bashr # for future shells
}

testGolangInstall() {
    echo "Testing GoLang Install"
    if go version ; then
        echo "GoLang succesfully installed "
    else
        echo "GoLang install failed"
        exit 77
    fi
}

installIPFS() {
    echo "Installing IPFS"
    wget https://dist.ipfs.tech/kubo/v0.18.0/kubo_v0.18.0_linux-amd64.tar.gz
    tar -xvzf kubo_v0.18.0_linux-amd64.tar.gz
    cd kubo
    sudo bash install.sh
    ipfs --version
    ipfs cat /ipfs/QmQPeNsJPyVWPFDVHb77w8G42Fvo15z4bG2X8D2GhfbSXc/readme
    export IPFS_CONNECT=/ip4/127.0.0.1/tcp/5001 # for current shell
    echo "export IPFS_CONNECT=/ip4/127.0.0.1/tcp/5001" >> ~/.bashrc # for future shells
}

testIPFSInstall() {
    echo "Testing IPFS Install"
    if ipfs --version ; then
        echo "IPFS succesfully installed "
    else
        echo "IPFS install failed"
        exit 77
    fi
}

runIPFS() {
    echo "Starting IPFS Daemon"
    ipfs init
    ipfs config Addresses.API /ip4/0.0.0.0/tcp/5001
    ipfs config Addresses.Gateway /ip4/0.0.0.0/tcp/8080
    ipfs config --json API.HTTPHeaders.Access-Control-Allow-Methods '["PUT", "POST"]'
    ipfs config Pinning.Recursive true
    screen -dmS ipfs ipfs daemon --routing=dhtclient
}

installBacalhau() {
    curl -sL https://get.bacalhau.org/install.sh | bash
}

testBacalhauInstall() {
    echo "Testing Bacalhau Install"
    if bacalhau version ; then
        echo "Bacalhau succesfully installed "
    else
        echo "Bacalhau install failed"
        exit 77
    fi
}

runBacalhau() {
    owner="labdaostage"
    if [ $PLEX_ENV = "prod" ]; then
        owner="labdao"
    fi
    screen -dmS bacalhau bacalhau serve --node-type compute,requester --ipfs-connect $IPFS_CONNECT --limit-total-gpu 1 --limit-job-memory 12gb --job-selection-accept-networked --job-selection-data-locality anywhere --labels owner=$owner
}

installConda() {
    # maybe switch to miniconda
    echo "Installing Conda"
    wget https://repo.anaconda.com/archive/Anaconda3-2022.10-Linux-x86_64.sh
    bash Anaconda3-2022.10-Linux-x86_64.sh -b
    export PATH=~/anaconda3/bin:$PATH # for current shell
    echo "export PATH=~/anaconda3/bin:$PATH" >> ~/.bashrc # for future shells
}

testCondaInstall() {
    if conda --version ; then
        echo "Conda succesfully installed "
    else
        echo "Conda install failed"
        exit 77
    fi
}

installJuypter() {
    conda install -c conda-forge jupyter-book --yes
    conda install -c conda-forge jupyterlab --yes
}

testJuypterInstall() {
    if jupyter --version ; then
        echo "Jupyter succesfully installed "
    else
        echo "Jupyter install failed"
        exit 77
    fi
}

runJuypter() {
    jupyter notebook password
    screen -dmS jupyter jupyter lab --ip=0.0.0.0
}

printLogo() {
logo="
                                        @
                                 @@@@@@@@@@@@@@@
                               @@@@@@@@@@@@@@@@@@@
                              @@@@@@@@@@@@@@@@@@@@@
             @@@@@@@@@@      @@@@@@@@@@@@@@@@@@@@@@@      @@@@@@@@@@
           @@@@@@@@@@@@      @@@@@@@@@@@@@@@@@@@@@@@      @@@@@@@@@@@@
         @@@@@@@@@@@@@@      @@@@@@@@@@@@@@@@@@@@@@@      @@@@@@@@@@@@@@
        *@@@@@@@@@@@@@      @@@@@@@@@@@@@@@@@@@@@@         @@@@@@@@@@@@@
         @@@@@@@@@@        @@@@@@@@@@@@@@@@@@@@@%            &@@@@@@@@@@
           @@@@           @@@@@@@@@@@@@@@@@@&                     @@@@
                        @@@@@@@@
                   @@@@@@@@@
      @@@@@@@@@@@@@@@@@@@@        ,@@@@@@@@@@@                 @@@@@@@@@@@@
   @@@@@@@@@@@@@@@@@@@@@@       @@@@@@@@@@@@@@@@@           @@@@@@@@@@@@@@@@@@
  @@@@@@@@@@@@@@@@@@@@@@      @@@@@@@@@@@@@@@@@@@@@       @@@@@@@@@@@@@@@@@@@@@
 @@@@@@@@@@@@@@@@@@@@@@@     @@@@@@@@@@@@@@@@@@@@@@@      @@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@     @@@@@@@@@@@@@@@@@@@@@@@     @@@@@@@@@@@@@@@@@@@@@@@
 @@@@@@@@@@@@@@@@@@@@@@      @@@@@@@@@@@@@@@@@@@@@@@     @@@@@@@@@@@@@@@@@@@@@@@
  @@@@@@@@@@@@@@@@@@@@@       @@@@@@@@@@@@@@@@@@@@@      @@@@@@@@@@@@@@@@@@@@@@
   @@@@@@@@@@@@@@@@@@           @@@@@@@@@@@@@@@@@       @@@@@@@@@@@@@@@@@@@@@@
      @@@@@@@@@@@@                 @@@@@@@@@@@         @@@@@@@@@@@@@@@@@@@@
                                                     @@@@@@@@@
                                                 @@@@@@@@
           @@@@                     &@@@@@@@@@@@@@@@@@@           @@@@
         @@@@@@@@@@             @@@@@@@@@@@@@@@@@@@@@        &@@@@@@@@@@
        *@@@@@@@@@@@@@        @@@@@@@@@@@@@@@@@@@@@@@      @@@@@@@@@@@@@
         @@@@@@@@@@@@@@      @@@@@@@@@@@@@@@@@@@@@@@      @@@@@@@@@@@@@@
           @@@@@@@@@@@@      @@@@@@@@@@@@@@@@@@@@@@@      @@@@@@@@@@@@
             @@@@@@@@@@      @@@@@@@@@@@@@@@@@@@@@@@      @@@@@@@@@@
                              @@@@@@@@@@@@@@@@@@@@@
                               @@@@@@@@@@@@@@@@@@@
                                 @@@@@@@@@@@@@@@
                                        @"
  echo "$logo"
  echo "Welcome to LabDAO! Documentation at https://github.com/labdao/"
}

setup() {
    installDocker
    testDockerInstall
    installNvidiaDrivers
    testNvidiainstall
    installNvidiaContainerToolKit
    testNvidiaContainerToolkitInstall
    installIPFS
    testIPFSInstall
    installGolang
    testGolangInstall
    installBacalhau
    testBacalhauInstall
    installConda
    testCondaInstall
    installJuypter
    testJuypterInstall
    printLogo
}

start() {
    case $PLEX_ENV in
        stage)
            echo "Starting for staging enviroment"
            ;;
        prod)
            echo "Starting for production enviroment"
            ;;
        *)
            echo "PLEX_ENV must be set to stage or prod"
            exit 77
            ;;
    esac
    runIPFS
    runBacalhau
    runJuypter
    echo "screen -ls"
}
