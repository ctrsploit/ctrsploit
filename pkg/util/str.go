package util

func Int8ToStr(arr []int8) string {
	byteSlice := make([]byte, len(arr))
	for i, v := range arr {
		byteSlice[i] = byte(v)
	}
	return string(byteSlice)
}

func StrToInt8(s string) []int8 {
	bs := []byte(s)
	res := make([]int8, len(bs))
	for i, v := range bs {
		res[i] = int8(v)
	}

	return res
}
