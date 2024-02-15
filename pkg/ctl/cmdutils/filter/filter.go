package filter

import (
	"fmt"
	"strings"

	"github.com/gobwas/glob"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	"k8s.io/apimachinery/pkg/util/sets"
)

// Filter holds filter configuration
type Filter struct {
	ExcludeAll bool // highest priority

	// include filters take precedence
	includeNames    sets.Set[string]
	includeGlobs    []glob.Glob
	rawIncludeGlobs []string

	excludeNames    sets.Set[string]
	excludeGlobs    []glob.Glob
	rawExcludeGlobs []string
}

// NewFilter returns a new initialized Filter
func NewFilter() Filter {
	return Filter{
		ExcludeAll:   false,
		includeNames: sets.New[string](),
		excludeNames: sets.New[string](),
	}
}

// AppendIncludeNames appends explicit names to the include filter
func (f *Filter) AppendIncludeNames(names ...string) { f.includeNames.Insert(names...) }

// AppendExcludeGlobs sets globs for exclusion rules
func (f *Filter) AppendExcludeGlobs(globExprs ...string) error {
	for _, expr := range globExprs {
		compiledExpr, err := glob.Compile(expr)
		if err != nil {
			return errors.Wrapf(err, "parsing glob filter %q", expr)
		}
		f.excludeGlobs = append(f.excludeGlobs, compiledExpr)
		f.rawExcludeGlobs = append(f.rawExcludeGlobs, expr)
	}
	return nil // exclude filter doesn't have to match anything, so we don't validate it
}

// AppendExcludeNames appends explicit names to the exclude filter
func (f *Filter) AppendExcludeNames(names ...string) { f.excludeNames.Insert(names...) }

func (*Filter) matchGlobs(name string, exprs []glob.Glob) bool {
	for _, compiledExpr := range exprs {
		if compiledExpr.Match(name) {
			return true
		}
	}
	return false
}

// hasIncludeRules returns true if the user has supplied inclusion globs or names
func (f *Filter) hasIncludeRules() bool {
	return len(f.includeGlobs) != 0
}

func (f *Filter) describeIncludeRules() string {
	rules := append(sets.List(f.includeNames), f.rawIncludeGlobs...)
	return strings.Join(rules, ",")
}

// hasExcludeRules returns true if the user has supplied exclusion globs or names
func (f *Filter) hasExcludeRules() bool {
	return len(f.excludeGlobs) != 0
}

func (f *Filter) describeExcludeRules() string {
	rules := append(sets.List(f.excludeNames), f.rawExcludeGlobs...)
	return strings.Join(rules, ",")
}

// Match given name against the filter and returns
// true or false if it has to be included or excluded
func (f *Filter) Match(name string) bool {
	if f.ExcludeAll {
		return false // force exclude
	}

	// Name overwrites
	if f.includeNames.Has(name) && !f.excludeNames.Has(name) {
		return true
	}

	if f.excludeNames.Has(name) {
		return false
	}

	hasIncludeRules := f.hasIncludeRules()
	hasExcludeRules := f.hasExcludeRules()

	if !hasIncludeRules && !hasExcludeRules {
		return true
	}

	// Exclusion override takes precedence
	if f.excludeNames.Has(name) {
		return false
	}

	if hasIncludeRules {
		if f.matchGlobs(name, f.includeGlobs) {
			if hasExcludeRules {
				// exclusion takes precedence
				if f.matchGlobs(name, f.excludeGlobs) {
					return false
				}
			}
			return true
		}

		// if there are include rules and it doesn't match then it must be excluded regardless of the exclude rules
		return false
	}

	// With only exclusion rules, everything that is not excluded is included
	if hasExcludeRules {
		// Overwrites by name take precedence
		if f.excludeNames.Has(name) {
			return false
		}
		if f.matchGlobs(name, f.excludeGlobs) {
			return false
		}
		return true
	}

	return false
}

// doMatchAll all names against the filter and return two sets of names - included and excluded
func (f *Filter) doMatchAll(names []string) (sets.Set[string], sets.Set[string]) {
	included, excluded := sets.New[string](), sets.New[string]()
	if f.ExcludeAll {
		for _, n := range names {
			excluded.Insert(n)
		}
		return included, excluded
	}
	for _, n := range names {
		if f.Match(n) {
			included.Insert(n)
		} else {
			excluded.Insert(n)
		}
	}
	return included, excluded
}

// doAppendIncludeGlobs sets globs for inclusion rules
func (f *Filter) doAppendIncludeGlobs(names []string, resource string, globExprs ...string) error {
	for _, expr := range globExprs {
		compiledExpr, err := glob.Compile(expr)
		if err != nil {
			return errors.Wrapf(err, "parsing glob filter %q", expr)
		}
		f.includeGlobs = append(f.includeGlobs, compiledExpr)
		f.rawIncludeGlobs = append(f.rawIncludeGlobs, expr)
	}
	return f.includeGlobsMatchAnything(names, resource)
}

func (f *Filter) doSetExcludeExistingFilter(names []string, resource string) error {
	uniqueNames := sets.List(sets.New[string](names...))
	f.excludeNames.Insert(uniqueNames...)
	for _, n := range uniqueNames {
		isAlsoIncluded := f.includeNames.Has(n) || f.matchGlobs(n, f.includeGlobs)
		if isAlsoIncluded {
			return fmt.Errorf("existing %s %q should be excluded, but matches include filter: %s", resource, n, f.describeIncludeRules())
		}
	}
	if len(uniqueNames) != 0 {
		logger.Info("%d existing %s(s) (%s) will be excluded", len(uniqueNames), resource, strings.Join(uniqueNames, ","))
	}
	return nil
}

func (f *Filter) includeGlobsMatchAnything(names []string, resource string) error {
	if len(f.includeGlobs) == 0 {
		return nil
	}
	for _, n := range names {
		if f.matchGlobs(n, f.includeGlobs) {
			return nil
		}
	}
	return fmt.Errorf("no %ss match include glob filter specification: %q", resource, strings.Join(f.rawIncludeGlobs, ","))
}

func (f *Filter) doLogInfo(resource string, included, excluded sets.Set[string]) {
	logMsg := func(subset sets.Set[string], status string) {
		count := subset.Len()
		list := strings.Join(sets.List(subset), ", ")
		subjectFmt := "%d %ss (%s) were %s"
		if count == 1 {
			subjectFmt = "%d %s (%s) was %s"
		}
		logger.Info(subjectFmt, count, resource, list, status+" (based on the include/exclude rules)")
	}

	if f.hasIncludeRules() {
		logger.Info("combined include rules: %s", f.describeIncludeRules())
		if included.Len() == 0 {
			logger.Info("no %ss present in the current set were included by the filter", resource)
		}
	}
	if included.Len() > 0 {
		logMsg(included, "included")
	}
	if f.hasExcludeRules() {
		logger.Info("combined exclude rules: %s", f.describeExcludeRules())
		if excluded.Len() == 0 {
			logger.Info("no %ss present in the current set were excluded by the filter", resource)
		}
	}
	if excluded.Len() > 0 {
		logMsg(excluded, "excluded")
	}
}
