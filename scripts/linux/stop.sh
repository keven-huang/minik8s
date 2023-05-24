#!/bin/bash

pids=$(ps -aux | grep kub | grep -v grep | awk '{print $2}')

for pid in $pids; do
  sudo kill "$pid"
done

docker ps -aq --filter "name=^my-replicaset|^test" | xargs docker stop
docker ps -aq --filter "name=^my-replicaset|^test" | xargs docker rm

