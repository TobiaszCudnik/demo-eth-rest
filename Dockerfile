FROM golang:1.19.3-alpine
# deps
RUN apk add --no-cache curl git
RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin
# src
WORKDIR /app
COPY . .
# run
CMD ["task",  "start"]
