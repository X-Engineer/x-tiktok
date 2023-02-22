package util

import "strconv"

func StrArrToInt64Arr(strArr []string) ([]int64, error) {
	int64Arr := make([]int64, 0, len(strArr))
	for _, str := range strArr {
		int64Val, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, err
		}
		int64Arr = append(int64Arr, int64Val)
	}
	return int64Arr, nil
}
