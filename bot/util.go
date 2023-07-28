package bot

import "strings"

func AreAddChannelArgsValid(args []string) bool {
	if len(args) < 1 {
		return false
	}

	return true
}

func RemoveDuplicate[T string | int](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func TrimElements(sliceList []string, toTrim string) []string {
	for i, item := range sliceList {
		sliceList[i] = strings.Trim(item, toTrim)
	}
	return sliceList
}
