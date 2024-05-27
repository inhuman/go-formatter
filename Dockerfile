FROM golang:1.22-alpine AS builder

ENV CGO_ENABLED 0

ENV GCI_VERSION=v0.8.2
ENV GOFUMPT_VERSION=v0.4.0

WORKDIR /app

COPY . .

RUN go install github.com/daixiang0/gci@$GCI_VERSION
RUN go install mvdan.cc/gofumpt@$GOFUMPT_VERSION
RUN go install golang.org/x/tools/cmd/goimports@latest

RUN go build -o /format ./formatter.go

FROM scratch

COPY --from=builder /format /format

COPY --from=builder /go/bin/gci /bin/gci

COPY --from=builder /go/bin/gofumpt /bin/gofumpt

COPY --from=builder /go/bin/goimports /bin/goimports

COPY --from=builder /usr/local/go/bin/gofmt /bin/gofmt

ENTRYPOINT ["/format"]