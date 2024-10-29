package proxys

type Provider interface {
	GetProxy(region string, asn string) Proxy
}
