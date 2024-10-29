import json
from time import sleep

import requests

project = "messenger"


# project = "facebook_lite"
# project = "instagram_lite"

def GetEmail():
    resp = requests.post("http://127.0.0.1:5588/get_email", json={
        "project": project,
        "provider": "yx1024"
    })
    print(resp.text)
    return resp.json()


def GetCode(email):
    resp = requests.post("http://127.0.0.1:5588/get_code", json={
        "project": project,
        "email": email
    })
    print(resp.text)
    return resp.json()


# , passwd:
e = GetEmail()["email"]
while True:
    resp = GetCode(e)
    if resp.get("code") != None:
        print(resp["code"])
        break
    sleep(2)
# GetCode()
