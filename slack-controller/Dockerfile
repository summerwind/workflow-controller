FROM ubuntu:18.04 as base

RUN apt update && apt install -y ca-certificates

#######################################

FROM golang:1.12 as builder

ENV GO111MODULE=on

WORKDIR /workspace

COPY go.mod go.sum .

RUN go mod download

COPY . .

WORKDIR /workspace/slack-controller

RUN go vet ./...
RUN go test -v ./...
RUN CGO_ENABLED=0 go build .

#######################################

FROM summerwind/whitebox-controller:0.4.0

COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /workspace/slack-controller/slack-controller /bin/slack-controller
COPY --from=builder /workspace/slack-controller/config.yaml /config.yaml

ENTRYPOINT ["/bin/whitebox-controller"]
