package common

//type ShortLinkCache struct {
//	Link string
//}
//
//type ShortLink struct {
//	client   *http.Client
//	cache    []*ShortLinkCache
//	lock     sync.Mutex
//	maxCache int
//	file     *os.File
//}
//
//func CreateShortLink(maxCache int, savePath string, _proxy *Proxy) *ShortLink {
//	var cache = make([]*ShortLinkCache, 0)
//	readFile, err := os.ReadFile(savePath)
//	if err == nil {
//		lines := strings.Split(string(readFile), "\n")
//		for _, line := range lines {
//			cache = append(cache, &ShortLinkCache{
//				Link: line,
//			})
//		}
//	}
//	log.Info("read short link count %d", len(cache))
//	file, err := os.OpenFile(savePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
//	if err != nil {
//		panic(err)
//	}
//
//	return &ShortLink{
//		client:   CreateGoHttpClient(EnableHttp2(), _proxy.GetProxy()),
//		cache:    cache,
//		maxCache: maxCache,
//		file:     file,
//	}
//}
//
//func (this *ShortLink) getFromCache() string {
//	if len(this.cache) > 0 {
//		ret := this.cache[RandIndex(this.cache)]
//		return ret.Link
//	}
//	return ""
//}
//
//func (this *ShortLink) GetShortLink(_url string) string {
//	this.lock.Lock()
//	defer this.lock.Unlock()
//
//	if len(this.cache) > this.maxCache {
//		return this.getFromCache()
//	}
//	short := this.getShortLink(_url)
//	if short != "" {
//		this.file.WriteString(short + "\n")
//		this.file.Sync()
//		this.cache = append(this.cache, &ShortLinkCache{Link: short})
//		return short
//	}
//	return this.getFromCache()
//}
//
//func (this *ShortLink) getShortLink(_url string) string {
//	do, err := HttpDo(this.client, &RequestOpt{
//		Params: map[string]string{
//			"url":      _url,
//			"shorturl": "",
//			"opt":      "0",
//		},
//		Header: map[string]string{
//			"referer":      "https://is.gd/create.php",
//			"content-type": "application/x-www-form-urlencoded",
//		},
//		IsPost: true,
//		ReqUrl: "https://is.gd/create.php",
//	})
//	if err != nil {
//		return ""
//	}
//	return GetMidString(do, "id=\"short_url\" value=\"", "\" onclick=\"select_text()")
//}
