package RedisUtil

type Conf struct {
	Addr     string
	Password string
	DB       int
}

func GetAddr() string {
	return "ip:6379"
}

func GetPwd() string {
	return "password"
}
func GetDb() int {
	return 0
}
