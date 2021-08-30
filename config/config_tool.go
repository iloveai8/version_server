package tool
 
import (
	"bufio"
	"encoding/json"
	"os"
)


type Config struct {
	RedisConfig RedisConfig `json:"redis_config"`
}

 
//Redis属性定义
type RedisConfig struct {
	Addr string `json:"addr"`
	Port string `json:"port"`
	Db int `json:"db"`
}
 
//方法或者变量首字母大写, 表示对外可见
func GetConfig() *Config  {
	return cfg
}
 
//全局变量 对外不可见
var cfg * Config = nil
 
func ParseConfig(path string)(*Config, error)  {
	file, err := os.Open(path)
	defer file.Close()
 
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)
	if  err = decoder.Decode(&cfg); err != nil{
		return nil, err
	}
	return cfg, nil
}