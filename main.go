package main

import (
	"errors"
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

	ls, err := get_path_entries(first_arg)

	if err != nil {
		os.Exit(1)
		return
	}

	pos := 0
	tm := terminal.Terminal_reader{}

	hide_cursor()
	clear()
	choosen_path := re_run(tm, ls, make([]render.Entry, 0), pos, "")
	show_cursor()

	joined_path := path.Join(first_arg, choosen_path)

	clear()
	os.Stdout.Write([]byte(joined_path))

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

func get_path_entries(path string) ([]render.Entry, error) {
	file_system := os.DirFS(path)
	ls, err := fs.ReadDir(file_system, ".")

	if err != nil {
		print(err)
		log.Fatal("No such file or directory")

		return nil, errors.New("No such file or directory")
	}

	var render_entries []render.Entry

	for _, dir_list_item := range ls {
		entry, err := render.MakeEntry(dir_list_item)

		if err != nil || entry.IsDotEntry {
			continue
		}

		render_entries = append(render_entries, entry)
	}

	return render_entries, nil
}

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
	max_render_size := 5
	render.RenderList(searched_list, pos, max_render_size)

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

	new_pos := terminal.GetNavKeyCalculatedValue(k, pos, len(searched_list))

	return re_run(tm, ls_name_list, searched_list, new_pos, formatted_search_str)
}
