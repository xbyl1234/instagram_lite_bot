package common

import (
	tls "github.com/refraction-networking/utls"
	"strconv"
	"strings"
)

type GenExtensionFunc = func(gen *TlsConfigGen) tls.TLSExtension

type TlsConfigGen struct {
	TlsVersion      []string
	CipherSuites    []uint16
	Curves          []tls.CurveID
	SupportedPoints []uint8
	ExtensionsCode  []string
	ExtensionConfig ExtensionConfig
}

type ExtensionConfig struct {
	SupportedSignatureAlgorithms_13 []tls.SignatureScheme
	Renegotiation_65281             tls.RenegotiationSupport
	KeyShare_51                     []tls.KeyShare
	Modes_45                        []uint8
	SupportedVersions_43            []uint16
}

var extensionsMap map[string]GenExtensionFunc

func GetInsLiteExtensionConfig() ExtensionConfig {
	return ExtensionConfig{
		SupportedSignatureAlgorithms_13: []tls.SignatureScheme{
			tls.ECDSAWithP256AndSHA256,
			tls.PSSWithSHA256,
			tls.PKCS1WithSHA256,
			tls.ECDSAWithP384AndSHA384,
			tls.PSSWithSHA384,
			tls.PKCS1WithSHA384,
			tls.PSSWithSHA512,
			tls.PKCS1WithSHA512,
			tls.PKCS1WithSHA1,
		},
		Renegotiation_65281: tls.RenegotiateNever,
		KeyShare_51: []tls.KeyShare{
			{Group: tls.X25519},
		},
		Modes_45: []uint8{tls.PskModeDHE},
		SupportedVersions_43: []uint16{
			tls.VersionTLS13,
			tls.VersionTLS12,
			tls.VersionTLS11,
			tls.VersionTLS10,
		},
	}
}

func init() {
	extensionsMap = make(map[string]GenExtensionFunc)
	extensionsMap["0"] = func(gen *TlsConfigGen) tls.TLSExtension {
		return &tls.SNIExtension{}
	}
	extensionsMap["23"] = func(gen *TlsConfigGen) tls.TLSExtension {
		return &tls.ExtendedMasterSecretExtension{}
	}
	extensionsMap["65281"] = func(gen *TlsConfigGen) tls.TLSExtension {
		return &tls.RenegotiationInfoExtension{Renegotiation: gen.ExtensionConfig.Renegotiation_65281}
	}
	extensionsMap["10"] = func(gen *TlsConfigGen) tls.TLSExtension {
		return &tls.SupportedCurvesExtension{Curves: gen.Curves}
	}
	extensionsMap["11"] = func(gen *TlsConfigGen) tls.TLSExtension {
		return &tls.SupportedPointsExtension{SupportedPoints: gen.SupportedPoints}
	}
	extensionsMap["35"] = func(gen *TlsConfigGen) tls.TLSExtension {
		return &tls.SessionTicketExtension{}
	}
	extensionsMap["5"] = func(gen *TlsConfigGen) tls.TLSExtension {
		return &tls.StatusRequestExtension{}
	}
	extensionsMap["13"] = func(gen *TlsConfigGen) tls.TLSExtension {
		return &tls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: gen.ExtensionConfig.SupportedSignatureAlgorithms_13}
	}
	extensionsMap["51"] = func(gen *TlsConfigGen) tls.TLSExtension {
		return &tls.KeyShareExtension{KeyShares: gen.ExtensionConfig.KeyShare_51}
	}
	extensionsMap["45"] = func(gen *TlsConfigGen) tls.TLSExtension {
		return &tls.PSKKeyExchangeModesExtension{Modes: gen.ExtensionConfig.Modes_45}
	}
	extensionsMap["43"] = func(gen *TlsConfigGen) tls.TLSExtension {
		return &tls.SupportedVersionsExtension{
			Versions: gen.ExtensionConfig.SupportedVersions_43,
		}
	}
	extensionsMap["21"] = func(gen *TlsConfigGen) tls.TLSExtension {
		return &tls.UtlsPaddingExtension{GetPaddingLen: tls.BoringPaddingStyle}
	}
}

func GenTlsConfig(ja3 string, extensionsConfig ExtensionConfig) *tls.ClientHelloSpec {
	var gen = createTlsConfigGen(ja3, extensionsConfig)
	return gen.gen()
}

func atoi(s string) int64 {
	parseInt, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return parseInt
}

func createTlsConfigGen(ja3 string, extensionsConfig ExtensionConfig) *TlsConfigGen {
	sp := strings.Split(ja3, ",")
	_tlsVersion := strings.Split(sp[1], "-")
	_ciphers := strings.Split(sp[1], "-")
	_rawExtensions := strings.Split(sp[2], "-")
	_curves := strings.Split(sp[3], "-")
	_pointFormats := strings.Split(sp[4], "-")

	config := &TlsConfigGen{}
	for _, item := range _curves {
		config.Curves = append(config.Curves, tls.CurveID(uint16(atoi(item))))
	}
	for _, item := range _pointFormats {
		config.SupportedPoints = append(config.SupportedPoints, uint8(atoi(item)))
	}
	for _, c := range _ciphers {
		config.CipherSuites = append(config.CipherSuites, uint16(atoi(c)))
	}
	config.ExtensionsCode = _rawExtensions
	config.TlsVersion = _tlsVersion
	config.ExtensionConfig = extensionsConfig
	return config
}

func MAX(a []uint16) uint16 {
	result := a[0]
	for _, v := range a {
		if v > result {
			result = v
		}
	}
	return result
}

func MIN(a []uint16) uint16 {
	result := a[0]
	for _, v := range a {
		if v < result {
			result = v
		}
	}
	return result
}

func (this *TlsConfigGen) gen() *tls.ClientHelloSpec {
	var Extensions = []tls.TLSExtension{}
	for _, item := range this.ExtensionsCode {
		Extensions = append(Extensions, extensionsMap[item](this))
	}
	var config = &tls.ClientHelloSpec{
		CipherSuites:       this.CipherSuites,
		Extensions:         Extensions,
		CompressionMethods: nil,
		TLSVersMin:         MIN(this.ExtensionConfig.SupportedVersions_43),
		TLSVersMax:         MAX(this.ExtensionConfig.SupportedVersions_43),
		GetSessionID:       nil,
	}
	return config
}
