FROM golang:1.22-alpine3.19 as builder

RUN apk update && apk add git make bash musl-dev gcc libwebp-dev

WORKDIR /app

COPY . ./

# Do dep installs outside, due to private git modules
# RUN make dep

RUN make build
RUN make build-ssh

FROM alpine:latest

WORKDIR /app

RUN apk update && apk add libwebp

COPY --from=builder /app/main /app/
COPY --from=builder /app/main-sh /app/
COPY --from=builder /app/wrapper.sh /app/

COPY --from=builder /app/public /app/public
COPY --from=builder /app/dist /app/dist
COPY --from=builder /app/pub /app/pub

RUN ls -la /app/

EXPOSE 4001

EXPOSE 23234

CMD ["sh","/app/wrapper.sh"]

