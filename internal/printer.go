package internal

import (
	"github.com/ctrsploit/sploit-spec/pkg/printer"
)

// Printer default equal to Text, will be overwritten if --colorful is set
var Printer = printer.GetPrinter(printer.TypeText)

func Print(printers ...printer.Printer) (s string) {
	s = printer.Print(Printer, printers...)
	return
}
