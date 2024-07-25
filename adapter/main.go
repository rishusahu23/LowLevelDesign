package main

import "fmt"

type Printer interface {
	Print(msg string)
}

type LegacyPrinter struct {
}

func (lp *LegacyPrinter) OldPrint(msg string) {
	fmt.Println("legacy: ", msg)
}

type Adapter struct {
	legacyPrinter *LegacyPrinter
}

func (a *Adapter) Print(msg string) {
	a.legacyPrinter.OldPrint(msg)
}

var _ Printer = &Adapter{}

func main() {
	lg := &LegacyPrinter{}
	adapter := &Adapter{legacyPrinter: lg}
	var pr Printer = adapter
	pr.Print("hello")
}
