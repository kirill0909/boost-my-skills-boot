FROM golang:1.19-alpine AS build
RUN apk --no-cache add gcc g++ make git

# create temp dir for build app
RUN mkdir "app"
ADD ./app /app
WORKDIR /app

ARG ACCESS_AX_USER
ARG ACCESS_AX_TOKEN
ARG ACCESS_USER
ARG ACCESS_TOKEN
RUN go env -w GOPRIVATE=${CI_SERVER_HOST} GOSUMDB=off
RUN echo -e "machine gitlab.axarea.ru\nlogin ${ACCESS_AX_USER}\npassword ${ACCESS_AX_TOKEN}" > ~/.netrc
RUN echo -e "machine gitlab.com\nlogin ${ACCESS_USER}\npassword ${ACCESS_TOKEN}" >> ~/.netrc

# download requirements && compile
RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build -o /bin/main ./cmd/api/main.go

FROM alpine:3.9
RUN apk --no-cache add ca-certificates
ADD https://github.com/golang/go/raw/master/lib/time/zoneinfo.zip /zoneinfo.zip
ENV ZONEINFO /zoneinfo.zip
COPY --from=build /bin/main /bin/main
ENTRYPOINT ["/bin/main"]
