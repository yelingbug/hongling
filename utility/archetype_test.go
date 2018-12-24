package utility

import (
	"os"
	"testing"
)

const _FILE = "./test_archetype"

func createTestFile(t *testing.T) {
	if _, err := os.Create(_FILE); err != nil {
		t.Error("创建" + _FILE + "文件 失败.")
	}
}

func createTestDir(t *testing.T) {
	if err := os.Mkdir(_FILE, os.ModePerm); err != nil {
		t.Error("创建" + _FILE + "目录失败.")
	}
}

func removeTestFile(t *testing.T) {
	os.Remove(_FILE)
}

func removeTestDir(t *testing.T) {
	removeTestFile(t)
}

func TestStatWithFile(t *testing.T) {
	createTestFile(t)

	if _, err := os.Stat(_FILE); err != nil {
		t.Error("文件存在时，不会发生error.")
	}

	removeTestFile(t)
}

func TestStatWithDir(t *testing.T) {
	createTestDir(t)

	if _, err := os.Stat(_FILE); err != nil {
		t.Error(err)
	}
	removeTestDir(t)
}
