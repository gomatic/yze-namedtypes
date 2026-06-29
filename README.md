# yze-namedtypes

A [`yze`](https://github.com/gomatic/yze) analyzer (category `types`) enforcing the gomatic Go standard that function parameters use named domain types rather than bare primitives.

v1 flags parameters of non-method function declarations whose type is a bare predeclared primitive identifier (`int`, `string`, `bool`, …). A named domain type — whether a defined type (`type Count int`) or an alias (`type Celsius = float64`) — is exempt, as are the predeclared `error`/`any` interfaces and imported types. Method receivers (interface-satisfaction carve-outs) and composite types (`[]string`, …) are deferred.

- **Rule:** `yze/namedtypes`
- **Library:** exports `Analyzer` and `Registration` for the [`yze`](https://github.com/gomatic/yze) aggregator and [`stickler`](https://github.com/gomatic/stickler) runner.
- **Binary:** `cmd/yze-namedtypes` runs it standalone (`text`/`-json`, and as a `go vet -vettool`).

Built on the [`go-yze`](https://github.com/gomatic/go-yze) framework.
