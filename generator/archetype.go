package generator

import (
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v2"
	"os"
	"time"
	"io/ioutil"
	"hongling/utility"
)

const (
	_FLAG_SPRING       = "spring"
	_FLAG_SPRING_WEB   = "spring-web"
	_FLAG_SPRING_BOOT   = "spring-boot"
	_POM_SPRING        = `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>

    <groupId>%s</groupId>
    <artifactId>%s</artifactId>
    <version>1.0.0</version>
    <packaging>pom</packaging>

    <properties>
        <project.build.sourceEncoding>UTF-8</project.build.sourceEncoding>
        <mockito.version>2.15.0</mockito.version>
        <powermock.version>2.0.0-beta.5</powermock.version>
        <junit.version>4.12</junit.version>
        <spring.version>5.1.3.RELEASE</spring.version>
        <druid.version>1.1.9</druid.version>
        <guava.version>23.5-jre</guava.version>
        <mysql.version>5.1.44</mysql.version>
    </properties>

    <dependencies>
        <dependency>
            <groupId>org.springframework</groupId>
            <artifactId>spring-context</artifactId>
            <version>${spring.version}</version>
        </dependency>

        <dependency>
            <groupId>org.springframework</groupId>
            <artifactId>spring-context-support</artifactId>
            <version>${spring.version}</version>
        </dependency>

        <dependency>
            <groupId>org.springframework</groupId>
            <artifactId>spring-core</artifactId>
            <version>${spring.version}</version>
        </dependency>

        <dependency>
            <groupId>org.springframework</groupId>
            <artifactId>spring-beans</artifactId>
            <version>${spring.version}</version>
        </dependency>

        <dependency>
            <groupId>org.springframework</groupId>
            <artifactId>spring-jdbc</artifactId>
            <version>${spring.version}</version>
        </dependency>

        <dependency>
            <groupId>org.springframework</groupId>
            <artifactId>spring-tx</artifactId>
            <version>${spring.version}</version>
        </dependency>

        <dependency>
            <groupId>org.springframework</groupId>
            <artifactId>spring-aop</artifactId>
            <version>${spring.version}</version>
        </dependency>

        <dependency>
            <groupId>org.springframework</groupId>
            <artifactId>spring-orm</artifactId>
            <version>${spring.version}</version>
        </dependency>

        <dependency>
            <groupId>org.springframework</groupId>
            <artifactId>spring-test</artifactId>
            <version>${spring.version}</version>
        </dependency>

        <dependency>
            <groupId>com.alibaba</groupId>
            <artifactId>druid</artifactId>
            <version>${druid.version}</version>
        </dependency>

        <dependency>
            <groupId>com.google.guava</groupId>
            <artifactId>guava</artifactId>
            <version>${guava.version}</version>
        </dependency>

        <dependency>
            <groupId>mysql</groupId>
            <artifactId>mysql-connector-java</artifactId>
            <version>${mysql.version}</version>
        </dependency>

        <dependency>
            <groupId>junit</groupId>
            <artifactId>junit</artifactId>
            <version>${junit.version}</version>
        </dependency>

        <dependency>
            <groupId>org.mockito</groupId>
            <artifactId>mockito-core</artifactId>
            <version>${mockito.version}</version>
        </dependency>

        <dependency>
            <groupId>org.powermock</groupId>
            <artifactId>powermock-module-junit4</artifactId>
            <version>${powermock.version}</version>
        </dependency>

        <dependency>
            <groupId>org.powermock</groupId>
            <artifactId>powermock-api-mockito2</artifactId>
            <version>${powermock.version}</version>
        </dependency>
    </dependencies>
</project>    
`
	_POM_SPRING_WEB    = `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>

    <groupId>%s</groupId>
    <artifactId>%s</artifactId>
    <version>1.0.0</version>
    <packaging>pom</packaging>

    <properties>
        <project.build.sourceEncoding>UTF-8</project.build.sourceEncoding>
        <mockito.version>2.15.0</mockito.version>
        <powermock.version>2.0.0-beta.5</powermock.version>
        <junit.version>4.12</junit.version>
        <spring.version>5.1.3.RELEASE</spring.version>
        <druid.version>1.1.9</druid.version>
        <guava.version>23.5-jre</guava.version>
        <mysql.version>5.1.44</mysql.version>
    </properties>

    <dependencies>
        <dependency>
            <groupId>org.springframework</groupId>
            <artifactId>spring-context</artifactId>
            <version>${spring.version}</version>
        </dependency>

        <dependency>
            <groupId>org.springframework</groupId>
            <artifactId>spring-context-support</artifactId>
            <version>${spring.version}</version>
        </dependency>

        <dependency>
            <groupId>org.springframework</groupId>
            <artifactId>spring-core</artifactId>
            <version>${spring.version}</version>
        </dependency>

        <dependency>
            <groupId>org.springframework</groupId>
            <artifactId>spring-beans</artifactId>
            <version>${spring.version}</version>
        </dependency>

        <dependency>
            <groupId>org.springframework</groupId>
            <artifactId>spring-jdbc</artifactId>
            <version>${spring.version}</version>
        </dependency>

        <dependency>
            <groupId>org.springframework</groupId>
            <artifactId>spring-tx</artifactId>
            <version>${spring.version}</version>
        </dependency>

        <dependency>
            <groupId>org.springframework</groupId>
            <artifactId>spring-aop</artifactId>
            <version>${spring.version}</version>
        </dependency>

        <dependency>
            <groupId>org.springframework</groupId>
            <artifactId>spring-orm</artifactId>
            <version>${spring.version}</version>
        </dependency>

        <dependency>
            <groupId>org.springframework</groupId>
            <artifactId>spring-test</artifactId>
            <version>${spring.version}</version>
        </dependency>

        <dependency>
            <groupId>org.springframework</groupId>
            <artifactId>spring-web</artifactId>
            <version>${spring.version}</version>
        </dependency>

        <dependency>
            <groupId>org.springframework</groupId>
            <artifactId>spring-webmvc</artifactId>
            <version>${spring.version}</version>
        </dependency>

        <dependency>
            <groupId>com.alibaba</groupId>
            <artifactId>druid</artifactId>
            <version>${druid.version}</version>
        </dependency>

        <dependency>
            <groupId>com.google.guava</groupId>
            <artifactId>guava</artifactId>
            <version>${guava.version}</version>
        </dependency>

        <dependency>
            <groupId>mysql</groupId>
            <artifactId>mysql-connector-java</artifactId>
            <version>${mysql.version}</version>
        </dependency>

        <dependency>
            <groupId>junit</groupId>
            <artifactId>junit</artifactId>
            <version>${junit.version}</version>
        </dependency>

        <dependency>
            <groupId>org.mockito</groupId>
            <artifactId>mockito-core</artifactId>
            <version>${mockito.version}</version>
        </dependency>

        <dependency>
            <groupId>org.powermock</groupId>
            <artifactId>powermock-module-junit4</artifactId>
            <version>${powermock.version}</version>
        </dependency>

        <dependency>
            <groupId>org.powermock</groupId>
            <artifactId>powermock-api-mockito2</artifactId>
            <version>${powermock.version}</version>
        </dependency>
    </dependencies>
</project>
`
	_POM_SPRING_BOOT    = `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>

    <groupId>%s</groupId>
    <artifactId>%s</artifactId>
    <version>1.0.0</version>
    <packaging>pom</packaging>

    <properties>
        <project.build.sourceEncoding>UTF-8</project.build.sourceEncoding>
        <mockito.version>2.15.0</mockito.version>
        <powermock.version>2.0.0-beta.5</powermock.version>
        <junit.version>4.12</junit.version>
        <springboot.version>2.1.1.RELEASE</springboot.version>
        <druid.version>1.1.9</druid.version>
        <guava.version>23.5-jre</guava.version>
        <mysql.version>5.1.44</mysql.version>
    </properties>

    <dependencies>
        <parent>
	        <groupId>org.springframework.boot</groupId>
	        <artifactId>spring-boot-starter-parent</artifactId>
	        <version>${springboot.version}</version>
        </parent>

        <dependency>
            <groupId>com.alibaba</groupId>
            <artifactId>druid</artifactId>
            <version>${druid.version}</version>
        </dependency>

        <dependency>
            <groupId>com.google.guava</groupId>
            <artifactId>guava</artifactId>
            <version>${guava.version}</version>
        </dependency>

        <dependency>
            <groupId>mysql</groupId>
            <artifactId>mysql-connector-java</artifactId>
            <version>${mysql.version}</version>
        </dependency>

        <dependency>
            <groupId>junit</groupId>
            <artifactId>junit</artifactId>
            <version>${junit.version}</version>
        </dependency>

        <dependency>
            <groupId>org.mockito</groupId>
            <artifactId>mockito-core</artifactId>
            <version>${mockito.version}</version>
        </dependency>

        <dependency>
            <groupId>org.powermock</groupId>
            <artifactId>powermock-module-junit4</artifactId>
            <version>${powermock.version}</version>
        </dependency>

        <dependency>
            <groupId>org.powermock</groupId>
            <artifactId>powermock-api-mockito2</artifactId>
            <version>${powermock.version}</version>
        </dependency>
    </dependencies>

    <build>
	    <plugins>
		    <plugin>
			    <groupId>org.springframework.boot</groupId>
			    <artifactId>spring-boot-maven-plugin</artifactId>
		    </plugin>
	    </plugins>
    </build>
</project>
`
	_CONFIG_SPRING     = `<?xml version="1.0" encoding="UTF-8"?>
<beans xmlns="http://www.springframework.org/schema/beans"
       xmlns:context="http://www.springframework.org/schema/context"
       xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
       xmlns:aop="http://www.springframework.org/schema/aop"
       xsi:schemaLocation="http://www.springframework.org/schema/beans
                          http://www.springframework.org/schema/beans/spring-beans-4.3.xsd
                          http://www.springframework.org/schema/context
                          http://www.springframework.org/schema/context/spring-context-4.3.xsd
                          http://www.springframework.org/schema/aop
                          http://www.springframework.org/schema/aop/spring-aop-4.3.xsd">
    <aop:aspectj-autoproxy/>
    <context:annotation-config/>
    <context:component-scan base-package="*" use-default-filters="true"/>
</beans>
`
	_CONFIG_SPRING_WEB = `<?xml version="1.0" encoding="UTF-8"?>
<beans xmlns="http://www.springframework.org/schema/beans"
       xmlns:context="http://www.springframework.org/schema/context"
       xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
       xmlns:aop="http://www.springframework.org/schema/aop"
       xsi:schemaLocation="http://www.springframework.org/schema/beans
                          http://www.springframework.org/schema/beans/spring-beans-4.3.xsd
                          http://www.springframework.org/schema/context
                          http://www.springframework.org/schema/context/spring-context-4.3.xsd
                          http://www.springframework.org/schema/mvc
                          http://www.springframework.org/schema/mvc/spring-mvc-4.3.xsd
                          http://www.springframework.org/schema/aop
                          http://www.springframework.org/schema/aop/spring-aop-4.3.xsd">
    <aop:aspectj-autoproxy/>
    <context:annotation-config/>
    <bean id="conversionService" class="org.springframework.format.support.FormattingConversionServiceFactoryBean"/>
    <mvc:annotation-driven
            conversion-service="conversionService"
            validator="validator"
            content-negotiation-manager="contentNegotiationManager">
    <mvc:interceptors>
        <bean class="org.springframework.web.servlet.i18n.LocaleChangeInterceptor">
            <property name="paramName" value="language"/>
        </bean>
        <bean class="org.springframework.web.servlet.handler.ConversionServiceExposingInterceptor">
            <constructor-arg ref="conversionService"/>
        </bean>
    </mvc:interceptors>
    <mvc:default-servlet-handler/>
</beans>
`
	_CONFIG_SPRING_BOOT = ``
)

var ArchetypeCommand = &cli.Command{
	Name:     "archetype",
	Category: "项目模版生成",
	Aliases:  []string{"template"},
	Usage:    "hl [global options] archetype/template [command options] [arguments...]",
	Action:   archetype,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "name",
			Usage: "指定项目名称.",
		},
		&cli.BoolFlag{
			Name:  _FLAG_SPRING,
			Usage: "生成基于spring-framework(5.1.3.RELEASE)项目模版(包含core/beans/context/context-support/jdbc/transaction/aop/aspects/orm/test).",
		},
		&cli.BoolFlag{
			Name:  _FLAG_SPRING_WEB,
			Usage: "生成基于spring-framework(5.1.3.RELEASE)项目模版(包含core/beans/context/context-support/jdbc/transaction/aop/aspects/orm/web/webmvc/test).",
		},
		&cli.BoolFlag{
			Name:  _FLAG_SPRING_BOOT,
			Usage: "生成基于spring-boot(2.1.1)项目模版.",
		},
	},
}

type archetype_args struct {
	name string //项目名称
	t    string //项目类型spring/spring-web
}

type archetype_action interface {
	isSpring() bool
	isSpringWeb() bool
	isSpringBoot() bool
	getRoot() string
	generatePom() error
	generateApplicationConfig() error
	generateApplicationProperties() error
	generate() error
}

func (aa *archetype_args) getRoot() string {
	return utility.CacheDir + aa.name
}

func (aa *archetype_args) isSpring() bool {
	if aa.t == _FLAG_SPRING {
		return true
	} else {
		return false
	}
}

func (aa *archetype_args) isSpringWeb() bool {
	if aa.t == _FLAG_SPRING_WEB {
		return true
	} else {
		return false
	}
}

func (aa *archetype_args) isSpringBoot() bool {
	if aa.t == _FLAG_SPRING_BOOT {
		return true
	} else {
		return false
	}
}

func (aa *archetype_args) generatePom() error {
	if aa.isSpring() {
		utility.Logger.Info(fmt.Sprintf("生成%s类型项目的pom文件.", _FLAG_SPRING))
		return generateSpringPom(aa)
	} else if aa.isSpringWeb() {
		utility.Logger.Info(fmt.Sprintf("生成%s类型项目的pom文件.", _FLAG_SPRING_WEB))
		return generateSpringWebPom(aa)
	} else if aa.isSpringBoot() {
		utility.Logger.Info(fmt.Sprintf("生成%s类型项目的pom文件.", _FLAG_SPRING_BOOT))
		return generateSpringBootPom(aa)
	} else {
		utility.Logger.Warn(fmt.Sprintf("不支持%s类型项目的创建，略过.", aa.t))
		return nil
	}
}

func generateSpringPom(aa *archetype_args) error {
	return createFile(aa.getRoot()+"/pom.xml", _POM_SPRING)
}

func generateSpringWebPom(aa *archetype_args) error {
	return createFile(aa.getRoot()+"/pom.xml", _POM_SPRING_WEB)
}

func generateSpringBootPom(aa *archetype_args) error {
	return createFile(aa.getRoot()+"/pom.xml", _POM_SPRING_BOOT)
}

func (aa *archetype_args) generateApplicationConfig() error {
	utility.Logger.Info(fmt.Sprintf("生成src/test的application配置."))
	for _, entry := range []map[string]string{
		{
			"fileName": aa.getRoot() + "/src/main/resources/application.xml",
			"content":  generateApplicationConfig(aa),
		},
		{
			"fileName": aa.getRoot() + "/src/test/resources/application.xml",
			"content":  generateApplicationConfig(aa),
		},
	} {
		if err := createFile(entry["fileName"], entry["content"]); err != nil {
			return err
		}
	}
	return nil
}

func generateApplicationConfig(aa *archetype_args) string {
	if aa.isSpringWeb() {
		return _CONFIG_SPRING_WEB
	} else if aa.isSpring() {
		return _CONFIG_SPRING
	} else if aa.isSpringBoot() {
		return _CONFIG_SPRING_BOOT
	}
	return ""
}

func (aa *archetype_args) generateApplicationProperties() error {
	utility.Logger.Info(fmt.Sprintf("生成src/test的项目属性配置."))
	for _, entry := range []map[string]string{
		{
			"fileName": aa.getRoot() + "/src/main/resources/application.properties",
			"content":  "#Generated by template at " + time.Now().Format("2006/01/02 15:04:05"),
		},
		{
			"fileName": aa.getRoot() + "/src/test/resources/application.properties",
			"content":  "#Generated by template at " + time.Now().Format("2006/01/02 15:04:05"),
		},
	} {
		if err := createFile(entry["fileName"], entry["content"]); err != nil {
			return err
		}
	}
	return nil
}

func (aa *archetype_args) generate() error {
	//初始化项目目录
	name_ := utility.CacheDir + aa.name
	if err := createDir(name_); err != nil {
		return err
	}

	//初始化maven样式的目录结构
	utility.Logger.Info(fmt.Sprintf("初始化项目%s的maven样式目录结构.", aa.name))
	for _, dir := range []string{
		"src/main/java",
		"src/main/resources",
		"src/test/java",
		"src/test/resources",
	} {
		if err := createDir(name_ + "/" + dir); err != nil {
			return err
		}
	}

	//生成pom文件
	if err := aa.generatePom(); err != nil {
		return err
	}

	//生成application.xml配置
	if err := aa.generateApplicationConfig(); err != nil {
		return err
	}

	//生成application.properties配置
	if err := aa.generateApplicationProperties(); err != nil {
		return err
	}

	utility.Logger.Info(fmt.Sprintf("生成项目%s模版完成.", aa.name))
	return nil
}

func archetype(c *cli.Context) error {
	name := c.String("name")
	if name == "" {
		return errors.New("需要指定项目名称.")
	}

	aa := &archetype_args{
		name: name,
	}
	if c.Bool(_FLAG_SPRING) {
		aa.t = _FLAG_SPRING
	} else if c.Bool(_FLAG_SPRING_WEB) {
		aa.t = _FLAG_SPRING_WEB
	} else if c.Bool(_FLAG_SPRING_BOOT) {
		aa.t = _FLAG_SPRING_BOOT
	} else {
		utility.Logger.Info("没有指定项目类型,默认为spring-boot.")
		aa.t = _FLAG_SPRING_BOOT
	}

	return aa.generate()
}

func createDir(name string) error {
	f, err := os.Stat(name)
	if err == nil && f.IsDir() {
		utility.Logger.Warn("目录" + name + "存在.")
	} else if err != nil && os.IsNotExist(err) {
		utility.Logger.Info("创建目录" + name + ".")
		if err_ := os.MkdirAll(name, os.ModePerm); err_ != nil {
			return err_
		} else {
			return nil
		}
	} else if err != nil {
		return err
	}
	return nil
}

func createFile(name, content string) error {
	return ioutil.WriteFile(name, []byte(content), 0666)
}
