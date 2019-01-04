#!/bin/bash

ABSPATH=$(cd "$(dirname "$0")"; pwd)
echo $ABSPATH

cd $ABSPATH/../convert-ws
docker build -t cqrs/convert-ws:latest .
cd -

cd $ABSPATH/../write-ws
docker build -t cqrs/write-ws:latest .
cd -

cd $ABSPATH/../read-ws
docker build -t cqrs/read-ws:latest .
cd -

cd $ABSPATH/../kong-init
docker build -t cqrs/kong-init:latest .
cd -