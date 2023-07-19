#!/bin/bash
usage(){
    echo "Usage: sh multi_nodes.sh [grpc|http]"
    exit 0
}
TARGET=""
case $1 in
    "grpc" )
        TARGET="multi_grpc_nodes.go" ;;
    "http" )
        TARGET="multi_http_nodes.go" ;;
    * )
        usage;;    
esac

trap "rm server;kill 0" EXIT

echo "Starting multiple cupcake_cache nodes"

go build -o server $TARGET 
./server -port 8080 & 
./server -port 8081 &
./server -port 8082 -proxy 1 &

sleep 3

echo "Sending request to cupcake_cache"
echo "request 1 key:Test, val:" "$(curl -s http://localhost:8888/api?key=Test)" &
echo "request 2 key:Test, val:" "$(curl -s http://localhost:8888/api?key=Test)" &
echo "request 3 key:Test, val:" "$(curl -s http://localhost:8888/api?key=Test)" &

wait