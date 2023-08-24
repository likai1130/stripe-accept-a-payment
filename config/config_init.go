package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

// Setup 载入配置文件
func Setup(configFile string) {

	v := viper.New()
	//自动获取全部的env加入到viper中。默认别名和环境变量名一致。（如果环境变量多就全部加进来）
	v.AutomaticEnv()

	//替换读取格式。默认a.b.c.d格式读取env，改为a_b_c_d格式读取
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 本地配置文件位置
	v.SetConfigFile(configFile)

	//支持 "yaml", "yml", "json", "toml", "hcl", "tfvars", "ini", "properties", "props", "prop", "dotenv", "env":
	v.SetConfigType("yaml")

	//读配置文件到viper配置中
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// 系列化成config对象
	if err = v.Unmarshal(&Cfg); err != nil {
		panic(err)
	}
}
