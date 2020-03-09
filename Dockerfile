FROM golang:alpine
COPY go.mod *.go /source/
WORKDIR /source
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o docker-hosts

FROM scratch
COPY --from=0 /source/docker-hosts /bin/docker-hosts
COPY hosts /etc/hosts
CMD ["/bin/docker-hosts"]
