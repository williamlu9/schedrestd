## RESTFUL API Daemon for HPC scheduler

Provide RESTFUL interface to run HPC scheduler commands

### Build ###

It requires golang and make.

Run "make" to build the binary schedrestd.

### Deploy ###

```
make install
systemctl enable schedrestd
systemctl start schedrestd
```

This will install schedrestd to /usr/sbin, schedrestd.yaml to /etc/schedrestd,
and install the service
