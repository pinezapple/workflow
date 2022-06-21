# BUILD
> Contains anything for build app. Etc dockerfile, docker-compose,....

## TABLE OF CONTENTS
- [BUILD](#build)
  - [TABLE OF CONTENTS](#table-of-contents)
  - [DOCKER BUILD:](#docker-build)
  - [DOCKER COMPOSE:](#docker-compose)

## DOCKER BUILD: 

- Docker build based on context:
  - -> need to change the context to parent directory (../)

  ```bash
  docker build -t heimdall -f ./Dockerfile ..
  ```
- Tag the docker images -> docker repo:
    ```bash
    docker tag heimdall <repo>
    ```
    - github repo: docker.pkg.github.com/internvinbdi/heimdall/heimdall
      - will tag :latest to the repo
      - Need to store the [package docs](./PackageDocs.md)

## DOCKER COMPOSE: