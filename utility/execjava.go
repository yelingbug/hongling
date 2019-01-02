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
	"gopkg.in/AlecAivazis/survey.v1"
	"gopkg.in/AlecAivazis/survey.v1/core"
	"reflect"
	"math/rand"
	"net/http"
	"io/ioutil"
)

const (
	_whoami = iota
	_whoareu
	_whoisit
)

const (
	_user     = "yelin.g"
	_password = "hi, hongling"
)

var modules = map[string]map[string]interface{}{
	"common": {
		"orderBy":      1,
		"name":         "common",
		"relativePath": "common",
		"forFix":       false,
		"priority":     _whoami,
	},

	/*"mcommon": {
		"relativePath": "message-center/mc-common",
		"forFix": false,
	},*/
	"mclient": {
		"orderBy":      2,
		"name":         "mc-hsf-client",
		"relativePath": "message-center/mc-hsf-client",
		"forFix":       false,
		"priority":     _whoareu,
	},
	/*"ucommon": {
		"relativePath": "user-center/uc-common",
		"forFix": false,
	},*/
	"uclient": {
		"orderBy":      3,
		"name":         "uc-hsf-client",
		"relativePath": "user-center/uc-hsf-client",
		"forFix":       false,
		"priority":     _whoareu,
	},
	/*"uacommon": {
		"relativePath": "user-account-center/uac-common",
		"forFix": false,
	},*/
	"uaclient": {
		"orderBy":      4,
		"name":         "uac-hsf-client",
		"relativePath": "user-account-center/uac-hsf-client",
		"forFix":       false,
		"priority":     _whoareu,
	},
	"ssoclient": {
		"orderBy":      5,
		"name":         "sso-client",
		"relativePath": "single-sign-on/sso-client",
		"forFix":       false,
		"priority":     _whoareu,
	},
	"tcbidclient": {
		"orderBy":      6,
		"name":         "trans-bidding-hsf-client",
		"relativePath": "transaction/trans-bidding-hsf-client",
		"forFix":       false,
		"priority":     _whoareu,
	},
	"tcothersclient": {
		"orderBy":      7,
		"name":         "trans-others-client",
		"relativePath": "transaction/trans-others-client",
		"forFix":       false,
		"priority":     _whoareu,
	},
	"tcrepayclient": {
		"orderBy":      8,
		"name":         "trans-repayment-hsf-client",
		"relativePath": "transaction/trans-repayment-hsf-client",
		"forFix":       false,
		"priority":     _whoareu,
	},
	"tctransclient": {
		"orderBy":      9,
		"name":         "trans-transfer-hsf-client",
		"relativePath": "transaction/trans-transfer-hsf-client",
		"forFix":       false,
		"priority":     _whoareu,
	},
	"ptclient": {
		"orderBy":      10,
		"name":         "portal-hsf-client",
		"relativePath": "portal/portal-hsf-client",
		"forFix":       false,
		"priority":     _whoareu,
	},
	"yxclient": {
		"orderBy":      11,
		"name":         "youxuan-hsf-client",
		"relativePath": "youxuan/youxuan-hsf-client",
		"forFix":       false,
		"priority":     _whoareu,
	},
	"tcore": {
		"orderBy":      12,
		"name":         "trans-core",
		"relativePath": "transaction/trans-core",
		"forFix":       false,
		"priority":     _whoareu,
	},

	"uc": {
		"orderBy":      13,
		"name":         "uc-hsf-service",
		"relativePath": "user-center/uc-hsf-service",
		"priority":     _whoisit,
		"forFix": map[string][]string{
			TEST: {"10.139.51.136@22", "10.139.54.223@22"},
			PROD: {"10.253.43.53@22", "10.139.51.37@22", "10.139.54.60@22"},
		},
	},
	"uac": {
		"orderBy":      14,
		"name":         "uac-hsf-service",
		"relativePath": "user-account-center/uac-hsf-service",
		"priority":     _whoisit,
		"forFix": map[string][]string{
			TEST: {"10.139.39.111@22", "10.139.55.25@22"},
			PROD: {"10.139.48.208@22", "10.139.51.147@22", "10.253.42.231@22"},
		},
	},
	"schd": {
		"orderBy":      15,
		"name":         "scheduler",
		"relativePath": "scheduler",
		"priority":     _whoisit,
		"forFix": map[string][]string{
			TEST: {"10.139.49.117@22", "10.139.52.127@22"},
			PROD: {"10.139.55.170@22", "10.253.43.49@22"},
		},
	},
	"tcbid": {
		"orderBy":      16,
		"name":         "trans-bidding",
		"relativePath": "transaction/trans-bidding",
		"priority":     _whoisit,
		"forFix": map[string][]string{
			TEST: {"10.139.38.6@22", "10.139.52.11@22"},
			PROD: {"10.139.49.84@22", "10.253.43.12@22"},
		},
	},
	"tctrans": {
		"orderBy":      17,
		"name":         "trans-transfer",
		"relativePath": "transaction/trans-transfer",
		"priority":     _whoisit,
		"forFix": map[string][]string{
			TEST: {"10.139.48.224@22", "10.139.52.74@22"},
			PROD: {"10.253.43.6@22", "10.139.51.178@22"},
		},
	},
	"tc": {
		"orderBy":      18,
		"name":         "trans-bidding-hsf-service",
		"relativePath": "transaction/trans-bidding-hsf-service",
		"priority":     _whoisit,
		"forFix": map[string][]string{
			TEST: {"10.139.38.220@22", "10.139.53.143@22"},
			PROD: {"10.139.38.106@22", "10.253.43.38@22"},
		},
	},
	"mc": {
		"orderBy":      19,
		"name":         "mc-hsf-service",
		"relativePath": "message-center/mc-hsf-service",
		"priority":     _whoisit,
		"forFix": map[string][]string{
			TEST: {"10.139.51.215@22", "10.139.53.109@22"},
			PROD: {"10.253.43.37@22", "10.139.52.162@22"},
		},
	},
	"pt": {
		"orderBy":      20,
		"name":         "portal-hsf-service",
		"relativePath": "portal/portal-hsf-service",
		"priority":     _whoisit,
		"forFix": map[string][]string{
			TEST: {"10.139.48.247@22", "10.139.54.34@22"},
			PROD: {"10.253.43.59@22", "10.139.55.140@22"},
		},
	},
	"tcrepay": {
		"orderBy":      21,
		"name":         "trans-repayment",
		"relativePath": "transaction/trans-repayment",
		"priority":     _whoisit,
		"forFix": map[string][]string{
			TEST: {"10.139.39.85@22", "10.139.48.194@22"},
			PROD: {"10.139.55.16@22", "10.139.52.96@22", "10.253.43.41@22"},
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
			Name:  "r",
			Usage: "远程服务器上项目部署的根目录,格式:<ip>:<目录>,如果所有服务器的目录相同,可以忽略ip,即:<目录>.",
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
    public void fix(FixLogger logger) {
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

      public void fix(FixLogger logger) {
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
var _remoteDirs = map[string]string{
	TEST: "/home/admin/taobao-tomcat-7.0.59/deploy/ROOT/WEB-INF/classes/",
	PROD: "/home/admin/taobao-tomcat-7.0.59/deploy/ROOT/WEB-INF/classes/",
}

func execJava(c *cli.Context) error {
	env, err := Verify(c.String("environment"), TEST)
	if err != nil {
		return err
	}

	if env == DEV {
		env = TEST
	}

	rdirs := c.String("r")
	if rdirs == "" {
		rdirs = _remoteDirs[env]
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

	rip2dir, err := parse(rdirs)
	if err != nil {
		return err
	}

	if err := upload(c.Command.Name, env, rip2dir, c.Args().First()); err != nil {
		return err
	}

	if err := execute(c.Command.Name, env, rip2dir, c.Args().First()); err != nil {
		return err
	}

	//connect("104.128.237.136", 57896)
	//uploadClassFileToHost("./main", "104.128.237.136@57896", "/root/main")

	/*var sftpClient *sftp.Client
	if sshClient, err := createSshTunnel("104.128.237.136", "57896"); err != nil {
		return err
	} else {
		var err_ error
		if sftpClient, err_ = createSftpClient(sshClient); err_ != nil {
			return err_
		}
	}
	defer sftpClient.Close()

	timer := time.NewTicker(2 * time.Second)
	for {
		if f, err := sftpClient.Open("/root/abc"); err != nil && os.IsNotExist(err) {
			fmt.Println("文件不存在.....")
		} else {
			if s, sss := ioutil.ReadAll(f); sss != nil {
				fmt.Println("读取文件失败...")
			} else {
				fmt.Println(fmt.Sprintf("文件%s存在.", s))
				break
			}
		}
		<- timer.C
	}
	defer timer.Stop()*/

	return nil
}

// 通过内网调用远程服务http接口执行命令.
func execute(command, env string, rip2dir map[string]string, class string) error {
	// 随机选择一台服务器.
	cluster := modules[command]["forFix"].(map[string][]string)[env]
	rand.Seed(time.Now().Unix())
	iport := cluster[rand.Intn(len(cluster))]

	classFilePath := strings.Replace(class, ".", string(os.PathSeparator), -1)
	classFileRemotePaths := getRemotePathForClass(modules[command], env, rip2dir, classFilePath)

	iportAfterSplitted := strings.Split(iport, "@")
	if len(iportAfterSplitted) != 2 {
		return errors.New(fmt.Sprintf("预定义服务器ip端口%s格式错误.", iport))
	}

	// 准备http post请求,送出class类路径字符串
	ip, port := strings.TrimSpace(iportAfterSplitted[0]), strings.TrimSpace(iportAfterSplitted[1])
	url := fmt.Sprintf("http://%s:8080/fix", ip)

	client := &http.Client{
		Timeout: 1 * time.Minute,
	}

	// 发出请求,服务立即返回
	resp, err := client.Post(url, "text/plain", strings.NewReader(class))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if result, err := ioutil.ReadAll(resp.Body); err != nil {
		return err
	} else {
		resultAsString := string(result)
		if resultAsString == "OK" {
			Logger.Info("请求成功,坐等服务进程处理结果.")
		} else {
			status := fmt.Sprintf("请求成功,但是服务进程预处理失败,检查参数:[%s].", resultAsString)
			return errors.New(status)
		}
	}

	var sftpClient *sftp.Client
	if sshClient, err := createSshTunnel(ip, port); err != nil {
		return err
	} else {
		var err_ error
		if sftpClient, err_ = createSftpClient(sshClient); err_ != nil {
			return err_
		}
	}
	defer sftpClient.Close()

	// 2s一次检查任务是否执行完成.
	timer := time.NewTicker(2 * time.Second)
	timeout := 0
	for {
		if f, err := sftpClient.Open(classFileRemotePaths[iport] + "/" + class + ".doing"); err != nil && os.IsNotExist(err) {
			Logger.Info("远程服务进程还未开始,等待... ... ...")
		} else if err == nil {
			if _, err_ := sftpClient.Open(classFileRemotePaths[iport] + "/" + class + ".done"); err_ != nil && os.IsNotExist(err_) {
				Logger.Info("远程服务进程已经开始,处理日志正在刷新,还未完成... ... ...")
			} else if err_ == nil {
				if result, err__ := ioutil.ReadAll(f); err__ != nil {
					Logger.Warn(fmt.Sprintf("读取远程服务进程处理结果失败,请自行前往服务器%s相应目录下查看结果.", ip))
				} else {
					Logger.Info(fmt.Sprintf(`远程服务进程处理完成,处理结果:
%s
`, result))
					break
				}

			}

		} else {
			return errors.New(fmt.Sprintf("尝试打开远程服务进程处理日志失败:%s", err))
		}
		<- timer.C
		timeout++
		if timeout == 30 {
			Logger.Info(fmt.Sprintf("等待远程服务进程处理结果超过1分钟,请自行前往服务器%s相应目录下查看结果.", ip))
			break
		}
	}
	defer timer.Stop()

	return nil
}

// 拉取hl.main的production分支.
func pullProductionBranch() error {
	// clone仓库,如果已经克隆,忽略
	execCommand(CacheDir, "git", "clone", "--single-branch", "--branch", "production", "git@172.16.0.100:java/hl.main.git")

	// 对clone的仓库,做一次更新,保证代码是最新的
	if err := execCommand(CacheDir+"hl.main", "git", "pull"); err != nil {
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

// 上载文件.
func upload(command string, env string, rip2dir map[string]string, classAs string) error {
	// 上载class文件到对应服务的集群部署目录
	if err := uploadClassFileToCluster(modules[command], env, rip2dir, classAs); err != nil {
		return err
	}

	return nil
}

// 解析目录参数.
func parse(rdirs string) (map[string]string, error) {
	// 如果不同于缺省的远程目录,明确指定服务器对应的目录结构,有两种可能:
	// 1,如果有部分不同于缺省目录的服务器,格式为"<ip1|ip2...>:<dir>,<ip3|ip4...>:<dir>",比如"172.16.0.20:/home/opt/tomcat,172.16.0.11|172.16.0.12:/home/admin/tomcat"
	// 2,如果所有的服务器的目录都不同于缺省目录,格式为"<dir>",比如"/home/opt/tomcat"
	rdirsAfterSplitted := strings.Split(rdirs, ",")

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
	classFilePath := strings.Replace(classAs, ".", string(os.PathSeparator), -1) + ".class"
	classFileLocalPath := CacheDir + "hl.main/" + module["relativePath"].(string) + "/src/main/java/" + classFilePath
	for k, v := range getRemotePathForClass(module, env, rdirs, classFilePath) {
		if err := uploadClassFileToHost(classFileLocalPath, k, v); err != nil {
			return err
		}
	}

	return nil
}

// 如果classFilePath是以.class结尾的文件,返回的就是每个服务器对应的远程class文件放置的绝对路径;
// 如果classFilePath不以.class结尾,返回的就是每个服务器对应的远程class文件所在的目录.
func getRemotePathForClass(module map[string]interface{}, env string, rdirs map[string]string, classFilePath string) map[string]string {
	// 如果不以.class文件结尾,要的结果就是class文件所在的远程服务器目录.
	if !strings.HasSuffix(classFilePath, ".class") {
		index := strings.LastIndex(classFilePath, string(os.PathSeparator))
		if index == -1 {
			classFilePath = ""
		} else {
			classFilePath = classFilePath[:index+1]
		}
	}
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
				classFileRemotePaths[ip] = _remoteDirs[env] + classFilePath
			}
		}
	}
	return classFileRemotePaths
}

// 通过ssh传送文件到远程服务器指定目录.
func uploadClassFileToHost(localPath, iport, remotePath string) error {
	iportAfterSplitted := strings.Split(iport, "@")
	if len(iportAfterSplitted) != 2 {
		return errors.New(fmt.Sprintf("预定义服务器ip端口%s格式错误.", iport))
	}

	ip, port := iportAfterSplitted[0], iportAfterSplitted[1]

	var sftpClient *sftp.Client
	if sshClient, err := createSshTunnel(ip, port); err != nil {
		return err
	} else {
		var _err error
		sftpClient, _err = createSftpClient(sshClient)
		if _err != nil {
			return _err
		}
	}
	defer sftpClient.Close()

	l, err_ := os.Open(localPath)
	if err_ != nil {
		return err_
	}
	defer l.Close()

	r, err__ := sftpClient.Create(remotePath)
	if err__ != nil {
		return err__
	}
	defer r.Close()

	size, err___ := io.Copy(bufio.NewWriter(r), bufio.NewReader(l))
	if err___ != nil && err___ != io.EOF {
		return err___
	}
	Logger.Info(fmt.Sprintf(" 上传本地文件%s到远程主机%s的指定目录%s成功,传输大小%d字节.", localPath, iport, remotePath, size))
	return nil
}

// 连接远程主机
func createSshTunnel(ip string, port string) (*ssh.Client, error) {
	var (
		sshClient *ssh.Client
		err       error
	)

	// 连接
	addr := fmt.Sprintf("%s:%s", ip, port)

	auth := struct {
		// survey内部使用反射,首字母必须大写.
		User     string
		Password string
	}{_user, _password}

	for {
		sshClient, err = ssh.Dial("tcp", addr, buildConnectContext(auth.User, auth.Password))
		if err == nil {
			return sshClient, nil
		}

		// 只对鉴权失败继续尝试.
		if !strings.Contains(err.Error(), "handshake failed: ssh: unable to authenticate") {
			return nil, errors.New(fmt.Sprintf("连接%s时候发生异常:%s.", addr, err))
		}

		Logger.Warn(fmt.Sprintf("连接%s鉴权失败,重新输入用户名密码.", addr))

		// 自定义错误输出模板
		core.ErrorTemplate = `{{color "red"}}{{ ErrorIcon }} 抱歉, 输入无效: {{.Error}}{{color "reset"}}
`
		inputs := []*survey.Question{
			{
				Name: "user",
				Prompt: &survey.Input{
					Message: "用户名:",
					Help:    "用户名不能过长,长度一般在3到20之间.",
				},
				Validate: requiredWithMessage("用户名"),
			},
			{
				Name: "password",
				Prompt: &survey.Password{
					Message: "密码:",
					Help:    "密码稍微复杂一点,严格保密.",
				},
				Validate: requiredWithMessage("密码"),
			},
		}
		if err_ := survey.Ask(inputs, &auth); err_ != nil {
			return nil, errors.New(fmt.Sprintf("连接时候交互输入处理失败:%s.", err_))
		}
	}
}

func createSftpClient(client *ssh.Client) (*sftp.Client, error) {
	// create sftp client
	if sftpClient, err := sftp.NewClient(client); err != nil {
		return nil, errors.New(fmt.Sprintf("连接成功但是创建ftp连接时失败:%s.", err))
	} else {
		return sftpClient, nil
	}
}

func buildConnectContext(user, password string) *ssh.ClientConfig {
	// 授权
	auth := make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	return &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
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

// isZero returns true if the passed value is the zero object
func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Slice, reflect.Map:
		return v.Len() == 0
	}

	// compare the types directly with more general coverage
	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}

func requiredWithMessage(message string) func(val interface{}) error {
	return func(val interface{}) error {
		// the reflect value of the result
		value := reflect.ValueOf(val)

		// if the value passed in is the zero value of the appropriate type
		if isZero(value) && value.Kind() != reflect.Bool {
			return errors.New(message + "不能为空.")
		}
		return nil
	}
}
