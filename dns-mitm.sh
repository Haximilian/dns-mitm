#!/bin/bash
trap "networksetup -setdnsservers Wi-Fi 8.8.8.8" INT
networksetup -setdnsservers Wi-Fi 127.0.0.1
killall -HUP mDNSResponder
go run dns-server.go