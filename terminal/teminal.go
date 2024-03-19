package terminal

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

type Terminal_reader struct {
}

func (tr Terminal_reader) Reader() (Key, error) {
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
			// fmt.Printf(">>> %v %q %T\n", item, string(item), item)
			formattedByte = append(formattedByte, item)
		}

		// println("lenght", len(formattedByte))

		if len(formattedByte) == 1 {
			return newKey(formattedByte[0]), nil
		} else {
			if k, err := hotKeys(formattedByte); err == nil {
				return k, err
			} else {
				println("Not supported hot key")
			}
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
	PayloadByte byte
	desc        string
}

func (k Key) To_string() string {
	return fmt.Sprintln(k)
}

func newKey(b byte) Key {
	k := new(Key)

	k.payload = fmt.Sprintf("%q", string(b))
	k.PayloadByte = b
	return *k
}

func hotKeys(keys []byte) (Key, error) {

	var cleanKeys []byte

	for _, b := range keys {
		if b > 0 {
			cleanKeys = append(cleanKeys, b)
		}
	}

	ckLen := len(cleanKeys)

	if ckLen == 3 {
		return buildHotkey(cleanKeys)
	} else if ckLen == 6 {
		return buildExtraHotkey(cleanKeys)
	} else {
		return *new(Key), errors.New("No hot keys")
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

	if len(keys) != 6 {
		return 0, errors.New("Invaild key sequence(Not allowed)")
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

	key = &Key{alt: false, ctrl: false, esc: true, payload: keyValue, PayloadByte: keys[2]}

	if nav := basicNav[keys[2]]; nav != nil {
		key.desc = nav[1]
	}

	return *key, nil
}

func buildExtraHotkey(keys []byte) (Key, error) {

	key := new(Key)

	if _, err := checkHotkey(keys, true); err != nil {
		return *key, err
	}

	meta := metaKey[keys[4]]

	// println("META ====", keys[4], meta)

	keyValue := string(keys[5])

	key = &Key{alt: false, ctrl: false, esc: true, payload: keyValue, PayloadByte: keys[2]}

	for _, m := range meta {
		// println("meta", m, m == "shift")

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
var basicNav = map[byte][]string{
	65: {"A", "Up"},
	66: {"B", "Down"},
	67: {"C", "Right"},
	68: {"D", "Left"},
	10: {"\n", "Enter"},
	9:  {"\t", "Tab"},
	90: {"Z", "Shift+Tab"},
}

var basicNavValue = map[byte]int{
	65: -1, //{"A", "Up"},
	66: 1,  //{"B", "Down"},
	67: 0,  //{"C", "Right"},
	68: 0,  //{"D", "Left"},
	10: 0,  //{"\n", "Enter"},
	9:  0,  //{"\t", "Tab"},
	90: 0,  //{"Z", "Shift+Tab"},
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
