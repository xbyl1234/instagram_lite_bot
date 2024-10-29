package main

//func TestDB(t *testing.T) {
//	log.InitDefaultLog("", true, false)
//	db.InitMysql()
//	outlook := email.CreateImapEmail(&email.Config{
//		Server: "outlook.office365.com",
//		Port:   "993",
//		RetryConfig: &email.RetryConfig{
//			LastReqTime:          time.Time{},
//			RetryTimeoutDuration: 60 * time.Second,
//			RetryDelayDuration:   3 * time.Second,
//		},
//	}, &email.Account{
//		Username: "mitannobui3@hotmail.com",
//		Password: "WHaERh11",
//	})
//	_ = outlook
//	db.SaveNewEmailDb(outlook, "facebook", "1")
//	db.UpdateEmailProject(outlook.Username, "facebook", "on_used")
//	db.BlackEmail(outlook.Username, "")
//	db.SelectAvailableEmail("facebook", 10)
//}

//func TestEmail(t *testing.T) {
//	log.InitDefaultLog("", true, false)
//	outlook := email.CreateImapEmail(&email.Config{
//		Server: "outlook.office365.com",
//		Port:   "993",
//		RetryConfig: &email.RetryConfig{
//			LastReqTime:          time.Time{},
//			RetryTimeoutDuration: 60 * time.Second,
//			RetryDelayDuration:   3 * time.Second,
//		},
//	}, &email.Account{
//		Username: "thossmosifau@hotmail.com",
//		Password: "ZQA6Pb71",
//	})
//
//	err := outlook.Test()
//	log.Info("%v", err)
//	//code, err := ctrl.WaitForInstagram(outlook.EmailClientData)
//	//code, err := ctrl.WaitForCodeMessenger(outlook.EmailClientData)
//	//outlook.WaitForCode("registration@facebookmail.com",
//	//	outlook.GetAccount().Username,
//	//	func(_email *email.EmailClientData, resp *email.RespEmail) (string, error) {
//	//		log.Info("find")
//	//		return "", nil
//	//	})
//	//log.Info("%v %v", code, err)
//}
//
//func TestJson(t *testing.T) {
//	input := []byte(`{"key1":[{},{"key2":{"key3":[1,2,3]}}]}`)
//	root, err := sonic.Get(input)
//	root.Set("key5", ast.NewBool(false))
//	root.Set("key4", ast.NewBool(false))
//
//	println("%v", root.Get("as212"))
//	buf, err := root.MarshalJSON()
//	println(string(buf))
//	_ = err
//}
//
//func TestHttp(t *testing.T) {
//	// Ja3 is by default HelloChrome_105
//	var bot gostruct.BotData                                    // Creating struct
//	bot.HttpRequest.Request.URL = "https://tls.peet.ws/api/all" // URL
//	bot.HttpRequest.Request.Method = "GET"                      // All method that are supported by fhttp is here supported too
//	bot.HttpRequest.Request.Protocol = "2"                      // 2.0 and 1.1 supported
//	bot.HttpRequest.Request.ReadResponseBody = true             // Used to read the body from the server.
//	bot.HttpRequest.Request.ReadResponseCookies = true          // Used to read the cookies from the server.
//	bot.HttpRequest.Request.ReadResponseHeaders = true          // Used to read the headers from the server.
//	gorequest.HttpRequest(&bot)                                 // Request
//	for k, v := range bot.HttpRequest.Response.Headers {        // Printing the Headers from the server.
//		println(k, ":	", v)
//	}
//	for k, v := range bot.HttpRequest.Response.Cookies { // Printing the Cookies from the server.
//		println(k, ":	", v)
//	}
//	println(bot.HttpRequest.Response.Source) // Printing the source from the server.request
//}
//
//func TestHttpClient(t *testing.T) {
//	jar := tls_client.NewCookieJar()
//	options := []tls_client.HttpClientOption{
//		tls_client.WithTimeoutSeconds(30),
//		tls_client.WithClientProfile(tls_client.Chrome_105),
//		tls_client.WithNotFollowRedirects(),
//		tls_client.WithCookieJar(jar), // create cookieJar instance and pass it as argument
//	}
//
//	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//
//	req, err := http.NewRequest(http.MethodGet, "https://tls.peet.ws/api/all", nil)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//
//	req.Header = http.Header{
//		"accept":          {"*/*"},
//		"accept-language": {"de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7"},
//		"user-agent":      {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36"},
//		http.HeaderOrderKey: {
//			"accept",
//			"accept-language",
//			"user-agent",
//		},
//	}
//
//	resp, err := client.Do(req)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//
//	defer resp.Body.Close()
//
//	log.Println(fmt.Sprintf("status code: %d", resp.StatusCode))
//
//	readBytes, err := io.ReadAll(resp.Body)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//
//	log.Println(string(readBytes))
//}
