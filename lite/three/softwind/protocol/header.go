package protocol

type Header struct {
	ProxyAddress    string
	SNI             string
	GrpcServiceName string
	Cipher          string
	Password        string
	IsClient        bool
}
