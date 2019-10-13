#!/bin/bash
#
# please export HELLO_REPOSITORY_URL environment vaiable before build
#
BUILD_NUMBER=`git rev-parse --short HEAD`
TAG=$HELLO_REPOSITORY_URL:$BUILD_NUMBER
PORT=5000
sed -i '/CMD/d' Dockerfile
echo CMD [\"$PORT\", \"$BUILD_NUMBER\"] >> Dockerfile
$(aws ecr get-login --no-include-email --region $TF_VAR_region)
docker build -t $TAG --build-arg port=$PORT .
docker push $TAG
