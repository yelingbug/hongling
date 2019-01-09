package execjava

import (
	"github.com/c-bata/go-prompt"
	"hongling/menu"
	"strings"
	"fmt"
	"hongling/utility"
	"github.com/pkg/errors"
)

type ExecJavaSuggestion struct {
	rValue, bValue, module, class string
}

func init() {
	menu.AddMenuCommand(_command, "开发/测试/生产环境远程紧急修复.", &ExecJavaSuggestion{})
}

func (ej *ExecJavaSuggestion) Export(in prompt.Document) []prompt.Suggest {
	args := strings.Split(in.TextBeforeCursor(), " ")

	p1 := ""
	if len(args) > 1 {
		p1 = args[1]
	}

	switch p1 {
	case "-" + _flags[0]:
		// execjava -r ??? -b
		if len(args) == 4 {

			// 有四个参数,说明-r的值已经实锤了.
			ej.rValue = args[2]

			if args[3] == "-" + _flags[1] {
				return prompt.FilterHasPrefix(appendBranchs(nil), "", true)
			} else {
				return prompt.FilterHasPrefix(appendModules(appendArgs(nil, 1)), in.GetWordBeforeCursor(), true)
			}
		}

		if len(args) == 5 {
			if strings.HasPrefix(args[3], "-") {
				ej.bValue = args[4]
			} else {
				ej.module = args[3]
				ej.class = args[4]
				return []prompt.Suggest{}
			}
			if args[3] == "-" + _flags[1] {
				// 有五个参数,说明-b的值在输入当中,但是肯定已经实锤了.
				ej.bValue = args[4]
				return prompt.FilterHasPrefix(appendBranchs(nil), in.GetWordBeforeCursor(), true)
			}
		}

		if len(args) == 6 {
			ej.bValue = args[4]//刷新,因为有可能是通过tab选择的值
			ej.module = args[5]
			return prompt.FilterHasPrefix(appendModules(nil), in.GetWordBeforeCursor(), true)
		}

		if len(args) == 7 {
			ej.module = args[5]//刷新,因为有可能是通过tab选择的值
			ej.class = args[6]
		}

		return []prompt.Suggest{}
	case "-" + _flags[1]:
		if len(args) == 3 {
			ej.bValue = args[2]
			return prompt.FilterHasPrefix(appendBranchs(nil), in.GetWordBeforeCursor(), true)
		}

		// execjava -b ??? -r
		if len(args) == 4 {
			ej.bValue = args[2]//刷新,因为有可能是通过tab选择的值

			if args[3] == "-" + _flags[0] {
				return []prompt.Suggest{}
			} else {
				return prompt.FilterHasPrefix(appendModules(appendArgs(nil, 0)), in.GetWordBeforeCursor(), true)
			}
		}

	 	// execjava -b ??? -r ???
		if len(args) == 5 {
			if strings.HasPrefix(args[3], "-") {
				ej.rValue = args[4]
			} else {
				ej.module = args[3]
				ej.class = args[4]
				return []prompt.Suggest{}
			}
		}

		if len(args) == 6 {
			ej.module = args[5]
			return prompt.FilterHasPrefix(appendModules(nil), in.GetWordBeforeCursor(), true)
		}

		if len(args) == 7 {
			ej.module = args[5]//刷新,因为有可能是通过tab选择的值
			ej.class = args[6]
		}

		return []prompt.Suggest{}
	default:
		if len(args) == 3 {
			ej.module = args[1]
			ej.class = args[2]
			return []prompt.Suggest{}
		}
		return prompt.FilterHasPrefix(appendModules(appendArgs(appendArgs(nil, 0), 1)), in.GetWordBeforeCursor(), true)
	}
}

func appendBranchs(suggests []prompt.Suggest) []prompt.Suggest {
	if suggests == nil {
		suggests = make([]prompt.Suggest, 0)
	}
	return append(suggests, prompt.Suggest{
		Text:        "pre-production",
		Description: "预发布分支.",
	}, prompt.Suggest{
		Text:        "production",
		Description: "生产分支.",
	})
}

func appendArgs(suggests []prompt.Suggest, index int) []prompt.Suggest {
	if suggests == nil {
		suggests = make([]prompt.Suggest, 0)
	}
	return append(suggests, prompt.Suggest{
		Text:        "-" + _flags[index],
		Description: _flagsDetail[_flags[index]],
	})
}

func appendModules(suggests []prompt.Suggest) []prompt.Suggest {
	if suggests == nil {
		suggests = make([]prompt.Suggest, 0)
	}
	for key, value := range _modules {
		if _, ok := value["forFix"].(map[string][]string); ok {
			suggests = append(suggests, prompt.Suggest{
				Text:        key,
				Description: value["description"].(string),
			})
		}
	}
	return suggests
}

func (ej *ExecJavaSuggestion) Get() string {
	return _command
}

func (ej *ExecJavaSuggestion) Process() {
	if ej.bValue != "pre-production" && ej.bValue != "production" {
		utility.Logger.Info("指定的分支不是pre-production或者production,默认为production")
		ej.bValue = "production"
	}

	env := utility.PROD
	if ej.bValue == "pre-production" {
		env = utility.TEST
	}

	if ej.rValue == "" {
		ej.rValue = _remoteDirs[env]
		utility.Logger.Info(fmt.Sprintf("没有指定远程服务器项目的部署根目录,缺省为[%s].", ej.rValue))
	}

	if err := execJava_(ej.module, env, ej.bValue, ej.rValue, ej.class, true); err != nil {
		fmt.Println(err)
	}

	// 完成之后要清理,准备一下次继续.
	ej.rValue = ""
	ej.bValue = ""
	ej.module = ""
	ej.class = ""
}

func (ej *ExecJavaSuggestion) String() string {
	return fmt.Sprintf("rValue=[%s], bValue=[%s], module=[%s], class=[%s].", ej.rValue, ej.bValue, ej.module, ej.class)
}

func (ej *ExecJavaSuggestion) Validate() error {
	if _, ok := _modules[ej.module]; !ok {
		return errors.New(fmt.Sprintf("%s不是合法的模块.", ej.module))
	}

	return nil
}
