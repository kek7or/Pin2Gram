package bot

import (
	"pinterest-tg-autopost/dbtypes"
	"pinterest-tg-autopost/pinterest"
	"strings"
)

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

func InSlicePin(value pinterest.Pin, slice []pinterest.Pin) bool {
	for _, e := range slice {
		if e.ID == value.ID {
			return true
		}
	}

	return false
}

func NonImplicationPins(sliceA []pinterest.Pin, sliceB []dbtypes.Post) []pinterest.Pin {
	if len(sliceB) == 0 {
		return sliceA
	}

	newSlice := make([]pinterest.Pin, 0)

	for _, a := range sliceA {
		isAinB := false
		for _, b := range sliceB {
			if a.ID == b.PinId {
				isAinB = true
			}
		}

		if !isAinB {
			newSlice = append(newSlice, a)
		}
	}

	return newSlice
}

func TrimElements(sliceList []string, toTrim string) []string {
	for i, item := range sliceList {
		sliceList[i] = strings.Trim(item, toTrim)
	}
	return sliceList
}
