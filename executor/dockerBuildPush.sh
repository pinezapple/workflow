#!/bin/sh
DockerTag=ghcr.io/vfluxus/executor:stage-$(date +%s)
docker build -t $DockerTag .
docker push $DockerTag