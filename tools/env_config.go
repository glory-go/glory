package tools

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/glory-go/glory/common"

	"github.com/rs/xid"
	//"github.com/glory-go/glory/log" 不能用这个
)

// GetEnv 获取从环境变量读出的值
func GetEnv(envStr string) string {
	return os.Getenv(envStr)
}

// ReadFromEnvIfNeed 传入rawConfig为指向配置结构的指针
// 保证了在配置文件中可以混用真实配置值和环境变量值，要求环境变量与真实配置数据不冲突
/*例如：
"master-mysql":
   config_source: env # 环境变量读入
   host: MYSQL_USER_HOST
   port: 3306
   username: MYSQL_USER_USERNAME
   password: MYSQL_USER_PASSWORD
   dbname: MYSQL_DBNAME
*/
func ReadFromEnvIfNeed(rawConfig interface{}) error {
	rawVal := reflect.ValueOf(rawConfig)
	if rawVal.Kind() != reflect.Ptr { //保证是指针
		fmt.Println("ReadFromEnvIfNeed param rawConfig should be pointer to struct")
		return errors.New("input rawConfig kind invalid")
	}
	val := rawVal.Elem()
	typ := reflect.TypeOf(rawConfig).Elem()
	kd := val.Kind() //获取到a对应的类别
	if kd != reflect.Struct {
		fmt.Println("ReadAllConfFromEnv expect struct but got ", kd)
		return errors.New("ReadAllConfFromEnv expect struct")
	}

	num := val.NumField()
	for i := 0; i < num; i++ {
		tagVal := typ.Field(i).Tag.Get("yaml")
		if tagVal == "config_source" {
			val, ok := val.Field(i).Interface().(string)
			if !ok {
				//不符合要求，直接跳过
				fmt.Println("yaml tag 'config_source' type should be string")
				continue
			}
			if val == "env" {
				//符合环境变量读入要求
				readAllConfFromEnv(rawConfig)
				return nil
			}
		}
	}
	return nil
}


// readAllConfFromEnv 提供从环境变量拉取的服务
// 对于传入的rawConfig结构所有yaml标签不为config_source 的字段，尝试从环境变量中拉取，如果拉取到则替换原有值
func readAllConfFromEnv(rawConfig interface{}) {
	val := reflect.ValueOf(rawConfig) //获取reflect.Type类型
	typ := reflect.TypeOf(rawConfig)
	//获取到该结构体有几个字段
	num := val.Elem().NumField()
	//遍历结构体的所有字段
	for i := 0; i < num; i++ {
		//fmt.Printf("Field %d:值=%v\n", i, val.Elem().Field(i))
		tagVal := typ.Elem().Field(i).Tag.Get("yaml")
		//如果该字段有tag标签就显示，否则就不显示
		if envValue := GetEnv(val.Elem().Field(i).String()); tagVal != "config_source" && envValue != "" {
			// 非config字段，并且环境变量里面有定义
			val.Elem().Field(i).SetString(envValue)
		}
	}
}

func SetTimeClickFunction(t time.Duration, f func()) {
	ticker := time.NewTicker(t)
	for {
		<-ticker.C
		f()
	}
}

func PrometheusParseToSupportMetricsName(input string) string {
	input = strings.Replace(input, "-", "_", -1)
	input = strings.Replace(input, ".", "_", -1)
	input = strings.Replace(input, "$", "_", -1)
	return input
}

func Addr2AddrLabel(localAddress common.Address) string {
	str := localAddress.GetUrl()
	str = strings.Replace(str, ":", "-", -1)
	str = strings.Replace(str, ".", "_", -1)
	return str
}

func AddrLabel2Addr(addrLabel string) common.Address {
	addrLabel = strings.Replace(addrLabel, "_", ".", -1)
	res := strings.Split(addrLabel, "-")
	addr := common.Address{}
	if len(res) != 2 {
		fmt.Println("error: AddrLabel2Ip error!")
		return addr
	}

	addr.Host = res[0]
	var err error
	addr.Port, err = strconv.Atoi(res[1])
	if err != nil {
		fmt.Println("error: AddrLabel2Ip error!")
	}
	return addr
}

// GenerateXID generate unique id
func GenerateXID() string {
	return xid.New().String()
}
