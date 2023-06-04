package docker

import "strconv"

func TranslateCPU(str string) int64 {
	len := len(str)
	res := 0.0
	if str[len-1] == 'm' { //
		res, _ = strconv.ParseFloat(str[:len-1], 32)
		res *= 1e6
	} else {
		res, _ = strconv.ParseFloat(str[:len], 32)
		res *= 1e9
	}
	return int64(res)
}

func TranslateMem(str string) int64 {
	len := len(str)
	res, _ := strconv.Atoi(str[:len-1])
	if str[len-1] == 'M' || str[len-1] == 'm' { // mb
		res = res * 1024 * 1024
	} else if str[len-1] == 'K' || str[len-1] == 'k' { // kb
		res = res * 1024
	} else if str[len-1] == 'G' || str[len-1] == 'g' { // gb
		res = res * 1024 * 1024 * 1024
	}
	return int64(res)
}
