package android

type MccMncInfo struct {
	Mcc string   `json:"mcc"`
	Mnc string   `json:"mnc"`
	Apn []string `json:"apn"`
}

type SimInfo struct {
	Id                string       `json:"id"`
	MccMnc            []MccMncInfo `json:"mccmnc"`
	CarrierId         int          `json:"carrier_id"`
	SimOperatorNameCn string       `json:"sim_operator_name_cn"`
	SimOperatorNameEn string       `json:"sim_operator_name_en"`
	SimCountryIso     string       `json:"sim_country_iso"`
	AreaCode          string       `json:"area_code"`
	PhoneNumberLength int          `json:"phone_number_length"`
	PhoneNumberPref   []string     `json:"phone_number_pref"`
}

type Country struct {
	Language  []string `json:"language"`
	Country   string   `json:"country"`
	CountryCn string   `json:"country_cn"`
	Timezone  string   `json:"timezone"`
}

type CountryResource struct {
	Sim     map[string][]*SimInfo `json:"sim"`
	Country map[string]*Country   `json:"country"`
}
