# Tagged URN - Go Implementation

Go implementation of Tagged URN with strict validation, pattern matching, and graded specificity comparison.

## Features

- **Strict Rule Enforcement** - Follows exact same rules as Rust, JavaScript, and Objective-C implementations
- **Case Insensitive** - All input normalized to lowercase (except quoted values)
- **Tag Order Independent** - Canonical alphabetical sorting
- **Special Pattern Values** - `*` (must-have-any), `?` (unspecified), `!` (must-not-have)
- **Value-less Tags** - Tags without values (`tag`) mean must-have-any (`tag=*`)
- **Graded Specificity** - Exact values score higher than wildcards
- **JSON Serialization** - Full JSON marshal/unmarshal support
- **Zero Dependencies** - Only standard library (testify for tests only)

## Installation

```bash
go get github.com/machinefabric/tagged-urn-go
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"

    taggedurn "github.com/machinefabric/tagged-urn-go"
)

func main() {
    // Parse a URN
    urn, err := taggedurn.NewTaggedUrnFromString("cap:generate;ext=pdf")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Has generate marker:", urn.HasMarkerTag("generate"))  // true
    fmt.Println("Canonical:", urn.ToString())                            // "cap:ext=pdf;generate"

    // Build a URN
    built := taggedurn.NewTaggedUrnBuilder("cap").
        Marker("extract").
        Tag("format", "pdf").
        Build()

    // Check matching
    pattern, _ := taggedurn.NewTaggedUrnFromString("cap:generate")
    conforms, err := urn.ConformsTo(pattern)
    if err != nil {
        log.Fatal(err)
    }
    if conforms {
        fmt.Println("URN conforms to pattern")
    }

    // Get specificity
    fmt.Println("Specificity:", urn.Specificity())
}
```

## API Reference

### TaggedUrn

| Function/Method | Description |
|-----------------|-------------|
| `NewTaggedUrnFromString(s)` | Parse URN from string |
| `NewTaggedUrnFromTags(prefix, tags)` | Create from prefix and tag map |
| `Empty(prefix)` | Create empty URN with prefix |
| `GetTag(key)` | Get value for a tag key |
| `HasTag(key, value)` | Check if tag exists with value |
| `WithTag(key, value)` | Return new URN with tag added/updated |
| `WithoutTag(key)` | Return new URN with tag removed |
| `ConformsTo(pattern)` | Check if URN conforms to a pattern |
| `Accepts(instance)` | Check if URN (as pattern) accepts an instance |
| `CanHandle(request)` | Check if URN can handle a request |
| `Specificity()` | Get graded specificity score |
| `SpecificityTuple()` | Get (exact, mustHaveAny, mustNot) counts |
| `IsMoreSpecificThan(other)` | Compare specificity with another URN |
| `ToString()` | Get canonical string representation |
| `Hash()` | Get SHA256 hash of canonical form |

### TaggedUrnBuilder

| Method | Description |
|--------|-------------|
| `NewTaggedUrnBuilder(prefix)` | Create builder with prefix |
| `Tag(key, value)` | Add or update a tag (chainable) |
| `Build()` | Build the URN |
| `BuildWithValidation()` | Build with validation (returns error) |

## Matching Semantics

Every form means the set of states its key may be in, on either side of a
comparison, and two URNs are compared by those sets. Three questions:

| Question | Go |
|---|---|
| Is everything `a` describes described by `b`? (a guarantee) | `a.ConformsTo(b)` |
| Could `a` and `b` be about the same thing? (a possibility) | `a.Meets(b)` |
| Does `a`, a complete thing — what it does not mention it does not have — fit `b`? | `a.Satisfies(b)` |

What an instance must say to be guaranteed to fit a pattern:

| Pattern | Instance omits K | Instance `K=v` | Instance `K=x` (x≠v) | Instance `K` (any value) |
|---------|------------------|----------------|----------------------|--------------------------|
| (missing) or `?K` | fits | fits | fits | fits |
| `!K` | no — an omission promises nothing | no | no | no |
| `K` (=`K=*`) | no | fits | fits | fits |
| `K=v` | no | fits | no | no — "some value" is not `v` |

A complete thing that omits `K` does fit `!K`: use `Satisfies` where the
left side is what something is — a value's media, a cap's own tags — rather than
what something is declared to take or give. "Some value" could be `v`:
`Meets` says so, and is the answer a search wants; a route is only ever
taken on a guarantee.

The rules are proved in `../formal` (Lean), and this package runs code generated
from them.

## Graded Specificity

| Value Type | Score |
|------------|-------|
| Exact value (`K=v`) | 3 |
| Must-have-any (`K=*`) | 2 |
| Must-not-have (`K=!`) | 1 |
| Unspecified (`K=?`) or missing | 0 |

## Error Codes

| Code | Constant | Description |
|------|----------|-------------|
| 1 | `ErrorInvalidFormat` | Empty or malformed URN |
| 2 | `ErrorEmptyTag` | Empty key or value component |
| 3 | `ErrorInvalidCharacter` | Disallowed character in key/value |
| 4 | `ErrorInvalidTagFormat` | Tag not in key=value format |
| 5 | `ErrorMissingPrefix` | URN does not start with prefix |
| 6 | `ErrorDuplicateKey` | Same key appears twice |
| 7 | `ErrorNumericKey` | Key is purely numeric |
| 8 | `ErrorUnterminatedQuote` | Quoted value never closed |
| 9 | `ErrorInvalidEscapeSequence` | Invalid escape in quoted value |
| 10 | `ErrorEmptyPrefix` | Prefix is empty |
| 11 | `ErrorPrefixMismatch` | Prefixes don't match in comparison |

## Testing

```bash
go test -v ./...
```

## Cross-Language Compatibility

This Go implementation produces identical results to:
- [Rust implementation](https://github.com/machinefabric/tagged-urn-rs)
- [JavaScript implementation](https://github.com/machinefabric/tagged-urn-js)
- [Objective-C implementation](https://github.com/machinefabric/tagged-urn-objc)

All implementations pass the same test cases and follow identical rules. See [Tagged URN RULES.md](https://github.com/machinefabric/tagged-urn-rs/blob/main/docs/RULES.md) for the complete specification.
