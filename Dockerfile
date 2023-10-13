# Build the project
FROM golang:1.19 as builder

WORKDIR /go/src/git.mytaxi.uz/mytaxi/billing-data-processor
ADD . .

RUN make build
# RUN make test

# Create production image for application with needed files
FROM golang:1.19-alpine3.17

EXPOSE 8000

RUN apk add --no-cache ca-certificates

COPY --from=builder /go/src/git.mytaxi.uz/mytaxi/billing-data-processor .

CMD ["./main"]
