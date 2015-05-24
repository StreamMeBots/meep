serve:
	cd client && npm run-script build

serve:
	cd client && npm run-script build

dev: dev-assets serve run

run:
	godep go build -race -a && \
		./meep -config-path="$(shell pwd)/config.json"

dev-assets:
	-rm bindata.go
	-rm bindata_assetfs.go
	go get github.com/jteeuwen/go-bindata/...
	go get github.com/elazarl/go-bindata-assetfs/...
	go-bindata -debug client/serve/...

prod: assets

assets:
	-rm bindata.go
	-rm bindata_assetfs.go
	go-bindata-assetfs client/serve/... 

