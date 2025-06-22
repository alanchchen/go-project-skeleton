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
BUILD_BIN_DIR := $(BUILD_DIR)/bin
DIRS := \
	$(BUILD_BIN_DIR)

# Define your docker repository
DOCKER_REPOSITORY ?= ghcr.io/alanchchen/go-project-skeleton
DOCKER_COMMAND ?= docker  # or podman
REV ?= $(shell git rev-parse --short HEAD 2> /dev/null)

define find-subdir
$(shell find $(1) -maxdepth 1 -mindepth 1 -type d -o -type l)
endef

APPS := $(sort $(notdir $(call find-subdir,cmd)))
PHONY += $(APPS)

all: $(APPS)

.SECONDEXPANSION:
$(APPS): $(addprefix $(BUILD_BIN_DIR)/,$$@)

$(DIRS):
	$(Q)mkdir -p $@

$(BUILD_BIN_DIR)/%: $(BUILD_BIN_DIR) FORCE
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
	$(Q)rm -fr $(BUILD_BIN_DIR)

.PHONY: help
help:
	@echo  'App targets:'
	@echo  '  all                         - Build all apps marked with [*]'
	@for app in $(APPS); do \
		printf "* %s\n" $$app; done
	@echo  ''
	@echo  'Docker targets:'
	@echo  '  docker                      - Build docker image which includes all executables'
	@echo  '  docker-push                 - Push $(DOCKER_REPOSITORY):$(REV)'
	@echo  ''
	@echo  'Test targets:'
	@echo  '  test                        - Run all tests'
	@echo  ''
	@echo  'Cleaning targets:'
	@echo  '  clean                       - Remove built executables'
	@echo  ''
	@echo  'Execute "make" or "make all" to build all apps marked with [*] '
	@echo  'For further info see the ./README.md file'

.PHONY: $(PHONY)

.PHONY: FORCE
FORCE:
