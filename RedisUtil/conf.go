package RedisUtil

type Conf struct {
	Addr     string
	Password string
	DB       int
}

func GetAddr() string {
	return "192.168.192.129:6379"
}

func GetPwd() string {
	return "caojinbo"
}
func GetDb() int {
	return 0
}
