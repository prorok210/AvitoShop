FROM golang:1.22

WORKDIR ${GOPATH}/AvitoShop/
COPY . ${GOPATH}/AvitoShop/

RUN go build -o /build ./cmd/avito_shop_service \
    && go clean -cache -modcache

EXPOSE 8080

CMD ["/build"]