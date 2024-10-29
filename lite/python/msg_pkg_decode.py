import json
from io import BytesIO
from urllib.parse import unquote
from parse_xml import *
import multipart as mp

from zstd import zip_decode, gzip_decode


def process_path(path):
    return path


def parse_mult_form(data):
    s = data[2:  data.find(b"\r\n")]
    p = mp.MultipartParser(BytesIO(data), s)
    return p.parts()


def prepare_request_params(pkg):
    pkg["params"] = []
    body = pkg["req_body"].decode()
    for item in body.split("&"):
        key = item[:item.find("=")]
        value = item[item.find("=") + 1:]
        pkg["params"].append(KeyValue(key, unquote(value)))


def prepare_params_variables(pkg):
    variables = kvlist_get(pkg["params"], "variables")
    if not variables:
        return

    pkg["variables"] = variables
    jsn = json.loads(variables)

    if jsn.get("params"):
        pkg["has_variables_params"] = True

        if jsn.get("params").get("params"):
            params = jsn["params"]["params"]
            params = params.replace(r"\"", "\r")
            params = params.replace("{params:", "")
            params = params[:len(params) - 2]
            pkg["variables_params_params"] = params
            pkg["has_variables_params_params"] = True
        else:
            pkg["has_variables_params_params"] = False

    else:
        pkg["has_variables_params_params"] = False
        pkg["has_variables_params"] = False


def get_params_template(pkg):
    params = pkg.get("params")
    params_template = {}
    if params is not None:
        for item in params:
            params_template[item.k] = item.v
    return params_template


def get_friendly_name(pkg):
    k = kvlist_get(pkg["params"], "fb_api_req_friendly_name")
    if k:
        return k
    k = kvlist_get(pkg["req_header"], "X-Fb-Friendly-Name")
    if k:
        return k
    return pkg["url"]


def make_api_json(pkg):
    info = {}
    info["host"] = pkg["host"]
    info["path"] = pkg["path"]
    info["method"] = pkg["method"]
    info["is_json_response"] = is_json(pkg["resp_body"])

    body = {}
    md5_keys = []
    body["header_template"] = get_header_template(pkg)
    if pkg.get("variables"):
        body["variables_template"] = pkg["variables"]
        md5_keys.extend(json_txt(json.loads(body["variables_template"])))

    body["is_json_response"] = is_json(pkg["resp_body"])
    body["is_block_setting"] = pkg.get("has_variables_params_params")
    body["head_seq"] = kvlist_keys(pkg["req_header"])
    md5_keys.extend(body["head_seq"])

    if pkg.get("params"):
        body["params_seq"] = kvlist_keys(pkg.get("params"))
        md5_keys.extend(body["params_seq"])
        body["params_template"] = get_params_template(pkg)

    info["body"] = {}
    info["body"][str_md5(str(md5_keys))] = body
    return info


def write_params_variables(name, pkg, api):
    if pkg["has_variables_params_params"]:
        f = open("./data/pkgs/" + name, encoding="utf-8", mode="w")
        try:
            f.write(json.dumps(json.loads(pkg["variables_params_params"]), indent=4, ensure_ascii=False))
        except Exception as e:
            print(pkg["url"], e)


def prepare_graphql_body(pkg):
    prepare_request_params(pkg)
    prepare_params_variables(pkg)
    # write_params_variables(str(index) + "_" + get_friendly_name(pkg), pkg, api)
    # index += 1


def prepare_logging_events(pkg):
    pass
    # mult = parse_mult_form(pkg["req_body"])
    # for item in mult:
    #     if item.filename == 'message':
    #         f = open("./msg_env_改机2/" + str(index) + ".json", "wb")
    #         data = json.dumps(json.loads(gzip_decode(item.raw).decode("utf-8")),
    #                           indent=4,
    #                           ensure_ascii=False,
    #                           sort_keys=False)
    #         f.write(data.encode())


def prepare_form_body(pkg):
    pass


def prepare_body(pkg):
    global index
    url = pkg["url"]
    if "graph.facebook.com/graphql" in url:
        prepare_graphql_body(pkg)
    elif "graph.facebook.com/logging_client_events" in url:
        prepare_logging_events(pkg)
    else:
        prepare_form_body(pkg)


def prepare_pkg(pkgs):
    for pkg in pkgs:
        if pkg["method"] == "POST":
            prepare_body(pkg)


def run_pkg(pkgs):
    apis = {}
    index = 0
    for pkg in pkgs:
        path = pkg["path"]
        if "graph.facebook.com/graphql" in pkg["url"]:
            path = get_friendly_name(pkg)
        if path =='FbBloksAppRootQuery-com.bloks.www.bloks.caa.reg.confirmation.fb.bottomsheet':
            pass
        api = make_api_json(pkg)

        if apis.get(path):
            new_key = list(api.get("body").keys())[0]
            new_value = list(api.get("body").values())[0]
            if apis.get(path).get("body").get(new_key):
                print("dup " + path)
            else:
                apis.get(path).get("body")[new_key] = new_value
        else:
            apis[path] = api

            index += 1
    return apis


noneed = ["xdrig.com", "graph.facebook.com/logging_client_events"]
need = []
pkgs = []
# pkgs.extend(read_pkg("./data/message改机邮箱注册失败2.xml", need, noneed, process_path))
# pkgs.extend(read_pkg("./data/message改机邮箱注册失败.xml", need, noneed, process_path))
# pkgs.extend(read_pkg("./data/message邮箱注册失败.xml", need, noneed, process_path))
pkgs.extend(read_pkg("./data/message邮箱注册失败.xml", need, noneed, process_path))

prepare_pkg(pkgs)

apis = run_pkg(pkgs)
print(json.dumps(apis))
