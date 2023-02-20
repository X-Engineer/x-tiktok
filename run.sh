#!/bin/bash
pkill -f './x-tiktok'
go build .
chmod +x ./x-tiktok
nohup ./x-tiktok > console.log 2>&1 &