#!/bin/bash

pids=$(ps -a | grep kub | awk '{print $1}')

for pid in $pids; do
  kill "$pid"
done