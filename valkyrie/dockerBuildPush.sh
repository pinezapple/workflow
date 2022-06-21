#!/bin/sh
DockerTag=ghcr.io/vfluxus/valkyrie:stage-$(date +%s)
docker build -t $DockerTag . &&
docker push $DockerTag
