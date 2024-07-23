#!/bin/bash

# get ip via: 
# nmap -sn 192.168.1.0/24
IP=192.168.1.222

#export CGO_ENABLED=1
export GOOS=linux
export GOARCH=arm64  

#go build -o SolarPi_arm64 ./cmd/SolarPi 

#scp ./SolarPi_arm64 konrad@$IP:~

#ssh konrad@$IP "./SolarPi_arm64" 

# send code
scp -r ./cmd konrad@192.168.1.222:~ 