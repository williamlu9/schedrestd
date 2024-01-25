#!/usr/bin/python3
import requests, getpass, json, sys

baseurl = 'http://localhost:8088/sa/v1/'

if len(sys.argv) < 2:
    print("Usage:", sys.argv[0], "username [token_valid_in_minute]")
    sys.exit(1)

duration = 120
if len(sys.argv) == 3:
    duration = int(sys.argv[2])

username = sys.argv[1]

try:
    # obtain user password
    password= getpass.getpass(prompt="Password: ")
except Exception as error:
    print('ERROR', error)
    sys.exit(1)

# call API to generate a token that is valid for 120 minutes
reply = requests.post(baseurl + 'login', json = {'username':username, 'password':password, 'duration': duration})
res = json.loads(reply.text)
if (not 'data' in res) or (res['data'] == None):
    print(res["msg"],"")
    sys.exit(1)

# print token
print(res["data"]["token"]["token"])
