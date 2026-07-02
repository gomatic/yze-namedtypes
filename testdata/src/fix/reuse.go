package fix

import "other"

// FilePath is an existing named domain type; the reuse fixes below adopt it
// instead of minting a skeleton.
type FilePath string

// writeRoot mirrors the gomatic/go-rewrite shape: its only call site already
// converts from FilePath, so the fix retypes the parameter to FilePath,
// converts the body uses back to the primitive, and drops the call-site
// conversion.
func writeRoot(name string) string { // want `parameter type string is a bare primitive; define a named domain type`
	if name == "" {
		return "."
	}
	return name + "/root"
}

func writeCaller(path FilePath) string {
	return writeRoot(string(path))
}

// Label is reused by tagged even though a second call site passes an untyped
// constant; the constant argument is wrapped in Label.
type Label string

func tagged(label string) string { // want `parameter type string is a bare primitive; define a named domain type`
	return "#" + label
}

func tagCaller(l Label) string {
	return tagged(string(l))
}

func tagDefault() string {
	return tagged("general")
}

// Retries is reused by two different functions in one pass: reuse mints no
// name, so the per-pass minted-name dedupe does not apply and both
// diagnostics carry fixes.
type Retries int

func waitFor(attempts int) int { // want `parameter type int is a bare primitive; define a named domain type`
	return attempts * 100
}

func retryLoop(times int) int { // want `parameter type int is a bare primitive; define a named domain type`
	return times + 1
}

func retriesCaller(r Retries) int {
	return waitFor(int(r)) + retryLoop(int(r))
}

// Owner and Repo feed slug's two flagged parameters through ONE call
// expression; each reuse fix edits only its own argument range, so the two
// fixes never conflict.
type Owner string

// Repo is the reuse candidate for slug's second parameter.
type Repo string

func slug(owner string, repo string) string { // want `string is a bare primitive` `string is a bare primitive`
	return owner + "/" + repo
}

func slugCaller(o Owner, r Repo) string {
	return slug(string(o), string(r))
}

// City and Town disagree on the reuse candidate for locate, so the skeleton
// fix applies instead.
type City string

// Town is the second, disagreeing reuse candidate for locate.
type Town string

func locate(place string) string { // want `parameter type string is a bare primitive; define a named domain type`
	return place
}

func cityCaller(c City) string {
	return locate(string(c))
}

func townCaller(t Town) string {
	return locate(string(t))
}

// Pathish is an alias, not a defined type; a conversion from it does not
// trigger reuse and the skeleton fix applies.
type Pathish = string

func fromAlias(p string) string { // want `parameter type string is a bare primitive; define a named domain type`
	return p
}

func aliasCaller(p Pathish) string {
	return fromAlias(string(p))
}

// Fahrenheit's underlying type is float64, not int, so the conversion int(f)
// does not trigger reuse and the skeleton fix applies.
type Fahrenheit float64

func chill(deg int) int { // want `parameter type int is a bare primitive; define a named domain type`
	return deg - 1
}

func chillCaller(f Fahrenheit) int {
	return chill(int(f))
}

// Tag is generic; its instances cannot be spelled as a bare parameter type,
// so reuse is skipped and the skeleton fix applies.
type Tag[T any] string

func withTag(g string) string { // want `parameter type string is a bare primitive; define a named domain type`
	return g
}

func tagUser(t Tag[int]) string {
	return withTag(string(t))
}

// span's only conversion argument comes from a type declared in another
// package; reuse is same-package only, so the skeleton fix applies.
func span(width int) int { // want `parameter type int is a bare primitive; define a named domain type`
	return width * 2
}

func spanCaller(d other.Distance) int {
	return span(int(d))
}

// rawMixed's second call passes a non-constant expression that is not a
// conversion from a named type, so reuse cannot hold and the skeleton fix
// wraps every call site.
func rawMixed(s string) string { // want `parameter type string is a bare primitive; define a named domain type`
	return s
}

func render(fp FilePath) string {
	return string(fp)
}

func rawCaller(fp FilePath) string {
	return rawMixed(string(fp)) + rawMixed(render(fp))
}

// Marker would be reused by mark, but the constant call site below sits where
// a local declaration shadows Marker, so the wrap cannot reference it and the
// skeleton fix applies.
type Marker string

func mark(m string) string { // want `parameter type string is a bare primitive; define a named domain type`
	return m
}

func markCaller(v Marker) string {
	return mark(string(v))
}

func markShadowed() string {
	type Marker struct{}
	_ = Marker{}
	return mark("x")
}
