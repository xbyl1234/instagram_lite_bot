import json


class Devices:
    def __init__(self, line):
        self.json = json.loads(line)

    def getDeviceName(self):
        name = (self.getPropNoThrow('ro.product.brand') + "_" +
                self.getPropNoThrow('ro.product.name') + "_" +
                self.getPropNoThrow("ro.product.device") + "_" +
                self.getPropNoThrow('ro.build.version.sdk')).lower()
        return name

    def getDeviceType(self):
        name2 = self.getDeviceName()
        if "xiaomi" in name2:
            return 0
        if "huawei" in name2:
            return 1
        if "honor" in name2:
            return 2
        if "oppo" in name2:
            return 3
        if "vivo" in name2:
            return 4
        if "realme" in name2:
            return 5
        if "samsung" in name2:
            return 6
        if "redmi" in name2:
            return 7
        if "sony" in name2:
            return 9
        if "asus" in name2:
            return 10
        print("unknow device:" + name2)
        return -1

    def getProp(self, key):
        prop_info = self.json["prop_info"]
        return prop_info[key]

    def getPropNoThrow(self, key):
        try:
            prop_info = self.json["prop_info"]
            return prop_info[key]
        except Exception as e:
            return ""

    def getFeature(self):
        features = self.json["features"]["systemAvailableFeatures"]
        return features

    def getShareLib(self):
        libs = self.json["features"]["systemSharedLibraryNames"]
        return libs

    def getConfiguration(self):
        configuration = self.json["configuration"]
        return {
            "colorMode": configuration["colorMode"],
            "densityDpi": configuration["densityDpi"],
            "fontScale": configuration["fontScale"],
            "fontWeightAdjustment": configuration["fontWeightAdjustment"],
            "hardKeyboardHidden": configuration["hardKeyboardHidden"],
            "keyboard": configuration["keyboard"],
            "keyboardHidden": configuration["keyboardHidden"],
            "navigation": configuration["navigation"],
            "navigationHidden": configuration["navigationHidden"],
            "orientation": configuration["orientation"],
            "screenHeightDp": configuration["screenHeightDp"],
            "screenLayout": configuration["screenLayout"],
            "screenWidthDp": configuration["screenWidthDp"],
            # "semMobileKeyboardCovered": configuration["semMobileKeyboardCovered"],
            # "smallestScreenWidthDp": configuration["smallestScreenWidthDp"],
            "touchscreen": configuration["touchscreen"],
            "uiMode": configuration["uiMode"],
        }

    def getConfigurationInfo(self):
        configuration_info = self.json["configuration_info"]
        return {
            "reqTouchScreen": configuration_info["reqTouchScreen"],
            "reqKeyboardType": configuration_info["reqKeyboardType"],
            "reqNavigation": configuration_info["reqNavigation"],
            "reqInputFeatures": configuration_info["reqInputFeatures"],
            "reqGlEsVersion": configuration_info["reqGlEsVersion"],
        }

    def getCpus(self):
        cpu = self.json["cpu"]
        result = []
        cpuCount = self.json["cpu_count"]

        def parseInt(s):
            return int(s.replace("\n", ""), 10)

        for i in range(0, cpuCount):
            cpuName = "cpu" + str(i)
            item = {
                "min_frequency": parseInt(cpu["min"][cpuName]),
                "max_frequency": parseInt(cpu["max"][cpuName]),
            }
            result.append(item)
        return result

    def getActivityMemory(self):
        activity_memory = self.json["activity_memory"]
        return activity_memory

    def getPkgs(self):
        return {
            "tts_pkg_name": "com.google.android.tts",
            "launcher_pkg_name": self.json["resolve_package_no_ref"]["launcher"],
            "home_pkg_name": self.json["resolve_package_no_ref"]["home"]
        }

    def getEnv(self):
        http_agent = "Dalvik/" + "2.1.0" + " (Linux; U; Android " + self.getProp(
            "ro.build.version.release_or_codename")
        if "REL" == self.getProp("ro.build.version.codename"):
            http_agent += "; "
            http_agent += self.getProp("ro.product.model")
        http_agent += " Build/"
        http_agent += self.getProp("ro.build.id")
        http_agent += ")"
        return {
            "user_agent": self.json["user_agent"],
            "http_agent": http_agent,
            "assets_locales": self.json["locale"]["assetsLocales"],
            "available_locales": self.json["locale"]["availableLocales"]
        }

    def getPhone(self):
        return {
            "phone_type": self.json["extract"]["phone"]["phone"]["PhoneType"]
        }

    def getWindowSize(self):
        return {
            "density_dpi": self.json["screen_density_dpi"],
            "real_width": self.json["screen_stable_display_size"]["x"],
            "real_height": self.json["screen_stable_display_size"]["y"],
            "app_width": self.json["display_info"]["getWidth"],
            "app_height": self.json["display_info"]["getHeight"],
            "status_bar": 0,
            "navigation_bar": 0,
            "refresh_rate": self.json["display_info"]["getRefreshRate"],
            "scale": self.json["screen_scaled_density"],
            "xdpi": self.json["screen_xdpi"],
            "ydpi": self.json["screen_ydpi"]
        }

    def getRuntimeMemory(self):
        return {
            "max_memory": self.json["runtime_memory"]["maxMemory"],
            "total_memory": self.json["runtime_memory"]["totalMemory"],
            "free_memory": self.json["runtime_memory"]["freeMemory"],
        }

    def getOpenGl(self):
        return {
            "egl_extensions": self.json["other"]["opengl"]
        }

    def getMemory(self):
        return {
            "memory_class": int(self.json["runtime_memory"]["maxMemory"] / 1024 / 1024),
            "activity_memory": self.getActivityMemory(),
            "runtime_memory": self.getRuntimeMemory(),
        }

    def getBattery(self):
        return {
            "scale": self.json["battery"]["battery.scale"]
        }

    def toInsJson(self):
        return {
            "props": {
                "ro.build.version.codename": self.getProp("ro.build.version.codename"),
                "ro.build.version.incremental": self.getProp("ro.build.version.incremental"),
                "ro.build.version.sdk": self.getProp("ro.build.version.sdk"),
                "ro.build.version.security_patch": self.getProp("ro.build.version.security_patch"),
                "ro.product.first_api_level": self.getProp("ro.product.first_api_level"),
                "ro.build.date.utc": self.getProp("ro.build.date.utc"),
                "ro.product.cpu.abilist": self.getProp("ro.product.cpu.abilist"),
                "ro.product.cpu.abilist32": self.getProp("ro.product.cpu.abilist32"),
                "ro.product.cpu.abilist64": self.getProp("ro.product.cpu.abilist64"),
                "ro.product.board": self.getProp("ro.product.board"),
                "ro.product.model": self.getProp("ro.product.model"),
                "ro.product.brand": self.getProp("ro.product.brand"),
                "ro.product.manufacturer": self.getProp("ro.product.manufacturer"),
                "ro.build.version.release": self.getProp("ro.build.version.release"),
                "ro.build.product": self.getProp("ro.build.product"),
                "ro.build.id": self.getProp("ro.build.id"),
                "ro.hardware": self.getProp("ro.hardware"),
                "ro.product.device": self.getProp("ro.product.device"),
                "ro.boot.hardware": self.getProp("ro.boot.hardware"),
                "ro.mediatek.platform": self.getPropNoThrow("ro.mediatek.platform"),
                "ro.board.platform": self.getProp("ro.board.platform"),
                "ro.chipname": self.getPropNoThrow("ro.chipname"),
            },
            "device_name": self.getDeviceName(),
            "device_type": self.getDeviceType(),
            "features": self.getFeature(),
            "shared_library": self.getShareLib(),
            "configuration": self.getConfiguration(),
            "configuration_info": self.getConfigurationInfo(),
            "cpus": self.getCpus(),
            "memory": self.getMemory(),
            "opengl": self.getOpenGl(),
            "pkg_infos": self.getPkgs(),
            "env": self.getEnv(),
            "phone": self.getPhone(),
            "battery": self.getBattery(),
            "window_size": self.getWindowSize(),
            "sdk_int": int(self.getProp("ro.build.version.sdk")),
        }


devices = []
data = open("./raw_devices.txt", encoding="utf8").read()
sp = data.split("\n")
for line in sp:
    if not sp:
        continue
    try:
        devices.append(Devices(line))
    except Exception as e:
        print(e)

dumps = []
for d in devices:
    try:
        dumps.append(d.toInsJson())
    except Exception as e:
        print(e)

f = open("../android.json", "w")
f.write(json.dumps(dumps))
