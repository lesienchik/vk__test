ARG GO_IMAGE=golang
ARG GO_IMAGE_VERSION=1.22-alpine3.20
ARG ALPINE_IMAGE_TAG=alpine:3.20

FROM ${GO_IMAGE}:${GO_IMAGE_VERSION} AS builder

LABEL stage=gobuilder

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /vktest/app

ADD go.mod .
ADD go.sum .
RUN go mod download

COPY . /vktest/app/

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o vktest ./cmd/vktest/vktest.go
FROM ${ALPINE_IMAGE_TAG}



WORKDIR /root/
COPY --from=builder /vktest/app/vktest .
COPY --from=builder /vktest/app/local_files/config.json ./local_files/config.json
COPY --from=builder /vktest/app/docs/ ./docs/

CMD ["./vktest"]
