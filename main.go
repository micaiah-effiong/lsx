package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"

	"github.com/micaiah-effiong/lsx/render"
	"github.com/micaiah-effiong/lsx/terminal"
)

// func main() {
// 	tm := terminal.TerminalReader{}
//
// 	k, _ := tm.Reader()
//
// 	fmt.Println(k.ToStrint())
// }

func main() {

	handle_termination()

	a_flag := flag.Bool("a", false, "")
	s_flag := flag.Int("s", 5, "")

	flag.Parse()
	args := flag.Args()

	var first_arg string

	if len(args) < 1 {
		first_arg = "."
	} else {
		first_arg = args[0]
	}

	tm := terminal.Terminal_reader{}

	config_flags := map[string]int{"s": *s_flag, "a": bool2Int(*a_flag)}
	config := LsxConfig{tm, first_arg, 0, *s_flag, config_flags}
	choosen_path := find_lsx_path(config)

	clear()
	os.Stdout.Write([]byte(choosen_path))

	return
}

func clear() {
	print("\033[0;0H\033[2J")
}

// make cursor invisible
func hide_cursor() {
	println("\033[?25l")
}

// make cursor visible
func show_cursor() {
	println("\033[?25h")
}

func handle_termination() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		for range c {
			// log.Printf("captured %v, stopping profiler and exiting..", sig)
			show_cursor()
			clear()
			os.Exit(0)
			return
		}
	}()
}

func render_get_choice(config LsxConfig, ls_name_list []render.Entry, searched_list []render.Entry, search_str string) string {

	clear()
	println(config.SearchPath)
	println(fmt.Sprintf("Search: %v \n", search_str))

	var filtered_searched_list []render.Entry
	if search_str != "" {
		for _, entry := range ls_name_list {
			if strings.Contains(strings.ToLower(entry.Name), strings.ToLower(search_str)) {
				filtered_searched_list = append(filtered_searched_list, entry)
			}
		}
	} else {
		filtered_searched_list = ls_name_list
	}

	searched_list = filtered_searched_list
	render.RenderList(searched_list, config.Position, config.ListSize)

	k, err := config.Terminal.Reader()
	if err != nil {
		panic(err)
	}

	clear()

	if k.PayloadByte == 9 && len(searched_list) > 0 { // tab
		choice_item := searched_list[config.Position].Name
		return handle_list_tab(config, choice_item)
	}

	if k.PayloadByte == 90 && k.IsHotKey() { // shift+tab
		return handle_list_tab(config, "..")
	}

	if k.PayloadByte == 10 { // <Enter>
		if len(searched_list) <= 0 {
			return ""
		}
		choice_item := searched_list[config.Position].Name
		return path.Join(config.SearchPath, choice_item)
	}

	// println(k.ToString())

	formatted_search_str := search_str
	if !k.IsHotKey() {
		formatted_search_str += string(k.PayloadByte)
	}

	if k.PayloadByte == 127 && !k.IsHotKey() { // backspace
		fstr_len := len(formatted_search_str)
		if fstr_len > 1 {
			formatted_search_str = formatted_search_str[:fstr_len-2]
		}

		if fstr_len == 1 {
			formatted_search_str = ""
		}
	}

	config.Position = terminal.GetNavKeyCalculatedValue(k, config.Position, len(searched_list))

	return render_get_choice(config, ls_name_list, searched_list, formatted_search_str)
}

func find_lsx_path(config LsxConfig) string {

	ls, err := terminal.GetPathEntries(config.SearchPath, config.Flags["a"] != 0)

	if err != nil {
		os.Exit(1)
		return ""
	}

	hide_cursor()
	choosen_path := render_get_choice(config, ls, make([]render.Entry, 0), "")
	show_cursor()

	return choosen_path
	// return path.Join(config.SearchPath, choosen_path)
}

type LsxConfig struct {
	Terminal   terminal.Terminal_reader
	SearchPath string
	Position   int
	ListSize   int
	Flags      map[string]int
}

func handle_list_tab(config LsxConfig, choice_item string) string {
	joined_path := path.Join(config.SearchPath, choice_item)
	config.SearchPath = joined_path
	config.Position = 0
	println("tab-into:", config.SearchPath, choice_item, joined_path)
	return find_lsx_path(config)
}

func bool2Int(n bool) int {
	if n {
		return 1
	}
	return 0
}
