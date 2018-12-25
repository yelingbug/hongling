package utility

import (
	"gopkg.in/urfave/cli.v2"
	"fmt"
)

var modules = map[string]map[string]interface{}{
	"uc": {
		"relativePath":"user-center/uc-hsf-service",
	},
	"uac": {
		"relativePath":"user-account-center/uac-hsf-service",
	},
	"schd": {
		"relativePath":"scheduler",
	},
	"tcbid": {
		"relativePath":"transaction/trans-bidding",
	},
	"tctrans": {
		"relativePath":"transaction/trans-transfer",
	},
	"tcore": {
		"relativePath":"transaction/trans-bidding-hsf-service",
	},
	"mc": {
		"relativePath":"message-center/mc-hsf-service",
	},
	"pt": {
		"relativePath":"portal/portal-hsf-service",
	},
	"tcrepay": {
		"relativePath":"transaction/trans-repayment/",
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
	for key := range modules {
		ExecjavaCommand.Subcommands = append(ExecjavaCommand.Subcommands, &cli.Command{
			Name:key,
			Action:execJava,
		})
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

	exec(c.Command.Name)

	return nil
}

// 拉取hl.main的production分支.
func pullProductionBranch() error  {
	return nil
}

// 编译项目,拉取依赖.
func compile() error {
	return nil
}

// 执行修复.
func exec(command string) {
	ss := modules[command]["relativePath"]
}