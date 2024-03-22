package terminal

func GetNavKeyCalculatedValue(key Key, pos int, size int) int {
	val := get_nav_key_value(key)
	assm := val + pos

	size_indexing := size - 1

	if assm > size_indexing {
		pos = max(size_indexing, 0)
	} else if assm < 0 {
		pos = 0
	} else {
		pos = assm
	}

	return pos
}

func get_nav_key_value(k Key) int {
	if k.alt || k.ctrl || k.shift {
		return 0
	}

	return basicNavValue[k.PayloadByte]
}
