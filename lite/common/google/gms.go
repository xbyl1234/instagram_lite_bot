package google

import (
	"CentralizedControl/common/http_helper"
	"errors"
	"net/http"
)

type GmsClient struct {
	host string
	h    *http.Client
}

type TokenDetails struct {
	Auth   string `json:"auth"`
	Expiry int    `json:"expiry"`
	Scopes string `json:"scopes"`
}

type AccountToken struct {
	AuthAccount  string       `json:"authAccount"`
	Error        string       `json:"Error"`
	TokenDetails TokenDetails `json:"tokenDetails"`
	AccountType  string       `json:"accountType"`
	AuthToken    string       `json:"authtoken"`
}

func (this *GmsClient) GetAccountToken(account, scope, pkg_name string) *AccountToken {
	var resp = &AccountToken{}
	err := http_helper.HttpDoJson(this.h, &http_helper.RequestOpt{
		IsPost: true,
		ReqUrl: this.host + "/get_account_token",
		JsonData: map[string]string{
			"account":  account,
			"scope":    "audience:server:client_id:" + scope,
			"pkg_name": pkg_name,
		},
	}, resp)
	if err != nil {
		return nil
	}
	return resp
}

type GetAccountResp struct {
	Error    string   `json:"error"`
	Accounts []string `json:"accounts"`
}

func (this *GmsClient) GetAccount() (string, error) {
	var resp = &GetAccountResp{}
	err := http_helper.HttpDoJson(this.h, &http_helper.RequestOpt{
		IsPost: true,
		ReqUrl: this.host + "/get_accounts",
	}, resp)
	if err != nil || resp.Accounts == nil {
		return "", err
	}
	return resp.Accounts[0], err
}

type SafetynetAttestResp struct {
	Error     string `json:"error"`
	JwsResult string `json:"jws_result"`
}

func (this *GmsClient) SafetyNetAttest(apiKey, nonce string) (string, error) {
	var resp = &SafetynetAttestResp{}
	err := http_helper.HttpDoJson(this.h, &http_helper.RequestOpt{
		IsPost: true,
		ReqUrl: this.host + "/safetynet_attest",
		JsonData: map[string]string{
			"api_key": apiKey,
			"nonce":   nonce,
		},
	}, resp)
	if err != nil {
		return "", err
	}
	if resp.Error != "ok" {
		return "", errors.New(resp.Error)
	}
	return resp.JwsResult, nil
}
