#!/bin/bash

export MAX_NODES=5
export SAMPLE_SIZE=10
export QUORUM_SIZE=14
export DECISION_THRESHOLD=20

for ((i = 0; i < MAX_NODES; i++)); do
  export PORT=$((9000 + $i))
  export NODE_ID=$((9000 + $i))
  go run ./cmd/srv/. &
done

wait
sleep 10
