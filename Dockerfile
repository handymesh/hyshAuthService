FROM golang:1.20.5-alpine as builder

# Install dep
RUN apk add --update ca-certificates git 
    # && go get -u github.com/golang/dep/cmd/dep

# Build project
WORKDIR /go/src/github.com/handymesh/go_authService
COPY . .
# RUN dep ensure
RUN go mod init github.com/handymesh/handy_authService && go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o go-auth cmd/go-auth/main.go

RUN go build -o myapp

FROM scratch

# RUN addgroup -S 997 && adduser -S 997 -G 997
# USER 997

WORKDIR /app/
COPY --from=builder /go/src/github.com/handymesh/go_authService/go-auth .
CMD ["./go-auth"]
