FROM golang:1.22-alpine3.19 as builder

RUN apk update && apk add git make bash

WORKDIR /app

COPY . ./

COPY .env /app/.env

ENV CGO_ENABLED=1

RUN uname -m

RUN apk add musl-dev gcc git make vips vips-dev vips-tools
# Do dep installs outside, due to private git modules
# RUN make dep

RUN make build

FROM alpine:latest

WORKDIR /app

RUN apk add curl vips vips-dev vips-tools \
    && vips --version

COPY --from=builder /app/main /app/
COPY --from=builder /app/public /app/public
COPY --from=builder /app/dist /app/dist
COPY --from=builder /app/pub /app/pub
COPY --from=builder /app/.env /app/.env

RUN mkdir /app/store

EXPOSE 4001

CMD [ "/app/main" ]
