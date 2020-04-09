IMPORT_PATH = github.com/cloustone/pandas
# V := 1 # When V is set, print cmd and build progress.
Q := $(if $V,,@)

VERSION          := $(shell git describe --tags --always --dirty="-dev")
DATE             := $(shell date -u '+%Y-%m-%d-%H%M UTC')
VERSION_FLAGS    := -ldflags='-X "main.Version=$(VERSION)" -X "main.BuildTime=$(DATE)"'

# Space separated patterns of packages to skip in list, test, format.
DOCKER_NAMESPACE := cloustone

# Space separated patterns of packages to skip in list, test, format.
IGNORED_PACKAGES := /vendor/

MAINFLUX_SERVICES = http ws coap lora influxdb-writer influxdb-reader mongodb-writer \
	mongodb-reader cassandra-writer cassandra-reader postgres-writer postgres-reader cli \
  opcua  mqtt

UNAME = $(shell uname)
DOCKER_REPO = docker.io
IMAGES = apimachinery dmms pms rulechain headmast lbs authn
IMAGE_NAME_PREFIX := pandas-
IMAGE_DIR := $(IMAGE_NAME)
ifeq ($(IMAGE_NAME),bridge)
    IMAGE_DIR := edge/$(IMAGE_NAME)
else ifneq (,$(filter $(IMAGE_NAME), apimachinery dmms pms rulechain headmast lbs authn))
    IMAGE_DIR := cmd/$(IMAGE_NAME)
else ifeq ($(IMAGE_NAME),cabinet)
    IMAGE_DIR := security/$(IMAGE_NAME)
endif

GCFLAGS  := -gcflags="-N -l"

DOCKERS_DEV = $(addprefix docker_dev_,$(IMAGES))
define make_docker_dev
	$(eval svc=$(subst docker_dev_,,$(1)))
	@echo IMAGE_DIR is $(IMAGE_DIR1) 
	@echo svc is $(svc)
	@echo building $(IMAGE_NAME_PREFIX)$(svc) image ...
	@if [ ! -d "cmd/$(svc)/bin/" ]; then mkdir cmd/$(svc)/bin/ ; fi
	@cp scripts/dockerize cmd/$(svc)/bin/
	cp bin/$(svc) cmd/$(svc)/bin/main
	@full_img_name=$(IMAGE_NAME_PREFIX)$(svc); \
		cd ./cmd/$(svc)/ && \
			docker build -t $(DOCKER_REPO)/$(DOCKER_NAMESPACE)/$$full_img_name ../../../. -f Dockerfile.dev 
	@rm -rf cmd/$(svc)/bin
endef

.PHONY: all
all: build

.PHONY: docker
docker: export GOOS=linux
docker: pandas-base $(addprefix docker-build-, $(IMAGES)) 
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

pandas-base:
	@echo building $(IMAGE_NAME_PREFIX)pandas-base image ...
	docker build -t $(DOCKER_REPO)/$(DOCKER_NAMESPACE)/pandas-base . -f docker/base/Dockerfile

.PHONY: dockers_dev
$(DOCKERS_DEV):
	$(call make_docker_dev,$(@))
dockers_dev: $(DOCKERS_DEV)

.PHONY: deploy
deploy:
	@helm install .

.PHONY: upgrade
upgrade:
	@existing=$$(helm list | grep pandas | awk '{print $$1}' | head -n 1); \
		(if [ ! -z "$$existing" ]; then echo "Upgrade the stack via helm. This may take a while."; helm upgrade "$$existing"; echo "The stack has been upgraded."; fi) > /dev/null;

.PHONY: undeploy
undeploy:
	@existing=$$(helm list | grep pandas | awk '{print $$1}' | head -n 1); \
		(if [ ! -z "$$existing" ]; then echo "Undeploying the stack via helm. This may take a while."; helm del --purge "$$existing"; echo "The stack has been undeployed."; fi) > /dev/null;

.PHONY: all
all: build

.PHONY: build
build: apimachinery  dmms  pms rulechain lbs headmast  authn  users bootstrap realms

.PHONY: apimachinery 
apimachinery: 
	@echo "building api server (apimachinery)..."
	$Q CGO_ENABLED=0 go build -o bin/$@ $(GCFLAGS) $(if $V,-v) $(VERSION_FLAGS) $(IMPORT_PATH)/cmd/apimachinery 

.PHONY: dmms 
dmms: cmd/dmms 
	@echo "building device management server (dmms)..."
	$Q CGO_ENABLED=1 go build -o bin/$@ $(GCFLAGS) $(if $V,-v) $(VERSION_FLAGS) $(IMPORT_PATH)/cmd/dmms

.PHONY: pms 
pms: cmd/pms 
	@echo "building project management server (pms)..."
	$Q CGO_ENABLED=1 go build -o bin/$@ $(GCFLAGS) $(if $V,-v) $(VERSION_FLAGS) $(IMPORT_PATH)/cmd/pms

.PHONY: rulechain 
rulechain: cmd/rulechain
	@echo "building rulechain server (rulechain)..."
	$Q CGO_ENABLED=0 go build -o bin/$@ $(GCFLAGS) $(if $V,-v) $(VERSION_FLAGS) $(IMPORT_PATH)/cmd/rulechain

.PHONY: lbs 
lbs: cmd/lbs
	@echo "building location based service (lbs)..."
	$Q CGO_ENABLED=0 go build -o bin/$@ $(GCFLAGS) $(if $V,-v) $(VERSION_FLAGS) $(IMPORT_PATH)/cmd/lbs

.PHONY: headmast 
headmast: cmd/headmast
	@echo "building headmast service (headmast)..."
	$Q CGO_ENABLED=0 go build -o bin/$@ $(GCFLAGS) $(if $V,-v) $(VERSION_FLAGS) $(IMPORT_PATH)/cmd/headmast

.PHONY: authn 
authn: 
	@echo "building key management service (authn)..."
	$Q CGO_ENABLED=1 go build -o bin/$@ $(GCFLAGS) $(if $V,-v) $(VERSION_FLAGS) $(IMPORT_PATH)/cmd/authn

.PHONY: authz
authz: 
	@echo "building authorization service (authn)..."
	$Q CGO_ENABLED=1 go build -o bin/$@ $(GCFLAGS) $(if $V,-v) $(VERSION_FLAGS) $(IMPORT_PATH)/cmd/authz


.PHONY: users 
users: 
	@echo "building user management service (users)..."
	$Q CGO_ENABLED=1 go build -o bin/$@ $(GCFLAGS) $(if $V,-v) $(VERSION_FLAGS) $(IMPORT_PATH)/cmd/users

.PHONY: realms 
realms: 
	@echo "building identify authentiation management service (realms)..."
	$Q CGO_ENABLED=1 go build -o bin/$@ $(GCFLAGS) $(if $V,-v) $(VERSION_FLAGS) $(IMPORT_PATH)/cmd/realms



.PHONY: bootstrap 
bootstrap: cmd/bootstrap 
	@echo "building bootstrap service (bootstrap)..."
	$Q CGO_ENABLED=1 go build -o bin/$@ $(GCFLAGS) $(if $V,-v) $(VERSION_FLAGS) $(IMPORT_PATH)/cmd/bootstrap


.PHONY: mainflux 
mainflux: 
	@echo "building backend service (mainflux)..."
	$Q CGO_ENABLED=1 go build -o bin/$@ $(GCFLAGS) $(if $V,-v) $(VERSION_FLAGS) $(IMPORT_PATH)/cmd/things
	$Q CGO_ENABLED=1 go build -o bin/$@ $(GCFLAGS) $(if $V,-v) $(VERSION_FLAGS) $(IMPORT_PATH)/cmd/twins
	$Q CGO_ENABLED=1 go build -o bin/$@ $(GCFLAGS) $(if $V,-v) $(VERSION_FLAGS) $(IMPORT_PATH)/mainflux/cmd/ws
	$Q CGO_ENABLED=1 go build -o bin/$@ $(GCFLAGS) $(if $V,-v) $(VERSION_FLAGS) $(IMPORT_PATH)/mainflux/cmd/postgres-writer
	$Q CGO_ENABLED=1 go build -o bin/$@ $(GCFLAGS) $(if $V,-v) $(VERSION_FLAGS) $(IMPORT_PATH)/mainflux/cmd/postgres-reader
	$Q CGO_ENABLED=1 go build -o bin/$@ $(GCFLAGS) $(if $V,-v) $(VERSION_FLAGS) $(IMPORT_PATH)/mainflux/cmd/opcua
	$Q CGO_ENABLED=1 go build -o bin/$@ $(GCFLAGS) $(if $V,-v) $(VERSION_FLAGS) $(IMPORT_PATH)/mainflux/cmd/mqtt
	$Q CGO_ENABLED=1 go build -o bin/$@ $(GCFLAGS) $(if $V,-v) $(VERSION_FLAGS) $(IMPORT_PATH)/mainflux/cmd/mongodb-writer
	$Q CGO_ENABLED=1 go build -o bin/$@ $(GCFLAGS) $(if $V,-v) $(VERSION_FLAGS) $(IMPORT_PATH)/mainflux/cmd/mongodb-reader
	$Q CGO_ENABLED=1 go build -o bin/$@ $(GCFLAGS) $(if $V,-v) $(VERSION_FLAGS) $(IMPORT_PATH)/mainflux/cmd/lora
	$Q CGO_ENABLED=1 go build -o bin/$@ $(GCFLAGS) $(if $V,-v) $(VERSION_FLAGS) $(IMPORT_PATH)/mainflux/cmd/influxdb-writer
	$Q CGO_ENABLED=1 go build -o bin/$@ $(GCFLAGS) $(if $V,-v) $(VERSION_FLAGS) $(IMPORT_PATH)/mainflux/cmd/influxdb-reader
	$Q CGO_ENABLED=1 go build -o bin/$@ $(GCFLAGS) $(if $V,-v) $(VERSION_FLAGS) $(IMPORT_PATH)/mainflux/cmd/http
	$Q CGO_ENABLED=1 go build -o bin/$@ $(GCFLAGS) $(if $V,-v) $(VERSION_FLAGS) $(IMPORT_PATH)/mainflux/cmd/coap
	$Q CGO_ENABLED=1 go build -o bin/$@ $(GCFLAGS) $(if $V,-v) $(VERSION_FLAGS) $(IMPORT_PATH)/mainflux/cmd/cassandra-writer
	$Q CGO_ENABLED=1 go build -o bin/$@ $(GCFLAGS) $(if $V,-v) $(VERSION_FLAGS) $(IMPORT_PATH)/mainflux/cmd/cassandra-reader


.PHONY: test
test: 
	$Q go test  ./...





