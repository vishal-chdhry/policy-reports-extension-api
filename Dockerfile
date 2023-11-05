ARG BUILDER_IMAGE="golang:1.21.3-alpine3.18"
FROM $BUILDER_IMAGE as builder

WORKDIR /
COPY ./server .

# Build Go binary
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o prext .

# Create image
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder prext prext
ENTRYPOINT ["/prext"]
