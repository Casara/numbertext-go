# Dev tools, version-pinned via `go run` instead of a globally installed
# binary: numbertext-go is a library, so it doesn't make sense to pull
# linter/vuln-scanner transitive dependencies into go.mod/go.sum for
# whoever depends on it.
GOLANGCI_LINT := github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2
GO_ARCH_LINT  := github.com/fe3dback/go-arch-lint@v1.16.0
GREMLINS      := github.com/go-gremlins/gremlins/cmd/gremlins@v0.6.0
# Unpinned on purpose: govulncheck is only useful when it knows about the
# vulnerabilities disclosed since the last release, and it reads the
# database over the network anyway. Pinning it would freeze the analyzer,
# not the data.
GOVULNCHECK   := golang.org/x/vuln/cmd/govulncheck@latest

# The one non-Go dev tool here, run via npx instead of a global install.
# Requires Node.js >= 20.
MARKDOWNLINT := markdownlint-cli2@0.23.2

.PHONY: help build run fmt lint lint-fix lint-md arch-lint vuln \
	test test-race coverage test-select test-mutation test-mutation-dry-run \
	verify-docs doc-sync gen-locales sync-data check clean

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'

build: ## Build the CLI (bin/numbertext)
	@go build -o bin/numbertext ./cmd/numbertext

run: ## Run the CLI (pass flags via ARGS, e.g. make run ARGS="-lang en -cardinal 42")
	@go run ./cmd/numbertext $(ARGS)

fmt: ## Format the code (gofmt/goimports, via golangci-lint)
	@go run $(GOLANGCI_LINT) fmt

lint: ## Run the linter (golangci-lint v2, default:all - see .golangci.yml)
	@go run $(GOLANGCI_LINT) run

lint-fix: ## Run the linter and apply available auto-fixes
	@go run $(GOLANGCI_LINT) run --fix

lint-md: ## Lint markdown files (markdownlint-cli2, requires Node.js >= 20)
	@npx --yes $(MARKDOWNLINT) "**/*.md" "#NOTES.md" "#TEMP.md"

arch-lint: ## Check the package dependency graph (.go-arch-lint.yml)
# internal/soros must stay a self-contained interpreter with no
# knowledge of numbertext/cmd - this is what actually enforces the
# invariant CLAUDE.md/AGENTS.md only document.
	@go run $(GO_ARCH_LINT) check

vuln: ## Scan for known vulnerabilities reachable from this code (govulncheck)
	@go run $(GOVULNCHECK) ./...

test: ## Run the tests
	@go test ./...

test-race: ## Run the tests with the race detector
	@go test -race ./...

coverage: ## Generate coverage.out and coverage.html with the coverage report
	@go test -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Report at coverage.html"

test-select: ## Sanity-check the build-tag selective-embedding feature (see README)
# A build choosing only "en"+"pt" must still pass the tests exercising
# those two languages - proof that opting in doesn't silently break the
# default-everything build's behavior for the languages you did keep.
	@go build -tags numbertext_select ./...
	@go test -tags "numbertext_select numbertext_lang_en numbertext_lang_pt" ./... -run 'English|Portuguese|Gender'

test-mutation: ## Run mutation tests (gremlins), writes mutation.json
# gremlins doesn't handle Go's "./..." pattern well (silently reports
# nothing, verified); passing "." makes it recurse through the whole
# module on its own.
	@go run $(GREMLINS) unleash -o mutation.json .

test-mutation-dry-run: ## List the mutants without running the tests (much faster)
	@go run $(GREMLINS) unleash --dry-run .

verify-docs: ## Compile every self-contained Go program embedded in the Markdown
# Complements `make check`: that proves the code is correct, not that the
# documentation still describes it. Illustrative fragments (a bare
# statement, a snippet meant to be read alongside surrounding prose)
# can't compile standalone and are reported as unverified rather than
# silently assumed correct.
	@bash scripts/verify-docs.sh

doc-sync: ## List documentation a change to the exported API has left behind
# Deliberately never fails: a diff can touch an exported symbol without
# any of the listed files needing an edit, and a check that cries wolf
# gets ignored. It covers the residue verify-docs cannot - fragments and
# prose.
	@bash scripts/doc-sync-check.sh

gen-locales: ## Regenerate locale_<code>.go after data/*.sor gains or loses a file
	@python3 scripts/gen-locale-embeds.py
	@go run $(GOLANGCI_LINT) fmt

sync-data: ## Refresh data/*.sor from upstream libnumbertext (see data/SOURCE.md)
	@scripts/sync-data.sh

check: lint arch-lint verify-docs test-race test-select ## Run lint + arch-lint + verify-docs + tests with the race detector (the minimum CI should run)

clean: ## Remove build/test artifacts
	@rm -rf bin coverage.out coverage.html mutation.json
