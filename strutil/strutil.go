package strutil

// Excel表格头中的英文字母转索引，例如 A -> 0, B -> 1, Z -> 25, AA -> 26, AB -> 27
func ExcelHeadChar2Index(header string) int {
	var result int
	for idx, it := range header {
		intChar := int(it) - 65
		if idx == len(header)-1 {
			result = result + (len(header)-1-idx)*(intChar+1)*26 + intChar
		} else {
			result = result + (len(header)-1-idx)*(intChar+1)*26
		}
	}
	return result
}
