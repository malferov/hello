FROM golang AS build
WORKDIR /go/src
COPY . .
RUN go get github.com/gin-gonic/gin
RUN go build -o app

FROM centos
ARG port
EXPOSE $port
COPY --from=build /go/src/app .
ENTRYPOINT ["./app"]
CMD ["5000", "24adf75"]
