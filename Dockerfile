ARG appName=app

FROM golang:1.20.1-alpine3.17 AS build
ARG appName
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app
COPY /tmp tmp
COPY go.mod go.sum main.go ./
COPY /src src
COPY /docs docs
COPY /vendor vendor

RUN pwd
RUN ls docs

ENV CGO_ENABLED=0
RUN go build -o ${appName}

FROM scratch
ARG appName

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Copy timezone data, which avaiable after installing tzdata
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=Asia/Jakarta

COPY --from=build /app/${appName} /app
COPY --from=build /app/docs/index.html /index.html
COPY --from=build /app/tmp/.add /tmp/.add
COPY --from=build /app/src/config/location.json /location.json


ENTRYPOINT [ "/app" ]