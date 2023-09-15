package syserrno

const (
	OK = "0"
	// 通用系統錯誤
	Undefined = "9999"

	HTTP           = "10"
	ValidParameter = "11"
	RMQ            = "13"

	// storage
	Mongo = "20"
	MySQL = "21"
	Redis = "22"
	AWSS3 = "23"
)
