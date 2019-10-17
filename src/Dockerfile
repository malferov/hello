FROM golang AS build
WORKDIR /go/src
COPY . .
RUN go get github.com/gin-gonic/gin
RUN go get github.com/go-redis/redis
RUN go build -o hello

FROM centos
ARG port
EXPOSE $port
COPY --from=build /go/src/hello .
ENTRYPOINT ["./hello"]
CMD ["5000", "2fbeb67"]
