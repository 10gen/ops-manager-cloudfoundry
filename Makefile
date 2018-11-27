#!/bin/bash
PROJECT_ROOT := $(shell pwd)
export PROJECT_ROOT
GOPATH := $(PROJECT_ROOT)
PATH := $(PATH):$(GOPATH)/bin
OS_NAME := $(shell uname -s)
OS_ARCH := $(shell uname -p)

export GOPATH
export PATH
#
#OSFLAG 				:=
#ifeq ($(OS),Windows_NT)
#	OSFLAG :=$(OSFLAG)_windows
#	ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
#		OSFLAG:=$(OSFLAG)_amd64
#	endif
#	ifeq ($(PROCESSOR_ARCHITECTURE),x86)
#		OSFLAG:=$(OSFLAG)_386
#	endif
#else
#	UNAME_S := $(shell uname -s)
#	ifeq ($(UNAME_S),Linux)
#		OSFLAG:=$(OSFLAG)_linux
#	endif
#	ifeq ($(UNAME_S),Darwin)
#		OSFLAG:=$(OSFLAG)_darwin
#	endif
#		UNAME_P := $(shell uname -p)
#	ifeq ($(UNAME_P),x86_64)
#		OSFLAG:=$(OSFLAG)_amd64
#	endif
#		ifneq ($(filter %86,$(UNAME_P)),)
#	OSFLAG:=$(OSFLAG)_386
#		endif
#	ifneq ($(filter arm%,$(UNAME_P)),)
#		OSFLAG:=$(OSFLAG)_arm
#	endif
#endif
ifeq ($(OS),Windows_NT)
    CCFLAGS += -D WIN32
    ifeq ($(PROCESSOR_ARCHITEW6432),AMD64)
        CCFLAGS += -D AMD64
    else
        ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
            CCFLAGS += -D AMD64
        endif
        ifeq ($(PROCESSOR_ARCHITECTURE),x86)
            CCFLAGS += -D IA32
        endif
    endif
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
        CCFLAGS += -D LINUX
    endif
    ifeq ($(UNAME_S),Darwin)
        CCFLAGS += -D OSX
    endif
    UNAME_P := $(shell uname -p)
    ifeq ($(UNAME_P),x86_64)
        CCFLAGS += -D AMD64
    endif
    ifneq ($(filter %86,$(UNAME_P)),)
        CCFLAGS += -D IA32
    endif
    ifneq ($(filter arm%,$(UNAME_P)),)
        CCFLAGS += -D ARM
    endif
endif
TERRAFORM_DOWNLOAD_URL := https://releases.hashicorp.com/terraform/0.11.10/terraform_0.11.10

#TODO:
#	Add brew install for MacOS
test:
	@echo $(OSFLAG)
	@echo $(shell uname -s)
	@echo $(shell uname -p)
	@echo $(TERRAFORM_DOWNLOAD_URL)$(OSFLAG).zip

#Installing Ops-Manager 1.8+ https://github.com/pivotal-cf/om
install/om:
ifeq ($(shell uname -s), Linux)
	wget "https://github.com/pivotal-cf/om/releases/download/0.44.0/om-linux"
endif
ifeq ($(shell uname -s), Windows)
	wget "https://github.com/pivotal-cf/om/releases/download/0.44.0/om-windows.exe"
endif
ifeq ($(shell uname -s), Darwin)
	wget "https://github.com/pivotal-cf/om/releases/download/0.44.0/om-darwin"
endif
	#	OSFLAG 				:=
	#	ifeq ($(OS),Windows_NT)
	#		OSFLAG :=$(OSFLAG)_windows
	#		ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
	#			OSFLAG:=$(OSFLAG)_amd64
	#		endif
	#		ifeq ($(PROCESSOR_ARCHITECTURE),x86)
	#			OSFLAG:=$(OSFLAG)_386
	#		endif
	#	endif
#	OS_NAME :=
#	ifeq ( $(shell uname -s),Linux)
#		ops_url:=$(ops_url)linux
#	endif
#	ifeq ( $(shell uname -s),Darwin)
#		ops_url:=$(ops_url)darwin
#	endif
#	@echo $(ops_url)
#	wget -O ops_manager_$(shell uname -s) $(ops_url)
	#apt-key add -echo "deb http://apt.starkandwayne.com stable main" 
	#tee /etc/apt/sources.list.d/starkandwayne.list
	#apt-get update
	#apt-get install om
	#go get -u github.com/pivotal-cf/texplate



#Installing the Cloud Foundary CLI https://docs.cloudfoundry.org/cf-cli/install-go-cli.html
install/pcf-cli:
	wget -q -O - https://packages.cloudfoundry.org/debian/cli.cloudfoundry.org.key | sudo apt-key add -
	echo "deb https://packages.cloudfoundry.org/debian stable main" | sudo tee /etc/apt/sources.list.d/cloudfoundry-cli.list
	sudo apt-get update
	sudo apt-get install cf-cli
	#dpkg -i path/to/cf-cli-*.deb && apt-get install -f

install/terraform:
	#https://releases.hashicorp.com/terraform/0.11.10/terraform_0.11.10_linux_amd64.zip
	#https://releases.hashicorp.com/terraform/0.11.10/terraform_0.11.10_darwin_amd64.zip
	#https://releases.hashicorp.com/terraform/0.11.10/terraform_0.11.10_windows_amd64.zip
	TERRAFORM_DOWNLOAD_URL := "https://releases.hashicorp.com/terraform/0.11.10/terraform" + OSFLAG
	export PATH="$PATH:/path/to/dir"

pfc-product/download:

pfc-product/upload:

install/pas-tile:

install/mongodb-tile:
