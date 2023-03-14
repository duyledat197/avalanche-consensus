#!/bin/bash

export PORT=9090
export MAX_NODES=50
export SAMPLE_SIZE=10
export QUORUM_SIZE=14
export DECISION_THRESHOLD=20

go run ./cmd/srv/.