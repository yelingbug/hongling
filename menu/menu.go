package menu

import (
	"github.com/c-bata/go-prompt"
	"fmt"
	"gopkg.in/urfave/cli.v2"
	"strings"
)

type Suggestion interface {
	Get() string
	Export(document prompt.Document, args...string) []prompt.Suggest
}

var menuCommands = make(map[string]Suggestion)
var suggests []prompt.Suggest

func AddMenuCommand(name, description string, suggestion Suggestion) {
	suggests = append(suggests, prompt.Suggest{Text:name, Description:description})
	menuCommands[name] = suggestion
}

func BuildMenu(context *cli.Context) error {
	prompt := prompt.New(
		process,
		completer,
		prompt.OptionTitle("Utility by Yelin.G"),
		prompt.OptionHistory([]string{"SELECT * FROM users;"}),
		prompt.OptionPrefix(">>> "),
		prompt.OptionLivePrefix(prefix),
		prompt.OptionPrefixTextColor(prompt.Yellow),
		prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionBGColor(prompt.DarkGray))
	prompt.Run()
	return nil
}

func completer(in prompt.Document) []prompt.Suggest {
	text := strings.TrimSpace(in.TextBeforeCursor())
	// 初始为空.
	if text == "" {
		return []prompt.Suggest{}
	}

	// 找到第一个命令
	args := strings.Split(text, " ")
	rootCommand := strings.TrimSpace(args[0])
	if _, _ok := menuCommands[rootCommand]; (_ok && len(args) == 1 && in.GetWordBeforeCursor() == ""/*增加这个判断处理输入命令之后再输入一个空格才会出现后面的提示*/) || len(args) > 1 {
		// 命令存在,指派到具体的接口实现中去获取Suggest;不存在,返回空.
		if value, ok := menuCommands[rootCommand]; ok {
			return value.Export(in, args[1:]...)
		}
	}
	return prompt.FilterHasPrefix(suggests, in.GetWordBeforeCursor(), true)
}

func prefix() (string, bool) {
	return "", false
}

func process(text string) {
	fmt.Println(text)
	return
}
