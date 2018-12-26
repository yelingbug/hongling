package utility

import (
	"gopkg.in/urfave/cli.v2"
	"fmt"
	"os/exec"
	"bufio"
	"io"
	"os"
	"strings"
	"sort"
	"sync"
	"github.com/pkg/errors"
)

const (
	whoami = iota
	whoareu
	whoisit
)

var modules = map[string]map[string]interface{}{
	"common": {
		"orderBy": 1,
		"name": "common",
		"relativePath": "common",
		"forFix":       false,
		"priority":     whoami,
	},

	/*"mcommon": {
		"relativePath": "message-center/mc-common",
		"forFix": false,
	},*/
	"mclient": {
		"orderBy": 2,
		"name": "mc-hsf-client",
		"relativePath": "message-center/mc-hsf-client",
		"forFix":       false,
		"priority":     whoareu,
	},
	/*"ucommon": {
		"relativePath": "user-center/uc-common",
		"forFix": false,
	},*/
	"uclient": {
		"orderBy": 3,
		"name": "uc-hsf-client",
		"relativePath": "user-center/uc-hsf-client",
		"forFix":       false,
		"priority":     whoareu,
	},
	/*"uacommon": {
		"relativePath": "user-account-center/uac-common",
		"forFix": false,
	},*/
	"uaclient": {
		"orderBy": 4,
		"name": "uac-hsf-client",
		"relativePath": "user-account-center/uac-hsf-client",
		"forFix":       false,
		"priority":     whoareu,
	},
	"ssoclient": {
		"orderBy": 5,
		"name": "sso-client",
		"relativePath": "single-sign-on/sso-client",
		"forFix":       false,
		"priority":     whoareu,
	},
	"tcbidclient": {
		"orderBy": 6,
		"name": "trans-bidding-hsf-client",
		"relativePath": "transaction/trans-bidding-hsf-client",
		"forFix":       false,
		"priority":     whoareu,
	},
	"tcothersclient": {
		"orderBy": 7,
		"name": "trans-others-client",
		"relativePath": "transaction/trans-others-client",
		"forFix":       false,
		"priority":     whoareu,
	},
	"tcrepayclient": {
		"orderBy": 8,
		"name": "trans-repayment-hsf-client",
		"relativePath": "transaction/trans-repayment-hsf-client",
		"forFix":       false,
		"priority":     whoareu,
	},
	"tctransclient": {
		"orderBy": 9,
		"name": "trans-transfer-hsf-client",
		"relativePath": "transaction/trans-transfer-hsf-client",
		"forFix":       false,
		"priority":     whoareu,
	},
	"ptclient": {
		"orderBy": 10,
		"name": "portal-hsf-client",
		"relativePath": "portal/portal-hsf-client",
		"forFix":       false,
		"priority":     whoareu,
	},
	"yxclient": {
		"orderBy": 11,
		"name": "youxuan-hsf-client",
		"relativePath": "youxuan/youxuan-hsf-client",
		"forFix":       false,
		"priority":     whoareu,
	},
	"tcore": {
		"orderBy": 12,
		"name": "trans-core",
		"relativePath": "transaction/trans-core",
		"forFix":       false,
		"priority":     whoareu,
	},

	"uc": {
		"orderBy": 13,
		"name": "uc-hsf-service",
		"relativePath": "user-center/uc-hsf-service",
		"forFix":       true,
		"priority":     whoisit,
	},
	"uac": {
		"orderBy": 14,
		"name": "uac-hsf-service",
		"relativePath": "user-account-center/uac-hsf-service",
		"forFix":       true,
		"priority":     whoisit,
	},
	"schd": {
		"orderBy": 15,
		"name": "scheduler",
		"relativePath": "scheduler",
		"forFix":       true,
		"priority":     whoisit,
	},
	"tcbid": {
		"orderBy": 16,
		"name": "trans-bidding",
		"relativePath": "transaction/trans-bidding",
		"forFix":       true,
		"priority":     whoisit,
	},
	"tctrans": {
		"orderBy": 17,
		"name": "trans-transfer",
		"relativePath": "transaction/trans-transfer",
		"forFix":       true,
		"priority":     whoisit,
	},
	"tc": {
		"orderBy": 18,
		"name": "trans-bidding-hsf-service",
		"relativePath": "transaction/trans-bidding-hsf-service",
		"forFix":       true,
		"priority":     whoisit,
	},
	"mc": {
		"orderBy": 19,
		"name": "mc-hsf-service",
		"relativePath": "message-center/mc-hsf-service",
		"forFix":       true,
		"priority":     whoisit,
	},
	"pt": {
		"orderBy": 20,
		"name": "portal-hsf-service",
		"relativePath": "portal/portal-hsf-service",
		"forFix":       true,
		"priority":     whoisit,
	},
	"tcrepay": {
		"orderBy": 21,
		"name": "trans-repayment",
		"relativePath": "transaction/trans-repayment",
		"forFix":       true,
		"priority":     whoisit,
	},
}

var ExecjavaCommand = &cli.Command{
	Name:     "execjava",
	Category: "应用修复",
	Aliases:  []string{"ej"},
	Usage:    "hl [global options] execjava/ej [command options] [arguments...]",
	Action:   execJavaUsage,
}

func init() {
	for key, value := range modules {
		if value["forFix"].(bool) {
			ExecjavaCommand.Subcommands = append(ExecjavaCommand.Subcommands, &cli.Command{
				Name:   key,
				Action: execJava,
			})
		}
	}
}

func execJavaUsage(c *cli.Context) error {
	fmt.Println(`使用方法:
  1,写一个public的java类,类名随意,必须包含默认构造函数,只要包含方法签名:
    public void fix() {
      //修复逻辑
    }
  2,类可以定义在任意模块中,支持的模块有` + `,类实例在初始化的时候会被对应模块的spring上下文autowire,所以放心的声明任意依赖的dao/service等...
  3,毫无疑问,类在对应的模块中必须能被编译通过,当然测试必须OK.
  4,搞定之后,要同步到production分支,修复的意义也就在此.
  5,经过我或者开发的确认,去找运维的同事执行如下命令(他们懂得):
    g execjava/ej <模块>..<类的全路径>

  例子:
    在uc服务上写一个类example:

    package com.hongling.abc;
    public class example {
      @Resource
      protected UserDAO userDAO;
      @Resource
      protected AccountService accountService;

      public void fix() {
        System.out.println("准备开工");
        ...
        System.out.println("一共执行了100行.")
        ...
      }
    }
    本地编译测试通过,同步奥production环境,审核确认后,告诉运维在生产环境执行:
    g execjava/ej uc.com.hongling.abc.example
    
    THAT'S IT!
`)
	return nil
}

func execJava(c *cli.Context) error {
	if err := pullProductionBranch(); err != nil {
		Logger.Error("拉取production分之失败:[%s].", err)
		return err
	}

	if err := compile(); err != nil {
		Logger.Error("编译项目失败:[%s].", err)
		return err
	}

	execute(c.Command.Name)

	return nil
}

// 拉取hl.main的production分支.
func pullProductionBranch() error {
	// clone仓库
	execCommand(CacheDir, "git", "clone", "--single-branch", "--branch", "production", "git@172.16.0.100:java/hl.main.git")
	return nil
}

// 编译项目,拉取依赖.
func compile() error {
	// 排序
	var keysOrderBy []int
	keysOrderByMap := make(map[int]string, len(modules))
	for key, value := range modules {
		keysOrderBy = append(keysOrderBy, value["orderBy"].(int))
		keysOrderByMap[value["orderBy"].(int)] = key
	}

	sort.Ints(keysOrderBy)

	// 编译
	Logger.Info("开始前置项目串行编译...")
	for _, item := range keysOrderBy {
		value := modules[keysOrderByMap[item]]
		if !value["forFix"].(bool) {
			Logger.Info(fmt.Sprintf("编译项目[%s]...", value["name"]))
			r, n := value["relativePath"].(string), value["name"].(string)
			if r[0:(len(r) - len(n))] == "" {
				if err := execCommand(CacheDir+"/hl.main/"+value["name"].(string), "mvn", "install"); err != nil {
					return err
				}
			} else {
				if err := execCommand(CacheDir+"/hl.main/"+r[0:(len(r) - len(n))], "mvn", "-am", "-pl", value["name"].(string), "install"); err != nil {
					return nil
				}
			}
		}
	}

	Logger.Info("开始项目并行编译...")
	var wg sync.WaitGroup
	ch := make(chan error)
	for _, item := range keysOrderBy {
		value := modules[keysOrderByMap[item]]
		if value["forFix"].(bool) {
			r, n := value["relativePath"].(string), value["name"].(string)
			Logger.Info(fmt.Sprintf("编译项目[%s]在目录[%s]中...", value["name"], r[0:(len(r) - len(n))]))
			wg.Add(1)
			go func(dir string, v map[string]interface{}) {
				defer wg.Done()
				if dir == "" {
					if err := execCommand(CacheDir+"/hl.main/", "mvn", "install", "-f", CacheDir+"/hl.main/"+v["name"].(string)+"/pom.xml"); err != nil {
						ch <- err
					}
				} else {
					if err := execCommand(CacheDir+"/hl.main/", "mvn", "-am", "-pl", v["name"].(string), "install", "-f", CacheDir+"/hl.main/"+dir+"pom.xml"); err != nil {
						ch <- err
					}
				}
			}(r[0:(len(r) - len(n))], value)
		}
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var errs []string
	for err := range ch {
		errs = append(errs, err.Error())
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}

// 执行修复.
func execute(command string) {
	//ss := modules[command]["relativePath"]
}

func execCommand(toDir string, command string, params ...string) error {
	if err := os.Chdir(toDir); err != nil {
		Logger.Fatal(err)
		return err
	}

	commandAs := exec.Command(command, params...)
	Logger.Info("执行", commandAs.Args)

	out, err := commandAs.StdoutPipe()
	defer out.Close()
	if err != nil {
		Logger.Fatal("执行命令%s时候获取标准输出通道失败.", err)
		return err
	}

	outForErr, err_ := commandAs.StderrPipe()
	defer out.Close()
	if err_ != nil {
		Logger.Fatal("执行命令%s时候获取错误输出通道失败.", err)
		return err_
	}

	commandAs.Start()

	//go func() {
	reader := bufio.NewReader(out)
	for {
		line, err__ := reader.ReadString('\n')
		if err__ != nil || err_ == io.EOF {
			break
		}
		Logger.Info(line)
	}
	//}()

	//go func() {
	readerAsErr := bufio.NewReader(outForErr)
	for {
		line, err__ := readerAsErr.ReadString('\n')
		if err__ != nil || err_ == io.EOF {
			break
		}
		Logger.Info(line)
	}
	//}()

	if err__ := commandAs.Wait(); err__ != nil {
		if strings.HasSuffix(err__.Error(), "not started") {
			Logger.Error(fmt.Sprintf("命令%s不存在.", command))
		}
		return err__
	}

	return nil
}
