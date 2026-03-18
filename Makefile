download-docs:
	bash ./scripts/update-docs.sh

download-blog:
	bash ./scripts/update-blog.sh

update-docs: download-docs

update-blog: download-blog

update-resources: update-docs update-blog

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

dev-playground: build
	VM_INSTANCE_ENTRYPOINT=https://play.victoriametrics.com \
	VM_INSTANCE_TYPE=cluster \
	MCP_SERVER_MODE=http \
	MCP_LISTEN_ADDR=:8080 \
	./mcp-victoriametrics

all: test check lint build
