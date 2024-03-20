package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"

	"github.com/micaiah-effiong/lsx/render"
	"github.com/micaiah-effiong/lsx/terminal"
)

func re_run(tm terminal.Terminal_reader, ls_name_list []render.Entry, searched_list []render.Entry, pos int, search_str string) string {

	println(fmt.Sprintf("Search: %v \n", search_str))

	var _searched_list []render.Entry
	if search_str != "" {
		for _, entry := range ls_name_list {
			if strings.Contains(strings.ToLower(entry.Name), strings.ToLower(search_str)) {
				_searched_list = append(_searched_list, entry)
			}
		}
	} else {
		_searched_list = ls_name_list
	}

	searched_list = _searched_list

	render_list(searched_list, pos)

	k, err := tm.Reader()
	if err != nil {
		panic(err)
	}

	clear()

	if k.PayloadByte == 10 {
		return searched_list[pos].Name
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

	new_pos := terminal.GetNavKeyCalculatedValue(k, pos, len(searched_list)-1)

	return re_run(tm, ls_name_list, searched_list, new_pos, formatted_search_str)
}

// func main() {
// 	tm := terminal.TerminalReader{}
//
// 	k, _ := tm.Reader()
//
// 	fmt.Println(k.ToStrint())
// }

func main() {

	handle_termination()

	flag.Parse()
	args := flag.Args()

	var first_arg string

	if len(args) < 1 {
		first_arg = "."
	} else {
		first_arg = args[0]
	}

	_fs := os.DirFS(first_arg)
	ls, err := fs.ReadDir(_fs, ".")

	if err != nil {
		print(err)
		log.Fatal("No such file or directory")
		// panic("An error occurred while reading directory")

		os.Exit(1)
		return
	}

	clear()

	pos := 0
	var ls_name_list []render.Entry

	for _, dir_list_item := range ls {

		entry, err := render.MakeEntry(dir_list_item)

		if err != nil {
			continue
		}

		if entry.IsDotEntry {
			continue
		}

		ls_name_list = append(ls_name_list, entry)
	}

	tm := terminal.Terminal_reader{}

	hide_cursor()
	clear()
	choosen_path := re_run(tm, ls_name_list, make([]render.Entry, 0), pos, "")
	show_cursor()

	joined_path := path.Join(first_arg, choosen_path)

	clear()
	os.Stdout.Write([]byte(joined_path))

	return
}

func render_list(list []render.Entry, pos int) {
	const SIZE = 5

	var new_list []render.Entry
	list_len := len(list)
	// println("pos: ", pos, len(list))

	var new_pos = pos

	if list_len <= SIZE {
		new_list = list
	} else if pos > SIZE-1 {
		new_list = list[pos-SIZE+1 : pos+1]
		pos = len(new_list) - 1
	} else {
		new_list = list[0:SIZE]
	}

	for line_pos, line := range new_list {
		if line_pos == pos {
			println("\033[38;5;220m → ", line.RenderName, "\033[38;5;231m")
		} else {
			println(" ", line.RenderName)
		}
	}

	println(fmt.Sprintf("%v/%v", new_pos+1, list_len))
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
