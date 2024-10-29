package ctrl

import (
	"CentralizedControl/common/android"
	"CentralizedControl/common/log"
	"CentralizedControl/facebook"
	"CentralizedControl/instagram"
	"CentralizedControl/messenger"
	"encoding/json"
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
)

type MessengerCookies struct {
	messenger.Cookies
	Error string
}

func FillMessengerAccount(reqJson *ast.Node) *MessengerCookies {
	var session android.Session
	var account android.Account

	_session, _ := sonic.Marshal(reqJson.GetString("session"))
	json.Unmarshal(_session, &session)

	_account, _ := sonic.Marshal(reqJson.GetString("account"))
	json.Unmarshal(_account, &account)

	return &MessengerCookies{
		Cookies: messenger.Cookies{},
		Error:   reqJson.GetString("error"),
	}
}

func SaveMessengerAccount(reqJson *ast.Node) *MessengerCookies {
	msg := FillMessengerAccount(reqJson)
	saveMessengerAccount(msg)
	return msg
}

func saveMessengerAccount(ck *MessengerCookies) error {
	result := accountDB.Table(MessengerAccountTable).Create(ck)
	if result.Error != nil {
		log.Error("update messenger db error: %v", result.Error)
	}
	return result.Error
}

func saveFacebookAccount(ck *facebook.Cookies) error {
	result := accountDB.Table(FacebookAccountTable).Create(ck)
	if result.Error != nil {
		log.Error("update messenger db error: %v", result.Error)
	}
	return result.Error
}

func saveInstagramAccount(ck *instagram.Cookies) error {
	result := accountDB.Table(InstagramAccountTable).Create(ck)
	if result.Error != nil {
		log.Error("update messenger db error: %v", result.Error)
	}
	return result.Error
}

func UpdateMessenger2FA(ck *messenger.Cookies) error {
	result := accountDB.Table(MessengerAccountTable).Where("email = ?", ck.Email).
		Update("two_fa_key", ck.TwoFAKey).
		Update("qr_code_uri", ck.QrCodeUri).
		Update("had_2fa", ck.Had2FA)
	if result.Error != nil {
		log.Error("update messenger db error: %v", result.Error)
	}
	return result.Error
}

func UpdateFacebook2FA(ck *facebook.Cookies) error {
	result := accountDB.Table(MessengerAccountTable).Where("email = ?", ck.Email).
		Update("two_fa_key", ck.TwoFAKey).
		Update("qr_code_uri", ck.QrCodeUri).
		Update("had_2fa", ck.Had2FA)
	if result.Error != nil {
		log.Error("update messenger db error: %v", result.Error)
	}
	return result.Error
}
