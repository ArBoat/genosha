package config

import (
  "flag"
  "fmt"
  "github.com/spf13/viper"
  "log"
)

var (
  V *viper.Viper
  env = flag.String("env", "dev", "environment")
)

func init()  {
  flag.Parse()
  V = viper.New()
  V.SetConfigName(*env)
  V.SetConfigType("yaml")
  V.AddConfigPath("./configure")
  err := V.ReadInConfig()
  if err != nil {                // 读取配置信息失败
    log.Println(fmt.Errorf("Fatal error config file: %s \n", err))
  }
}
