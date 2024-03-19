package terminal

func GetNavKeyCalculatedValue(key Key, pos int, size int) int {
	val := get_nav_key_value(key)

	if (val + pos) > size {
		pos = size
	} else if (val + pos) < 0 {
		pos = 0
	} else {
		pos += val
	}

	return pos
}

func get_nav_key_value(k Key) int {
	if k.alt || k.ctrl || k.shift {
		return 0
	}

	return basicNavValue[k.PayloadByte]
}
