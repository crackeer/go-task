# build image

```
export DOCKER_BUILDKIT=1
export HTTPS_PROXY=http://127.0.0.1:7897
export HTTP_PROXY=http://127.0.0.1:7897

docker build --platform linux/amd64 -t go-task .
```

# docker-compose.yml

```
services:
  app:
    container_name: go-task
    image: go-task:latest
    restart: always
    ports:
    - "7501:80"
```
