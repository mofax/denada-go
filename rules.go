package denada

import "fmt"
import "strings"

type Cardinality int

const (
	Zero = iota
	Optional
	ZeroOrMore
	Singleton
	OneOrMore
)

type RuleInfo struct {
	ContextPath []string
	Context     RuleContext
	Name        string
	Cardinality Cardinality
}

func (r RuleInfo) checkCount(count int) error {
	switch r.Cardinality {
	case Zero:
		if count != 0 {
			return fmt.Errorf("Expected zero of rule %s, found %d", r.Name, count)
		}
	case Optional:
		if count > 1 {
			return fmt.Errorf("Expected at most 1 of rule %s, found %d", r.Name, count)
		}
	case Singleton:
		if count != 1 {
			return fmt.Errorf("Expected at exactly 1 of rule %s, found %d", r.Name, count)
		}
	case OneOrMore:
		if count == 0 {
			return fmt.Errorf("Expected at least 1 of rule %s, found %d", r.Name, count)
		}
	}
	return nil
}

func ParseRuleName(desc string) (rule RuleInfo, err error) {
	return ParseRule(desc, NullContext())
}

func ParseRule(desc string, context RuleContext) (rule RuleInfo, err error) {
	rule = RuleInfo{Cardinality: Zero}

	// Note a rule is of the form "myrule>childrule".  If no ">" is present,
	// child rules are assumed to be indicated by the "contents" of the current
	// rule.

	parts := strings.Split(desc, ">")
	str := desc
	path := []string{"."}

	if len(parts) == 0 {
		err = fmt.Errorf("Empty rule string")
		return
	} else if len(parts) == 2 {
		str = parts[0]
		path = strings.Split(parts[1], "/")
	} else if len(parts) > 2 {
		err = fmt.Errorf("Rule contains multiple child rule indicators (>)")
		return
	}

	// Shorthand notation
	if str[0] == '^' {
		path = []string{".."}
		str = str[1:]
	}

	rctxt, ferr := context.Find(path...)
	if ferr != nil {
		err = fmt.Errorf("Error finding %v: %v", path, ferr)
		return
	}
	rule.Context = rctxt
	rule.ContextPath = path

	l := len(str) - 1
	lastchar := str[l]
	if lastchar == '-' {
		rule.Cardinality = Zero
		str = str[0:l]
	} else if lastchar == '+' {
		rule.Cardinality = OneOrMore
		str = str[0:l]
	} else if lastchar == '*' {
		rule.Cardinality = ZeroOrMore
		str = str[0:l]
	} else if lastchar == '?' {
		rule.Cardinality = Optional
		str = str[0:l]
	} else {
		rule.Cardinality = Singleton
	}
	rule.Name = str

	return
}
