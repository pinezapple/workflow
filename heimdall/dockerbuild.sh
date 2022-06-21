DockerTag=ghcr.io/vfluxus/heimdall:stage-$(date +%s)
docker build -t $DockerTag -f ./Dockerfile ..
docker push $DockerTag
