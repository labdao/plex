#TODO #20
wget https://dist.ipfs.tech/kubo/v0.18.0/kubo_v0.18.0_linux-amd64.tar.gz
tar -xvzf kubo_v0.18.0_linux-amd64.tar.gz
cd kubo
sudo bash install.sh
ipfs --version

# post installation
ipfs init
ipfs cat /ipfs/QmQPeNsJPyVWPFDVHb77w8G42Fvo15z4bG2X8D2GhfbSXc/readme


# port forwarding
# 127.0.0.1:5001/webui 
# name files with ipfs get QmWLy5XUKuSVJiYirLfw3xSxD42BDPZEVvCv7f2mGanUiD -o test.png

# NOTRUN screen -ls
screen -dm ipfs daemon -D
