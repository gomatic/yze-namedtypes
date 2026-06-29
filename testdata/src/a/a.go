package a

import "io"

// Count is a defined named domain type.
type Count int

// Celsius is a type alias to a primitive; an alias is itself a named domain type
// (the exact thing this analyzer exists to encourage) and must not be flagged.
type Celsius = float64

// withBare takes a bare primitive parameter and must be flagged.
func withBare(a int) {} // want `parameter type int is a bare primitive; define a named domain type`

// withBareString is flagged on its bare string primitive.
func withBareString(s string) {} // want `parameter type string is a bare primitive; define a named domain type`

// withBareBool is flagged on its bare bool primitive.
func withBareBool(b bool) {} // want `parameter type bool is a bare primitive; define a named domain type`

// withBareByte is flagged: byte is a predeclared primitive alias and its
// underlying type is a primitive.
func withBareByte(b byte) {} // want `parameter type byte is a bare primitive; define a named domain type`

// withNamed takes a defined named domain type and must not be flagged.
func withNamed(c Count) {}

// withAlias takes an alias to a primitive; the alias is a named domain type and
// must not be flagged.
func withAlias(t Celsius) {}

// withError takes the predeclared error interface (not a primitive); not flagged.
func withError(e error) {}

// withAny takes the predeclared any (an interface, not a primitive); not flagged.
func withAny(v any) {}

// withReader takes an imported interface type (a selector, not an identifier);
// not flagged.
func withReader(r io.Reader) {}

// withSlice takes a composite (non-identifier) type; deferred in v1, not flagged.
func withSlice(s []string) {}

// noParams has no parameters and must not be flagged.
func noParams() {}

// T carries a method below.
type T struct{}

// method has a pointer/value receiver; methods are deferred in v1 and the bare
// primitive parameter must not be flagged.
func (T) method(a int) {}
