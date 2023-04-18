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

        LATEST_RELEASE=$(curl -s https://api.github.com/repos/labdao/plex/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
        RELEASE_WITHOUT_V=${LATEST_RELEASE#v}

        if [ "$OS" = "darwin" ]
        then
            if [ "$ARCH" = "amd64" ] || [ "$ARCH" = "x86_64" ]
            then
                curl -sSL https://github.com/labdao/plex/releases/download/$LATEST_RELEASE/plex_${RELEASE_WITHOUT_V}_darwin_amd64.tar.gz | tar xvz
            elif [ "$ARCH" = "arm64" ]
            then
                curl -sSL https://github.com/labdao/plex/releases/download/$LATEST_RELEASE/plex_${RELEASE_WITHOUT_V}_darwin_arm64.tar.gz | tar xvz
            else
                echo "Cannot install Plex. Unsupported architecture for Darwin OS: $ARCH"
            fi
        elif [ "$OS" = "linux" ]
        then
            if [ "$ARCH" = "amd64" ] || [ "$ARCH" = "x86_64" ]
            then
                curl -sSL https://github.com/labdao/plex/releases/download/$LATEST_RELEASE/plex_${RELEASE_WITHOUT_V}_linux_amd64.tar.gz | tar xvz
            else
                echo "Cannot install Plex. Unsupported architecture for Linux: $ARCH"
            fi
        elif [ "$OS" = "windows" ]
        then
            if [ "$ARCH" = "amd64" ] || [ "$ARCH" = "x86_64" ]
            then
                curl -sSL https://github.com/labdao/plex/releases/download/$LATEST_RELEASE/plex_${RELEASE_WITHOUT_V}_windows_amd64.tar.gz
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

getTools() {
    echo "Downloading tools..."
    mkdir -p tools
    curl -sL -O https://github.com/labdao/plex/archive/refs/heads/main.zip
    unzip -qq main.zip && mv plex-main/tools/*.json tools/ && rm -rf plex-main
    rm -f main.zip
}

getTestData() {
    echo "Downloading test data..."
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
getTools
displayLogo

echo "Installation complete. Welcome to LabDAO! Documentation at https://github.com/labdao/plex"
echo "To get started, please run the following steps:"
echo "1. Please change the permissions of plex on your system:"
echo "chmod +x ./plex"
echo "2. Please add your access token to use plex. To request an access token, visit https://whe68a12b61.typeform.com/to/PpbO2HYf"
echo "export PLEX_ACCESS_TOKEN=<your access token>"
echo "3. [Linux only] If you recieve a warning about download speeds on Linux you can optionally run:"
echo "sudo sysctl -w net.core.rmem_max=2500000"

echo "After these steps, you're ready to generate computational biology data! Run the following command to run Equibind on test data:"
echo "./plex -app equibind -input-dir ./testdata/binding/pdbbind_processed_size1/"
