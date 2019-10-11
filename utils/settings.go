package utils

type qiniuCfg struct {
	accessKey string
	secretKey string
	bucket    string
	zone      string

	domain string
}

func QiniuCfg() qiniuCfg {
	cfg := qiniuCfg{
		accessKey: "rfHlsZIoIAw69VVfW0vgVwtG876MkXd_Skr_kNtE",
		secretKey: "WxvytzvTSPbBULfA6GCxuH1CWzvekykclDxJ9MAQ",
		bucket:    "tccnu-ups",
		zone:      "ZoneHuanan",
		domain:    "https://tccnu-ups.bigdata.starclink.com",
	}

	return cfg
}
