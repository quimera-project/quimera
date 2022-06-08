package storage

import (
	"fmt"
	"os"
	"os/user"
	"runtime"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/quimera-project/quimera/internal/check"
	"github.com/quimera-project/quimera/internal/env"
)

type storage struct {
	checks map[string]check.Checks
	statistics
	status int
}

var store *storage

// init initializes the Storage struct.
func init() {
	store = &storage{checks: make(map[string]check.Checks)}
	store.Time = time.Now()
	user, err := user.Current()
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
	store.User = user.Username
	store.Id = uuid.New().ID()
	store.Os = runtime.GOOS
	store.Arch = runtime.GOARCH
}

// Checks returns a map of Checks stored by their corresponding category.
func Checks() map[string]check.Checks {
	return store.checks
}

// ChecksByCategory returns a sorted Checks slice of the supplied category.
func ChecksByCategory(cat string) (checks check.Checks) {
	checks = store.checks[cat]
	sort.Sort(checks)
	return
}

// SortedCategories returns a sorted slice of the Checks categories names.
func SortedCategories() (categories []string) {
	for k := range store.checks {
		categories = append(categories, k)
	}
	return
}

// Status returns the current Quimera status.
func Status() int {
	return store.status
}

// Stats returns the current Quimera statistics.
func Stats() statistics {
	return store.statistics
}

// Categories returns the different number of the executed categories checks.
func Categories() map[string]int {
	count := map[string]int{}
	for cat, checks := range store.checks {
		count[cat] = checks.Len()
	}
	return count
}

// Add adds a check to the Storage.
//
// If the check failed and the SkipFailed is provided by the user, the check is not saved and ignored instead.
func Add(c *check.Check) {
	store.Total++
	if c.Failed {
		store.Failed++
		if env.Config.SkipFailed {
			return
		}
	}
	store.checks[c.Category] = append(store.checks[c.Category], c)
}

// AddError adds an error to the Storage.
func AddError() {
	store.Total++
	store.Errors++
}

// RenderStats returns the rendered string of the statistics.
func RenderStats() string {
	return store.statistics.render()
}

// finish calculates the Stats duration.
func Finish() {
	store.Duration = time.Since(store.Time)
}
