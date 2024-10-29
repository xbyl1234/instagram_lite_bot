import time

import requests

serverIp = "192.168.123.209"
serverPort = "11187"


def devices():
    resp = requests.post("http://" + serverIp + ":" + serverPort + "/devices")
    print(resp.status_code)
    print(resp.text)
    return resp


devices()