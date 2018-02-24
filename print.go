package goat

import "github.com/fatih/color"

func PrintSuccess(s string) {
	color.Green(s)
}

func PrintInfo(s string) {
	color.Blue(s)
}

func PrintFail(s string) {
	color.Red(s)
}

func PrintWarning(s string) {
	color.Red(s)
}
