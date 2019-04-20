FROM ubuntu:18.04 as base

RUN apt update && apt install -y ca-certificates

#######################################

FROM golang:1.12 as builder

ENV GO111MODULE=on

WORKDIR /workspace

COPY go.mod go.sum .

RUN go mod download

COPY . .

RUN go vet ./...
RUN go test -v ./...
RUN CGO_ENABLED=0 go build -o ./bin/feed-controller ./feed-controller
RUN CGO_ENABLED=0 go build -o ./bin/slack-controller ./slack-controller
RUN CGO_ENABLED=0 go build -o ./bin/github-controller ./github-controller

#######################################

FROM summerwind/whitebox-controller:0.4.0 as feed-controller

COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /workspace/bin/feed-controller /bin/feed-controller
COPY --from=builder /workspace/feed-controller/config.yaml /config.yaml

ENTRYPOINT ["/bin/whitebox-controller"]

#######################################

FROM summerwind/whitebox-controller:0.4.0 as slack-controller

COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /workspace/bin/slack-controller /bin/slack-controller
COPY --from=builder /workspace/slack-controller/config.yaml /config.yaml

ENTRYPOINT ["/bin/whitebox-controller"]

#######################################

FROM summerwind/whitebox-controller:0.4.0 as github-controller

COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /workspace/bin/github-controller /bin/github-controller
COPY --from=builder /workspace/github-controller/config.yaml /config.yaml

ENTRYPOINT ["/bin/whitebox-controller"]
