FROM golang:1.20-alpine AS build
WORKDIR /src
COPY . .
RUN apk add --no-cache git \
    && go build -o /bin/mockauth ./

FROM alpine:3.18
COPY --from=build /bin/mockauth /bin/mockauth
EXPOSE 3000
ENTRYPOINT ["/bin/mockauth"]
