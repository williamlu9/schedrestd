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

### Configuration ###

By default, the servide listens to the port of 8088 on HTTP. The file /etc/scheddrestd/schedrestd.yaml allows changing some configurations like adding SSL certificates, changing port, and etc. Refer to the file comment area for the description of each parameter.

After modifying the configuration file, restart the service is required.

```
systemctl restart schedrestd
```
#### Configuration parameters in /etc/schedrestd/schedrest.yaml ####

- ssl: 1 = Enable, 0 = Disable (default)
- http_port: Speficying listening port, if ssl = 0. The default value is 8088
- https_port: Specifying listening port in HTTPS, if ssl = 1. The default value is 8043.
- cert: SSL cert file path. The default is: /opt/cert/server.crt
- key: SSL cert key path. The default is: /opt/cert/server.key
- timeout: The token valid duration in minute. The default value is 30. User can change this value upon individual login.
- log_level: Log level in /var/log/schedrestd.log.hostname. The default is "info".
- web_url_path: Any prefix path in front of /sa/v1

### Example ###

The example code is in the directory of test.

1. Login to generate a token
```
#!/usr/bin/python3
import requests, getpass, json, sys

baseurl = 'http://localhost:8088/sa/v1/'

if len(sys.argv) < 2:
    print("Usage:", sys.argv[0], "username")
    sys.exit(1)

username = sys.argv[1]

try:
    # obtain user password
    password = getpass.getpass(prompt="Password: ")
except Exception as error:
    print('ERROR', error)
    sys.exit(1)

# call API to generate a token that is valid for 120 minutes
reply = requests.post(baseurl + 'login', json = {'username':username, 'password':password, 'duration':120})
res = json.loads(reply.text)
if not 'data' in res:
    print(res["msg"],"")
    sys.exit(1)

# print token
print(res["data"]["token"]["token"])
```

2. Run a command
```
#!/usr/bin/python3
import request, json, sys

baseurl = 'http://localhost:8088/sa/v1/'

if len(sys.argv) < 3:
    print("Usage:", sys.argv[0], "token", "command ...")
    sys.exit(1)

token = sys.argv[1]
command = ' '.join(sys.argv[2:])

headers = {'Authorization': 'Bearer ' + token}

# specify command, current working directory for the command to run (optional),
# and environment variables (optional):
inputvar = {'command':command,
            'cwd':'/var/tmp',
            'envs':['aaa=aaa',
                    'bbb=bbb']
           }

reply = requests.post(baseurl + 'cmd/run', json = inputvar, headers=headers)

res = json.loads(reply.text)

if res['data'] == None:
    print(res['msg'], "")
    sys.exit(1)

print(res['data']['Output'] + res['data']['Error'], end='')
```
