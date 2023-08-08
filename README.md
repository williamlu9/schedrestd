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

### API Document ###

After the schedrestd service is running, in a browser, access the document at
http://IP_ADDR:8088/sa/v1/swagger/index.html

### Example ###

The example code to use the API is: test/cmd.py

### Configuration ###

By default, the servide listens to the port of 8088 on HTTP. The file /etc/scheddrestd/schedrestd.yaml allows changing some configurations like adding SSL certificates, changing port, and etc. Refer to the file comment area for the description of each parameter.

After modifying the configuration file, restart the service is required.

```
systemctl restart schedrestd
```
