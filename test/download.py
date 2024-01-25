#!/usr/bin/python3 -u
import requests, sys

baseurl = 'http://localhost:8088/sa/v1/'

if len(sys.argv) < 3:
    print("Usage:", sys.argv[0], "token remote_file_path ")
    sys.exit(1)

token = sys.argv[1]

headers = {'Authorization': 'Bearer ' + token}

reply = requests.get(baseurl + 'file/download/' + sys.argv[2], headers=headers)

print (reply.text)
