package internal

import (
	"github.com/ctrsploit/sploit-spec/pkg/printer"
)

// Printer default equal to Text, will be overwritten if --colorful is set
var Printer = printer.GetPrinter(printer.TypeText)

func Print(printers ...printer.Printer) (s string) {
	for _, i := range printers {
		if !i.IsEmpty() {
			s += Printer(i) + "\n"
		}
	}
	return
}
