package denada

import "fmt"

import "github.com/bitly/go-simplejson"

type Modifications map[string]*simplejson.Json

type Element struct {
	/* Common to all elements */
	Qualifiers    []string
	Name          string
	Description   string
	Modifications Modifications
	Contents      ElementList      // Used by definitions
	Value         *simplejson.Json // Used by declarations

	rulepath   string
	rule       string
	definition bool
}

func NewDefinition(name string, desc string, qualifiers ...string) *Element {
	return &Element{
		Qualifiers:    qualifiers,
		Name:          name,
		Description:   desc,
		Modifications: Modifications{},
		Contents:      ElementList{},
		Value:         nil,
		definition:    true,
	}
}

func NewDeclaration(name string, desc string,
	qualifiers ...string) *Element {
	return &Element{
		Qualifiers:    qualifiers,
		Name:          name,
		Description:   desc,
		Modifications: Modifications{},
		Contents:      ElementList{},
		Value:         nil,
		definition:    false,
	}
}

// This checks whether a given element has EXACTLY the listed qualifiers (in the exact order)
func (e Element) HasQualifiers(quals ...string) bool {
	if len(quals) != len(e.Qualifiers) {
		return false
	}
	for i, q := range e.Qualifiers {
		if quals[i] != q {
			return false
		}
	}
	return true
}

func (e Element) Unparse(rules bool) string {
	return UnparseElement(e, rules)
}

func (e Element) Clone() *Element {
	// TODO: Clone modifications and qualifiers

	children := []*Element{}
	if e.Contents != nil {
		children = append(children, e.Contents...)
	} else {
		children = nil
	}

	return &Element{
		Qualifiers:    e.Qualifiers,
		Name:          e.Name,
		Description:   e.Description,
		Modifications: e.Modifications,
		Contents:      children,
		Value:         e.Value,
		rule:          e.rule,
		definition:    e.definition,
	}
}

func (e Element) String() string {
	ret := ""
	for _, q := range e.Qualifiers {
		ret += q + " "
	}
	ret += e.Name

	if e.IsDefinition() {
		return fmt.Sprintf("%s { ... }", ret)
	} else {
		if e.Value != nil {
			return fmt.Sprintf("%s = %v;", ret, e.Value)
		} else {
			return fmt.Sprintf("%s;", ret)
		}
	}
}

func (e Element) Rule() string {
	return e.rule
}

func (e Element) RulePath() string {
	return e.rulepath
}

func (e Element) IsDefinition() bool {
	return e.definition
}

func (e Element) IsDeclaration() bool {
	return !e.definition
}

func equalValues(l *simplejson.Json, r *simplejson.Json) (bool, error) {
	lbytes, err := l.Encode()
	if err != nil {
		return false, err
	}
	rbytes, err := r.Encode()
	if err != nil {
		return false, err
	}
	return string(lbytes) == string(rbytes), nil
}

func (e *Element) Append(children ...*Element) error {
	if e.IsDefinition() {
		e.Contents = append(e.Contents, children...)
		return nil
	} else {
		return fmt.Errorf("Attempted to append elements to a declaration")
	}
}

func (e Element) Equals(o Element) error {
	// Check that they have the same number of qualifiers
	if len(e.Qualifiers) != len(o.Qualifiers) {
		return fmt.Errorf("Length mismatch (%d vs. %d)",
			len(e.Qualifiers), len(o.Qualifiers))
	}

	// Then check that each qualifier is identical (and in identical order)
	for i, q := range e.Qualifiers {
		if q != o.Qualifiers[i] {
			return fmt.Errorf("Qualifier mismatch: %s vs %s", q, o.Qualifiers[i])
		}
	}

	// Next, check that they have the same name
	if e.Name != o.Name {
		return fmt.Errorf("Name mismatch: %s vs %s", e.Name, o.Name)
	}

	// And then the same description
	if e.Description != o.Description {
		return fmt.Errorf("Description mismatch: %s vs %s", e.Description, o.Description)
	}

	// Now we check the modifications to make sure that all keys match and that the
	// value for each key is identical between both sets of modifications
	for k, v := range e.Modifications {
		ov, exists := o.Modifications[k]
		if !exists {
			return fmt.Errorf("Mismatch in modification for key %s missing from argument", k)
		}
		eq, err := equalValues(v, ov)
		if err != nil {
			return err
		}
		if !eq {
			return fmt.Errorf("Mismatch in value for key %s: %v vs %v", k, v, ov)
		}
	}
	for k, _ := range o.Modifications {
		_, exists := e.Modifications[k]
		if !exists {
			return fmt.Errorf("Mismatch in modification for key %s missing from object", k)
		}
	}

	err := e.Contents.Equals(o.Contents)
	if err != nil {
		return fmt.Errorf("Error in child elements comparing %v with %v: %v",
			e, o, err)
	}

	if e.Value != nil && o.Value != nil {
		// If they both have values, make sure they are equal
		eq, err := equalValues(e.Value, o.Value)
		if err != nil {
			return err
		}
		if !eq {
			return fmt.Errorf("Mismatch in values: %v vs %v", e.Value, o.Value)
		}
	} else {
		// If they don't both have values, then make sure that both are nil
		if e.Value != nil || o.Value != nil {
			return fmt.Errorf("Mismatch in value (one has a value, the other doesn't")
		}
	}
	// If we get here, nothing was unequal
	return nil
}

func (e *Element) SetModification(key string, data interface{}) {
	e.Modifications[key] = makeJson(data)
}

func (e *Element) SetValue(data interface{}) *simplejson.Json {
	val := makeJson(data)
	e.Value = val
	return val
}

func (e *Element) StringValueOf(defval string) string {
	if e == nil {
		return defval
	}
	if e.Value == nil {
		return defval
	}
	s, err := e.Value.String()
	if err != nil {
		return defval
	}
	return s
}

func (e *Element) FirstNamed(name string) *Element {
	if e.definition {
		return e.Contents.FirstNamed(name)
	} else {
		return nil
	}
}
