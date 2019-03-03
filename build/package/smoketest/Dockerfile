FROM golang:1.10.4-stretch

COPY .  /go/src/github.com/fairyhunter13/tax-calculator
COPY ./configs/config.ini /configs/config.ini
WORKDIR /go/src/github.com/fairyhunter13/tax-calculator/test/smoke
COPY ./scripts/smoketestwrapper.sh .

RUN apt-get update; \
    apt-get install -y --no-install-recommends curl; \
    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh ; \
    dep ensure -v --vendor-only; \
    apt-get clean autoclean ;\
    apt-get autoremove --yes; \
    rm -rf /var/lib/{apt,dpkg,cache,log}/; \
    chmod +x smoketestwrapper.sh

CMD ./smoketestwrapper.sh