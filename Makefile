build:
	go build -o ai-check-guardrails ./cmd/ai-check-guardrails

gen-docs:
	go run ./cmd/gendocs

preview-docs: gen-docs
	uvx --with mkdocs-material mkdocs serve --dev-addr 0.0.0.0:8000

build-docs: gen-docs
	uvx --with mkdocs-material mkdocs build --strict
