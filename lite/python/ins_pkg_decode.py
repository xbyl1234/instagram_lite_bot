import json
import re
import urllib
from urllib.parse import unquote
from parse_xml import *


def match_template(string, template):
    pattern = template.replace(r'%d', r'\d+')
    pattern = pattern.replace(r'%s', r'\w+')
    match = re.match(pattern, string)
    if match:
        return True
    else:
        return False


def process_path(path):
    paths = [
        "/api/v1/highlights/%d/highlights_tray/",
        "/api/v1/music/profile/%d/",
        "/api/v1/feed/user/%d/story/",
        "/api/v1/feed/user/%d/",
        "/api/v1/users/%d/info/",
        "/api/v1/fundraiser/%d/standalone_fundraiser_info/",
        "/api/v1/fundraiser/%d/standalone_fundraiser_info/",
        "/api/v1/highlights/%d/highlights_tray/",
        "/api/v1/feed/user/%d/",
        "/api/v1/users/%d/info/",
        "/api/v1/highlights/%d/highlights_tray/",
        "/api/v1/users/%d/info/",
        "/api/v1/friendships/create/%d/",
        "/api/v1/nametag/nametag_lookup_by_name/%s/"
    ]
    for item in paths:
        if match_template(path, item):
            return (item.replace("%d", "%v")
                    .replace("%s", "%v"))
    return path


def make_api_json(pkg):
    info = {}
    info["host"] = pkg["host"]
    info["path"] = pkg["path"]
    info["method"] = pkg["method"]
    info["is_json_response"] = is_json(pkg["resp_body"])

    bodys = {}
    md5_keys = []
    if pkg.get("query"):
        bodys["query"] = pkg["query"]
    if pkg.get("query_seq"):
        bodys["query_seq"] = pkg["query_seq"]
        md5_keys.extend(bodys["query_seq"])

    bodys["head_seq"] = kvlist_keys(pkg["req_header"])
    bodys["header_template"] = get_header_template(pkg)
    md5_keys.extend(bodys["head_seq"])

    if pkg.get("body_type"):
        bodys["body_type"] = pkg["body_type"]

    if pkg.get("form_seq"):
        bodys["form_seq"] = pkg["form_seq"]
        bodys["form_template"] = pkg["form_template"]
        md5_keys.extend(bodys["form_seq"])

    if pkg.get("json_template"):
        bodys["json_template"] = pkg["json_template"]
        md5_keys.extend(json_txt(json.loads(bodys["json_template"])))

    bodys["had_nav_chain"] = kvlist_get(pkg["req_header"], "X-Ig-Nav-Chain") is not None
    info["body"] = {}
    info["body"][str_md5(str(md5_keys))] = bodys
    return info


# signed_body=SIGNATURE.{"bool_opt_policy":"0","mobileconfigsessionless":"","api_version":"3","unit_type":"1","query_hash":"e1faa64a4a2408ba55531b85db97d0a6664f9dfa3a579dd56e946ed57849db75","device_id":"9ff0fad8-c663-47cc-93cc-8da3c06caf7a","fetch_type":"ASYNC_FULL","family_device_id":"9FFF5C83-7F4B-409A-9424-2C18C9667ED4"}
def prepare_signed_body(pkg, body):
    body = urllib.parse.unquote(body)
    body = body.replace("signed_body=SIGNATURE.", "")
    pkg["json_template"] = body
    pkg["body_type"] = "signed_body"


def prepare_form_body(pkg, body):
    form = decode_query(body)
    pkg["form_seq"] = kvlist_keys(form)
    pkg["form_template"] = kvlist2map(form)
    pkg["body_type"] = "form_body"


def prepare_json_body(pkg, body):
    body = urllib.parse.unquote(body)
    pkg["json_template"] = body
    pkg["body_type"] = "json_body"


def prepare_body(pkg):
    body = pkg["req_body"].decode()
    if body.startswith("signed_body="):
        prepare_signed_body(pkg, body)
    else:
        if is_json(unquote(body)):
            prepare_json_body(pkg, body)
        else:
            prepare_form_body(pkg, body)


def prepare_pkg(pkgs):
    for pkg in pkgs:
        if pkg["method"] == "POST":
            prepare_body(pkg)


def run_pkg(pkgs):
    apis = {}
    index = 0
    for pkg in pkgs:
        if "i.instagram.com" == pkg["host"]:
            path = pkg["path"]
            api = make_api_json(pkg)

            if apis.get(path):
                new_key = list(api.get("body").keys())[0]
                new_value = list(api.get("body").values())[0]
                if apis.get(path).get("body").get(new_key):
                    print("dup " + path)
                else:
                    apis.get(path).get("body")[new_key] = new_value
            else:
                apis[pkg["path"]] = api

            index += 1
    return apis


def analyse_heade(pkgs):
    keyset = {}
    for pkg in pkgs:
        if "i.instagram.com" == pkg["host"]:
            keyset[kvlist_md5(pkg["req_header"])] = kvlist_keys(pkg["req_header"])
            # for header in pkg["req_header"]:
            #     keyset[header.k] = header.v

    print(json.dumps(keyset))


noneed = []
need = ["i.instagram.com"]
# need = ["https://graph.facebook.com/graphql"]
pkgs = []
# pkgs.extend(read_pkg("./data/2023-08-30-facebook跳Instagram登录.xml", need, noneed))
# pkgs.extend(read_pkg("./data/ins1.xml", need, noneed, process_path))
# pkgs.extend(read_pkg("./data/扫描二维码关注.xml", need, noneed, process_path))
pkgs.extend(read_pkg("./data/开启app.xml", need, noneed, process_path))

prepare_pkg(pkgs)
# analyse_heade(pkgs)

apis = run_pkg(pkgs)
print(json.dumps(apis))

# def write_pkgs_xls(pkgs):
#     f = open("./data/pkgs.xls", "w")
#
#     def write_value(pkg, k):
#         if pkg.get(k) is not None:
#             f.write(pkg[k])
#         f.write("\t")
#
#     def write_listkey(pkg, k):
#         if pkg.get(k) is not None:
#             f.write(kvlist2str(pkg[k]))
#         f.write("\t")
#
#     for pkg in pkgs:
#         try:
#             write_value(pkg, "url")
#             write_value(pkg, "fb_api_req_friendly_name")
#             # f.write(kvlist_get(pkg["req_header"], "x-fb-privacy-context"))
#             # write_value(pkg, "params_key_md5")
#             f.write(kvlist_get(pkg["params"], "client_trace_id"))
#             # write_listkey(pkg, "params")
#             # write_value(pkg, "variables_params")
#             # write_value(pkg, "variables")
#         except Exception as e:
#             print("write_pkgs_xls", pkg["url"], e)
#         finally:
#             f.write("\r")
# print(pkgs)
