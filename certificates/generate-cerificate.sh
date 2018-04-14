#!/bin/bash
openssl genrsa -out proxy-ca.key 2048
openssl req -x509 -new -nodes -key proxy-ca.key -sha256 -days 1825 -out proxy-ca.pem
