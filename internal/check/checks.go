package check

// Checks represents a checks pointers slice.
type Checks []*Check

// Len returns the corresponding checks slice length.
func (cs Checks) Len() int { return len(cs) }

// Less returns true if the check title of the checks slice on the i position is lower than the one in the j position. Otherwise it returns false.
func (cs Checks) Less(i, j int) bool { return cs[i].Title < cs[j].Title }

// Swap swaps the checks on the i and j checks slice positions.
func (cs Checks) Swap(i, j int) { cs[i], cs[j] = cs[j], cs[i] }

// String returns the title of the corresponding check on the i position.
func (cs Checks) String(i int) string { return cs[i].Title }
