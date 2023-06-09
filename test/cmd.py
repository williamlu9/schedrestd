#!/usr/bin/python3 -u
import requests, sys, json, getpass

baseurl = 'http://localhost:8088/sa/v1/'

if len(sys.argv) < 3:
    print("Usage:", sys.argv[0], "username command")
    sys.exit(1)

username = sys.argv[1]
command = sys.argv[2]

try:
    password = getpass.getpass()
except Exception as error:
    print('ERROR', error)
    sys.exit(1)

reply = requests.post(baseurl + 'login', json = {'username': username, 'password':password})

res = json.loads(reply.text)

if res["data"] == None:
    print(res["msg"],"")
    sys.exit(1)

token = res["data"]["token"]["token"]
user = res["data"]["token"]["userName"]

headers = {'Authorization': 'Bearer ' + token}

inputvar = {'command':command,
            'envs':["aaa=aaa"],
            'cwd':'/var/tmp'
           }

reply = requests.post(baseurl + 'cmd/run', json = inputvar, headers=headers)

res = json.loads(reply.text)

if res['data'] == None:
    print(res['msg'],"")
    sys.exit(1)

print(res['data']['Output'], end='')
