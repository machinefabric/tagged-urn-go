# Tagged URN - Go Implementation

Go implementation of Tagged URN with strict validation, pattern matching, and graded specificity comparison.

## Features

- **Strict Rule Enforcement** - Follows exact same rules as Rust, JavaScript, and Objective-C implementations
- **Six Tag Forms** - exact `k=v`, present `k`, absent `!k`, unconstrained `?k`, `k?=v` (absent or not v), `k!=v` (present and not v)
- **Typed Tags** - Each tag is stored as the constraint it writes (`taggedurn.Constraint`: `ConstraintExact`, `ConstraintPresent`, …)
- **Identifiers** - The prefix and every key are an ASCII letter followed by ASCII letters, digits and `-`
- **Case Insensitive** - The prefix, keys and unquoted values are lowered; quoted values keep their case
- **Literal Quoting** - A quoted value is always the literal value: `x="*"` is the exact value `*`
- **Tag Order Independent** - Canonical alphabetical sorting
- **Graded Specificity** - `?k` 0, `k?=v` 1, `k` 2, `k!=v` 3, `k=v` 4, `!k` 5
- **JSON Serialization** - Full JSON marshal/unmarshal support
- **Proved Relations** - Refinement, equivalence and specificity are code generated from the proved model

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
    built, err := taggedurn.NewTaggedUrnBuilder("cap").
        Marker("extract").
        Tag("format", "pdf").
        Build()
    if err != nil {
        log.Fatal(err)
    }
    format, _ := built.GetTag("format")
    fmt.Println(format == taggedurn.ConstraintExact{Value: "pdf"}) // true

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
| `NewTaggedUrnFromTags(prefix, tags)` | Create from prefix and a map of constraints (checked) |
| `Empty(prefix)` | Create empty URN with prefix |
| `GetTag(key)` | The constraint for a key, and whether the URN mentions it |
| `Constraint(key)` | The constraint for a key, `ConstraintMissing` when it is not mentioned |
| `HasTag(key, value)` | Whether the URN writes the exact value `key=value` |
| `HasMarkerTag(key)` | Whether the URN writes the marker `key` |
| `WithTag(key, constraint)` | Return new URN with tag added/replaced (checked) |
| `WithoutTag(key)` | Return new URN with tag removed |
| `ConformsTo(pattern)` | Check if URN conforms to a pattern |
| `Accepts(instance)` | Check if URN (as pattern) accepts an instance |
| `Specificity()` | Get graded specificity score |
| `SpecificityTuple()` | Count tags by form: (absent, exact, presentNot, present, optionalNot) |
| `IsMoreSpecificThan(other)` | Compare specificity with another URN |
| `Compare(other)` | Order by parts: prefix, then each tag by key, constraint form, value |
| `LeastUpperBound(prefix, urns)` | The most specific URN every one of `urns` conforms to, made by the proved model |
| `ToString()` | Get canonical string representation |
| `Hash()` | Get SHA256 hash of canonical form |

### TaggedUrnBuilder

| Method | Description |
|--------|-------------|
| `NewTaggedUrnBuilder(prefix)` | Create builder with prefix |
| `Tag(key, value)` | Add the exact value, taken literally (chainable) |
| `Marker(key)` | Add a marker (chainable) |
| `Constraint(key, constraint)` | Add a tag of any form (chainable) |
| `Build()` | Build the URN, or the first error a tag met |
| `BuildAllowEmpty()` | Build, allowing no tags |

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

| Form | Score |
|------|-------|
| Unconstrained (`?K`) or missing | 0 |
| Absent or not v (`K?=v`) | 1 |
| Present (`K`) | 2 |
| Present and not v (`K!=v`) | 3 |
| Exact value (`K=v`) | 4 |
| Absent (`!K`) | 5 |

## Error Codes

| Code | Constant | Description |
|------|----------|-------------|
| 1 | `ErrorInvalidFormat` | Empty or malformed URN |
| 2 | `ErrorEmptyTag` | Empty key or value component |
| 3 | `ErrorInvalidCharacter` | Disallowed character in prefix, key or value |
| 4 | `ErrorInvalidTagFormat` | Tag not in key=value format |
| 5 | `ErrorMissingPrefix` | URN does not start with prefix |
| 6 | `ErrorDuplicateKey` | Same key appears twice |
| 8 | `ErrorUnterminatedQuote` | Quoted value never closed |
| 9 | `ErrorInvalidEscapeSequence` | Invalid escape in quoted value |
| 10 | `ErrorEmptyPrefix` | Prefix is empty |
| 11 | `ErrorPrefixMismatch` | Prefixes don't match in comparison |
| 12 | `ErrorWhitespaceInInput` | Input has leading or trailing whitespace |

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
