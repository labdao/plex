Function setOSandArch {
    if($env:OS -like "*windows*") {
        $global:OS = "windows"
    }
    else {
        $global:OS = "linux"
    }

    if($env:PROCESSOR_ARCHITECTURE -eq "AMD64") {
        $global:ARCH = "amd64"
    }
    elseif($env:PROCESSOR_ARCHITECTURE -eq "x86") {
        $global:ARCH = "386"
    }
    elseif($env:PROCESSOR_ARCHITECTURE -eq "ARM64") {
        $global:ARCH = "arm64"
    }
    else {
        $global:ARCH = "unknown"
    }
}


Function makeParentFolder() {
    New-Item -ItemType Directory -Force -Path plex
    Set-Location -Path plex
}


Function makeConfigFolder() {
    New-Item -ItemType Directory -Force -Path config
}


Function downloadPlex {
    setOSandArch

    $LATEST_RELEASE=$(Invoke-RestMethod -Uri "https://api.github.com/repos/labdao/plex/releases/latest" | Select-Object -ExpandProperty tag_name)
    $RELEASE_WITHOUT_V=$LATEST_RELEASE.TrimStart('v')

    Write-Host "OS: $OS"
    Write-Host "ARCH: $ARCH"
    Write-Host "Latest Release: $LATEST_RELEASE"

    if($OS -eq "windows") {
        if($ARCH -eq "amd64" -or $ARCH -eq "x86_64") {
            $url = "https://github.com/labdao/plex/releases/download/$LATEST_RELEASE/plex_${RELEASE_WITHOUT_V}_windows_amd64.tar.gz"
            Invoke-WebRequest -Uri $url -OutFile "plex.tar.gz"
            tar -xzvf "plex.tar.gz"
            Remove-Item "plex.tar.gz"
        }
        else {
            Write-Host "Cannot install Plex. Unsupported architecture for Windows: $ARCH"
        }
    }
}


Function getTools() {
    Write-Host "Downloading tools..."
    New-Item -ItemType Directory -Force -Path tools

    $toolsUrl = "https://raw.githubusercontent.com/labdao/plex/main/tools/"
    $toolFiles = @("equibind.json", "diffdock.json", "bam2fastq.json", "colabfold-large.json", "colabfold-mini.json", "colabfold-standard.json", "oddt.json")

    foreach($tool in $toolFiles) {
        $url = $toolsUrl + $tool
        Invoke-WebRequest -Uri $url -OutFile ("tools\" + $tool)
    }
}


Function getTestData() {
    Write-Host "Downloading test data..."

    $testDataUrl = "https://raw.githubusercontent.com/labdao/plex/main/testdata/"

    $testDataFiles = @{
        "binding/abl/" = @("7n9g.pdb", "ZINC000003986735.sdf", "ZINC000019632618.sdf");
        "folding/" = @("test.fasta");
        "design/" = @("insulin_target.pdb");
    "scoring/pdbbind_processed_size1_equibind/" = @("6d08_ligand.sdf", "6d08_protein_processed.pdb", "6d08_protein_processed_6d08_ligand_docked.sdf")
    }

    foreach($testDataPath in $testDataFiles.Keys) {
        New-Item -ItemType Directory -Force -Path ("testdata\" + $testDataPath.Replace('/', '\'))

        foreach($fileName in $testDataFiles[$testDataPath]) {
            $url = $testDataUrl + $testDataPath + $fileName
            Invoke-WebRequest -Uri $url -OutFile ("testdata\" + $testDataPath.Replace('/', '\') + $fileName)
        }
    }
}

makeParentFolder
makeConfigFolder
downloadPlex
getTestData
getTools

Write-Host "Installation complete. Welcome to LabDAO! Documentation at https://github.com/labdao/plex"
Write-Host "To get started, please run the following steps:"
Write-Host "1. To request access to the Jupyter Hub please visit: https://try.labdao.xyz"
Write-Host "2. Run the following command to run Equibind on test data:"
Write-Host "./plex init -t tools/equibind.json -i '{"protein": "testdata/binding/abl/7n9g.pdb"], "small_molecule": ["testdata/binding/abl/ZINC000003986735.sdf"]}' --scatteringMethod=dotProduct --autoRun=true"