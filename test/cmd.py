#!/usr/bin/python3 -u
import requests, sys, json, getpass

baseurl = 'http://localhost:8088/sa/v1/'

if len(sys.argv) < 3:
    print("Usage:", sys.argv[0], "token command")
    sys.exit(1)

token = sys.argv[1]
command = ' '.join(sys.argv[2:])

headers = {'Authorization': 'Bearer ' + token}

inputvar = {'command':command,
            'cwd':'/var/tmp',
            'envs':['aaa=aaa','bbb=bbb']
           }

reply = requests.post(baseurl + 'cmd/run', json = inputvar, headers=headers)

res = json.loads(reply.text)

if (not 'data' in res) or (res['data'] == None):
    print(res['msg'],"")
    sys.exit(1)

print(res['data']['Output'] + res['data']['Error'], end='')
sys.exit(res['data']['ExitCode'])
