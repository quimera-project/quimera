package check

import (
	"github.com/quimera-project/quimera/internal/utils/live"
)

// Task represents the Crafter channel type.
type Task struct {
	Check *Check
	Err   error
}

// Crafter represents a crafter used to craft checks in a concurrency way.
type Crafter struct {
	Ch chan *Task
}

// NewCrafter creates a new crafter and returns its pointer.
func NewCrafter() *Crafter {
	return &Crafter{Ch: make(chan *Task, 1)}
}

// Craft sends through the crafter channel the crafted error and any error encountered.
func (cr *Crafter) Craft(c *Check) {
	defer close(cr.Ch)
	defer live.SpinnerDone()
	live.SpinnerAdd()
	cr.Ch <- &Task{Check: c, Err: c.Craft()}
}
