#!/bin/bash
basepath=$(dirname "$0")

cd "$basepath/../certificates"
./generate-certificate.sh
