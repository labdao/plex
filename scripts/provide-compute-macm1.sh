#!/bin/bash
# This script is tested to run on a fresh Ubuntu 20.04 LTS Nvidia GPU compute instance
# Specifically tested on an AWS P3, G4, and G5 instance

# this exception allows exit 77 to exit the whole script within subshells
# https://unix.stackexchange.com/questions/48533/exit-shell-script-from-a-subshell
set -E
trap '[ "$?" -ne 77 ] || exit 77' ERR

testDockerInstall() {
    echo "Testing Docker Install"
    if docker run hello-world ; then
        echo "Docker succesfully installed"
    else
        echo "Docker not installed - please install Docker on your machine first"
        exit 77
    fi
EONG
}

testGolangInstall() {
    echo "Testing GoLang Install"
    if go version ; then
        echo "GoLang succesfully installed "
    else
        echo "GoLang not installed - please install GoLang on your machine first"
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
    sudo sysctl -w net.core.rmem_max=2500000
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
    export owner=$owner // so that we have the variable on restarts
    echo "export owner=$owner" >> ~/.bashrc # for future shells
    screen -dmS bacalhau bacalhau serve --node-type compute,requester --ipfs-connect $IPFS_CONNECT --limit-total-gpu 1 --limit-job-memory 12gb --job-selection-accept-networked --job-selection-data-locality anywhere --labels owner=$owner --private-internal-ipfs=false
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
    testDockerInstall
    installIPFS
    testIPFSInstall
    testGolangInstall
    installBacalhau
    testBacalhauInstall
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
    screen -ls
}
