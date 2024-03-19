package main

import (
	"flag"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"

	"github.com/micaiah-effiong/lsx/terminal"
)

func re_run(tm terminal.Terminal_reader, ls_name_list []string, pos int) string {

	render_list(ls_name_list, pos)
	k, err := tm.Reader()
	if err != nil {
		panic(err)
	}

	clear()

	// fmt.Println(k.To_string())

	if k.PayloadByte == 10 {
		return ls_name_list[pos]
	}

	new_pos := terminal.GetNavKeyCalculatedValue(k, pos, len(ls_name_list)-1)

	return re_run(tm, ls_name_list, new_pos)
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

	// for _, arg := range args {
	// 	println("all args ", arg)
	// }

	var first_arg string

	if len(args) < 1 {
		first_arg = "."
	} else {
		first_arg = args[0]
	}

	// if first_arg == "~" {
	// 	home_dir, err := os.UserHomeDir()
	// 	if err != nil {
	// 		print(err)
	// 		print("Error getting user home directory")
	// 		// panic("An error occurred while reading directory")
	// 	}
	//
	// 	first_arg = home_dir
	// }

	_fs := os.DirFS(first_arg)
	ls, err := fs.ReadDir(_fs, ".")

	if err != nil {
		print(err)
		print("No such file or directory")
		// panic("An error occurred while reading directory")
	}

	clear()

	pos := 0
	var ls_name_list []string

	for _, dir_list_item := range ls {
		info, err := dir_list_item.Info()

		if err != nil {
			continue
		}

		if strings.HasPrefix(info.Name(), ".") {
			continue
		}

		var name string
		// if info.IsDir() {
		// 	name = "\033[38;5;112m" + info.Name() + "\033[38;5;231m"
		// } else {
		name = dir_list_item.Name()
		// }

		ls_name_list = append(ls_name_list, name)
	}

	tm := terminal.Terminal_reader{}

	hide_cursor()
	clear()
	choosen_path := re_run(tm, ls_name_list, pos)
	show_cursor()

	if err != nil {
		log.Fatal(err)
		panic("Failed to get working directory")
	}

	joined_path := path.Join(first_arg, choosen_path)

	clear()
	println(">> ", first_arg)
	println(">> ", choosen_path)
	os.Stdout.Write([]byte(joined_path))

	return
}

func render_list(list []string, pos int) {
	const SIZE = 5

	var new_list []string
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
			println("\033[38;5;220m → ", line, "\033[38;5;231m")
		} else {
			println(" ", line)
		}
	}

	println(new_pos+1, "/", list_len)
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
			os.Exit(0)
			return
		}
	}()

}
