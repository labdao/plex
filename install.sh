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
            if [ "$ARCH" = "amd64" ] || [ "$ARCH" = "x86_64" ] || [ "$ARCH" = "i386" ]
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
getTestData
getTools
displayLogo

echo "Installation complete. Welcome to LabDAO! Documentation at https://github.com/labdao/plex"
echo "You're ready to generate computational biology data! Run the following command to run Equibind on test data:"
echo './plex init -t tools/equibind.json -i '\''{"protein": "testdata/binding/abl/7n9g.pdb"], "small_molecule": ["testdata/binding/abl/ZINC000003986735.sdf"]}'\'' --scatteringMethod=dotProduct --autoRun=true'