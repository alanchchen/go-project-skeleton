VERBOSE_ORIGINS := "command line" "environment"
ifdef V
  ifeq ($(filter $(VERBOSE_ORIGINS),$(origin V)),)
    BUILD_VERBOSE := $(V)
  endif
endif

ifndef BUILD_VERBOSE
  BUILD_VERBOSE := 0
endif

ifeq ($(BUILD_VERBOSE),1)
  Q :=
else
  Q := @
endif

PHONY += all test clean docker docker-push
CURDIR := $(shell pwd)
BUILD_DIR ?= $(CURDIR)/build
GOBIN_DIR := $(BUILD_DIR)/bin
DIRS := \
	$(GOBIN_DIR)

HOST_OS := $(shell uname -s)

# Define your docker repository
DOCKER_REPOSITORY ?= ghcr.io/alanchchen/go-project-skeleton
DOCKER_COMMAND ?= docker  # or podman
REV ?= $(shell git rev-parse --short HEAD 2> /dev/null)
GOPATH ?= $(shell go env GOPATH)

export PATH:=$(GOTOOLS_DIR):$(PATH)
export REV

define find-subdir
$(shell find $(1) -maxdepth 1 -mindepth 1 -type d -o -type l)
endef

APPS := $(sort $(notdir $(call find-subdir,cmd)))
PHONY += $(APPS)

all: $(APPS)

.SECONDEXPANSION:
$(APPS): $(addprefix $(GOBIN_DIR)/,$$@)

$(DIRS) :
	$(Q)mkdir -p $@

$(GOBIN_DIR)/%: $(GOBIN_DIR) FORCE
	$(Q)go build -o $@ ./cmd/$(notdir $@)
	@echo "Done building."
	@echo "Run \"$(subst $(CURDIR),.,$@)\" to launch $(notdir $@)."

docker:
	$(Q)$(DOCKER_COMMAND) build -t $(DOCKER_REPOSITORY):$(REV) .

docker-push:
	$(Q)$(DOCKER_COMMAND) push $(DOCKER_REPOSITORY):$(REV)

test:
	$(Q)go test -v ./...

clean:
	$(Q)rm -fr $(GOBIN_DIR) $(HOST_DIR)

.PHONY: help
help:
	@echo  'Generic targets:'
	@echo  '  all                         - Build all targets marked with [*]'
	@for app in $(APPS); do \
		printf "* %s\n" $$app; done
	@echo  ''
	@echo  'Docker targets:'
	@echo  '  docker                      - Build docker image which includes all executables'
	@echo  ''
	@echo  '  docker-push                 - Push $(DOCKER_REPOSITORY):$(REV)'
	@echo  ''
	@echo  'Test targets:'
	@echo  '  test                        - Run all tests'
	@echo  ''
	@echo  'Cleaning targets:'
	@echo  '  clean                       - Remove built executables'
	@echo  ''
	@echo  'Execute "make" or "make all" to build all targets marked with [*] '
	@echo  'For further info see the ./README.md file'

.PHONY: $(PHONY)

.PHONY: FORCE
FORCE:
