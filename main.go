package main

import (
	"github.com/fzdwx/gh-sp/cmd"
)

func main() {
	cmd.Execute()
}

//package main
//
//import (
//"fmt"
//"github.com/fzdwx/gh-sp/model/table"
//)
//
//func main() {
//	headers := table.Headers{
//		&table.Header{Text: "repo name", Ratio: 5, MinWidth: 25},
//		&table.Header{Text: "description", Ratio: 10},
//		&table.Header{Text: "start count", Ratio: 3, MinWidth: 12},
//		&table.Header{Text: "visibility", Ratio: 2, MinWidth: 10},
//		&table.Header{Text: "issues", Ratio: 2, MinWidth: 10},
//	}
//
//	m := table.NewModel(headers)
//
//	m.SetRows([]string{"fzdwx/x", "golang pkg", "ð 0", "public", "ð¯ 0"})
//	m.SetRows([]string{"fzdwx/goshelldemo", "go shell demo", "ð 0", "private", "ð¯ 0"})
//	m.SetRows([]string{"fzdwx/burst ", ":technologist: åç½ç©¿é (Proxy intranet to in", "ð 46", "public ", "ð¯ 8"})
//	m.SetRows([]string{"index/worker", "the tauri showsace with vue3", "ð 0", "public  ", "ð¯ 0"})
//	m.SetRows([]string{"fzdwx/sky   ", "â¡ å¿«éåºäºNettyåå»ºåºä½ èªå·±çæå¡ï¼HTTPï¼Web", "ð 5", "public", "ð¯ 5"})
//
//	fmt.Println(m.View())
//}
