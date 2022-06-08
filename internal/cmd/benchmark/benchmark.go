package benchmark

import (
	"fmt"
	"sort"

	"github.com/quimera-project/quimera/internal/check"
	"github.com/quimera-project/quimera/internal/storage"
	"github.com/quimera-project/quimera/internal/utils/live"
	"github.com/quimera-project/quimera/internal/utils/yaml"
)

func Test() {
	var (
		benchmark = map[float64]*check.Check{}
	)
	live.SppinerStart()
	checks, err := yaml.AllChecksByPriority()
	if err != nil {
		live.Printer.Fatalf("getting all checks: %v", err)
	}
	for i := range checks {
		live.SpinnerSuffixf("Crafting priority %d", checks[i][0].Priority)
		n := len(checks[i])
		arr := make([]check.Crafter, n)

		for j := 0; j < n; j++ {
			arr[j] = *check.NewCrafter()
			go arr[j].Craft(checks[i][j])
		}
		for j := 0; j < n; j++ {
			for {
				tick, ok := <-arr[j].Ch
				if ok {
					benchmark[tick.Check.Duration.Seconds()] = tick.Check
				} else {
					break
				}
			}
		}
	}
	live.SpinnerStopf("All checks have been crafted")
	storage.Finish()
	var times []float64
	for k := range benchmark {
		times = append(times, k)
	}
	sort.Float64s(times)
	for i, t := range times {
		if t > 10 {
			fmt.Printf("(%d) [%f] %s\n", live.Printer.Color.Gray(12-1, i), live.Printer.Color.Red(t), live.Printer.Color.Red(benchmark[t].Title))
		} else if t > 5 {
			fmt.Printf("(%d) [%f] %s\n", live.Printer.Color.Gray(12-1, i), live.Printer.Color.Yellow(t), live.Printer.Color.Yellow(benchmark[t].Title))
		} else {
			fmt.Printf("(%d) [%f] %s\n", live.Printer.Color.Gray(12-1, i), live.Printer.Color.Green(t), live.Printer.Color.Green(benchmark[t].Title))
		}
	}
}
