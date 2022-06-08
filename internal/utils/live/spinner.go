package live

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/quimera-project/quimera/internal/env"
	"github.com/theckman/yacspin"
)

type spinner struct {
	model    *yacspin.Spinner
	mu       sync.Mutex
	crafting uint
}

// newSpinner creates a new spinner and returns its pointer.
func newSpinner() *spinner {
	rand.Seed(time.Now().UnixNano())
	var colors, stopColors, stopFailColors []string
	if !env.Config.Uncolorized {
		colors, stopColors, stopFailColors = []string{"fgYellow"}, []string{"fgGreen"}, []string{"fgRed"}
	} else {
		colors, stopColors, stopFailColors = nil, nil, nil
	}
	model, _ := yacspin.New(yacspin.Config{
		Writer:            os.Stderr,
		Frequency:         150 * time.Millisecond,
		Colors:            colors,
		ColorAll:          !env.Config.Uncolorized,
		CharSet:           yacspin.CharSets[rand.Intn(90)],
		SuffixAutoColon:   true,
		Message:           "...",
		Suffix:            " ",
		StopCharacter:     "✓",
		StopFailCharacter: "✗",
		StopColors:        stopColors,
		StopFailColors:    stopFailColors,
		StopFailMessage:   "Something failed",
	})
	return &spinner{model: model}
}

// Infof formats according to a format specifier and writes to stderr trough the spinner as an info message.
func (s *spinner) Infof(format string, a ...any) {
	if !env.Config.Uncolorized {
		s.model.StopFailCharacter("⊡")
		s.model.StopFailColors("fgBlue")
	}
	s.model.StopFailMessage(fmt.Sprintf(format, a...))
	s.model.StopFail()
	s.model.Start()
}

// Warningf formats according to a format specifier and writes to stderr trough the spinner as a warning message.
func (s *spinner) Warningf(format string, a ...any) {
	if !env.Config.Uncolorized {
		s.model.StopFailCharacter("☢")
		s.model.StopFailColors("fgYellow")
	}
	s.model.StopFailMessage(fmt.Sprintf(format, a...))
	s.model.StopFail()
	s.model.Start()
}

// Errorf formats according to a format specifier and writes to stderr trough the spinner as an error message.
func (s *spinner) Errorf(format string, a ...any) {
	if !env.Config.Uncolorized {
		s.model.StopFailCharacter("✗")
		s.model.StopFailColors("fgRed")
	}
	s.model.StopFailMessage(fmt.Sprintf(format, a...))
	s.model.StopFail()
	s.model.Start()
}

// Fatalf formats according to a format specifier and writes to stderr trough the spinner as an error message.
//
// Quimera exits with a status code of 1.
func (s *spinner) Fatalf(format string, a ...any) {
	if !env.Config.Uncolorized {
		s.model.StopFailCharacter("✗")
		s.model.StopFailColors("fgRed")
	}
	s.model.StopFailMessage(fmt.Sprintf(format, a...))
	s.model.StopFail()
	os.Exit(1)
}

// suffixf formats according to a format specifier and change the sppiner suffix.
func (s *spinner) suffixf(format string, a ...any) {
	s.model.Suffix(" " + fmt.Sprintf(format, a...))
}

// stopf formats according to a format specifier and change the sppiner suffix.
//
// It also stops the spinner.
func (s *spinner) stopf(format string, a ...any) {
	s.model.Suffix(" " + fmt.Sprintf(format, a...))
	s.model.Stop()
}

// chao stops the spinner with a "Chao!" message.
func (s *spinner) chao() {
	s.model.Suffix(" ")
	s.model.StopFailMessage("Chao!")
	s.model.StopFail()
}

// addCheck adds one check to the crafting variable and calls Infof to reflect the remaining checks.
func (s *spinner) addCheck() {
	s.mu.Lock()
	s.crafting++
	s.mu.Unlock()
	s.model.Message(fmt.Sprintf("%d checks remaining", s.crafting))
}

// doneCheck subtracts one check to the crafting variable and calls Infof to reflect the remaining checks.
func (s *spinner) doneCheck() {
	s.mu.Lock()
	s.crafting--
	s.mu.Unlock()
	s.model.Message(fmt.Sprintf("%d checks remaining", s.crafting))
}
