BINARY_NAME=schedrestd
GOFLAGS=-v
MODULE_NAME=schedrestd

all: schedrestd sched_auth

init:
	@if [ ! -f go.mod ]; then \
		go mod init $(MODULE_NAME); \
		go mod tidy; \
	fi

schedrestd: init
	go build $(GOFLAGS) -o $(BINARY_NAME) .

sched_auth: src/sched_auth/sched_auth.c
	gcc -o sched_auth src/sched_auth/sched_auth.c -g -lpam

clean:
	go clean
	rm -f $(BINARY_NAME) sched_auth

distclean: clean
	rm -f go.mod go.sum
	GOPATH=$(GOPATH) go clean -cache

install: schedrestd
	install schedrestd /usr/sbin/schedrestd
	install sched_auth /usr/sbin/sched_auth
	install schedrestd.service /lib/systemd/system
	if [ ! -d /etc/schedrestd ];then mkdir /etc/schedrestd;fi
	install schedrestd.yaml /etc/schedrestd/schedrestd.yaml
