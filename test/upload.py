#!/usr/bin/python3 -u
import requests, sys, json

baseurl = 'http://localhost:8088/sa/v1/'

if len(sys.argv) < 3:
    print("Usage:", sys.argv[0], "token local_file_path ")
    sys.exit(1)

token = sys.argv[1]

headers = {'Authorization': 'Bearer ' + token}
files = {'file': open(sys.argv[2], 'rb')}

reply = requests.post(baseurl + 'file/upload', files=files, headers=headers)

res = json.loads(reply.text)

if (not 'data' in res) or (res['data'] == None):
    print(res["msg"],"")
    sys.exit(1)

# print token
print(res["data"]['file']['path'])
