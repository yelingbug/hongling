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
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"time"
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
		"priority":     whoisit,
		"forFix": map[string][]string{
			TEST: {"10.0.0.1"},
			PROD: {},
		},
	},
	"uac": {
		"orderBy": 14,
		"name": "uac-hsf-service",
		"relativePath": "user-account-center/uac-hsf-service",
		"priority":     whoisit,
		"forFix": map[string][]string{
			TEST: {"10.0.0.1"},
			PROD: {},
		},
	},
	"schd": {
		"orderBy": 15,
		"name": "scheduler",
		"relativePath": "scheduler",
		"priority":     whoisit,
		"forFix": map[string][]string{
			TEST: {"10.0.0.1"},
			PROD: {},
		},
	},
	"tcbid": {
		"orderBy": 16,
		"name": "trans-bidding",
		"relativePath": "transaction/trans-bidding",
		"priority":     whoisit,
		"forFix": map[string][]string{
			TEST: {"10.0.0.1"},
			PROD: {},
		},
	},
	"tctrans": {
		"orderBy": 17,
		"name": "trans-transfer",
		"relativePath": "transaction/trans-transfer",
		"priority":     whoisit,
		"forFix": map[string][]string{
			TEST: {"10.0.0.1"},
			PROD: {},
		},
	},
	"tc": {
		"orderBy": 18,
		"name": "trans-bidding-hsf-service",
		"relativePath": "transaction/trans-bidding-hsf-service",
		"priority":     whoisit,
		"forFix": map[string][]string{
			TEST: {"10.0.0.1"},
			PROD: {},
		},
	},
	"mc": {
		"orderBy": 19,
		"name": "mc-hsf-service",
		"relativePath": "message-center/mc-hsf-service",
		"priority":     whoisit,
		"forFix": map[string][]string{
			TEST: {"10.0.0.1"},
			PROD: {},
		},
	},
	"pt": {
		"orderBy": 20,
		"name": "portal-hsf-service",
		"relativePath": "portal/portal-hsf-service",
		"priority":     whoisit,
		"forFix": map[string][]string{
			TEST: {"10.0.0.1"},
			PROD: {},
		},
	},
	"tcrepay": {
		"orderBy": 21,
		"name": "trans-repayment",
		"relativePath": "transaction/trans-repayment",
		"priority":     whoisit,
		"forFix": map[string][]string{
			TEST: {"10.0.0.1"},
			PROD: {},
		},
	},
}

var ExecjavaCommand = &cli.Command{
	Name:     "execjava",
	Category: "应用修复",
	Aliases:  []string{"ej"},
	Usage:    "hl [global options] execjava/ej [command options] [arguments...]",
	Action:   execJavaUsage,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "r",
			Usage:"远程服务器上项目部署的根目录,格式:<ip>:<目录>,如果所有服务器的目录相同,可以忽略ip,即:<目录>.",
		},
	},
}

func init() {
	for key, value := range modules {
		if _, ok := value["forFix"].(map[string][]string); ok {
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

// 不同环境下远程服务器部署目录
var remoteDirs = map[string]string{
	TEST: "/home/admin/taobao-tomcat-7.0.59/deploy/ROOT/WEB-INF/classes/",
	PROD: "/home/admin/taobao-tomcat-7.0.59/deploy/ROOT/WEB-INF/classes/",
}

func execJava(c *cli.Context) error {
	env, err := Verify(c.String("environment"), TEST)
	if err != nil {
		return err
	}

	rdirs := c.String("r")
	if rdirs == "" {
		rdirs = remoteDirs[env]
		Logger.Info(fmt.Sprintf("没有指定远程服务器项目的部署根目录,缺省为[%s].", rdirs))
	}

	if !c.Args().Present() {
		Logger.Info("缺少类全路径参数:hl [global options] execjava/ej [uc/uac...] <类全路径>")
		return nil
	}

	if err := pullProductionBranch(); err != nil {
		Logger.Error(fmt.Sprintf("拉取production分之失败:[%s].", err))
		return err
	}

	if err := compile(); err != nil {
		Logger.Error(fmt.Sprintf("编译项目失败:[%s].", err))
		return err
	}

	if err := execute(c.Command.Name, env, rdirs, c.Args().First()); err != nil {
		return err
	}

	return nil
}

// 拉取hl.main的production分支.
func pullProductionBranch() error {
	// clone仓库,如果已经克隆,忽略
	execCommand(CacheDir, "git", "clone", "--single-branch", "--branch", "production", "git@172.16.0.100:java/hl.main.git")

	// 对clone的仓库,做一次更新,保证代码是最新的
	if err := execCommand(CacheDir + "hl.main", "git", "pull"); err != nil {
		return err
	}
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
		if v, ok := value["forFix"].(bool); ok && !v {
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
		if _, ok := value["forFix"].(map[string][]string); ok {
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

	// 编译不通过是走不下去的
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}

// 执行修复.
func execute(command string, env string, rdirs string, classAs string) error {
	// 如果不同于缺省的远程目录,明确指定服务器对应的目录结构,有两种可能:
	// 1,如果有部分不同于缺省目录的服务器,格式为"<ip1|ip2...>:<dir>,<ip3|ip4...>:<dir>",比如"172.16.0.20:/home/opt/tomcat,172.16.0.11|172.16.0.12:/home/admin/tomcat"
	// 2,如果所有的服务器的目录都不同于缺省目录,格式为"<dir>",比如"/home/opt/tomcat"
	rdirsAfterSplitted := strings.Split(rdirs, ",")
	rip2dir, err_ := parseDirs(rdirsAfterSplitted)
	if err_ != nil {
		return err_
	}

	// 上载class文件到对应服务的集群部署目录
	if err := uploadClassFileToCluster(modules[command], env, rip2dir, classAs); err != nil {
		return err
	}

	return nil
}

// 解析目录参数.
func parseDirs(rdirsAfterSplitted []string) (map[string]string, error) {
	rip2dir := make(map[string]string)
	// 清理前后空格
	rdirsAfterCleaned := make([]string, 0)
	for _, v := range rdirsAfterSplitted {
		v := strings.TrimSpace(v)
		if v == "" {
			continue
		}

		rdirsAfterCleaned = append(rdirsAfterCleaned, v)
	}

	if len(rdirsAfterCleaned) == 0 {
		return nil, errors.New("指定的远程目录结构格式不合适,应该是<ip|ip...>:<dir>,<ip|ip...>:<dir>或者<dir>.")
	} else if len(rdirsAfterCleaned) == 1 {
		// 如果只有1个元素,说明是上面第2种情况
		rip2dir["*"] = rdirsAfterCleaned[0]
	} else {
		// 否则,说明是上面第1种情况
		for _, v := range rdirsAfterCleaned {
			v_ := strings.Split(v, ":")

			if len(v_) != 2 {
				return nil, errors.New(fmt.Sprintf("指定的远程目录结构%s不合适,应该是<ip|ip...>:<dir>.", v))
			} else {
				ips, dir := v_[0], v_[1]
				ipsAfterSplitted := strings.Split(ips, "|")
				for _, v__ := range ipsAfterSplitted {
					v__ = strings.TrimSpace(v__)
					if v__ == "" {
						continue
					} else if v__ == "*" {
						return nil, errors.New("IP不能为*.")
					}

					if vx, ok := rip2dir[v__]; ok {
						return nil, errors.New(fmt.Sprintf("存在相同ip[%s]指向不同的目录[%s/%s]存在.", v__, vx, dir))
					} else {
						rip2dir[v__] = dir
					}
				}
			}
		}
	}

	return rip2dir, nil
}

// 上载类文件到集群指定的目录.
func uploadClassFileToCluster(module map[string]interface{}, env string, rdirs map[string]string, classAs string) error {
	// 将类路径转换为class文件的绝对路径
	classFilePath := strings.Replace(classAs, ".", string(os.PathSeparator), -1)
	classFileLocalPath := CacheDir + "hl.main/" + module["relativePath"].(string) + "/src/main/java/" + classFilePath

	clusters := module["forFix"].(map[string][]string)[env]
	classFileRemotePaths := make(map[string]string, len(clusters))
	if rp, ok := rdirs["*"]; ok {
		if !strings.HasSuffix(rp, string(os.PathSeparator)) {
			rp = rp + string(os.PathSeparator)
		}
		for _, ip := range clusters {
			classFileRemotePaths[ip] = rp + classFilePath
		}
	} else {
		for _, ip := range clusters {
			if v, ok := rdirs[ip]; ok {
				if !strings.HasSuffix(v, string(os.PathSeparator)) {
					rp = rp + string(os.PathSeparator)
				}
				classFileRemotePaths[ip] = v + classFilePath
			} else {
				classFileRemotePaths[ip] = remoteDirs[env] + classFilePath
			}
		}
	}

	for k, v := range classFileRemotePaths {
		if err := uploadClassFileToHost(classFileLocalPath, k, v); err != nil {
			return err
		}
	}

	return nil
}

// 通过ssh传送文件到远程服务器指定目录.
func uploadClassFileToHost(localPath, ip, remotePath string) error {
	client, err := connect(ip, 22)
	if err != nil {
		return err
	}
	defer client.Close()

	l, err_ := os.Open(localPath)
	if err_ != nil {
		return err_
	}
	defer l.Close()

	r, err__ := client.Create(remotePath)
	if err__ != nil {
		return err__
	}
	defer r.Close()

	size, err___ := io.Copy(bufio.NewWriter(r), bufio.NewReader(l))
	if err___ != nil && err___ != io.EOF {
		return err___
	}
	Logger.Info(fmt.Sprintf(" 上传本地文件%s到远程主机%s的指定目录%s成功,传输大小%d字节.", localPath, ip, remotePath, size))
	return nil
	/*buf := make([]byte, 1024)
	for {
		n, err___ := l.Read(buf)
		if err___ != nil && err___ != io.EOF {
			return err___
		}

		if n == 0 {
			break
		}

		_, err____ := r.Write(buf)
		if err____ != nil {
			return err____
		}
	}*/

}

// 连接远程主机
func connect(ip string, port int) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)

	user := ""
	password := ""

	// 授权
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
	}

	// 连接
	addr = fmt.Sprintf("%s:%d", ip, port)

	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}

	return sftpClient, nil
}

// 执行命令.
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
		Logger.Fatal(fmt.Sprintf("执行命令%s时候获取标准输出通道失败.", err))
		return err
	}

	outForErr, err_ := commandAs.StderrPipe()
	defer out.Close()
	if err_ != nil {
		Logger.Fatal(fmt.Sprintf("执行命令%s时候获取错误输出通道失败.", err))
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
