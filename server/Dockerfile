ARG BUILDER_IMAGE="golang:1.21.3-alpine3.18"
FROM $BUILDER_IMAGE as builder

WORKDIR /
COPY . ./

# Build Go binary
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o policy-reports-aggregation .

# Create image
FROM gcr.io/distroless/static:nonroot
COPY --from=builder policy-reports-aggregation policy-reports-aggregation
ENTRYPOINT ["/policy-reports-aggregation"]
