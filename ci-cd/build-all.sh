#!/bin/bash

cd ../convert-ws
docker build -t cqrs/convert-ws:latest .
cd -

cd ../write-ws
docker build -t cqrs/write-ws:latest .
cd -

cd ../read-ws
docker build -t cqrs/read-ws:latest .
cd -
