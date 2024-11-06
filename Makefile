GOPATH=`pwd`
GITCOMMIT=`git rev-parse HEAD`
BUILDTIME=`date`
RESTD_VERSION=0.1
RESTD_COPYRIGHT=2023

PROJECTS=schedrestd tokengen

sbin_SCRIPTS = schedrestd sched_auth

all: $(PROJECTS) sched_auth

$(PROJECTS): clean
	go env -w GO111MODULE=off
	GOPATH=$(GOPATH) RESTD_VERSION=@RESTD_VERSION@ RESTD_COPYRIGHT=@RESTD_COPYRIGHT@ go build -x -ldflags "-X '$@/version.GITCOMMIT=$(GITCOMMIT)' -X '$@/version.BUILDTIME=$(BUILDTIME)' -X '$@/version.VERSION=$(RESTD_VERSION)' -X '$@/version.RESTD_COPYRIGHT=$(RESTD_COPYRIGHT)'" $@

sched_auth: src/sched_auth/sched_auth.c
	gcc -o sched_auth src/sched_auth/sched_auth.c -g -lpam

clean: $(DELETE_PROJECTS)
	rm -f $(GOPATH)/schedrestd

distclean: clean
	GOPATH=$(GOPATH) go clean -cache

install: schedrestd
	install schedrestd /usr/sbin/schedrestd
	install sched_auth /usr/sbin/sched_auth
	isntall tokengen /usr/bin/tokengen
	install schedrestd.service /lib/systemd/system
	if [ ! -d /etc/schedrestd ];then mkdir /etc/schedrestd;fi
	install schedrestd.yaml /etc/schedrestd/schedrestd.yaml
