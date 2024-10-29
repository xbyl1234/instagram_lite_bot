import requests
import sys

# device xiaomixxx
# app=com.device,com.device2
# net_type wifi/5g
# sim 中国电信/
# country cn

serverIp = "192.168.123.193"
# serverIp = "192.168.123.224"
# serverIp = "192.168.123.189"
# serverIp = "192.168.123.229"
serverPort = "10086"


# honor_cdy-an90_hwcdy-h_29
# honor_hlk-al10_hwhlk-hp_29
# honor_jsn-tl00_hwjsn-h_29
# honor_koz-al00_hnkoz-s_29
# honor_lra-al00_hwlra-h_29
# honor_pct-al10_hwpct_29
# honor_yal-al00_hwyal_29
# huawei_alp-al00_hwalp_29
# huawei_ana-an00_hwana_29
# huawei_clt-al01_hwclt_29
# huawei_oce-an50_hwoce-ml_29
# huawei_pot-al00a_hwpot-hf_29
# huawei_sea-al10_hwsea-a_29
# huawei_spn-al00_hwspn_29
# huawei_stk-al00_hwstk-hf_29
# huawei_tas-al00_hwtas_29
# huawei_tas-an00_hwtas_29
# huawei_vog-al00_hwvog_29
# huawei_wkg-an00_hwwkg-m_29
# oppo_pbem00_pbem00_29
# oppo_pcam00_op46b1_29
# oppo_pcam10_op46f1_29
# oppo_pdym20_op4e21_29
# realme_rmx1901_rmx1901cn_29
# vivo_pd1838_pd1838_29
# vivo_pd1921_pd1921_29
# vivo_pd1962_pd1962_29
# vivo_pd2002_pd2002_29
# vivo_pd2031ea_pd2031ea_29
# vivo_pd2061_pd2061_29
# xiaomi_ginkgo_ginkgo_29

# http://127.0.0.1:10086/new_device?app=com.facebook.orca&enable_gid=false&network=wifi
def new_device(data):
    resp = requests.get("http://" + serverIp + ":" + serverPort + "/new_device", params=data)
    print("new_device", resp.text)
    return resp


def restore_device(data):
    resp = requests.get("http://" + serverIp + ":" + serverPort + "/restore_device", params=data)
    print("restore_device", resp.text)
    return resp


def stop_fake_device():
    resp = requests.get("http://" + serverIp + ":" + serverPort + "/stop_device")
    print("stop_fake_device", resp.text)
    return resp


def get_all_device_id():
    resp = requests.get("http://" + serverIp + ":" + serverPort + "/get_all_device_id")
    print("get_all_device_id", resp.text)
    return resp


def get_all_device():
    resp = requests.get("http://" + serverIp + ":" + serverPort + "/get_all_device")
    print("get_all_device", resp.text)
    return resp


def delete_device(data):
    resp = requests.get("http://" + serverIp + ":" + serverPort + "/delete_device", params=data)
    print("delete_device", resp.text)
    return resp


def cleanup_device(data):
    resp = requests.get("http://" + serverIp + ":" + serverPort + "/cleanup_device", params=data)
    print("cleanup_device", resp.text)
    return resp


def get_current_device():
    resp = requests.get("http://" + serverIp + ":" + serverPort + "/get_current_device")
    print("get_current_device", resp.text)
    return resp


def get_device_config(id):
    resp = requests.get("http://" + serverIp + ":" + serverPort + "/get_device_config", params={"id": id})
    print(resp.text)
    return resp

def get_device_env_config(id):
    resp = requests.get("http://" + serverIp + ":" + serverPort + "/get_device_env_config", params={"id": id})
    print(resp.text)
    return resp


resp = get_all_device()
devices = resp.text.split(",")
devices.remove("")
idx = 1

while True:
    resp = new_device({
        # "app": "com.quark.browser",
        # "app": "com.UCMobile",
        # "app": "com.instagram.android",
        # "app": "com.facebook.katana",
        # "app": "com.facebook.lite",
        # "app": "com.devices1",
        # "app": "com.instagram.android,com.facebook.katana,com.facebook.orca",
        "app": "com.facebook.orca",
        # "enable_boottime": "false",
        # "enable_boottime": "false",
        # "enable_fake_sysinfo": "false",
        # "enable_accessibility": "true",
        # "enable_gid": "true",
        "enable_dpi": "false",
        "enable_global_dpi": "false",
        # "device": "realme_rmx1901_rmx1901cn_29",
        # "device": devices[idx]
        "country": "us",
        "language": "en",
        "network": "wifi"
    })

    # resp = restore_device({
    #     "id": "1691940397_0_7001a693-8236-bd8e-38fc-7f74c020be62",
    # })
    # resp = stop_fake_device()
    # resp = cleanup_device({
    #     "id": "2",
    #     "keep_file_path": "app/com.quark.browser/files/dy_lib,app/com.quark.browser/shared_prefs",
    #     "ignore_file_type": ".so,.jpg",
    #     "ignore_file_size_greater": "1048576"
    # })
    # resp = get_all_device_id()
    # resp = delete_device({
    #     "id": "5"
    # })
    # get_current_device()
    # get_device_config("1694531793_0_65aee05e-f95e-8adb-cb5b-4eafd5a2cc52")
    # get_device_env_config("1691944985_5_0c633309-d2d8-2629-b828-b0c5280e7647")
    idx += 1
    t = input()
