preview-docs:
	uvx --with mkdocs-material mkdocs serve --dev-addr 0.0.0.0:8000

build-docs:
	uvx --with mkdocs-material mkdocs build --strict
