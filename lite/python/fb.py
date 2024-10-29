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


SubParamsTypeParamsNotRealJsonStr = open("./fb/SubParamsTypeParamsNotRealJsonStr", mode="w", encoding="utf-8")
SubParamsTypeParamsJsonStr = open("./fb/SubParamsTypeParamsJsonStr", mode="w", encoding="utf-8")
SubParamsTypeParamsParamsJson = open("./fb/SubParamsTypeParamsParamsJson", mode="w", encoding="utf-8")
SubParamsTypeParamsNo1Params = open("./fb/SubParamsTypeParamsNo1Params", mode="w", encoding="utf-8")
SubParamsTypeParamsNo2Params = open("./fb/SubParamsTypeParamsNo2Params", mode="w", encoding="utf-8")


def prepare_params_variables(pkg):
    variables = pkg["form_template"].get("variables")
    if not variables:
        print("no variables")
        return

    pkg["variables"] = variables
    jsn = json.loads(variables)

    if jsn.get("params"):
        if jsn.get("params").get("params"):
            params = jsn["params"]["params"]
            if params.startswith("{params:"):
                params = params.replace("{params:", "")
                params = params[:len(params) - 2]
                SubParamsTypeParamsNotRealJsonStr.write(pkg["friend_name"] + "\n")
                SubParamsTypeParamsNotRealJsonStr.write(json.dumps(json.loads(pkg["variables"])) + "\n")
                SubParamsTypeParamsNotRealJsonStr.write(params + "\n")
            else:
                if isinstance(params, str):
                    paramsJson = json.loads(params)
                    if paramsJson.get("params"):
                        SubParamsTypeParamsParamsJson.write(pkg["friend_name"] + "\n")
                        SubParamsTypeParamsParamsJson.write(json.dumps(json.loads(pkg["variables"])) + "\n")
                        SubParamsTypeParamsParamsJson.write(paramsJson.get("params") + "\n")
                        pass
                    else:
                        SubParamsTypeParamsJsonStr.write(pkg["friend_name"] + "\n")
                        SubParamsTypeParamsJsonStr.write(json.dumps(json.loads(pkg["variables"])) + "\n")
                        SubParamsTypeParamsJsonStr.write(json.dumps(paramsJson) + "\n")
                        pass
                    pass
                else:
                    # print(json.dumps(json.loads(pkg["variables"])))
                    # print(params)
                    pass
            # pkg["variables_params_params"] = params
            # print(json.dumps(json.loads(pkg["variables_params_params"])))
        else:
            SubParamsTypeParamsNo1Params.write(pkg["friend_name"] + "\n")
            SubParamsTypeParamsNo1Params.write(json.dumps(json.loads(pkg["variables"])) + "\n")
            SubParamsTypeParamsNo1Params.write(json.dumps(jsn.get("params")) + "\n")
    else:
        SubParamsTypeParamsNo2Params.write(pkg["friend_name"] + "\n")
        SubParamsTypeParamsNo2Params.write(json.dumps(json.loads(pkg["variables"])) + "\n")
        SubParamsTypeParamsNo2Params.write(json.dumps(jsn) + "\n")


def get_params_template(pkg):
    params = pkg.get("form_template")
    form_template = {}
    if params is not None:
        for item in params:
            form_template[item.k] = item.v
    return form_template


def get_friendly_name(pkg):
    k = pkg["form_template"].get("fb_api_req_friendly_name")
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
    body["head_seq"] = kvlist_keys(pkg["req_header"])
    md5_keys.extend(body["head_seq"])

    if pkg.get("variables"):
        md5_keys.extend(json_txt(json.loads(pkg["variables"])))

    if pkg.get("form_template"):
        body["form_seq"] = pkg["form_seq"]
        md5_keys.extend(body["form_seq"])
        body["form_template"] = pkg["form_template"]

    body["body_type"] = pkg["body_type"]

    if pkg.get("query"):
        body["query"] = pkg["query"]
        body["query_seq"] = pkg["query_seq"]

    info["body"] = {}
    info["body"][str_md5(str(md5_keys))] = body
    return info


def write_params_variables(name, pkg, api):
    if pkg.get("variables_params_params"):
        f = open("./data/pkgs/" + name, encoding="utf-8", mode="w")
        try:
            f.write(json.dumps(json.loads(pkg["variables_params_params"]), indent=4, ensure_ascii=False))
        except Exception as e:
            print(pkg["url"], e)


def prepare_form_body(pkg):
    body = pkg["req_body"].decode()
    form = decode_query(body)
    pkg["form_seq"] = kvlist_keys(form)
    pkg["form_template"] = kvlist2map(form)
    pkg["body_type"] = "form_body"


def prepare_graphql_body(pkg):
    body = pkg["req_body"].decode()
    form = decode_query(body)
    pkg["form_seq"] = kvlist_keys(form)
    pkg["form_template"] = kvlist2map(form)
    # write_params_variables(str(index) + "_" + get_friendly_name(pkg), pkg, api)
    # index += 1
    pkg["body_type"] = "graphql_body"


def prepare_logging_events(pkg):
    pkg["body_type"] = "events_body"
    # mult = parse_mult_form(pkg["req_body"])
    # for item in mult:
    #     if item.filename == 'message':
    #         f = open("./msg_env_改机2/" + str(index) + ".json", "wb")
    #         data = json.dumps(json.loads(gzip_decode(item.raw).decode("utf-8")),
    #                           indent=4,
    #                           ensure_ascii=False,
    #                           sort_keys=False)
    #         f.write(data.encode())


def prepare_body(pkg):
    global index
    url = pkg["url"]
    if "graph.facebook.com/graphql" in url:
        prepare_graphql_body(pkg)
    elif "graph.facebook.com/logging_client_events" in url:
        prepare_logging_events(pkg)
    else:
        prepare_form_body(pkg)

    if "graph.facebook.com/graphql" in url:
        pkg["friend_name"] = get_friendly_name(pkg)
        prepare_params_variables(pkg)


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
        if path == 'FbBloksAppRootQuery-com.bloks.www.bloks.caa.reg.confirmation.fb.bottomsheet':
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
need = ["facebook.com"]
# need = ["zero_headers_ping_params_v2"]
pkgs = []
# pkgs.extend(read_pkg("./data/message改机邮箱注册失败2.xml", need, noneed, process_path))
# pkgs.extend(read_pkg("./data/message改机邮箱注册失败.xml", need, noneed, process_path))
# pkgs.extend(read_pkg("./data/message邮箱注册失败.xml", need, noneed, process_path))
pkgs.extend(read_pkg("./data/2023-09-18-facebook邮箱注册成功.xml", need, noneed, process_path))

prepare_pkg(pkgs)

apis = run_pkg(pkgs)
print(json.dumps(apis))
