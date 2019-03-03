FROM golang:1.10.4-stretch AS builder

COPY .  /go/src/github.com/fairyhunter13/tax-calculator

WORKDIR /go/src/github.com/fairyhunter13/tax-calculator

RUN apt-get update; \
    apt-get install -y --no-install-recommends curl; \
    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh ; \
    dep ensure -v --vendor-only; \
    # Running the unit tests
    GOMAXPROCS=16 go test ./... -test.v -race -tags=unit; \
    cd cmd/taxcalculator; \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags="-w -s" -o /taxcalculator . ; \
    apt-get purge curl -y; \
    apt-get clean autoclean ;\
    apt-get autoremove --yes; \
    rm -rf /var/lib/{apt,dpkg,cache,log}/

FROM alpine

COPY --from=builder /taxcalculator .
COPY ./configs/config.ini /configs/config.ini
COPY ./scripts/applicationwrapper.sh .

RUN apk update ; \
    apk add --no-cache bash postgresql-client; \
    rm -rf /var/cache/apk/*; \ 
    chmod +x applicationwrapper.sh 

EXPOSE 9000

CMD ./applicationwrapper.sh


