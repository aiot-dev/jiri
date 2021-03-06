// Copyright 2017 The Fuchsia Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package color

import (
	"fmt"
	"os"
)

type Colorfn func(format string, a ...interface{}) string

const (
	escape = "\033["
	clear  = escape + "0m"
)

type ColorCode int

// Foreground text colors
const (
	BlackFg ColorCode = iota + 30
	RedFg
	GreenFg
	YellowFg
	BlueFg
	MagentaFg
	CyanFg
	WhiteFg
	DefaultFg
)

type Color interface {
	Black(format string, a ...interface{}) string
	Red(format string, a ...interface{}) string
	Green(format string, a ...interface{}) string
	Yellow(format string, a ...interface{}) string
	Blue(format string, a ...interface{}) string
	Magenta(format string, a ...interface{}) string
	Cyan(format string, a ...interface{}) string
	White(format string, a ...interface{}) string
	DefaultColor(format string, a ...interface{}) string
}

type color struct{}

func (color) Black(format string, a ...interface{}) string { return colorString(BlackFg, format, a...) }
func (color) Red(format string, a ...interface{}) string   { return colorString(RedFg, format, a...) }
func (color) Green(format string, a ...interface{}) string { return colorString(GreenFg, format, a...) }
func (color) Yellow(format string, a ...interface{}) string {
	return colorString(YellowFg, format, a...)
}
func (color) Blue(format string, a ...interface{}) string { return colorString(BlueFg, format, a...) }
func (color) Magenta(format string, a ...interface{}) string {
	return colorString(MagentaFg, format, a...)
}
func (color) Cyan(format string, a ...interface{}) string  { return colorString(CyanFg, format, a...) }
func (color) White(format string, a ...interface{}) string { return colorString(WhiteFg, format, a...) }
func (color) DefaultColor(format string, a ...interface{}) string {
	return colorString(DefaultFg, format, a...)
}

func colorString(c ColorCode, format string, a ...interface{}) string {
	if c == DefaultFg {
		return fmt.Sprintf(format, a...)
	}
	return fmt.Sprintf("%v%vm%v%v", escape, c, fmt.Sprintf(format, a...), clear)
}

type monochrome struct{}

func (monochrome) Black(format string, a ...interface{}) string   { return fmt.Sprintf(format, a...) }
func (monochrome) Red(format string, a ...interface{}) string     { return fmt.Sprintf(format, a...) }
func (monochrome) Green(format string, a ...interface{}) string   { return fmt.Sprintf(format, a...) }
func (monochrome) Yellow(format string, a ...interface{}) string  { return fmt.Sprintf(format, a...) }
func (monochrome) Blue(format string, a ...interface{}) string    { return fmt.Sprintf(format, a...) }
func (monochrome) Magenta(format string, a ...interface{}) string { return fmt.Sprintf(format, a...) }
func (monochrome) Cyan(format string, a ...interface{}) string    { return fmt.Sprintf(format, a...) }
func (monochrome) White(format string, a ...interface{}) string   { return fmt.Sprintf(format, a...) }
func (monochrome) DefaultColor(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

func NewColor(enableColor bool) Color {
	if enableColor {
		term := os.Getenv("TERM")
		switch term {
		case "dumb", "":
			enableColor = false
			fmt.Println("Warning: your terminal doesn't support colors")
		}
	}
	if enableColor {
		return color{}
	} else {
		return monochrome{}
	}
}
