FROM golang:1.12 as builder

ENV GO111MODULE=on

WORKDIR /workspace

COPY go.mod go.sum .

RUN go mod download

COPY . .

WORKDIR /workspace/feed-controller

RUN go vet ./...
RUN go test -v ./...
RUN CGO_ENABLED=0 go build .

#######################################

FROM summerwind/whitebox-controller:0.4.0

COPY --from=builder /workspace/feed-controller/feed-controller /bin/feed-controller
COPY --from=builder /workspace/feed-controller/config.yaml /config.yaml

ENTRYPOINT ["/bin/whitebox-controller"]
