#!/bin/bash

pids=$(ps -aux | grep kub | grep -v grep | awk '{print $2}')

for pid in $pids; do
  kill "$pid"
done