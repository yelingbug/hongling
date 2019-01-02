package utility

import (
	"github.com/lestrrat-go/file-rotatelogs"
	"time"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/rifflock/lfshook"
	"os/user"
	"runtime"
	"os"
	"net"
)

var Logger *logrus.Entry

var CacheDir = getCacheDir()

func getCacheDir() string {
	userdir, err := homedir.Dir()
	if err != nil {
		userdir = "." + string(os.PathSeparator)
	}

	return fmt.Sprintf("%s%s.utility%s", userdir, string(os.PathSeparator), string(os.PathSeparator))
}

// 初始化logger.
func init() {
	os.Mkdir(CacheDir, os.ModeDir|os.ModePerm)

	//info和error分开
	//info保存debug, info和warn level, error保存error和fatal level.
	forInfo, _ := rotatelogs.New(
		fmt.Sprintf("%s%sinfo", string(os.PathSeparator), CacheDir) + ".%Y%m%d",
		rotatelogs.WithRotationTime(1*time.Minute),
		rotatelogs.WithMaxAge(7 * 24 *time.Hour))
	forError, _ := rotatelogs.New(
		fmt.Sprintf("%s%serror", string(os.PathSeparator), CacheDir) + ".%Y%m%d",
		rotatelogs.WithRotationTime(1*time.Minute),
		rotatelogs.WithMaxAge(7 * 24 *time.Hour))

	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006/01/02 15:04:05",
		DisableTimestamp: false,
		FullTimestamp:true,
	})
	logger.SetOutput(os.Stdout)
	logger.AddHook(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.DebugLevel: forInfo,
			logrus.InfoLevel:  forInfo,
			logrus.WarnLevel: forError,
			logrus.FatalLevel: forError,
			logrus.ErrorLevel: forError,
		},
		&logrus.JSONFormatter{
			TimestampFormat: "2006/01/02 15:04:05",
			DisableTimestamp: false,
		},
	))

	var username, reasonAsUsernameFailed, hostname string
	if u, err := user.Current(); err != nil {
		username = "Unknown"
		if hostname_, err_ := os.Hostname(); err_ != nil {
			hostname = "Unknown"
		} else {
			hostname = hostname_
		}

		reasonAsUsernameFailed = fmt.Sprintf("Get current user failed in the OS(%s)/Arch(%s)/Hostname(%s).", runtime.GOOS, runtime.GOARCH, hostname)
	} else {
		username = u.Username
	}

	var ip, reasonAsAddressFailed string
	if addresses, err := net.InterfaceAddrs(); err != nil {
		ip = "Unknown"
		reasonAsAddressFailed = fmt.Sprintf("Get ip address failed in the OS(%s)/Arch(%s)/Hostname(%s).", runtime.GOOS, runtime.GOARCH, hostname)
	} else {
		for _, address := range addresses {
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ip = ipnet.IP.String()
					break
				} else {
					ip = "Unknown"
					reasonAsAddressFailed = fmt.Sprintf("Get ip4 address failed in the OS(%s)/Arch(%s)/Hostname(%s).", runtime.GOOS, runtime.GOARCH, hostname)
				}
			}
		}
	}

	Logger = logger.WithFields(logrus.Fields{"user": username, "ip": ip})

	if reasonAsUsernameFailed != "" {
		Logger.Error(reasonAsUsernameFailed)
	}

	if reasonAsAddressFailed != "" {
		Logger.Error(reasonAsAddressFailed)
	}

}
