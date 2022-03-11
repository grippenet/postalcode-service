# Args to make this more templateable
ARG SVC=grippenet/postalcode-service
ARG CMD_DIR=cmd/server

##########################
# Build stage
##########################

FROM golang:1.17-alpine as builder

RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
RUN mkdir -p /go/src/github.com/${SVC}
ENV GO111MODULE=on
ADD . /go/src/github.com/${SVC}
WORKDIR /go/src/github.com/${SVC}
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
WORKDIR /go/src/github.com/${SVC}/${CMD_DIR}
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o app .

##########################
# Image 
##########################
FROM scratch
COPY --from=builder /go/src/github.com/${SVC}/${CMD_DIR}/app /app/
COPY --from=builder /go/src/github.com/${SVC}/data /app/
WORKDIR /app
ENV GIN_MODE=release
ENV POSTAL_FILE=postal.json
EXPOSE 3251:3251
CMD ["./app"]
