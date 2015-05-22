serve:
	cd client && npm run-script build

dev: dev-assets serve
	godep go build -race -a && ./meep \
		-bot-key="14c5f71d-bf9f-491a-a940-bc5c14f1744a" \
		-bot-secret="85584c70bf15cbd89d20ab7a6438fefc227d712b5cfb77e7" \
		-client-id="a749c174-00b6-472a-8017-e59e85c85eea" \
		-client-secret="c662cab59262a71dd809a6a2ea6e983f072960b3" \
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
	go-bindata -debug client/serve/...

prod: assets

assets:
	-rm bindata.go
	-rm bindata_assetfs.go
	go-bindata-assetfs client/serve/... 

