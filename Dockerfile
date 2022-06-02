# This Dockerfile builds the lister sourcecode inside a golang docker
# container and sets the CMD to be executed when a container is launched
# from the built image.
FROM golang:1.18

WORKDIR /app/lister
COPY . .
RUN go build -v -o /usr/local/bin/lister ./...

CMD ["/usr/local/bin/lister"]