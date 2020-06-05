package confs

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func InitConfig(path string) map[string]string {
	myMap := make(map[string]string)

	f, err := os.Open(path)

	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		s := strings.TrimSpace(string(b))
		//fmt.Println(s)

		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}
		key := strings.TrimSpace(s[:index])
		if len(key) == 0 {
			continue
		}
		value := strings.TrimSpace(s[index+1:])
		if len(value) == 0 {
			continue
		}
		myMap[key] = value
	}
	return myMap
}

func GetEnv() string {
	var config string
	env := os.Getenv("ENV_CLUSTER")
	if env == "local" {
		config = "./configure/config.conf"
	} else if env == "prod" {
		config = "./configure_docker/config_prod.conf"
	} else if env == "test" {
		config = "./configure_docker/config_test.conf"
	} else {
		config = "../configure/config.conf"
	}
	return config
}

var ConfigMap = InitConfig(GetEnv())
