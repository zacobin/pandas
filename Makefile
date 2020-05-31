IMPORT_PATH = github.com/cloustone/pandas
BUILD_DIR = bin
# V := 1 # When V is set, print cmd and build progress.
Q := $(if $V,,@)

VERSION          := $(shell git describe --tags --always --dirty="-dev")
DATE             := $(shell date -u '+%Y-%m-%d-%H%M UTC')
VERSION_FLAGS    := -ldflags='-X "main.Version=$(VERSION)" -X "main.BuildTime=$(DATE)"'

# Space separated patterns of packages to skip in list, test, format.
DOCKER_NAMESPACE := cloustone

# Space separated patterns of packages to skip in list, test, format.
IGNORED_PACKAGES := /vendor/

SERVICES = swagger authn authz things bootstrap twins users vms realms swagger rulechain lbs 
	
ADAPTOR_SERVICE = http ws coap lora opcua mqtt cli
	
ADDONE_SERVICE = influxdb-writer influxdb-reader mongodb-writer mongodb-reader \
				cassandra-writer cassandra-reader postgres-writer postgres-reader

UNAME = $(shell uname)
#DOCKER_REPO = 10.4.47.129:5000
IMAGES = rulechain lbs authn authz things bootstrap twins users vms pms realms swagger
ADAPTOR_IMAGES = http ws coap lora opcua mqtt cli
ADDONE_IMAGES = influxdb-writer influxdb-reader mongodb-writer mongodb-reader \
				cassandra-writer cassandra-reader postgres-writer postgres-reader
IMAGE_NAME_PREFIX := pandas-
IMAGE_DIR := $(IMAGE_NAME)
ifeq ($(IMAGE_NAME),bridge)
    IMAGE_DIR := edge/$(IMAGE_NAME)
else ifneq (,$(filter $(IMAGE_NAME),  rulechain lbs authn))
    IMAGE_DIR := cmd/$(IMAGE_NAME)
else ifeq ($(IMAGE_NAME),cabinet)
    IMAGE_DIR := security/$(IMAGE_NAME)
endif

GCFLAGS  := -gcflags="-N -l"
CGO_ENABLED ?= 0
GOARCH ?= amd64
GOOS ?= linux

define compile_service
	@echo building service $(1) ...
	$Q CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o bin/$@ $(GCFLAGS) $(if $V,-v) $(VERSION_FLAGS) $(IMPORT_PATH)/cmd/$(1)
endef

DOCKERS_DEV = $(addprefix docker_dev_,$(IMAGES)) 
DOCKERS_ADAPTOR = $(addprefix docker_dev_,$(ADAPTOR_IMAGES))
DOCKERS_ADDONE = $(addprefix docker_dev_,$(ADDONE_IMAGES))
define make_docker_dev
	$(eval svc=$(subst docker_dev_,,$(1)))
	@echo building $(IMAGE_NAME_PREFIX)$(svc) image ...
	@if [ ! -d "cmd/$(svc)/bin/" ]; then mkdir cmd/$(svc)/bin/ ; fi
	@cp scripts/dockerize cmd/$(svc)/bin/
	cp bin/$(svc) cmd/$(svc)/bin/main
	@full_img_name=$(IMAGE_NAME_PREFIX)$(svc); \
		cd ./cmd/$(svc)/ && \
			docker build -t $(DOCKER_NAMESPACE)/$$full_img_name ../../../. -f Dockerfile.dev 
	@rm -rf cmd/$(svc)/bin
endef


.PHONY: $(SERVICES) $(ADAPTOR_SERVICE) $(ADDONE_SERVICE)
all: $(SERVICES) $(ADAPTOR_SERVICE) $(ADDONE_SERVICE)
	@echo evironment is [CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH)]
service: $(SERVICES)
adaptor: $(ADAPTOR_SERVICE)
addone: $(ADDONE_SERVICE)
$(SERVICES):
	$(call compile_service,$(@))
$(ADAPTOR_SERVICE):
	$(call compile_service,$(@))
$(ADDONE_SERVICE):
	$(call compile_service,$(@))

clean: 
	rm -rf ${BUILD_DIR}

.PHONY: docker dockers_dev deploy upgrade test undeploy test
docker: export GOOS=linux
docker: $(addprefix docker-build-, $(IMAGES)) 
	docker images | grep '<none>' | awk '{print $3}' | xargs docker rmi
	@echo "docker building completed!" 

# Docker build targets
$(addprefix docker-build-, $(IMAGES)): docker-build-%: %
	@IMAGE_NAME=$< make .docker-build

.docker-build:
	@echo building $(IMAGE_NAME_PREFIX)$(IMAGE_NAME) image ...
	@if [ ! -d "$(IMAGE_DIR)/bin/" ]; then mkdir $(IMAGE_DIR)/bin/ ; fi
	@cp scripts/dockerize $(IMAGE_DIR)/bin/
#	@if [ "$(UNAME)" = "Linux" ]; then cp bin/$(IMAGE_NAME) $(IMAGE_DIR)/bin/main ; fi
#	@if [ "$(UNAME)" = "Darwin" ]; then cp bin/linux_amd64/$(IMAGE_NAME) $(IMAGE_DIR)/bin/main ; fi
	cp bin/$(IMAGE_NAME) $(IMAGE_DIR)/bin/main
	@full_img_name=$(IMAGE_NAME_PREFIX)$(IMAGE_NAME); \
		cd ./$(IMAGE_DIR)/ && \
			docker build -t $(DOCKER_REPO)/$(DOCKER_NAMESPACE)/$$full_img_name ../../../. -f Dockerfile.dev 
	@rm -rf $(IMAGE_DIR)/bin
	@"./scripts/push.sh" $(IMAGE_NAME)
	# @kubectl delete pod $$(kubectl get pod -n pandas | grep $(IMAGE_NAME) | awk '{print $$1}') -n pandas 

pandas-base: export GOOS=linux
pandas-base: 
	@echo building $(IMAGE_NAME_PREFIX)pandas-base image ...
	docker build -t $(DOCKER_REPO)/$(DOCKER_NAMESPACE)/pandas-base . -f docker/base/Dockerfile

$(DOCKERS_DEV):
	$(call make_docker_dev,$(@))
$(DOCKERS_ADAPTOR):
	$(call make_docker_dev,$(@))
$(DOCKERS_ADDONE):
	$(call make_docker_dev,$(@))
dockers_dev: $(DOCKERS_DEV)
dockers_adaptor: $(DOCKERS_ADAPTOR)
dockers_addone:	$(DOCKERS_ADDONE)

deploy:
	@helm install .

upgrade:
	@existing=$$(helm list | grep pandas | awk '{print $$1}' | head -n 1); \
		(if [ ! -z "$$existing" ]; then echo "Upgrade the stack via helm. This may take a while."; helm upgrade "$$existing"; echo "The stack has been upgraded."; fi) > /dev/null;

undeploy:
	@existing=$$(helm list | grep pandas | awk '{print $$1}' | head -n 1); \
		(if [ ! -z "$$existing" ]; then echo "Undeploying the stack via helm. This may take a while."; helm del --purge "$$existing"; echo "The stack has been undeployed."; fi) > /dev/null;

test: 
	$Q go test  ./...






