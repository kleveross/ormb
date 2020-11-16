package cmd

import (
	"fmt"
	"strings"
)

func convertRef(ref string) string {
	refSlice := strings.Split(ref, "/")
	if len(refSlice) == 2 {
		temp := strings.Split(refSlice[len(refSlice)-1], ":")
		if len(temp) == 2 {
			return fmt.Sprintf("%v/%v/%v:%v", refSlice[0], temp[0], temp[0], temp[1])
		}
		return ref
	}

	return ref
}
