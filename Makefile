dev: dev-assets
	godep go build -race -a && ./meep \
		-bot-key="" \
		-bot-secret="" \
		-client-id="" \
		-client-secret="" \
		-auth-url http://pds.dev.ifi.tv/api-auth/authorize \
		-token-url http://pds.dev.ifi.tv/api-auth/token \
		-server-host localhost \
		-server-port :8888 \
		-url http://pds.dev.ifi.tv \
		-chat-host pds.dev.ifi.tv:2020 \
		-debug

dev-assets:
	-rm bindata.go
	-rm bindata_assetfs.go
	go get github.com/jteeuwen/go-bindata/...
	go get github.com/elazarl/go-bindata-assetfs/...
	go-bindata -debug client/...

prod: assets

assets:
	-rm bindata.go
	-rm bindata_assetfs.go
	go-bindata-assetfs client/... 

