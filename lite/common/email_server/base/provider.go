package base

const (
	ProviderYX1024Name      = "yx1024"
	ProviderYouXiang555Name = "youxiang555"

	ProviderLinShiName = "linshiyou"
	ProviderTempMailIo = "tempmailio"
)

//SADEfb.y31531
//ddsdew323323s

type Provider interface {
	GetType() string
	GetEmail() (EmailClient, error)
}
