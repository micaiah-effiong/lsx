package main

import (
	"flag"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path"
	"syscall"
)

func main() {
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

	if first_arg == "~" {
		home_dir, err := os.Hostname()
		if err != nil {
			print(err)
			print("Error getting user home directory")
			// panic("An error occurred while reading directory")
		}

		first_arg = home_dir
	}

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

	for _, line := range ls {
		line_name := line.Name()
		if string(line_name[0]) == "." {
			continue
		}
		ls_name_list = append(ls_name_list, line.Name())
	}

	render_list(ls_name_list, pos)

	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	var b []byte = make([]byte, 1)
	var command_modifier_byte []byte
	for {
		os.Stdin.Read(b)

		if b[0] == 10 {

			cwd, err := syscall.Getwd()
			if err != nil {
				log.Fatal(err)
				panic("Failed to get working directory")
			}

			joined_path := path.Join(cwd, first_arg, ls_name_list[pos])
			// println(joined_path, cwd, first_arg, ls_name_list[pos])
			// println(joined_path)

			clear()

			os.Stdout.Write([]byte(joined_path))

			// err = syscall.Chdir(joined_path)
			//
			// if err != nil {
			// 	log.Fatal(err)
			// 	panic("Failed to get change directory")
			// }

			os.Exit(0)
			return
		}

		// if len(command_modifier_byte) > 0 && command_modifier_byte[0] == 27 {
		// 	continue
		// }

		// ESC
		if b[0] == 27 {
			command_modifier_byte = nil
			command_modifier_byte = append(command_modifier_byte, b[0])
		}

		var key_string string
		if len(command_modifier_byte) > 0 && b[0] != 27 {
			command_modifier_byte = append(command_modifier_byte, b[0])
		}

		if len(command_modifier_byte) == 3 {
			k, valid := nice_bytes(command_modifier_byte)

			if valid {
				key_string = k
				pos += get_nice_key_val(k)
				command_modifier_byte = nil
			} else {
				println("Invalid", string(command_modifier_byte))
			}
		}

		if b[0] != 27 && len(command_modifier_byte) < 1 && len(key_string) < 1 {
			key_string = string(b[0])
		}

		if len(key_string) > 0 {
			clear()
			render_list(ls_name_list, pos)
		}

	}
}

func render_list(list []string, pos int) {
	for line_pos, line := range list {
		if line_pos == pos {
			println("→ ", line)
		} else {
			println(" ", line)
		}
	}
}

func clear() {

	// print("\033[6n")

	print("\033[0;0H\033[2J")
	// print("\033[H\033[2J")
}

func nice_bytes(bytes_list []byte) (string, bool) {
	if len(bytes_list) != 3 {
		println("not long enough")
		return "", false
	}

	// esc 27 = ESC
	// b_braket 91 = [
	// val ? = ?
	esc, b_braket, val := bytes_list[0], bytes_list[1], bytes_list[2]

	if esc != 27 || b_braket != 91 {
		println("proper no escape sequence")
		return "", false
	}

	val_string := string(val)

	switch val_string {
	case "A":
		return "Up", true
	case "B":
		return "Down", true
	case "C":
		return "Left", true
	case "D":
		return "Right", true
	case "Z":
		return "<Shift>Tab", true
	default:
		return "", false
	}

	// return string(val), true
}

func get_nice_key_val(key string) int {

	nice_key_map := make(map[string]int)
	nice_key_map["<Shift>Tab"] = -1
	nice_key_map["Up"] = -1
	nice_key_map["Down"] = 1
	nice_key_map["Left"] = 0
	nice_key_map["Right"] = 0

	return nice_key_map[key]

	// return 0
}
