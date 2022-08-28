package config

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestParseYAML(t *testing.T) {
	err := Init("conf.yml") // 解析yaml文件
	if err != nil {
		t.Fatal(err)
	}

	c := Show()
	data, err := json.MarshalIndent(&c, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(data))
}

// 测试更新配置文件
func TestWatch(t *testing.T) {
	err := Init("conf.yml")
	if err != nil {
		t.Error(err)
		return
	}

	fs := []func(){
		func() {
			fmt.Println("更新字段1")
		},
		func() {
			fmt.Println("更新字段2")
		},
	}
	WatchConfig(fs...)

	for i := 0; i < 30; i++ {
		fmt.Println("port:", Show().ServicePort)
		time.Sleep(time.Second)
	}
}
