package utils

func Uint8ToAny(data []uint8) (res []any) {
	for _, d := range data {
		res = append(res, d)
	}
	return
}

func Uint16ToAny(data []uint16) (res []any) {
	for _, d := range data {
		res = append(res, d)
	}
	return
}

func Uint64ToAny(data []uint64) (res []any) {
	for _, d := range data {
		res = append(res, d)
	}
	return
}
