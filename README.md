# build image

```
export DOCKER_BUILDKIT=1
export HTTPS_PROXY=http://127.0.0.1:7897
export HTTP_PROXY=http://127.0.0.1:7897

docker build --platform linux/amd64 -t go-task .
```

# docker-compose.yml

```
version: '3.8'

services:
  app:
    image: go-task:latest
    ports:
      - 7600:80
    environment:
      - PORT=80 
```
