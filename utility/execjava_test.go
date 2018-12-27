package utility

import (
	"testing"
	"strings"
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
	_, err := parseDirs([]string{"172.16.0.199|172.16.0.198:/home/admin/tomcat","172.16.0.197:/home/opt/tomcat"})
	if err != nil {
		t.Error(err)
	}

}
