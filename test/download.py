#!/usr/bin/python3 -u
import requests, sys

baseurl = 'http://localhost:8088/sa/v1/'

if len(sys.argv) < 3:
    print("Usage:", sys.argv[0], "token remote_file_path [source_dir]")
    sys.exit(1)

token = sys.argv[1]

headers = {'Authorization': 'Bearer ' + token}

cmd = 'file/download/' + sys.argv[2]
if len(sys.argv) == 4:
    cmd = cmd + '?dir=' + sys.argv[3]

reply = requests.get(baseurl + cmd + sys.argv[2], headers=headers)

open(sys.argv[2], 'wb').write(reply.content)
