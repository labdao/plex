#!/bin/bash
# executing these steps in order makes sure that each larger set is a superset of each smaller set
./tools/unidock/fetch_zinc22.sh tools/unidock/ZINC22-downloader-3D-complete-pdbqt.tgz.wget "H04" E2
./tools/unidock/fetch_zinc22.sh tools/unidock/ZINC22-downloader-3D-complete-pdbqt.tgz.wget "H04|H05|H06" E3
./tools/unidock/fetch_zinc22.sh tools/unidock/ZINC22-downloader-3D-complete-pdbqt.tgz.wget "H07|H08" E4
./tools/unidock/fetch_zinc22.sh tools/unidock/ZINC22-downloader-3D-complete-pdbqt.tgz.wget "H09|H10|H11" E5
./tools/unidock/fetch_zinc22.sh tools/unidock/ZINC22-downloader-3D-complete-pdbqt.tgz.wget "H12|H13" E6