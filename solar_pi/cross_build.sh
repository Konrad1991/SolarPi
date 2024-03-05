#!/bin/bash
cross build --target=armv7-unknown-linux-gnueabihf
# cross build --target=arm-unknown-linux-gnueabihf

scp /home/konrad/Documents/GitHub/SolarPi/solar_pi/target/armv7-unknown-linux-gnueabihf/debug/solar_pi konrad@192.168.1.192:~

