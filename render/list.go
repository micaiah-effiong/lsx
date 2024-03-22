package render

import "fmt"

func RenderList(list []Entry, pos int, max_render_size int) {

	var new_list []Entry
	list_len := len(list)
	// println("pos: ", pos, len(list))

	new_pos := &pos
	if list_len <= max_render_size {
		new_list = list
	} else if pos > max_render_size-1 {
		new_list = list[pos-max_render_size+1 : pos+1]
		*new_pos = len(new_list) - 1
	} else {
		new_list = list[0:max_render_size]
	}

	for line_pos, line := range new_list {
		if line_pos == *new_pos {
			println("\033[38;5;220m → ", line.RenderName, "\033[38;5;231m")
		} else {
			println(" ", line.RenderName)
		}
	}

	var current_page = 0

	if list_len == 0 {
		current_page = 0
	} else {
		current_page = *new_pos + 1
	}

	println(fmt.Sprintf("%v/%v", current_page, list_len))
}
