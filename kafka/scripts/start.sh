#!/bin/bash
/etc/confluent/docker/run &
./scripts/create-topics.sh
wait
