package ctrl

import "CentralizedControl/common/email_server/base"

var (
	MessengerAccountTable = "messenger_cookies"
	FacebookAccountTable  = "facebook_cookies"
	InstagramAccountTable = "instagram_cookies"
)

var (
	ConfigProxyTypeSubscribe = 0
	ConfigProxyTypeLuminati  = 1
	ConfigProxyTypeOxylabs   = 2
	ConfigProxyTypeZenoo     = 2

	ConfigProxyUse = ConfigProxyTypeSubscribe
)

var (
	ConfigEmailUse = base.ProviderYX1024Name
)
