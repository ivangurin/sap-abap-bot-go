lint:
	go tool golangci-lint run \
	--new-from-rev=main \
	--config=.golangci.yaml \
	--max-issues-per-linter=1000 \
	--max-same-issues=1000 \
	./...

lint-full:
	go tool golangci-lint run \
	--config=.golangci.yaml \
	--max-issues-per-linter=1000 \
	--max-same-issues=1000 \
	./...

genmock: 
	go tool mockery

generate: genmock

test:
	go test -v -race -count=1 ./...

run:
	go run ./cmd/bot

pull-and-run:
	docker compose down || true
	docker container rm sap-abap-bot || true
	docker image rm ghcr.io/ivangurin/sap-abap-bot-go || true
	docker compose up -d