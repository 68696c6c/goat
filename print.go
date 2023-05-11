package goat

import (
	"fmt"

	"github.com/fatih/color"
)

func Print(s string, a ...any) {
	println(getPrintOutput(s, a...))
}

func PrintSuccess(s string, a ...any) {
	color.Green(s, a...)
}

func PrintInfo(s string, a ...any) {
	color.Blue(s, a...)
}

func PrintWarning(s string, a ...any) {
	color.Yellow(s, a...)
}

func PrintDanger(s string, a ...any) {
	color.Red(s, a...)
}

func PrintHeading(s string, a ...any) {
	d := color.New(color.Bold)
	d.Println(getPrintOutput(s, a...))
}

func PrintIndent(s string, a ...any) {
	println("    " + getPrintOutput(s, a...))
}

func getPrintOutput(s string, a ...any) string {
	if len(a) > 0 {
		return fmt.Sprintf(s, a...)
	}
	return s
}
