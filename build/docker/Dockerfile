##########################
# Build stage
##########################

ARG DIR=/go/src/github.com/grippenet/postalcode-service

FROM golang:1.17-alpine as builder
ARG DIR # Populate global ARG for this build
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
RUN mkdir -p ${DIR}
ENV GO111MODULE=on
ADD . ${DIR}
WORKDIR ${DIR}
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
WORKDIR $DIR/cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o app .

##########################
# Image 
##########################
FROM scratch
ARG DIR
COPY --from=builder $DIR/cmd/server/app /app/
COPY --from=builder $DIR/data /app/
WORKDIR /app
ENV GIN_MODE=release
ENV POSTAL_FILE=postal.json
ENV PORT=3251
EXPOSE 3251:3251
CMD ["./app"]
