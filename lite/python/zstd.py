import gzip
import io
import json
import zlib

import zstandard
import pathlib


def zst_decode(data):
    return zstandard.ZstdDecompressor().decompress(data)


def zip_decode(data):
    return zlib.decompress(data)


def gzip_decode(data):
    compressed_file = io.BytesIO(data)
    return gzip.GzipFile(fileobj=compressed_file).read()

# data = open(path, "rb").read()
# print(zstd(r"F:\desktop\CentralizedControl\python\msg_env\0"))
