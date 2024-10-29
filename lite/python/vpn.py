import json

import requests

serverIp = "192.168.123.168"
serverPort = "12123"


def get_all():
    resp = requests.post("http://" + serverIp + ":" + serverPort + "/get_all")
    print("get_all", resp.text)
    return resp


def get_all_group():
    resp = requests.post("http://" + serverIp + ":" + serverPort + "/get_all_group")
    print("get_all_group", resp.text)
    return resp


def delete_all():
    resp = requests.post("http://" + serverIp + ":" + serverPort + "/delete_all")
    print("delete_all", resp.text)
    return resp


def add_proxy(path):
    resp = requests.post("http://" + serverIp + ":" + serverPort + "/add_proxy",
                         data=open(path, encoding="utf8").read().encode("utf8"))
    print("add_proxy", resp.text)
    return resp


def start_vpn(id):
    resp = requests.post("http://" + serverIp + ":" + serverPort + "/start_vpn", data=json.dumps({"id": id}))
    print("start_vpn", resp.text)
    return resp


def stop_vpn():
    resp = requests.post("http://" + serverIp + ":" + serverPort + "/stop_vpn")
    print("stop_vpn", resp.text)
    return resp

def reload_vpn(id):
    resp = requests.post("http://" + serverIp + ":12123/reload_vpn", json={"id": id})
    print(resp.text)
    return resp.json()

# add_proxy("./vpn.json")
# get_all_group()
# delete_all()
# start_vpn(37)
# start_vpn(25)
# get_all()
# stop_vpn()
