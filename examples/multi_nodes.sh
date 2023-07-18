#!/bin/bash

echo "Starting multiple cupcake_cache nodes"

go run multi_nodes.go -port 8080 & 
go run multi_nodes.go -port 8081 &
go run multi_nodes.go -port 8082 -proxy 1 &

sleep 2

echo "Sending request to cupcake_cache"
echo "request 1 key:Test, val:" "$(curl -s http://localhost:8888/api?key=Test)" &
echo "request 2 key:Test, val:" "$(curl -s http://localhost:8888/api?key=Test)" &
echo "request 3 key:Test, val:" "$(curl -s http://localhost:8888/api?key=Test)" &

wait