package cmdutils

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
	includeNames    sets.String
	includeGlobs    []glob.Glob
	rawIncludeGlobs []string

	excludeNames    sets.String
	excludeGlobs    []glob.Glob
	rawExcludeGlobs []string
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

func (f *Filter) hasIncludeRules() bool {
	return f.includeNames.Len()+len(f.includeGlobs) != 0
}

func (f *Filter) describeIncludeRules() string {
	rules := append(f.includeNames.List(), f.rawIncludeGlobs...)
	return fmt.Sprintf("%s", strings.Join(rules, ","))
}

func (f *Filter) hasExcludeRules() bool {
	return f.excludeNames.Len()+len(f.excludeGlobs) != 0
}

func (f *Filter) describeExcludeRules() string {
	rules := append(f.excludeNames.List(), f.rawExcludeGlobs...)
	return fmt.Sprintf("%s", strings.Join(rules, ","))
}

// Match given name against the filter and returns
// true or false if it has to be included or excluded
func (f *Filter) Match(name string) bool {
	if f.ExcludeAll {
		return false // force exclude
	}

	hasIncludeRules := f.hasIncludeRules()
	hasExcludeRules := f.hasExcludeRules()

	if !hasIncludeRules && !hasExcludeRules {
		return true // empty rules - include
	}

	mustInclude := false // use this override when rules overlap

	if hasIncludeRules {
		mustInclude = f.includeNames.Has(name)
		if f.matchGlobs(name, f.includeGlobs) {
			mustInclude = true
		}
		if !hasExcludeRules {
			// empty exclusion rules - explicit inclusion mode
			return mustInclude
		}
	}

	if hasExcludeRules {
		exclude := f.excludeNames.Has(name)
		if f.matchGlobs(name, f.excludeGlobs) {
			exclude = true
		}
		if exclude && !mustInclude {
			// exclude, unless overridden by an inclusion rule
			return false
		}
	}

	return true // biased to include
}

// doMatchAll all names against the filter and return two sets of names - included and excluded
func (f *Filter) doMatchAll(names []string) (sets.String, sets.String) {
	included, excluded := sets.NewString(), sets.NewString()
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
	uniqueNames := sets.NewString(names...).List()
	f.excludeNames.Insert(uniqueNames...)
	for _, n := range uniqueNames {
		isAlsoIncluded := f.includeNames.Has(n)
		if f.matchGlobs(n, f.includeGlobs) {
			isAlsoIncluded = true
		}
		if isAlsoIncluded {
			return fmt.Errorf("existing %s %q should be excluded, but matches include filter: %s", resource, n, f.describeIncludeRules())
		}
	}
	if len(uniqueNames) != 0 {
		logger.Info("%d %s(s) that already exist (%s) will be excluded", len(uniqueNames), resource, strings.Join(uniqueNames, ","))
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

func (f *Filter) doLogInfo(resource string, names []string) {
	logMsg := func(subset sets.String, status string) {
		count := subset.Len()
		list := strings.Join(subset.List(), ", ")
		subjectFmt := "%d %ss (%s) were %s"
		if count == 1 {
			subjectFmt = "%d %s (%s) was %s"
		}
		logger.Info(subjectFmt, count, resource, list, status+" (based on the include/exclude rules)")
	}

	included, excluded := f.doMatchAll(names)
	if f.hasIncludeRules() {
		logger.Info("combined include rules: %s", f.describeIncludeRules())
		if included.Len() == 0 {
			logger.Info("no %ss were included by the filter", resource)
		}
	}
	if included.Len() > 0 {
		logMsg(included, "included")
	}
	if f.hasExcludeRules() {
		logger.Info("combined exclude rules: %s", f.describeExcludeRules())
		if excluded.Len() == 0 {
			logger.Info("no %ss were excluded by the filter", resource)
		}
	}
	if excluded.Len() > 0 {
		logMsg(excluded, "excluded")
	}
}
