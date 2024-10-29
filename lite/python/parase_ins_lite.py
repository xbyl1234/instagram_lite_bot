import json
import os

path = r"F:\desktop\inslite\分析\邮箱注册\pkg_10405_520161963"
data = open(path).read()

data = data.split("\n")

dirPath = r"F:\desktop\inslite\分析\邮箱注册\1"
try:
    os.makedirs(dirPath)
except Exception as e:
    pass

idx = 0
for line in data:
    if not line:
        continue
    j = json.loads(line)
    f = open(dirPath + "/" + str(j["msg_code"]) + "_" + str(idx), "wb")
    f.write(bytes.fromhex(j["data"].replace(" ", "")))
    idx += 1
