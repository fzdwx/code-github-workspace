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
//	m.AppendRow([]string{"fzdwx/x", "golang pkg", "🌟 0", "public", "🎯 0"})
//	m.AppendRow([]string{"fzdwx/goshelldemo", "go shell demo", "🌟 0", "private", "🎯 0"})
//	m.AppendRow([]string{"fzdwx/burst ", ":technologist: 内网穿透 (Proxy intranet to in", "🌟 46", "public ", "🎯 8"})
//	m.AppendRow([]string{"index/worker", "the tauri showsace with vue3", "🌟 0", "public  ", "🎯 0"})
//	m.AppendRow([]string{"fzdwx/sky   ", "⚡ 快速基于Netty创建出你自己的服务（HTTP，Web", "🌟 5", "public", "🎯 5"})
//
//	fmt.Println(m.View())
//}
