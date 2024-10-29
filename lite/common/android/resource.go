package android

import (
	"CentralizedControl/common/utils"
	"io/ioutil"
	"os"
	"strings"
)

type resource struct {
	ico            []string
	username       []string
	androidDevice  []Device
	deviceResource CountryResource
	allCountry     []string
	asn            map[string][]string
	AppConfig      AppConfig
}

type AppConfig struct {
	GmsVersionStr  string `json:"gms_version_str"`
	GmsVersionCode int    `json:"gms_version_code"`
	//GsfVersionStr       string `json:"gsf_version_str"`
	//GsfVersionCode      int    `json:"gsf_version_code"`
	VendingVersionStr     string `json:"vending_version_str"`
	VendingVersionCode    int    `json:"vending_version_code"`
	DroidguardVersionStr  string `json:"droidguard_version_str"`
	DroidguardVersionCode int    `json:"droidguard_version_code"`
	InsLiteVersionStr     string `json:"ins_lite_version_str"`
	InsLiteVersionCode    int    `json:"ins_lite_version_code"`
	McQueryHashBin        string `json:"mc_query_hash_bin"`
	BloksVersionId        string `json:"bloks_version_id"`
}

var Resource resource

func init() {
	err := loadUsername("./res/username.txt")
	if err != nil {
		panic(err)
	}
	err = loadIco("./res/ico")
	if err != nil {
		panic(err)
	}
	err = loadSimInfo("./devices/sim.json")
	if err != nil {
		panic(err)
	}
	err = loadCountry("./devices/country.json")
	if err != nil {
		panic(err)
	}
	err = loadDevices("./devices/android.json")
	if err != nil {
		panic(err)
	}
	err = loadAsn("./devices/country_asn.json")
	if err != nil {
		panic(err)
	}
	err = loadAppConfig("./devices/app_config.json")
	if err != nil {
		panic(err)
	}
}

func loadUsername(usernamePath string) error {
	data, err := ioutil.ReadFile(usernamePath)
	if err != nil {
		return err
	}
	sp := strings.Split(string(data), "\n")
	for idx := range sp {
		username := sp[idx]
		username = strings.ReplaceAll(username, " ", "")
		username = strings.ReplaceAll(username, "\n", "")
		username = strings.ReplaceAll(username, "\r", "")
		if len(username) > 3 {
			Resource.username = append(Resource.username, username)
		}
	}
	return nil
}

func loadIco(icoPath string) error {
	dir, err := ioutil.ReadDir(icoPath)
	if err != nil {
		return err
	}
	PthSep := string(os.PathSeparator)
	for _, fi := range dir {
		if fi.IsDir() {
			continue
		}
		Resource.ico = append(Resource.ico, icoPath+PthSep+fi.Name())
	}
	return nil
}

func loadDevices(androidDevicePath string) error {
	err := utils.LoadJsonFile(androidDevicePath, &Resource.androidDevice)
	if err != nil {
		return err
	}
	return nil
}

func loadSimInfo(simPath string) error {
	err := utils.LoadJsonFile(simPath, &Resource.deviceResource.Sim)
	if err != nil {
		return err
	}
	return nil
}

func loadCountry(path string) error {
	err := utils.LoadJsonFile(path, &Resource.deviceResource.Country)
	if err != nil {
		return err
	}
	return nil
}

func loadOneDevices(name string, device *Device) {
	utils.LoadJsonFile("./devices/"+name+".json", device)
}

func loadAsn(path string) error {
	err := utils.LoadJsonFile(path, &Resource.asn)
	if err != nil {
		return err
	}
	return nil
}

func loadAppConfig(path string) error {
	err := utils.LoadJsonFile(path, &Resource.AppConfig)
	if err != nil {
		return err
	}
	return nil
}

//func GenAsn(country string) string {
//	return utils.ChoseOne(Resource.asn[country])
//}

func GetAreaCode(country string) string {
	return Resource.deviceResource.Sim[country][0].AreaCode
}
