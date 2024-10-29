package email_server

//const (
//	EmailDBRegSuccess  = "reg_success"
//	EmailDBRegBanEmail = "reg_email_ban"
//
//	EmailDBLose    = "lose"
//	EmailOnUsed    = "on_used"
//	EmailDBRelease = "release"
//)
//
//type EmailAccount struct {
//	Username      string `gorm:"primaryKey"`
//	Password      string
//	EmailProvider string
//	UnAvailable   string            `gorm:"column:unavailable"`
//	CreateTime    common.CustomTime `gorm:"autoCreateTime"`
//	UsedTime      common.CustomTime `gorm:"autoCreateTime"`
//	Instagram     string
//	Facebook      string
//}
//
//type EmailCache struct {
//	cache     map[string]map[string]base.EmailClient // project->email
//	cacheLock sync.Mutex
//	DB        *gorm.DB
//}

//func cleanProject(project string) string {
//	if strings.Contains(project, "_lite") {
//		project = strings.ReplaceAll(project, "_lite", "")
//	}
//	if project == "messenger" {
//		project = "facebook"
//	}
//	return project
//}
//
//func (this *EmailCache) SaveNewEmailDb(e base.EmailClient, project string, provider string) error {
//	project = cleanProject(project)
//	data := &EmailAccount{
//		Username:      e.GetAccount().Username,
//		Password:      e.GetAccount().Password,
//		EmailProvider: provider,
//		UnAvailable:   "",
//		Instagram:     "",
//		Facebook:      "",
//	}
//	switch project {
//	case "facebook":
//		data.Facebook = EmailOnUsed
//	case "instagram":
//		data.Instagram = EmailOnUsed
//	}
//	return this.DB.Table("email").Create(data).Error
//}
//
//func (this *EmailCache) BlackEmail(email string, err string) error {
//	return this.DB.Exec("UPDATE email SET unavailable = ? WHERE username = ?", err, email).Error
//}
//
//func (this *EmailCache) UpdateEmailProject(e string, project string, opt string) error {
//	project = cleanProject(project)
//	return this.DB.Exec("UPDATE email SET "+project+" = ? WHERE username = ?", opt, e).Error
//}
//
//func (this *EmailCache) SelectAvailableEmail(project string, maxCount int) []*EmailAccount {
//	var e []*EmailAccount
//	project = cleanProject(project)
//
//	tx := this.DB.Begin()
//	now := time.Now()
//
//	timeout := now.Add(-10 * time.Minute)
//	oneHourAgo := now.Add(-time.Hour * 2)
//	tx.Debug().Table("email").
//		Where("create_time BETWEEN ? AND ? and ("+project+"= ? or (used_time < ? and "+project+"= ?)"+" )and (unavailable is null or unavailable ='')",
//			common.CustomTime(oneHourAgo), common.CustomTime(now),
//			EmailDBRelease,
//			common.CustomTime(timeout),
//			EmailOnUsed).
//		Order("create_time DESC").
//		Limit(maxCount).
//		Find(&e)
//	if len(e) == 0 {
//		tx.Commit()
//		return nil
//	}
//
//	for _, item := range e {
//		log.Info("db cache email: %s", item.Username)
//		err := tx.Exec("UPDATE email SET "+project+" = ? ,used_time= ? WHERE username = ?", EmailOnUsed, common.CustomTime(now), item.Username).Error
//		if err != nil {
//			tx.Rollback()
//			return nil
//		}
//	}
//
//	tx.Commit()
//	return e
//}
//
//func (this *EmailCache) GetEmailByProject(project string) (base.EmailClient, error) {
//	this.cacheLock.Lock()
//	defer this.cacheLock.Unlock()
//
//	pjCache := this.cache[project]
//	if pjCache == nil {
//		pjCache = make(map[string]base.EmailClient)
//		this.cache[project] = pjCache
//	}
//	if len(pjCache) > 0 {
//		var key string
//		var value base.EmailClient
//		for k, v := range pjCache {
//			key = k
//			value = v
//			break
//		}
//		delete(pjCache, key)
//		log.Info("dbcache get email: %s, passwd: %s", value.GetAccount().Username, value.GetAccount().Password)
//		return value, nil
//	}
//	go this.RefreshDbCache(project)
//	return nil, nil
//}
//
//func (this *EmailCache) RefreshDbCache(project string) {
//	log.Info("start RefreshDbCache")
//	cacheSize := 10
//	accounts := this.SelectAvailableEmail(project, cacheSize)
//	testChan := make(chan *EmailAccount, cacheSize/2)
//	wait := sync.WaitGroup{}
//	wait.Add(cacheSize / 2)
//
//	this.cacheLock.Lock()
//	pjCache := this.cache[project]
//	for idx, item := range accounts {
//		if pjCache[item.Username] != nil {
//			accounts[idx] = nil
//		}
//	}
//	this.cacheLock.Unlock()
//
//	for i := 0; i < 10; i++ {
//		go func() {
//			defer wait.Done()
//			for account := range testChan {
//				imap := base.CreateImapEmail(DefaultOutlookConfig, &base.Account{
//					Username: account.Username,
//					Password: account.Password,
//				})
//				err := imap.Test()
//				if err != nil {
//					this.BlackEmail(account.Username, err.Error())
//					continue
//				}
//
//				this.cacheLock.Lock()
//				pjCache := this.cache[project]
//				pjCache[account.Username] = imap
//				this.cacheLock.Unlock()
//			}
//		}()
//	}
//
//	for _, account := range accounts {
//		if account != nil {
//			testChan <- account
//		}
//	}
//	wait.Wait()
//}
//
//func CreateDbCacheEmailProvider(DB *gorm.DB) *EmailCache {
//	result := &EmailCache{
//		DB:    DB,
//		cache: make(map[string]map[string]base.EmailClient),
//	}
//	err := result.DB.Table("email").AutoMigrate(&EmailAccount{})
//	if err != nil {
//		panic(err)
//	}
//	return result
//}
