import json
import urllib
import xml.dom.minidom
import base64
import hashlib
from urllib.parse import unquote


def _json_txt(json_key, dic_json):
    if isinstance(dic_json, dict):
        for key in dic_json:
            if isinstance(dic_json[key], dict):
                _json_txt(json_key, dic_json[key])
                json_key.append(key)
            else:
                json_key.append(key)


def json_txt(dic_json):
    json_key = []
    _json_txt(json_key, dic_json)
    return json_key


def is_json(data):
    try:
        json.loads(data)
        return True
    except Exception as e:
        return False


def str_md5(s):
    m = hashlib.md5()
    m.update(s.encode())
    return m.hexdigest()


def get_xml(file_path):
    burp = xml.dom.minidom.parse(file_path)
    root = burp.documentElement
    items = root.getElementsByTagName('item')
    return items


class KeyValue:
    def __init__(self, k, v):
        self.k = k
        self.v = v


def kvlist_md5(l):
    s = ""
    for head in l:
        s += head.k + ","
    m = hashlib.md5()
    m.update(s.encode())
    return m.hexdigest()


def kvlist_2str(l):
    s = ""
    for head in l:
        s += head.k + ","
    return s


def kvlist_get(l, k):
    for item in l:
        if item.k == k:
            return item.v
    return None


def kvlist_keys(l):
    r = []
    for item in l:
        r.append(item.k)
    return r


def kvlist2map(l):
    result = {}
    for item in l:
        if item.k in result.keys():
            print("waring ", item.k, " in map!")
        result[item.k] = item.v
    return result


def get_http_body(data):
    p = data.find(b"\r\n\r\n")
    return data[p + 4:]


def get_http_header(data):
    result = {"req_header": []}
    p = data.find(b"\r\n\r\n")
    header = data[:p]
    for line in header.split(b"\r\n"):
        if line == "":
            continue
        if line.startswith(b"POST"):
            result["method"] = "POST"
        elif line.startswith(b"GET"):
            result["method"] = "GET"
        else:
            header_key = line[:line.find(b":")].decode()
            if header_key == "host":
                continue
            h = KeyValue(header_key, line[line.find(b":") + 1:].decode())
            if h.v.startswith(" "):
                h.v = h.v[1:]
            result["req_header"].append(h)
    return result


def decode_query(s):
    result = []
    for item in s.split("&"):
        key = item[:item.find("=")]
        value = item[item.find("=") + 1:]
        result.append(KeyValue(key, unquote(value)))
    return result

def get_header_template(pkg):
    header_template = {}
    header = pkg["req_header"]
    for item in header:
        header_template[item.k] = item.v
    return header_template

def read_pkg(path, need, noneed, path_process):
    def is_in_url(url_list, url):
        for item in url_list:
            if item in url:
                return True
        return False

    pkgs = []
    xmlData = get_xml(path)
    for item in xmlData:
        try:
            pkg = {}
            pkg["url"] = item.getElementsByTagName("url")[0].firstChild.data
            if len(need) > 0 and not is_in_url(need, pkg["url"]):
                continue
            if is_in_url(noneed, pkg["url"]):
                continue
            url = urllib.parse.urlparse(pkg["url"])
            pkg["host"] = url.hostname
            pkg["path"] = path_process(url.path)
            if url.query:
                query = decode_query(url.query)
                pkg["query"] = kvlist2map(query)
                pkg["query_seq"] = kvlist_keys(query)

            req_raw = base64.b64decode(item.getElementsByTagName("request")[0].firstChild.data)
            resp_raw = base64.b64decode(item.getElementsByTagName("response")[0].firstChild.data)

            header_result = get_http_header(req_raw)
            pkg["method"] = header_result["method"]
            pkg["req_header"] = header_result["req_header"]
            pkg["req_body"] = get_http_body(req_raw)
            pkg["resp_body"] = get_http_body(resp_raw)

            pkgs.append(pkg)
        except Exception as e:
            print("read pkg error", e)

    return pkgs
