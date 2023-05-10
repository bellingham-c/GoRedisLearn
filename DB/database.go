package DB

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type DatabaseConnection struct {
	Host     string
	Username string
	Password string
	DB       string
	Port     string
	Charset  string
}

var DB *gorm.DB

func InitDB() *gorm.DB {
	config := viper.New()
	//配置文件名（不带扩展名）
	config.SetConfigName("application")
	//在项目中查找配置文件的路径，可以使用相对路径，也可以使用绝对路径
	config.AddConfigPath("./config")
	//多次调用以添加多个搜索路径
	//viper.AddConfigPath("D:/go_project/src/github.com/ourlang/demo/utils")
	//设置文件类型，这里是yaml文件
	config.SetConfigType("yaml")
	//定义用于接收配置文件的变量
	var configData DatabaseConnection
	//查找并读取配置文件
	err := config.ReadInConfig()
	if err != nil { // 处理读取配置文件的错误
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	if err := config.Unmarshal(&configData); err != nil { // 读取配置文件转化成对应的结构体错误
		panic(fmt.Errorf("read config file to struct err: %s \n", err))
	}

	fmt.Println(111111)

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=True&multiStatements=true&loc=Local",
		configData.Username,
		configData.Password,
		configData.Host,
		configData.Port,
		configData.DB,
		configData.Charset,
	)
	fmt.Println("dsn", dsn)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DriverName: "mysql",
		DSN:        dsn, // Data Source Name，参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name
	}), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})

	if err != nil {
		fmt.Println(22222)
		panic(err)
	}
	DB = db
	return DB
}

func GetDB() *gorm.DB {
	return DB
}
