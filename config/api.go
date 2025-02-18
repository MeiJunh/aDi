package config

// GetSConfDsn 返回静态配置dsn
func GetSConfDsn() string {
	if SConf != nil {
		return SConf.StaticDBDsn
	}
	return ""
}
