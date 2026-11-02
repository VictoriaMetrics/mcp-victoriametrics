download-docs:
	bash ./scripts/update-docs.sh

download-blog:
	bash ./scripts/update-blog.sh

update-docs: download-docs

update-blog: download-blog

update-resources: update-docs update-blog

update-readme:
	go run ./scripts/readmegen

test:
	bash ./scripts/test-all.sh

check:
	bash ./scripts/check-all.sh

lint:
	bash ./scripts/lint-all.sh

build: ui-build
	bash ./scripts/build-binaries.sh

ui-install:
	cd web && npm install

ui-build: ui-install
	cd web && npm run build && touch ./dist/.gitkeep

ui-dev:
	cd web && npm run dev

dev-run: build
	MCP_SERVER_MODE=http \
	MCP_LISTEN_ADDR=:8080 \
	./mcp-victoriametrics

dev-logs: export VM_INSTANCE_BACKEND = logs
dev-logs: dev-run

dev-traces: export VM_INSTANCE_BACKEND = traces
dev-traces: dev-run

dev-playground: export VM_INSTANCE_ENTRYPOINT = https://play.victoriametrics.com
dev-playground: export VM_INSTANCE_TYPE = cluster
dev-playground: dev-run

dev-playground-logs: export VM_INSTANCE_ENTRYPOINT = https://play-vmlogs.victoriametrics.com
dev-playground-logs: dev-logs

dev-playground-traces: export VM_INSTANCE_ENTRYPOINT = https://play-vtraces.victoriametrics.com
dev-playground-traces: dev-traces

all: test check lint build
