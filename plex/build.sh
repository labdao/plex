buildMacOsBinaries() {
    echo "Building Mac OS binaries..."
    
    if test releases/macos-amd64/plex ; then
        rm releases/macos-amd64/plex
    fi

    if test releases/macos-arm64/plex ; then
        rm releases/macos-arm64/plex
    fi

    GOOS=darwin GOARCH=amd64 go build -o releases/macos-amd64
    GOOS=darwin GOARCH=arm64 go build -o releases/macos-arm64
}

buildLinuxBinaries() {
    echo "Building Linux binaries..."

    if test releases/linux-amd64/plex ; then
        rm releases/linux-amd64/plex
    fi

    GOOS=linux GOARCH=amd64 go build -o releases/linux-amd64
}

buildWindowsBinaries() {
    echo "Building Windows binaries..."

    if test releases/windows-amd64/plex.exe ; then
        rm releases/windows-amd64/plex.exe
    fi

    GOOS=windows GOARCH=amd64 go build -o releases/windows-amd64
}

buildMacOsBinaries
buildLinuxBinaries
buildWindowsBinaries

echo "All builds complete."