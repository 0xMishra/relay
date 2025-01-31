package utils

import (
	"fmt"
	"os"
)

func CheckErr(err error, fatal bool) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		if fatal {
			os.Exit(1)
		}
	}
}
