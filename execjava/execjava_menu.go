package execjava

import (
	"github.com/c-bata/go-prompt"
	"hongling/menu"
	"strings"
)

type ExecJavaSuggestion struct {
}

func init() {
	menu.AddMenuCommand(_command, "开发/测试/生产环境远程紧急修复.", &ExecJavaSuggestion{})
}

func (ej *ExecJavaSuggestion) Export(in prompt.Document, args...string) []prompt.Suggest {
	//fmt.Println("[CurrentLine=", in.CurrentLine(), "]--[CurrentLineAfterCursor=", in.CurrentLineAfterCursor(), "]--[CurrentLineBeforeCursor=", in.CurrentLineBeforeCursor(), "]--[CursorPositionCol=", in.CursorPositionCol(), "]--[CursorPositionRow=", in.CursorPositionRow(), "]--")
	//fmt.Println(in.DisplayCursorPosition(), "--", in.FindEndOfCurrentWord(), "--", in.FindEndOfCurrentWordWithSpace(), "--", in.FindStartOfPreviousWord(), "--", in.FindStartOfPreviousWordWithSpace())
	//fmt.Println(string(in.GetCharRelativeToCursor(0)), "--", in.CursorPositionCol(), "--", in.GetCursorLeftPosition(2), "--", in.GetCursorRightPosition(1), "--", in.GetEndOfLinePosition())
	//fmt.Println(in.GetWordBeforeCursor(), ",,,", in.GetWordBeforeCursorWithSpace(), ",,,", in.GetWordAfterCursor(), ",,,", in.GetWordAfterCursorWithSpace())
	//fmt.Println(",,,", in.GetWordBeforeCursorUntilSeparator(" "), ",,,", in.GetWordBeforeCursor(), ",,,", in.GetWordBeforeCursorWithSpace())

	if len(args) <= 1 {
		secondSubcommand := ""
		if len(args) == 1 {
			secondSubcommand = strings.TrimSpace(args[0])
		}

		suggest := make([]prompt.Suggest, 0)
		for key, value := range _modules {
			if _, ok := value["forFix"].(map[string][]string); ok {
				suggest = append(suggest, prompt.Suggest{
					Text:        key,
					Description: value["description"].(string),
				})
			}
		}
		return prompt.FilterHasPrefix(suggest, secondSubcommand, true)
	} else {
		return []prompt.Suggest{}
	}
}

func (ej *ExecJavaSuggestion) Get() string {
	return _command
}
