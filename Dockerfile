ARG GO_VERSION=1.14
FROM golang:${GO_VERSION}-alpine AS builder

RUN apk add --no-cache ca-certificates git

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /src
# and will therefore be cached for speeding up the next build
COPY ./go.mod ./go.sum ./
#RUN go mod download

# Import the code from the context.
COPY ./ ./

RUN CGO_ENABLED=0 go build -installsuffix 'static' -o /app .

FROM alpine AS final
#COPY --from=builder /user/group /user/passwd /etc/
# Import the Certificate-Authority certificates for enabling HTTPS.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Import the compiled executable from the first stage.
COPY --from=builder /app /app

EXPOSE 8080

ENTRYPOINT ["/app"]
