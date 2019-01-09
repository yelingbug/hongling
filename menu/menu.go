package menu

import (
	"github.com/c-bata/go-prompt"
	"fmt"
	"gopkg.in/urfave/cli.v2"
	"strings"
)

type Suggestion interface {
	Get() string
	Export(document prompt.Document) []prompt.Suggest
	Process()
	Validate() error
}

var menuCommands = make(map[string]Suggestion)
var suggests []prompt.Suggest

func AddMenuCommand(name, description string, suggestion Suggestion) {
	suggests = append(suggests, prompt.Suggest{Text: name, Description: description})
	menuCommands[name] = suggestion
}

func BuildMenu(context *cli.Context) error {
	prompt := prompt.New(
		process,
		completer,
		prompt.OptionTitle("Utility by Yelin.G"),
		prompt.OptionHistory([]string{"SELECT * FROM users;"}),
		prompt.OptionMaxSuggestion(6),
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
	if value, ok := menuCommands[rootCommand]; ok {
		return value.Export(in)
	}
	return prompt.FilterHasPrefix(suggests, in.GetWordBeforeCursor(), true)
}

func prefix() (string, bool) {
	return "", false
}

func process(text string) {
	textAfterFixed := strings.TrimSpace(text)
	if textAfterFixed == "" {
		fmt.Println("你是在搞笑吗?")
		return
	}
	args := strings.Split(textAfterFixed, " ")
	if value, ok := menuCommands[args[0]]; ok {
		if err := value.Validate(); err == nil {
			value.Process()
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println("输入无效.")
	}
	return
}
