package colors

import "fmt"

var (
	//Black tints colors black
	Black = color("\033[1;30m%s\033[0m")
	//Red tints colors red
	Red = color("\033[1;31m%s\033[0m")
	//Green tints colors green
	Green = color("\033[1;32m%s\033[0m")
	//Yellow tints colors yellow
	Yellow = color("\033[1;33m%s\033[0m")
	//Purple tints colors purple
	Purple = color("\033[1;34m%s\033[0m")
	//Magenta tints colors magenta
	Magenta = color("\033[1;35m%s\033[0m")
	//Teal tints colors teal
	Teal = color("\033[1;36m%s\033[0m")
	//White tints colors white
	White = color("\033[1;37m%s\033[0m")

	//ERROR prints a stylized error prefix
	ERROR = White("[") + Red("ERROR") + White("]") + " "

	//DOCS prints a stylized doc prefix
	DOCS = White("[") + Magenta("DOCS") + White("]") + " "

	//STATUS prints a stylized status prefix
	STATUS = White("[") + Teal("STATUS") + White("]") + " "
)

func color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}
