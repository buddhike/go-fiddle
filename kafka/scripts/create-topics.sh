#!/bin/bash

kafka-topics --create --if-not-exists --topic "request" --zookeeper zookeeper:32181 --partitions 1 --replication-factor 1
kafka-topics --create --if-not-exists --topic "response" --zookeeper zookeeper:32181 --partitions 1 --replication-factor 1

kafka-topics --create --if-not-exists --topic "requestsummary" --zookeeper zookeeper:32181 --partitions 1 --replication-factor 1
kafka-topics --create --if-not-exists --topic "responsesummary" --zookeeper zookeeper:32181 --partitions 1 --replication-factor 1
