FROM golang:1.19.4-alpine3.17 as builder

RUN apk update && apk add git make bash

WORKDIR /app

COPY . ./

# Do dep installs outside, due to private git modules
# RUN make dep

RUN make build

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main /app/
COPY --from=builder /app/public /app/public
COPY --from=builder /app/dist /app/dist

EXPOSE 4001

CMD [ "/app/main" ]
