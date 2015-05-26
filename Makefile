serve:
	cd client && npm run-script build

dev: dev-assets serve run

run:
	godep go build -race -a && \
		./meep -config-path="$(shell pwd)/config.json"

dev-assets: clean deps
	go-bindata -debug client/serve/...

prod: clean assets
	godep go build -a

clean:
	rm -f bindata.go
	rm -f bindata_assetfs.go

assets: clean deps serve
	go-bindata-assetfs client/serve/...

deps:
	go get github.com/jteeuwen/go-bindata/...
	go get github.com/elazarl/go-bindata-assetfs/...

.PHONY: serve dev run dev-assets prod clean assets deps
