package utility

import (
	"github.com/c-bata/go-prompt"
	"gopkg.in/urfave/cli.v2"
	"fmt"
	"strings"
)

var suggests = []prompt.Suggest{
	{Text: "usage", Description: "用法说明."},
	{Text: "execjava", Description: "开发/测试/生产环境远程紧急修复."},
}

func completer(in prompt.Document) []prompt.Suggest {
	//fmt.Println("[CurrentLine=", in.CurrentLine(), "]--[CurrentLineAfterCursor=", in.CurrentLineAfterCursor(), "]--[CurrentLineBeforeCursor=", in.CurrentLineBeforeCursor(), "]--[CursorPositionCol=", in.CursorPositionCol(), "]--[CursorPositionRow=", in.CursorPositionRow(), "]--")
	//fmt.Println(in.DisplayCursorPosition(), "--", in.FindEndOfCurrentWord(), "--", in.FindEndOfCurrentWordWithSpace(), "--", in.FindStartOfPreviousWord(), "--", in.FindStartOfPreviousWordWithSpace())
	//fmt.Println(string(in.GetCharRelativeToCursor(0)), "--", in.CursorPositionCol(), "--", in.GetCursorLeftPosition(2), "--", in.GetCursorRightPosition(1), "--", in.GetEndOfLinePosition())
	//fmt.Println(in.GetWordBeforeCursor(), ",,,", in.GetWordBeforeCursorWithSpace(), ",,,", in.GetWordAfterCursor(), ",,,", in.GetWordAfterCursorWithSpace())
	fmt.Println(",,,", in.GetWordBeforeCursorUntilSeparator(" "), ",,,", in.GetWordBeforeCursor(), ",,,", in.GetWordBeforeCursorWithSpace())
	in.
	if in.TextBeforeCursor() == "" {
		return []prompt.Suggest{}
	}
	args := strings.Split(in.TextBeforeCursor(), " ")
	l := len(args)
	var option string
	if l >= 2 {
		option = args[l-2]
	}
	if option == "execjava" {

		return []prompt.Suggest{
			{Text: "uc", Description: "asdf."},
			{Text: "uac", Description: "开发/asdfasdfasdf/生产环境远程紧急修复."},
		}
	}
	/*l := len(args)
	if l >= 2 {
		option = args[l-2]
	}
	if strings.HasPrefix(option, "-") {
		return args[0], option, true
	}*/

	return []prompt.Suggest{}
}

func optionCompleter(args []string, long bool) []prompt.Suggest {
	/*l := len(args)
	if l <= 1 {
		if long {
			return prompt.FilterHasPrefix(optionHelp, "--", false)
		}
		return optionHelp
	}*/
	return []prompt.Suggest{
		{Text: "uc", Description: "sss."},
		{Text: "uac", Description: "sadfasdfasdfsafd."},
	}
}

func Menu(context *cli.Context) error {
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

func prefix() (string, bool) {
	return "", false
}

func process(text string) {
	fmt.Println(text)
	return
}
