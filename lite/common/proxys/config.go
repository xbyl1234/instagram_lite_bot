package proxys

var DebugHttpProxy, _ = CreateHttpProxy("http://127.0.0.1:8080")
var DebugS5Proxy, _ = CreateSocks5Proxy("http://127.0.0.1:8889")

var LocalSsrSocksProxy, _ = CreateSocks5Proxy("socks://127.0.0.1:10808")
