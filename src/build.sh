#!/bin/bash
BUILD_NUMBER=`git rev-parse --short HEAD`
PORT=5000
REG=docker.io
REPO=malferov
TAG=$REG/$REPO/app:$BUILD_NUMBER
sed -i '/CMD/d' Dockerfile
echo CMD [\"$PORT\", \"$BUILD_NUMBER\"] >> Dockerfile
docker build -t $TAG --build-arg port=$PORT .
docker push $TAG
