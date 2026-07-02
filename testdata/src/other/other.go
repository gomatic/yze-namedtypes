// Package other declares a named type outside the fix package, so the fix
// corpus can prove that reuse is same-package only.
package other

// Distance is a defined type with a primitive underlying type, but it lives
// outside the flagged function's package, so it is never reused.
type Distance int
