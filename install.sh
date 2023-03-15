#!/bin/bash

setOSandArch() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(arch)
}

makeParentFolder() {
    mkdir plex
    cd plex
}

makeConfigFolder() {
    mkdir config
}

downloadPlex() {
    if [[ ! -x plex ]]; then
        echo "Downloading Plex..."
        setOSandArch
        
        if [ "$OS" = "darwin" ]
        then
            if [ "$ARCH" = "amd64" ] || [ "$ARCH" = "x86_64" ]
            then
                curl -sSL https://github.com/labdao/plex/releases/download/v0.4.0/plex_0.4.0_darwin_amd64.tar.gz | tar xvz
            elif [ "$ARCH" = "arm64" ]
            then
                curl -sSL https://github.com/labdao/plex/releases/download/v0.4.0/plex_0.4.0_darwin_arm64.tar.gz | tar xvz
            else
                echo "Cannot install Plex. Unsupported architecture for Darwin OS: $ARCH"
            fi
        elif [ "$OS" = "linux" ]
        then
            if [ "$ARCH" = "amd64" ] || [ "$ARCH" = "x86_64" ]
            then
                curl -sSL https://github.com/labdao/plex/releases/download/v0.4.0/plex_0.4.0_linux_amd64.tar.gz | tar xvz
            else
                echo "Cannot install Plex. Unsupported architecture for Linux: $ARCH"
            fi
        elif [ "$OS" = "windows" ]
        then
            if [ "$ARCH" = "amd64" ] || [ "$ARCH" = "x86_64" ]
            then
                curl -sSL https://github.com/labdao/plex/releases/download/v0.4.0/plex_0.4.0_windows_amd64.tar.gz
            else
                echo "Cannot install Plex. Unsupported architecture for Windows: $ARCH"
            fi
        fi
    fi
}

getAppJsonl() {
    cd config
    curl -sSL -O https://raw.githubusercontent.com/labdao/plex/main/config/app.jsonl
    cd ..
}

getInstructionsTemplateJsonl() {
    cd config
    curl -sSL -O https://raw.githubusercontent.com/labdao/plex/main/config/instruction_template.jsonl
    cd ..
}

getTestData() {
    mkdir -p testdata
    curl -sL -O https://github.com/labdao/plex/archive/refs/heads/main.zip
    unzip -qq main.zip && mv plex-main/testdata/* testdata/ && rm -rf plex-main
    rm -f main.zip
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

    printf "\e[?25l"

    for (( i=0; i<${#logo}; i++)); do
        char="${logo:$i:1}"
        if [[ "$char" != " " ]]; then
            echo -n "$char"
        else
            echo -n "$char"
            sleep 0.0005
        fi
    done
    printf "\e[?25h"
    echo
}

makeParentFolder
makeConfigFolder
downloadPlex
getAppJsonl
getInstructionsTemplateJsonl
getTestData
displayLogo

echo "Installation complete. Welcome to LabDAO! Documentation at https://github.com/labdao/plex"
echo "To get started, please run the following 3 steps:"
echo "1. Please change the permissions of plex on your system:"
echo "chmod +x ./plex"
echo "2. Please run the following command to download large bacalhau results:"
echo "sudo sysctl -w net.core.rmem_max=2500000"
echo "3. To start using Plex, run the following command to run Equibind on test data:"
echo "./plex -app equibind -input-dir ./testdata/binding/pdbbind_processed_size1/"
