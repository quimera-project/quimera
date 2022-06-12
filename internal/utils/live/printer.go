package live

import "fmt"

type printer interface {
	Infof(format string, a ...any)
	Warningf(format string, a ...any)
	Errorf(format string, a ...any)
	Fatalf(format string, a ...any)
}

var Printer struct {
	printer
	Color   *color
	Spinner *spinner
}

// Init initializes the Printer struct.
//
// It selects whether the printer will be a color or a spinner struct.
func Init(spinner bool) {
	Printer.Color = newColor()
	Printer.Spinner = newSpinner()
	if spinner {
		Printer.printer = Printer.Spinner
		// Printer.Spinner.model.Start()
	} else {
		Printer.printer = Printer.Color
	}
}

// SpinnerStart starts the spinner.
func SppinerStart() {
	Printer.Spinner.model.Start()
}

// SpinnerMsg formats according to a format specifier and change the sppiner message.
func SpinnerMsg(format string, a ...any) {
	Printer.Spinner.model.Message(fmt.Sprintf(format, a...))
}

// SpinnerSuffixf formats according to a format specifier and change the sppiner suffix.
func SpinnerSuffixf(format string, a ...any) {
	Printer.Spinner.suffixf(format, a...)
}

// SpinnerStopf formats according to a format specifier and change the sppiner suffix.
//
// It also stops the spinner.
func SpinnerStopf(format string, a ...any) {
	Printer.Spinner.stopf(format, a...)
	Printer.printer = Printer.Color
}

// SpinnerChao stops the spinner with a "Chao!" message.
func SpinnerChao() {
	Printer.Spinner.chao()
}

// SpinnerAdd adds one check to the crafting variable and calls Infof to reflect the remaining checks.
func SpinnerAdd() {
	Printer.Spinner.addCheck()
}

// SpinnerDone subtracts one check to the crafting variable and calls Infof to reflect the remaining checks.
func SpinnerDone() {
	Printer.Spinner.doneCheck()
}
