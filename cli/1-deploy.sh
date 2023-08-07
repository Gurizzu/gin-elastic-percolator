#!/usr/bin/env bash

./cli/3-parse-env.sh

echo "Do you want to commit? (y/n)"
read isCommit

if [ $isCommit == "y" ]
then
    echo "Please input your commit:"
    read commit
    echo $commit
    git add .
    git commit -m "$commit"
    git push
fi

source .env

SERVER_SERVICE=$SERVER_SERVICE
APPNAME=$SVC_PERCOLATOR_NAME
TAG=latest
SERVER_CONFIG=/home/ws/$PROJECT_NAME/$APPNAME/.env
REGISTRY=$REGISTRY/$PROJECT_NAME
PORT=$SVC_PERCOLATOR_PORT
DOCKER_RUN="\
docker pull $REGISTRY/$APPNAME:$TAG;\
docker rm -f $APPNAME; \
docker run -d \
-m 512m \
--log-opt max-size=50m \
-p $PORT:$PORT \
--restart always \
--env-file $SERVER_CONFIG \
--name $APPNAME \
$REGISTRY/$APPNAME:$TAG"

echo "Deploying $REGISTRY/$APPNAME:$TAG"
echo "Use Port: $PORT"

swag i -parseVendor

docker build -t $APPNAME . --build-arg appName=$APPNAME
docker tag $APPNAME $REGISTRY/$APPNAME:$TAG
docker push $REGISTRY/$APPNAME:$TAG

notify-send "WEES" "Remote docker run executing"
echo $DOCKER_RUN

ssh root@$SERVER_SERVICE $DOCKER_RUN

echo "API DOCS:
SWAGGER: http://$SERVER_SERVICE:$PORT/swagger/index.html
RAPIDOC: http://$SERVER_SERVICE:$PORT"
