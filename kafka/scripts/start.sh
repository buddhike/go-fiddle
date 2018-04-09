#!/bin/bash
LOGS=/var/log/kafkalogs
touch $LOGS

/etc/confluent/docker/run >> $LOGS &

sleep 2
./scripts/create-topics.sh

tail -f $LOGS
