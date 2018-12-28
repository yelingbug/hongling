package utility

import (
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestStringSplit(t *testing.T) {
	s := strings.Split("/home/admin", ",")

	if len(s) == 0 {
		t.Error("结果应该为1.")
	}

	if s[0] != "/home/admin" {
		t.Error("结果应该是/home/admin.")
	}
}

func TestStringSplitToMulti(t *testing.T) {
	s := strings.Split("172.16.0.12:/home/admin,172.16.0.13:/home/yelin.g", ",")

	if len(s) != 2 {
		t.Error("结果应该为2.")
	}

	if s[1] != "172.16.0.13:/home/yelin.g" {
		t.Error("结果应该是172.16.0.13:/home/yelin.g")
	}
}

func TestInCompleteSplit(t *testing.T) {
	s := strings.Split(",123", ",")
	if len(s) != 2 {
		t.Error("结果应该为2,因为有一个逗号.")
	}

	s = strings.Split(",123,", ",")
	if len(s) != 3 {
		t.Error("结果应该为3,因为有两个逗号.")
	}

	s = strings.Split("  ,12 3 , ", ",")
	var vs []string
	for _, v := range s {
		if strings.TrimSpace(v) == "" {
			continue
		}
		vs = append(vs, strings.TrimSpace(v))
	}

	if len(vs) != 1 {
		t.Error("结果应该为1个.")
	}

	if vs[0] != "12 3" {
		t.Error("结果应该为12 3.")
	}
}

func TestParseDirs(t *testing.T) {
	result, err := parseDirs([]string{"172.16.0.199|172.16.0.198:/home/admin/tomcat","172.16.0.197:/home/opt/tomcat"})
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(result, map[string]string{
		"172.16.0.199":"/home/admin/tomcat",
		"172.16.0.198":"/home/admin/tomcat",
		"172.16.0.197":"/home/opt/tomcat",}) {
		t.Error("两个结果应该相等.")
	}

	result, err = parseDirs([]string{"", ""})
	t.Log(err)
	if err == nil {
		t.Error("应该要报错因为len(result) == 0.")
	}

	result, err = parseDirs([]string{"/home/abc/def"})
	t.Log(result)
	if err != nil {
		t.Error("这个调用是正确的，结果会是[*]=/home/abc/def")
	}
	if !reflect.DeepEqual(result, map[string]string{"*": "/home/abc/def"}) {
		t.Error("结果必须相等.")
	}

	result, err = parseDirs([]string{"172.15.s:/a/b/c:/d/e/f", ",,", "18|19:/56"})
	t.Log(err)
	if err == nil {
		t.Error("这个格式是不正确的，必须要报错.")
	}
}

func TestCopyFile(t *testing.T) {
	f1, _ := os.Create("execjava_1")
	f2, _ := os.Create("execjava_2")

	ioutil.WriteFile("execjava_1", []byte("123中国人"), os.ModePerm)
	io.Copy(f2, f1)
	c, _ := ioutil.ReadFile("execjava_2")
	if string(c) != "123中国人" {
		t.Error("拷贝没有按照预想的方式进行.")
	}

	os.RemoveAll("execjava_1")
	os.RemoveAll("execjava_2")

}
