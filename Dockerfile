FROM golang:1.14-alpine AS build-env
RUN apk --no-cache add build-base git mercurial gcc make
ADD ./ /src
RUN cd /src && go build cmd/service/main.go

FROM alpine
WORKDIR /
RUN mkdir config
COPY --from=build-env /src/main .
COPY --from=build-env /src/config /config
RUN chmod +x main
EXPOSE 8081
ENTRYPOINT ["/main"]