FROM golang:alpine AS build
WORKDIR /go/src/github.com/jessicagreben/kademlia/
ADD . .
RUN go build -o kademlia

FROM alpine
USER nobody
WORKDIR /app
COPY --from=build /go/src/github.com/jessicagreben/kademlia/kademlia /app/
COPY --from=build /go/src/github.com/jessicagreben/kademlia/startup.sh /app/
ENTRYPOINT ["./startup.sh"]
CMD ["./kademlia", "-s", "8080"]
