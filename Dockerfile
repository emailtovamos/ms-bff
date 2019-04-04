# 1st stage: build Go binary

FROM golang:1.10
WORKDIR /go/src/github.com/emailtovamos/ms-bff/

# Copy only Go package directories each separately
COPY vendor ./vendor/
COPY ./cli ./cli/
COPY bff ./bff/


RUN CGO_ENABLED=0        \
    GOOS=linux           \
    go install           \
      -a                 \
      -installsuffix cgo \
      ./cli


# 2nd stage: embed Go binary in small Linux distro (== Alpine)

FROM alpine:latest
WORKDIR /app/

# Copy the binary from the first build stage
COPY --from=0 /go/bin/cli ./binary

CMD ./binary 