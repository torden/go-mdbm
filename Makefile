PKG_NAME=go-mdbm

VERSION				:= $(shell git describe --tags --always --dirty="-dev")
DATE				:= $(shell date -u '+%Y-%m-%d-%H%M UTC')
VERSION_FLAGS		:= -ldflags='-X "main.Version=$(VERSION)" -X "main.BuildTime=$(DATE)"'
PLATFORM        	:=$(shell uname -a)
CMD_RM          	:=$(shell which rm)
CMD_CC          	:=$(shell which gcc)
CMD_STRIP       	:=$(shell which strip)
CMD_DIFF        	:=$(shell which diff)
CMD_RM          	:=$(shell which rm)
CMD_BASH        	:=$(shell which bash)
CMD_CP          	:=$(shell which cp)
CMD_AR          	:=$(shell which ar)
CMD_RANLIB      	:=$(shell which ranlib)
CMD_MV          	:=$(shell which mv)
CMD_AWK				:=$(shell which awk)
CMD_SED				:=$(shell which sed)
CMD_TAIL        	:=$(shell which tail)
CMD_FIND        	:=$(shell which find)
CMD_LDD         	:=$(shell which ldd)
CMD_MKDIR       	:=$(shell which mkdir)
CMD_TEST        	:=$(shell which test)
CMD_SLEEP       	:=$(shell which sleep)
CMD_SYNC        	:=$(shell which sync)
CMD_LN          	:=$(shell which ln)
CMD_ZIP        		:=$(shell which zip)
CMD_MD5SUM      	:=$(shell which md5sum)
CMD_READELF     	:=$(shell which readelf)
CMD_GDB         	:=$(shell which gdb)
CMD_FILE        	:=$(shell which file)
CMD_ECHO        	:=$(shell which echo)
CMD_NM          	:=$(shell which nm)
CMD_GO				:=$(shell which go)
CMD_GOLINT			:=$(shell which golint)
CMD_GOMETALINTER	:=$(shell which gometalinter)
CMD_GOIMPORTS		:=$(shell which goimport)
CMD_MAKE2HELP		:=$(shell which make2help)
CMD_GLIDE			:=$(shell which glide)
CMD_GOVER			:=$(shell which gover)
CMD_GOVERALLS		:=$(shell which goveralls)

PATH_REPORT=report
PATH_RACE_REPORT=$(PKG_NAME).race.report
PATH_CONVER_PROFILE=$(PKG_NAME).coverprofile
PATH_PROF_CPU=$(PKG_NAME).cpu.prof
PATH_PROF_MEM=$(PKG_NAME).mem.prof
PATH_PROF_BLOCK=$(PKG_NAME).block.prof
PATH_PROF_MUTEX=$(PKG_NAME).mutex.prof

VER_GOLANG=$(shell go version | awk '{print $$3}' | sed -e "s/go//;s/\.//g")
GOLANGV18_OVER=$(shell [ "$(VER_GOLANG)" -ge "180" ] && echo 1 || echo 0)
GOLANGV16_OVER=$(shell [ "$(VER_GOLANG)" -ge "160" ] && echo 1 || echo 0)

all: clean setup build

## Setup Build Environment
setup: installpkgs metalinter

## Install Packages
installpkgs::
	@$(CMD_ECHO)  -e "\033[1;40;32mInstall Packages.\033[01;m\x1b[0m"
	@$(CMD_GO) get github.com/Masterminds/glide
	@$(CMD_GO) get github.com/Songmu/make2help/cmd/make2help
	@$(CMD_GO) get github.com/davecgh/go-spew/spew
	@$(CMD_GO) get github.com/k0kubun/pp
	@$(CMD_GO) get github.com/mattn/goveralls
	@$(CMD_GO) get golang.org/x/tools/cmd/cover
	@$(CMD_GO) get github.com/modocache/gover
	@$(CMD_GO) get github.com/boltdb/bolt
ifeq ($(GOLANGV16_OVER),1)
	@$(CMD_GO) get github.com/golang/lint/golint
	@$(CMD_GO) get github.com/alecthomas/gometalinter
endif
	@$(CMD_GO) get -u github.com/awalterschulze/gographviz
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Build the go-mdbm
build: lint
	@$(CMD_ECHO)  -e "\033[1;40;32mBuilding\033[01;m\x1b[0m"
	@$(CMD_GO) build -a -n -v
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Install GoMetaLinter 
metalinter::
	@$(CMD_ECHO)  -e "\033[1;40;32mInstall Go-metalineter.\033[01;m\x1b[0m"
ifeq ($(GOLANGV16_OVER),1)
	@$(shell which gometalinter) --install
else
	@$(CMD_ECHO) -e "\033[1;40;36mSKIP: your golang is older version $(shell go version)\033[01;m\x1b[0m"
endif
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Run a LintChecker (Normal)
lint: setup
	@$(CMD_ECHO)  -e "\033[1;40;32mRun a LintChecker (Normal).\033[01;m\x1b[0m"
ifeq ($(GOLANGV16_OVER),1)
	@$(CMD_GO) vet $$($(shell which glide) novendor)
	@for pkg in $$($(shell which glide) novendor -x); do \
		$(CMD_GOLINT) -set_exit_status $$pkg || exit $$?; \
	done
else
	@$(CMD_ECHO) -e "\033[1;40;36mSKIP: your golang is older version $(shell go version)\033[01;m\x1b[0m"
endif
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Run a LintChecker (Strict)
strictlint: setup
	@$(CMD_ECHO)  -e "\033[1;40;32mRun a LintChecker (Strict).\033[01;m\x1b[0m"
ifeq ($(GOLANGV16_OVER),1)
	@$(CMD_GOMETALINTER) $$($(shell which glide) novendor)
else
	@$(CMD_ECHO) -e "\033[1;40;36mSKIP: your golang is older version $(shell go version)\033[01;m\x1b[0m"
endif
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Run Go Test with Data Race Detection
test: clean
	@$(CMD_MKDIR) -p $(PATH_REPORT)/raw/ $(PATH_REPORT)/doc/
	@$(CMD_ECHO)  -e "\033[1;40;32mRun Go Test.\033[01;m\x1b[0m"
	@GORACE="log_path=$(PATH_REPORT)/doc/$(PATH_RACE_REPORT)" $(CMD_GO) test -v -test.parallel 4 -race -coverprofile=$(PATH_REPORT)/raw/$(PATH_CONVER_PROFILE)
	@$(CMD_ECHO) -e "\033[1;40;36mGenerated a report of data race detection in $(PATH_REPORT)/doc/$(PATH_RACE_REPORT).pid\033[01;m\x1b[0m"
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Send a report of coverage profile to coveralls.io
coveralls::
	@$(CMD_ECHO)  -e "\033[1;40;32mSend a report of coverage profile to coveralls.io.\033[01;m\x1b[0m"
	@$(shell which goveralls) -coverprofile=$(PATH_REPORT)/raw/$(PATH_CONVER_PROFILE) -service=travis-ci
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Generate a report about coverage
cover: test
	@$(CMD_ECHO)  -e "\033[1;40;32mGenerate a report about coverage.\033[01;m\x1b[0m"
	@$(CMD_GO) tool cover -func=$(PATH_CONVER_PROFILE) -o $(PATH_CONVER_PROFILE).txt
	@$(CMD_GO) tool cover -html=$(PATH_CONVER_PROFILE)  -o $(PATH_CONVER_PROFILE).html
	@$(CMD_ECHO) -e "\033[1;40;36mGenerated a report file : $(PATH_CONVER_PROFILE).html\033[01;m\x1b[0m"
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Profiling
pprof: clean
	@$(CMD_MKDIR) -p $(PATH_REPORT)/raw/ $(PATH_REPORT)/doc/
	@$(CMD_ECHO)  -e "\033[1;40;32mGenerate profiles.\033[01;m\x1b[0m"
	@$(CMD_ECHO)  -e "\033[1;40;33mGenerate a CPU profile.\033[01;m\x1b[0m"
	@$(CMD_GO) test -tags unittest -test.parallel 4 -bench . -benchmem -cpuprofile=$(PATH_REPORT)/raw/$(PATH_PROF_CPU)
	@$(CMD_ECHO)  -e "\033[1;40;33mGenerate a Memory profile.\033[01;m\x1b[0m"
	@$(CMD_GO) test -tags unittest -test.parallel 4 -bench . -benchmem -memprofile=$(PATH_REPORT)/raw/$(PATH_PROF_MEM)
	@$(CMD_ECHO)  -e "\033[1;40;33mGenerate a Block profile.\033[01;m\x1b[0m"
	@$(CMD_GO) test -tags unittest -test.parallel 4 -bench . -benchmem -blockprofile=$(PATH_REPORT)/raw/$(PATH_PROF_BLOCK)
	@$(CMD_ECHO)  -e "\033[1;40;33mGenerate a Mutex profile.\033[01;m\x1b[0m"
ifeq ($(GOLANGV18_OVER),1)
	@$(CMD_GO) test -tags unittest -test.parallel 4 -bench . -benchmem -mutexprofile=$(PATH_REPORT)/raw/$(PATH_PROF_MUTEX)
else
	@$(CMD_ECHO) -e "\033[1;40;36mSKIP: your golang is older version $(shell go version)\033[01;m\x1b[0m"
endif
	@$(CMD_MV) -f *.test $(PATH_REPORT)/raw/
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Generate the report for profiling
report: pprof
	@$(CMD_MKDIR) -p $(PATH_REPORT)/raw/ $(PATH_REPORT)/doc/
	@$(CMD_ECHO)  -e "\033[1;40;33mGenerate all report in text format.\033[01;m\x1b[0m"
	@$(CMD_GO) tool pprof -text $(PATH_REPORT)/raw/$(PKG_NAME).test $(PATH_REPORT)/raw/$(PATH_PROF_CPU) > $(PATH_REPORT)/doc/$(PATH_PROF_CPU).txt
	@$(CMD_GO) tool pprof -text $(PATH_REPORT)/raw/$(PKG_NAME).test $(PATH_REPORT)/raw/$(PATH_PROF_MEM) > $(PATH_REPORT)/doc/$(PATH_PROF_MEM).txt
	@$(CMD_GO) tool pprof -text $(PATH_REPORT)/raw/$(PKG_NAME).test $(PATH_REPORT)/raw/$(PATH_PROF_BLOCK) > $(PATH_REPORT)/doc/$(PATH_PROF_BLOCK).txt
ifeq ($(GOLANGV18_OVER),1)
	@$(CMD_GO) tool pprof -text $(PATH_REPORT)/raw/$(PKG_NAME).test $(PATH_REPORT)/raw/$(PATH_PROF_MUTEX) > $(PATH_REPORT)/doc/$(PATH_PROF_MUTEX).txt
endif
	@$(CMD_ECHO)  -e "\033[1;40;33mGenerate all report in pdf format.\033[01;m\x1b[0m"
	@$(CMD_GO) tool pprof -pdf $(PATH_REPORT)/raw/$(PKG_NAME).test $(PATH_REPORT)/raw/$(PATH_PROF_CPU) > $(PATH_REPORT)/doc/$(PATH_PROF_CPU).pdf
	@$(CMD_GO) tool pprof -pdf $(PATH_REPORT)/raw/$(PKG_NAME).test $(PATH_REPORT)/raw/$(PATH_PROF_MEM) > $(PATH_REPORT)/doc/$(PATH_PROF_MEM).pdf
	@$(CMD_GO) tool pprof -pdf $(PATH_REPORT)/raw/$(PKG_NAME).test $(PATH_REPORT)/raw/$(PATH_PROF_BLOCK) > $(PATH_REPORT)/doc/$(PATH_PROF_BLOCK).pdf
ifeq ($(GOLANGV18_OVER),1)
	@$(CMD_GO) tool pprof -pdf $(PATH_REPORT)/raw/$(PKG_NAME).test $(PATH_REPORT)/raw/$(PATH_PROF_MUTEX) > $(PATH_REPORT)/doc/$(PATH_PROF_MUTEX).pdf
endif
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Show Help
help::
	@$(CMD_MAKE2HELP) $(MAKEFILE_LIST)

## Clean-up
clean::
	@$(CMD_ECHO)  -e "\033[1;40;32mClean-up.\033[01;m\x1b[0m"
	@$(CMD_RM) -rfv *.coverprofile *.swp *.core *.html *.prof *.test *.report ./$(PATH_REPORT)/* ./tmp/* *.mdbm *.txt *.log *.out
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

.PHONY: clean cover coveralls help lint pprof report run setup strictlint test build 
