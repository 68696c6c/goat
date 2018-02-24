package goat

import (
	"github.com/fatih/color"
	"encoding/json"
)

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
	color.Yellow(s)
}

func PrintHeading(s string) {
	d := color.New(color.FgHiWhite, color.Bold)
	d.Println(s)
}

func PrintIndent(s string) {
	json.MarshalIndent(s, "", "\t")
}
