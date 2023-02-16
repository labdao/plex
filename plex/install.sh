#!/bin/bash

setOSandArch() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(arch)
}

installGo() {
    if ! command -v go &> /dev/null
    then
        echo "Downloading and installing Go..."
        setOSandArch
        
        if [ "$OS" = "darwin"]
        then
            if [ "$ARCH" = "amd64" ] || [ "$ARCH" = "x86_64" ]
            then
                curl -O https://go.dev/dl/go1.19.6.darwin-amd64.pkg
                sudo installer -pkg go1.19.6.darwin-amd64.pkg -target /
                export PATH=$PATH:/usr/local/go/bin                
            elif [ "$ARCH" = "arm64" ]
            then
                curl -O https://go.dev/dl/go1.19.6.darwin-arm64.pkg
                sudo installer -pkg go1.19.6.darwin-arm64.pkg -target /
                export PATH=$PATH:/usr/local/go/bin
            else
                echo "Cannot install Go. Unsupported architecture for Darwin OS: $ARCH"
            fi
        elif [ "$OS" = "linux" ]
        then
            if [ "$ARCH" = "amd64" ] || [ "$ARCH" = "x86_64" ]
            then
                wget https://go.dev/dl/go1.19.6.linux-amd64.tar.gz
                sudo tar -C /usr/local -xvzf go1.19.6.linux-amd64.tar.gz
                rm go1.19.6.linux-amd64.tar.gz
                export PATH=$PATH:/usr/local/go/bin
            else
                echo "Cannot install Go. Unsupported architecture for Linux: $ARCH"
            fi
        elif [ "$OS" = "windows" ]
        then
            if [ "$ARCH" = "amd64" ] || [ "$ARCH" = "x86_64" ]
            then
                curl -O https://go.dev/dl/go1.19.6.windows-amd64.msi
                msiexec /i go1.19.6.windows-amd64.msi /quiet /qn
            else
                echo "Cannot install Go. Unsupported architecture for Windows: $ARCH"
            fi
        fi
    fi
}

installBacalhau() {
    echo "Installing Bacalhau..."
    curl -sL https://get.bacalhau.org/install.sh | bash
}

setW3SToken() {
    read -p "Enter your web3.storage API token: " WEB3STORAGE_TOKEN
    if [ -z "$WEB3STORAGE_TOKEN" ]
    then
        echo "web3.storage API token cannot be empty."
        setW3SToken
    else
        export WEB3STORAGE_TOKEN
        echo "web3.storage API token set successfully."
    fi
}

getPlex() {
    git clone https://github.com/labdao/ganglia.git
}

buildPlexBinary() {
    echo "Building Plex binary..."
    getPlex
    cd ganglia/plex
    go build
    echo "Plex binary built successfully."
}

displayLogo() {
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
                                        @
    "
    echo "$logo"
}

installGo
installBacalhau
setW3SToken
buildPlexBinary
displayLogo

echo "Installation complete. Welcome to LabDAO! Documentation at https://github.com/labdao/ganglia"
echo "To start using Plex, run the following command:"
echo "./plex -app equibind -gpu false -input-dir ./testdata/pdbbind_processed_size1"
