-include Makefile.local

.PHONY: fe-dev


fe-dev:
	@cd frontend && npm run dev
