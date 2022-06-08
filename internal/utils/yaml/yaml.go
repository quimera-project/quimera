package yaml

import (
	"fmt"
	"sort"

	qfs "github.com/quimera-project/quimera/internal/utils/fs"

	"github.com/quimera-project/quimera/internal/check"
	"github.com/quimera-project/quimera/internal/env"
	"gopkg.in/yaml.v3"
)

// AllChecks returns a slice containing all unmarshalled checks and any error encountered.
func AllChecks() (all check.Checks, err error) {
	categories, err := qfs.Checks.GetCategories()
	if err != nil {
		return nil, fmt.Errorf("reading categories: %v", err)
	}
	for _, cat := range categories {
		meta, err := qfs.Checks.ReadMeta(cat)
		if err != nil {
			return nil, fmt.Errorf("reading meta.yaml file from %s category: %v", cat, err)
		}

		var checks map[string]*check.Check
		err = yaml.Unmarshal(meta, &checks)
		if err != nil {
			return nil, fmt.Errorf("unmarshalling %s meta.yaml: %v", cat, err)
		}

		for name, c := range checks {
			if err := c.Prepare(cat, name); err != nil {
				return nil, err
			}
			all = append(all, c)
		}
	}
	return
}

// AllChecksByPriority returns a slice containing all unmarshalled checks in slices sorted by priority and any error encountered.
//
// It drops all checks not containing the category, tags and/or operating system supplied by the user.
func AllChecksByPriority() (all []check.Checks, err error) {
	var priority = map[int]check.Checks{}

	categories, err := qfs.Checks.GetCategories()
	if err != nil {
		return nil, fmt.Errorf("reading categories: %v", err)
	}

	for _, cat := range categories {
		if !validCategory(cat) {
			continue
		}

		meta, err := qfs.Checks.ReadMeta(cat)
		if err != nil {
			return nil, fmt.Errorf("reading meta.yaml file from %s category: %v", cat, err)
		}

		var checks map[string]*check.Check
		err = yaml.Unmarshal(meta, &checks)
		if err != nil {
			return nil, fmt.Errorf("unmarshalling %s meta.yaml: %v", cat, err)
		}

		for name, c := range checks {
			if !validTags(c.Tags) || !validOS(c.Os) {
				continue
			}
			if err := c.Prepare(cat, name); err != nil {
				return nil, err
			}
			priority[c.Priority] = append(priority[c.Priority], c)
		}
	}

	var keys = []int{}
	for k := range priority {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, v := range keys {
		all = append(all, priority[v])
	}

	return all, nil
}

// validCategory returns true if the category is valid.
//
// A category is valid if it is specified by the user or if the user hasn't especified any.
func validCategory(selected string) bool {
	if env.Config.Categories == nil {
		return true
	}
	for _, cat := range env.Config.Categories {
		if selected == cat {
			return true
		}
	}
	return false
}

// validTags return true if the tag is valid.
//
// A tag is valid if it is specified by the user or if the user hasn't especified any.
func validTags(selected []string) bool {
	if env.Config.Tags == nil {
		return true
	}
	for _, tag := range env.Config.Tags {
		for _, t := range selected {
			if t == tag {
				return true
			}
		}
	}
	return false
}

// validOS return true if the tested OS is valid.
//
// A tested OS is valid if it is specified by the user or if the user hasn't especified any.
func validOS(selected []string) bool {
	if env.Config.OS == nil {
		return true
	}
	for _, os := range env.Config.OS {
		for _, o := range selected {
			if o == os {
				return true
			}
		}
	}
	return false
}
