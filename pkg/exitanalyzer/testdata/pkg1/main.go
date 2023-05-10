package main

import "os"

func main() {
	if 2 > 1 {
		if 3 > 2 {
			if 5 > 4 {
				// nested
				os.Exit(1) // want "cannot use \"os.Exit\" in package main \"main\" function"
			}
		}
	}
	// top level
	os.Exit(0) // want "cannot use \"os.Exit\" in package main \"main\" function"
}
