#!/bin/bash

setOSandArch() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(arch)
}

makeParentFolder() {
    mkdir plex
    cd plex
}

downloadPlex() {
    if [[ ! -x plex ]]; then
        echo "Downloading Plex..."
        setOSandArch
        
        if [ "$OS" = "darwin" ]
        then
            if [ "$ARCH" = "amd64" ] || [ "$ARCH" = "x86_64" ]
            then
                curl -O https://raw.githubusercontent.com/labdao/ganglia/main/plex/releases/macos-amd64/plex
            elif [ "$ARCH" = "arm64" ]
            then
                curl -O https://raw.githubusercontent.com/labdao/ganglia/main/plex/releases/macos-arm64/plex
            else
                echo "Cannot install Go. Unsupported architecture for Darwin OS: $ARCH"
            fi
        elif [ "$OS" = "linux" ]
        then
            if [ "$ARCH" = "amd64" ] || [ "$ARCH" = "x86_64" ]
            then
                curl -O https://raw.githubusercontent.com/labdao/ganglia/main/plex/releases/linux-amd64/plex
            else
                echo "Cannot install Go. Unsupported architecture for Linux: $ARCH"
            fi
        elif [ "$OS" = "windows" ]
        then
            if [ "$ARCH" = "amd64" ] || [ "$ARCH" = "x86_64" ]
            then
                curl -O https://raw.githubusercontent.com/labdao/ganglia/main/plex/releases/windows-amd64/plex.exe
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

# setW3SToken() {
#     read -p "Enter your web3.storage API token: " WEB3STORAGE_TOKEN
#     if [ -z "$WEB3STORAGE_TOKEN" ]
#     then
#         echo "web3.storage API token cannot be empty."
#         setW3SToken
#     else
#         export WEB3STORAGE_TOKEN
#         echo "web3.storage API token set successfully."
#     fi
# }

getAppJsonl() {
    curl -sL -O https://raw.githubusercontent.com/labdao/ganglia/main/plex/app.jsonl
}

getInstructionsTemplateJsonl() {
    curl -sL -O https://raw.githubusercontent.com/labdao/ganglia/main/plex/instruction_template.jsonl
}

getTestData() {
    mkdir testdata
    cd testdata
    curl -sL -O https://raw.githubusercontent.com/labdao/ganglia/main/plex/testdata/pdbbind_processed_size1/6d08/6d08_protein_processed.pdb
    curl -sL -O https://raw.githubusercontent.com/labdao/ganglia/main/plex/testdata/pdbbind_processed_size1/6d08/6d08_ligand.sdf
    cd ..
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

makeParentFolder
downloadPlex
installBacalhau
# setW3SToken
getAppJsonl
getInstructionsTemplateJsonl
getTestData
displayLogo

echo "Installation complete. Welcome to LabDAO! Documentation at https://github.com/labdao/ganglia"
echo "Please run the following command to set your web3.storage token:"
echo "export WEB3STORAGE_TOKEN=<your API token>"
echo "After you have set your API token, to start using Plex run the following command:"
echo "./plex -app equibind -gpu false -input-dir ./testdata"
