package terminal

import (
	"fmt"
	"os"
	"os/exec"
)

type TerminalReader struct {
}

func (tr TerminalReader) Reader() {
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	var b []byte = make([]byte, 10)
	for {
		os.Stdin.Read(b)

		var formattedByte []byte

		for _, item := range b {
			if item == 0 {
				continue
			}
			fmt.Printf(">>> %v %q %T\n", item, string(item), item)

			formattedByte = append(formattedByte, item)
		}

		clear(b)

		println("lenght", len(formattedByte))

		if len(formattedByte) == 1 {
			println("|>", string(formattedByte))
		} else {
			println("hot key")
		}
	}
}

// basic navigation
// ESC [ = 27 91
// A = up
// B = down
// C = right
// D = left
// Z = tab (actually shift tab)

// ESC [ <key|1> ; <meta_key> <key>
// the second part (from ;) only happens when we get 1 and not a key value

var metaKey = map[int][]string{
	2: {"shift"},
	3: {"alt"},
	4: {"shift", "alt"},
	5: {"ctrl"},
	6: {"ctrl", "shift"},
}

// 2 = shift
// 5 = ctrl
// 3 = alt
// 6 = ctrl + shift
// 4 = shift + alt
