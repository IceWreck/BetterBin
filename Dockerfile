# builder
FROM golang:1.19 AS builder

COPY ${PWD} /app
WORKDIR /app

# compile BetterBin
RUN CGO_ENABLED=1 go build -ldflags '-s -w -extldflags "-static"' -o /app/BetterBin .

# download and compile goose which is responsible for managing migrations
RUN go get -u github.com/pressly/goose/cmd/goose

# main container
FROM debian:stable-slim
LABEL MAINTAINER Anchit Bajaj (@IceWreck) <ab@abifog.com>

# add new user 'betterbin'
RUN adduser --home "/betterbin" --disabled-password betterbin \
    --gecos "betterbin,-,-,-"

# copy required files from builder
COPY --from=builder /app/static /home/betterbin/app/static
COPY --from=builder /app/templates /home/betterbin/app/templates
COPY --from=builder /app/db /home/betterbin/app/db
COPY --from=builder /app/BetterBin /home/betterbin/app/BetterBin
COPY --from=builder /app/files/container-entrypoint.sh /home/betterbin/app/container-entrypoint.sh
COPY --from=builder /go/bin/goose /home/betterbin/app/goose

RUN chown betterbin /home/betterbin/app
RUN chmod +x /home/betterbin/app/container-entrypoint.sh

USER betterbin

RUN mkdir -p /home/betterbin/app/drops
RUN mkdir -p /home/betterbin/app/data

WORKDIR /home/betterbin/app

EXPOSE 8963

CMD ["./container-entrypoint.sh"]