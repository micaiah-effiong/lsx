package terminal

func GetNavKeyCalculatedValue(key Key, pos int, size int) int {
	val := getNavKeyValue(key)

	if (val + pos) > size {
		pos = size
	} else if (val + pos) < 0 {
		pos = 0
	} else {
		pos += val
	}

	return pos
}

func getNavKeyValue(k Key) int {
	if k.alt || k.ctrl || k.shift {
		return 0
	}

	return basicNavValue[k.PayloadByte]
}
