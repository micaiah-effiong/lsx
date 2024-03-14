package terminal

import (
	"errors"
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

		println("lenght", len(formattedByte))

		if len(formattedByte) == 1 {
			println("|>", string(formattedByte))
		} else {
			// println("hot key")
			hotKeys(b)
		}

		clear(b)
	}
}

type Key struct {
	shift       bool
	ctrl        bool
	alt         bool
	esc         bool
	payload     string
	payloadByte byte
}

func (k Key) toStrint() string {
	return fmt.Sprintln("\nshift", k.shift, "\nctrl", k.ctrl, "\nalt", k.alt, "\nesc", k.esc, "\npayload", k.payload, "\npayloadByte", k.payloadByte)
}

func hotKeys(keys []byte) {

	var cleanKeys []byte

	for _, b := range keys {
		if b > 0 {
			cleanKeys = append(cleanKeys, b)
		}
	}

	if len(cleanKeys) == 3 {
		println("hot keys")
		if k, err := buildHotkey(cleanKeys); err == nil {
			println("hk > ", k.toStrint())
		}
	} else if len(cleanKeys) > 3 {
		println("extra hot keys")
		if k, err := buildExtraHotkey(cleanKeys); err == nil {
			println("ehk > ", k.toStrint())
		}
	} else {
		println("not a hot key")
	}
}

func checkHotkey(keys []byte, checkExtra bool) (int, error) {

	if len(keys) < 3 {
		return 0, errors.New("Invaild key sequence")
	}

	if keys[0] != 27 {
		return 0, errors.New("Invaild key sequence")
	}

	if keys[1] != 91 {
		return 0, errors.New("Invaild key sequence")
	}

	if len(keys) == 3 {
		return 1, nil
	}

	if !checkExtra {
		return 1, nil
	}

	if keys[1] != 91 {
		return 0, errors.New("Invaild key sequence")
	}

	if keys[2] != 49 {
		return 0, errors.New("Invaild key sequence")
	}

	if keys[3] != 59 {
		return 0, errors.New("Invaild key sequence")
	}

	return 1, nil
}

func buildHotkey(keys []byte) (Key, error) {

	key := new(Key)

	if _, err := checkHotkey(keys, false); err != nil {
		return *key, err
	}

	keyValue := string(keys[2])

	key = &Key{alt: false, ctrl: false, esc: true, payload: keyValue, payloadByte: keys[2]}

	return *key, nil
}

func buildExtraHotkey(keys []byte) (Key, error) {

	key := new(Key)

	if _, err := checkHotkey(keys, false); err != nil {
		return *key, err
	}

	meta := metaKey[keys[4]]

	println("META ====", keys[4], meta)

	keyValue := string(keys[5])

	key = &Key{alt: false, ctrl: false, esc: true, payload: keyValue, payloadByte: keys[2]}

	for _, m := range meta {
		println("meta", m, m == "shift")

		if m == "shift" {
			key.shift = true
		}
		if m == "ctrl" {
			key.ctrl = true
		}

		if m == "alt" {
			key.alt = true
		}
	}

	return *key, nil

}

// basic navigation
// ESC [ = 27 91
// A = up
// B = down
// C = right
// D = left
// Z = tab (actually shift tab)
var basicNav = map[int]string{
	65: "A",
	66: "B",
	67: "C",
	68: "D",
	90: "Z",
}

// ESC [ <key|1> ; <meta_key> <key>
// the second part (from ;) only happens when we get 1 and not a key value

var metaKey = map[byte][]string{
	50: {"shift"},
	51: {"alt"},
	52: {"shift", "alt"},
	53: {"ctrl"},
	54: {"ctrl", "shift"},
}

// 2 = shift
// 5 = ctrl
// 3 = alt
// 6 = ctrl + shift
// 4 = shift + alt
