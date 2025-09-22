package util

func Int8ToStr(arr []int8) string {
	byteSlice := make([]byte, len(arr))
	for i, v := range arr {
		byteSlice[i] = byte(v)
	}
	return string(byteSlice)
}
