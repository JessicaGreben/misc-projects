FROM golang:alpine AS build
ADD . /src
RUN cd /src && go build -o kademlia

FROM alpine
RUN apk --no-cache add curl
USER nobody
WORKDIR /app
COPY --from=build /src/kademlia /app/
COPY --from=build /src/startup.sh /app/
ENTRYPOINT ["./startup.sh"]
CMD ["./kademlia", "-s", "8080"]
