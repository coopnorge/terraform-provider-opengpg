# Build parameters
CGO_ENABLED=0
LD_FLAGS="-extldflags '-static'"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOBUILD=CGO_ENABLED=$(CGO_ENABLED) $(GOCMD) build -v -buildmode=exe -ldflags $(LD_FLAGS)
GO_PACKAGES=./...
GO_TESTS=^.*$

GOLANGCI_LINT_VERSION=v1.44.0

BIN_PATH=$$HOME/bin

COVERPROFILE=c.out
CC_TEST_REPORTER_ID=7b59c847e0a4db01a8464de8547df8e2c9875dbbf7e4e1e474825885f8a3c133

.PHONY: all
all: build build-test test lint

.PHONY: download
download:
	$(GOMOD) download

.PHONY: install-golangci-lint
install-golangci-lint:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(BIN_PATH) $(GOLANGCI_LINT_VERSION)

.PHONY: install-cc-test-reporter
install-cc-test-reporter:
	curl -L https://codeclimate.com/downloads/test-reporter/test-reporter-latest-linux-amd64 > $(BIN_PATH)/cc-test-reporter
	chmod +x $(BIN_PATH)/cc-test-reporter

.PHONY: install-ci
install-ci: install-golangci-lint install-cc-test-reporter

.PHONY: build
build:
	$(GOBUILD)

.PHONY: test
test: build-test
	$(GOTEST) -run $(GO_TESTS) $(GO_PACKAGES)

.PHONY: lint
lint:
	golangci-lint run $(GO_PACKAGES)

.PHONY: build-test
build-test:
	$(GOTEST) -run=nope $(GO_PACKAGES)

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(OUTPUT_FILE) || true
	rm -f $(OUTPUT_FILE).sig || true

.PHONY: update
update:
	$(GOGET) -u $(GO_PACKAGES)
	$(GOMOD) tidy

.PHONY: all-cover
all-cover: build build-test test-cover lint

.PHONY: test-cover
test-cover: build-test
	$(GOTEST) -run $(GO_TESTS) -coverprofile=$(COVERPROFILE) $(GO_PACKAGES)

.PHONY: test-cover-upload-codecov
test-cover-upload-codecov: SHELL=/bin/bash
test-cover-upload-codecov: test-cover
test-cover-upload-codecov:
	bash <(curl -s https://codecov.io/bash) -f $(COVERPROFILE)

.PHONY: test-cover-upload-codeclimate
test-cover-upload-codeclimate: test-cover
test-cover-upload-codeclimate:
	env CC_TEST_REPORTER_ID=$(CC_TEST_REPORTER_ID) cc-test-reporter after-build -t gocov -p $$(go list -m)

.PHONY: test-cover-upload
test-cover-upload: test-cover-upload-codecov test-cover-upload-codeclimate

.PHONY: codespell
codespell:
	codespell -S .git,state.yaml,go.sum,terraform.tfstate,terraform.tfstate.backup,./local-testing/resources -L decorder

.PHONY: codespell-pr
codespell-pr:
	git diff master..HEAD | grep -v ^- | codespell -
	git log master..HEAD | codespell -

.PHONY: test-working-tree-clean
test-working-tree-clean:
	@test -z "$$(git status --porcelain)" || (echo "Commit all changes before running this target"; exit 1)

.PHONY: test-changelog
test-changelog: test-working-tree-clean
	make format-changelog
	@test -z "$$(git status --porcelain)" || (echo "Please run 'make format-changelog' and commit generated changes."; git diff; exit 1)

.PHONY: format-changelog
format-changelog:
	changelog fmt -o CHANGELOG.md.fmt
	mv CHANGELOG.md.fmt CHANGELOG.md

.PHONY: install-changelog
install-changelog:
	go get github.com/rcmachado/changelog@0.7.0
	go mod tidy

.PHONY: update-linters
update-linters:
	# Remove all enabled linters.
	sed -i '/^  enable:/q0' .golangci.yml
	# Then add all possible linters to config.
	golangci-lint linters | grep -E '^\S+:' | cut -d: -f1 | sort | sed 's/^/    - /g' | grep -v -E "($$(grep '^  disable:' -A 100 .golangci.yml  | grep -E '    - \S+$$' | awk '{print $$2}' | tr \\n '|' | sed 's/|$$//g'))" >> .golangci.yml

.PHONY: test-update-linters
test-update-linters: test-working-tree-clean
	make update-linters
	@test -z "$$(git status --porcelain)" || (echo "Linter configuration outdated. Run 'make update-linters' and commit generated changes to fix."; exit 1)

.PHONY: test-tidy
test-tidy: test-working-tree-clean
	go mod tidy
	@test -z "$$(git status --porcelain)" || (echo "Please run 'go mod tidy' and commit generated changes."; git diff; exit 1)
