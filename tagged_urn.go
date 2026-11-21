// Package taggedurn provides the fundamental tagged URN system with flat tag-based
// naming, pattern matching, and graded specificity comparison.
//
// Each tag is stored as the constraint it writes (formal.Constraint, re-exported
// here as Constraint):
//   - K=v:  ConstraintExact — must have key K with exact value v
//   - K:    ConstraintPresent — must have key K with any value
//   - !K:   ConstraintAbsent — must NOT have key K
//   - ?K:   ConstraintUnconstrained — no constraint on key K
//   - K?=v: ConstraintOptionalNot — absent, or present and not v
//   - K!=v: ConstraintPresentNot — present and not v
//   - (missing): no constraint; never stored
//
// Graded specificity scoring: ?K 0, K?=v 1, K 2, K!=v 3, K=v 4, !K 5.
package taggedurn

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"unicode"

	lungo "github.com/machinefabric/lungo-go"

	"github.com/machinefabric/tagged-urn-go/formal"
)

// Constraint is what a tag writes about its key: one of the model's seven
// forms. ConstraintMissing is never stored; it stands for a key a URN omits.
type Constraint = formal.Constraint

// The model's forms, by name.
type (
	ConstraintMissing       = formal.ConstraintMissing
	ConstraintUnconstrained = formal.ConstraintUnconstrained
	ConstraintOptionalNot   = formal.ConstraintOptionalNot
	ConstraintPresent       = formal.ConstraintPresent
	ConstraintPresentNot    = formal.ConstraintPresentNot
	ConstraintExact         = formal.ConstraintExact
	ConstraintAbsent        = formal.ConstraintAbsent
)

// TaggedUrn represents a tagged URN using flat, ordered tags with a configurable prefix.
//
// Examples:
//   - cap:generate;ext=pdf;out=binary;target=thumbnail
//   - cap:format;!debug  (format required, debug forbidden)
//   - myapp:key="Value With Spaces"
type TaggedUrn struct {
	prefix string
	tags   map[string]Constraint
	// The same URN on the proved model's side (package formal, generated from
	// ../formal by lungo): its tags with the proof that their keys are strictly
	// increasing. Every semantic question — does one refine another, are they
	// equivalent, how specific is it — is asked of this, so what this package
	// answers is what the theorems there are about. Built once, by assemble,
	// from the same tags as the fields beside it.
	formal formal.Wf
}

// orMissing is the model's form for a key: its constraint, or ConstraintMissing
// for a key the URN omits.
func orMissing(c Constraint, exists bool) Constraint {
	if !exists {
		return ConstraintMissing{}
	}
	return c
}

// assemble is the one way a TaggedUrn is made: the fields and the model's handle,
// together, from the same tags. Every caller has already checked the prefix and
// the keys and refused ConstraintMissing. The model refuses keys that are not
// strictly increasing, and sorted map keys always are (Go orders strings by their
// UTF-8 bytes, which is the code-point order Lean's String < uses) — so a refusal,
// or the runtime failing, is a broken invariant, and it panics.
func assemble(prefix string, tags map[string]Constraint) *TaggedUrn {
	keys := make([]string, 0, len(tags))
	for k := range tags {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	pairs := make([]lungo.Pair[string, Constraint], 0, len(keys))
	for _, k := range keys {
		c := tags[k]
		if c == nil {
			panic(fmt.Sprintf("tagged-urn: key '%s' stored with no constraint", k))
		}
		if _, missing := c.(ConstraintMissing); missing {
			panic(fmt.Sprintf("tagged-urn: key '%s' stored as Missing; a key with no constraint is omitted", k))
		}
		pairs = append(pairs, lungo.Pair[string, Constraint]{First: k, Second: c})
	}
	wf, err := formal.Make(prefix, pairs)
	if err != nil {
		panic(fmt.Sprintf("tagged-urn: the model could not build %s:%v: %v", prefix, keys, err))
	}
	if !wf.Valid {
		panic(fmt.Sprintf("tagged-urn: the model refused %s:%v, whose keys are sorted", prefix, keys))
	}
	return &TaggedUrn{prefix: prefix, tags: tags, formal: wf.Value}
}

// Formal is this URN on the proved model's side. A package generated from a model
// that builds on tagged-urn's names this type for the model's URNs, so what it
// returns passes to that package's functions as it is.
func (c *TaggedUrn) Formal() formal.Wf {
	return c.formal
}

// TaggedUrnRelationKind classifies the order-theoretic relation between two
// tagged URNs. It is derived from Accepts/IsComparable/IsEquivalent and is
// attached to coordinate deltas so callers can distinguish same-point edits,
// same-chain edits, and cross-branch edits without pretending delta only
// exists for comparable pairs.
type TaggedUrnRelationKind string

const (
	TaggedUrnRelationEquivalent   TaggedUrnRelationKind = "equivalent"
	TaggedUrnRelationComparable   TaggedUrnRelationKind = "comparable"
	TaggedUrnRelationIncomparable TaggedUrnRelationKind = "incomparable"
)

// TaggedUrnCoordinateDelta is the coordinate-space edit from one tagged URN to
// another with the same prefix.
//
// Removed contains canonical coordinate entries present in the base but absent
// or changed in the target. Added contains canonical coordinate entries absent
// from the base or changed in the target.
type TaggedUrnCoordinateDelta struct {
	Prefix       string
	Removed      map[string]Constraint
	Added        map[string]Constraint
	RelationKind TaggedUrnRelationKind
}

// IsEmpty returns true when the delta makes no coordinate-space change.
func (d *TaggedUrnCoordinateDelta) IsEmpty() bool {
	if d == nil {
		return true
	}
	return len(d.Removed) == 0 && len(d.Added) == 0
}

// TaggedUrnError represents errors that can occur during tagged URN operations
type TaggedUrnError struct {
	Code    int
	Message string
}

func (e *TaggedUrnError) Error() string {
	return e.Message
}

// Error codes for tagged URN operations. Code 7 is not used.
const (
	ErrorInvalidFormat         = 1
	ErrorEmptyTag              = 2
	ErrorInvalidCharacter      = 3
	ErrorInvalidTagFormat      = 4
	ErrorMissingPrefix         = 5
	ErrorDuplicateKey          = 6
	ErrorUnterminatedQuote     = 8
	ErrorInvalidEscapeSequence = 9
	ErrorEmptyPrefix           = 10
	ErrorPrefixMismatch        = 11
	ErrorWhitespaceInInput     = 12
)

// Parser states for state machine.
//
// The parser reads the six tag forms — the canonical alphabet of the
// constraint truth table — into the constraint each writes:
//
//	| Authored                 | Canonical | Constraint     | Score | Reading                       |
//	|--------------------------|-----------|----------------|------:|-------------------------------|
//	| `?x` ≡ `x?` ≡ `x=?`      | `?x`      | Unconstrained  |     0 | no constraint                 |
//	| `?x=v` ≡ `x?=v` ≡ `x=?v` | `x?=v`    | OptionalNot v  |     1 | absent OR (present and not v) |
//	| `x` ≡ `x=*`              | `x`       | Present        |     2 | present with any value        |
//	| `!x=v` ≡ `x!=v` ≡ `x=!v` | `x!=v`    | PresentNot v   |     3 | present and not v             |
//	| `x=v`                    | `x=v`     | Exact v        |     4 | present and exactly v         |
//	| `!x` ≡ `x!` ≡ `x=!`      | `!x`      | Absent         |     5 | absent (must-not-have)        |
//
// The qualifier `?` or `!` may appear as a key prefix (`?x`, `!x`, `?x=v`,
// `!x=v`), as an infix immediately before `=` (`x?`, `x!`, `x?=v`, `x!=v`),
// or as the first character of an unquoted value (`x=?v`, `x=!v`). The three
// notations are exact aliases. A quoted value is literal: `x="?v"` is the
// exact value `?v`, and `x="*"` the exact value `*`, never the marker.
//
// Disallowed (hard parse errors): `?x?`, `?x?=v`, `!x!=v`, `?!x`, `!?x`,
// `?x=?v`, `?x=*`, `!x=*`, `?x=?`, `x=?*`, `x=!?`, and an empty value.
type parseState int

const (
	stateExpectingKey parseState = iota
	// After `?` at key position; next character must begin a key.
	stateAfterPrefixQuestion
	// After `!` at key position.
	stateAfterPrefixBang
	stateInKey
	// In key, saw `?` after key chars; awaiting `=` (→ infix x?=v)
	// or `;`/end (→ bare x? ≡ ?x).
	stateInKeyAfterQuestion
	stateInKeyAfterBang
	stateExpectingValue
	stateInUnquotedValue
	stateInQuotedValue
	stateInQuotedValueEscape
	stateExpectingSemiOrEnd
)

// isIdentifierStart checks if a character may begin a prefix or a key
func isIdentifierStart(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// isValidKeyChar checks if a character may stand in a prefix or key after its first
func isValidKeyChar(c rune) bool {
	return isIdentifierStart(c) || (c >= '0' && c <= '9') || c == '-'
}

// identifier is s as an identifier — an ASCII letter, then ASCII letters, digits
// and `-` — lowered; what names it in the error ("prefix", "key").
func identifier(what, s string) (string, error) {
	if s == "" {
		if what == "prefix" {
			return "", &TaggedUrnError{Code: ErrorEmptyPrefix, Message: "tagged URN prefix cannot be empty"}
		}
		return "", &TaggedUrnError{Code: ErrorEmptyTag, Message: fmt.Sprintf("empty %s", what)}
	}
	for pos, c := range []rune(s) {
		allowed := isValidKeyChar(c)
		if pos == 0 {
			allowed = isIdentifierStart(c)
		}
		if !allowed {
			return "", &TaggedUrnError{
				Code:    ErrorInvalidCharacter,
				Message: fmt.Sprintf("invalid character '%c' in %s at position %d", c, what, pos),
			}
		}
	}
	return strings.ToLower(s), nil
}

// checkConstraint refuses a constraint a URN cannot write for key: none at all,
// ConstraintMissing (a key with no constraint is omitted), or an empty value.
func checkConstraint(key string, c Constraint) error {
	switch v := c.(type) {
	case nil, ConstraintMissing:
		return &TaggedUrnError{
			Code:    ErrorInvalidTagFormat,
			Message: fmt.Sprintf("key '%s' given no constraint: a key with no constraint is omitted", key),
		}
	case ConstraintExact:
		if v.Value == "" {
			return emptyValue(key)
		}
	case ConstraintOptionalNot:
		if v.Value == "" {
			return emptyValue(key)
		}
	case ConstraintPresentNot:
		if v.Value == "" {
			return emptyValue(key)
		}
	}
	return nil
}

func emptyValue(key string) error {
	return &TaggedUrnError{Code: ErrorEmptyTag, Message: fmt.Sprintf("empty value for key '%s'", key)}
}

// isValidUnquotedValueChar checks if a character may stand in an unquoted value:
// anything but the separators `;` `=`, the quoting characters `"` `\`, and
// whitespace.
func isValidUnquotedValueChar(c rune) bool {
	return c != ';' && c != '=' && c != '"' && c != '\\' && !unicode.IsSpace(c)
}

// needsQuoting reports whether a value must be quoted to read back as itself: it
// holds a character a bare value cannot, an uppercase letter (a bare value is
// lowered), or begins with `*`, `?` or `!` (which a bare value reads as a sigil or
// a qualifier).
func needsQuoting(value string) bool {
	for _, c := range value {
		if !isValidUnquotedValueChar(c) || unicode.IsUpper(c) {
			return true
		}
	}
	return strings.HasPrefix(value, "*") || strings.HasPrefix(value, "?") || strings.HasPrefix(value, "!")
}

// quoteValue quotes a value for serialization
func quoteValue(value string) string {
	var result strings.Builder
	result.WriteRune('"')
	for _, c := range value {
		if c == '"' || c == '\\' {
			result.WriteRune('\\')
		}
		result.WriteRune(c)
	}
	result.WriteRune('"')
	return result.String()
}

// writeValue is a value as it is written: bare, or quoted when it must be.
func writeValue(value string) string {
	if needsQuoting(value) {
		return quoteValue(value)
	}
	return value
}

// qualified is the constraint a value writes under a qualifier (0 for none).
func qualified(qualifier rune, value string) Constraint {
	switch qualifier {
	case 0:
		return ConstraintExact{Value: value}
	case '?':
		return ConstraintOptionalNot{Value: value}
	case '!':
		return ConstraintPresentNot{Value: value}
	}
	panic(fmt.Sprintf("invalid qualifier rune: %q", qualifier))
}

// refuseSigilUnderQualifier refuses a qualifier on a sigil: they make
// incompatible claims.
func refuseSigilUnderQualifier(qualifier rune, key, value string) error {
	if value == "*" || value == "?" || value == "!" {
		return &TaggedUrnError{
			Code: ErrorInvalidCharacter,
			Message: fmt.Sprintf(
				"qualifier '%c' on key '%s' cannot combine with sigil value '%s': "+
					"use a real value or drop the qualifier",
				qualifier, key, value),
		}
	}
	return nil
}

// canonicalNoValue is the constraint of a value-less tag, given its qualifier.
func canonicalNoValue(qualifier rune) Constraint {
	switch qualifier {
	case 0:
		return ConstraintPresent{}
	case '?':
		return ConstraintUnconstrained{}
	case '!':
		return ConstraintAbsent{}
	}
	panic(fmt.Sprintf("invalid qualifier rune: %q", qualifier))
}

// canonicalizeValue is the constraint a parsed (qualifier, value) pair writes. A
// quoted value is literal: Exact, or OptionalNot/PresentNot under a qualifier.
// An unquoted `*`, `?` or `!` is a sigil; an unquoted value that begins with
// `?` or `!` is that qualifier (`x=?v` ≡ `x?=v`). A qualifier on an unquoted
// sigil, or a second qualifier, is refused.
func canonicalizeValue(qualifier rune, quoted bool, key, value string) (Constraint, error) {
	if value == "" {
		return nil, emptyValue(key)
	}
	if quoted {
		return qualified(qualifier, value), nil
	}
	if qualifier == 0 {
		switch value {
		case "*":
			return ConstraintPresent{}, nil
		case "?":
			return ConstraintUnconstrained{}, nil
		case "!":
			return ConstraintAbsent{}, nil
		}
		first := rune(value[0])
		if first == '?' || first == '!' {
			rest := value[1:]
			if err := refuseSigilUnderQualifier(first, key, rest); err != nil {
				return nil, err
			}
			return qualified(first, rest), nil
		}
		return ConstraintExact{Value: value}, nil
	}
	if len(value) > 1 && (value[0] == '?' || value[0] == '!') {
		return nil, &TaggedUrnError{
			Code: ErrorInvalidCharacter,
			Message: fmt.Sprintf(
				"qualifier '%c' on key '%s' cannot combine with a second qualifier in value '%s': "+
					"quote the value to make it literal",
				qualifier, key, value),
		}
	}
	if err := refuseSigilUnderQualifier(qualifier, key, value); err != nil {
		return nil, err
	}
	return qualified(qualifier, value), nil
}

// NewTaggedUrnFromString creates a tagged URN from a string
// Format: prefix:key1=value1;key2=value2;... or prefix:key1="value with spaces";key2=simple
// The prefix is required and ends at the first colon
// Trailing semicolons are optional and ignored
// Tags are automatically sorted alphabetically for canonical form
//
// Case handling:
// - Prefix: Normalized to lowercase
// - Keys: Always normalized to lowercase
// - Unquoted values: Normalized to lowercase
// - Quoted values: Case preserved exactly as specified
func NewTaggedUrnFromString(s string) (*TaggedUrn, error) {
	// Fail hard on leading/trailing whitespace
	if s != strings.TrimSpace(s) {
		return nil, &TaggedUrnError{
			Code:    ErrorWhitespaceInInput,
			Message: fmt.Sprintf("tagged URN has leading or trailing whitespace: '%s'", s),
		}
	}

	if s == "" {
		return nil, &TaggedUrnError{
			Code:    ErrorInvalidFormat,
			Message: "tagged URN cannot be empty",
		}
	}

	// Find the prefix (everything before the first colon)
	colonPos := strings.Index(s, ":")
	if colonPos == -1 {
		return nil, &TaggedUrnError{
			Code:    ErrorMissingPrefix,
			Message: "tagged URN must have a prefix followed by ':'",
		}
	}

	if colonPos == 0 {
		return nil, &TaggedUrnError{
			Code:    ErrorEmptyPrefix,
			Message: "tagged URN prefix cannot be empty",
		}
	}

	prefix, err := identifier("prefix", s[:colonPos])
	if err != nil {
		return nil, err
	}
	tagsPart := s[colonPos+1:]
	tags := make(map[string]Constraint)

	// Handle empty tagged URN (prefix: with no tags or just semicolon)
	if tagsPart == "" || tagsPart == ";" {
		return assemble(prefix, tags), nil
	}

	state := stateExpectingKey
	var currentKey strings.Builder
	var currentValue strings.Builder
	chars := []rune(tagsPart)
	pos := 0
	// Tracks the qualifier for the tag currently being parsed:
	//   0    — no qualifier seen yet
	//   '?'  — `?` qualifier (prefix `?x` or infix `x?=`)
	//   '!'  — `!` qualifier (prefix `!x` or infix `x!=`)
	// Reset to 0 on each finishTag.
	var qualifier rune

	finishTag := func(constraint Constraint) error {
		key := currentKey.String()
		if _, exists := tags[key]; exists {
			return &TaggedUrnError{
				Code:    ErrorDuplicateKey,
				Message: fmt.Sprintf("duplicate tag key: %s", key),
			}
		}
		tags[key] = constraint
		currentKey.Reset()
		currentValue.Reset()
		qualifier = 0
		return nil
	}

	finishValue := func(quoted bool) error {
		constraint, err := canonicalizeValue(qualifier, quoted, currentKey.String(), currentValue.String())
		if err != nil {
			return err
		}
		return finishTag(constraint)
	}

	for pos < len(chars) {
		c := chars[pos]

		switch state {
		case stateExpectingKey:
			if c == ';' {
				pos++
				continue
			} else if c == '?' {
				qualifier = '?'
				state = stateAfterPrefixQuestion
			} else if c == '!' {
				qualifier = '!'
				state = stateAfterPrefixBang
			} else if isIdentifierStart(c) {
				currentKey.WriteRune(unicode.ToLower(c))
				state = stateInKey
			} else {
				return nil, &TaggedUrnError{
					Code:    ErrorInvalidCharacter,
					Message: fmt.Sprintf("invalid character '%c' at position %d", c, pos),
				}
			}

		case stateAfterPrefixQuestion, stateAfterPrefixBang:
			if isIdentifierStart(c) {
				currentKey.WriteRune(unicode.ToLower(c))
				state = stateInKey
			} else {
				return nil, &TaggedUrnError{
					Code: ErrorInvalidCharacter,
					Message: fmt.Sprintf(
						"expected key character after '%c' qualifier, got '%c' at position %d",
						qualifier, c, pos),
				}
			}

		case stateInKey:
			if c == '=' {
				state = stateExpectingValue
			} else if c == '?' {
				if qualifier != 0 {
					return nil, &TaggedUrnError{
						Code: ErrorInvalidCharacter,
						Message: fmt.Sprintf(
							"duplicate qualifier '?' at position %d: prefix and infix qualifiers cannot be combined on key '%s'",
							pos, currentKey.String()),
					}
				}
				qualifier = '?'
				state = stateInKeyAfterQuestion
			} else if c == '!' {
				if qualifier != 0 {
					return nil, &TaggedUrnError{
						Code: ErrorInvalidCharacter,
						Message: fmt.Sprintf(
							"duplicate qualifier '!' at position %d: prefix and infix qualifiers cannot be combined on key '%s'",
							pos, currentKey.String()),
					}
				}
				qualifier = '!'
				state = stateInKeyAfterBang
			} else if c == ';' {
				if err := finishTag(canonicalNoValue(qualifier)); err != nil {
					return nil, err
				}
				state = stateExpectingKey
			} else if isValidKeyChar(c) {
				currentKey.WriteRune(unicode.ToLower(c))
			} else {
				return nil, &TaggedUrnError{
					Code:    ErrorInvalidCharacter,
					Message: fmt.Sprintf("invalid character '%c' in key at position %d", c, pos),
				}
			}

		case stateInKeyAfterQuestion, stateInKeyAfterBang:
			if c == '=' {
				state = stateExpectingValue
			} else if c == ';' {
				if err := finishTag(canonicalNoValue(qualifier)); err != nil {
					return nil, err
				}
				state = stateExpectingKey
			} else {
				return nil, &TaggedUrnError{
					Code: ErrorInvalidCharacter,
					Message: fmt.Sprintf(
						"expected '=' or ';' after '%s%c' suffix qualifier, got '%c' at position %d",
						currentKey.String(), qualifier, c, pos),
				}
			}

		case stateExpectingValue:
			if c == '"' {
				state = stateInQuotedValue
			} else if c == ';' {
				return nil, emptyValue(currentKey.String())
			} else if isValidUnquotedValueChar(c) {
				currentValue.WriteRune(asciiLower(c))
				state = stateInUnquotedValue
			} else {
				return nil, &TaggedUrnError{
					Code:    ErrorInvalidCharacter,
					Message: fmt.Sprintf("invalid character '%c' in value at position %d", c, pos),
				}
			}

		case stateInUnquotedValue:
			if c == ';' {
				if err := finishValue(false); err != nil {
					return nil, err
				}
				state = stateExpectingKey
			} else if isValidUnquotedValueChar(c) {
				currentValue.WriteRune(asciiLower(c))
			} else {
				return nil, &TaggedUrnError{
					Code:    ErrorInvalidCharacter,
					Message: fmt.Sprintf("invalid character '%c' in unquoted value at position %d", c, pos),
				}
			}

		case stateInQuotedValue:
			if c == '"' {
				state = stateExpectingSemiOrEnd
			} else if c == '\\' {
				state = stateInQuotedValueEscape
			} else {
				currentValue.WriteRune(c)
			}

		case stateInQuotedValueEscape:
			if c == '"' || c == '\\' {
				currentValue.WriteRune(c)
				state = stateInQuotedValue
			} else {
				return nil, &TaggedUrnError{
					Code:    ErrorInvalidEscapeSequence,
					Message: fmt.Sprintf("invalid escape sequence at position %d (only \\\" and \\\\ allowed)", pos),
				}
			}

		case stateExpectingSemiOrEnd:
			if c == ';' {
				if err := finishValue(true); err != nil {
					return nil, err
				}
				state = stateExpectingKey
			} else {
				return nil, &TaggedUrnError{
					Code:    ErrorInvalidCharacter,
					Message: fmt.Sprintf("expected ';' or end after quoted value, got '%c' at position %d", c, pos),
				}
			}
		}

		pos++
	}

	switch state {
	case stateInUnquotedValue:
		if err := finishValue(false); err != nil {
			return nil, err
		}
	case stateExpectingSemiOrEnd:
		if err := finishValue(true); err != nil {
			return nil, err
		}
	case stateExpectingKey:
		// Valid — trailing semicolon or empty input after prefix.
	case stateInQuotedValue, stateInQuotedValueEscape:
		return nil, &TaggedUrnError{
			Code:    ErrorUnterminatedQuote,
			Message: fmt.Sprintf("unterminated quote at position %d", pos),
		}
	case stateAfterPrefixQuestion, stateAfterPrefixBang:
		return nil, &TaggedUrnError{
			Code:    ErrorEmptyTag,
			Message: fmt.Sprintf("qualifier '%c' at end of input has no key", qualifier),
		}
	case stateInKey, stateInKeyAfterQuestion, stateInKeyAfterBang:
		if err := finishTag(canonicalNoValue(qualifier)); err != nil {
			return nil, err
		}
	case stateExpectingValue:
		return nil, emptyValue(currentKey.String())
	}

	return assemble(prefix, tags), nil
}

// asciiLower lowers an ASCII letter and leaves every other character as it is,
// as unquoted values are lowered.
func asciiLower(c rune) rune {
	if c >= 'A' && c <= 'Z' {
		return c + ('a' - 'A')
	}
	return c
}

// NewTaggedUrnFromTags creates a tagged URN from a prefix and its tags. The prefix
// and every key must be identifiers (an ASCII letter, then ASCII letters, digits
// and `-`) and are lowered; every constraint must be one a URN can write.
func NewTaggedUrnFromTags(prefix string, tags map[string]Constraint) (*TaggedUrn, error) {
	p, err := identifier("prefix", prefix)
	if err != nil {
		return nil, err
	}
	result := make(map[string]Constraint, len(tags))
	for k, v := range tags {
		key, err := identifier("key", k)
		if err != nil {
			return nil, err
		}
		if err := checkConstraint(key, v); err != nil {
			return nil, err
		}
		if _, exists := result[key]; exists {
			return nil, &TaggedUrnError{Code: ErrorDuplicateKey, Message: fmt.Sprintf("duplicate tag key: %s", key)}
		}
		result[key] = v
	}
	return assemble(p, result), nil
}

// Empty creates an empty tagged URN with the specified prefix (required)
func Empty(prefix string) (*TaggedUrn, error) {
	p, err := identifier("prefix", prefix)
	if err != nil {
		return nil, err
	}
	return assemble(p, make(map[string]Constraint)), nil
}

// GetPrefix returns the prefix of this tagged URN
func (c *TaggedUrn) GetPrefix() string {
	return c.prefix
}

// GetTag returns the constraint the URN writes for key, and whether it mentions
// it. Key is normalized to lowercase for lookup.
func (c *TaggedUrn) GetTag(key string) (Constraint, bool) {
	value, exists := c.tags[strings.ToLower(key)]
	return value, exists
}

// Constraint returns the constraint the URN writes for key, ConstraintMissing
// when it does not mention it: the model's form of the key.
func (c *TaggedUrn) Constraint(key string) Constraint {
	value, exists := c.tags[strings.ToLower(key)]
	return orMissing(value, exists)
}

// AllTags returns a copy of all tags in this URN
func (c *TaggedUrn) AllTags() map[string]Constraint {
	result := make(map[string]Constraint, len(c.tags))
	for k, v := range c.tags {
		result[k] = v
	}
	return result
}

// HasTag reports whether this URN writes the exact value value for key
// (`key=value`). Key is normalized to lowercase; value comparison is case-sensitive.
func (c *TaggedUrn) HasTag(key, value string) bool {
	exact, ok := c.tags[strings.ToLower(key)].(ConstraintExact)
	return ok && exact.Value == value
}

// HasMarkerTag reports whether tagName is present as a marker (ConstraintPresent,
// written as just the key). Example: "cap:constrained;..." has marker tag
// "constrained".
func (c *TaggedUrn) HasMarkerTag(tagName string) bool {
	_, ok := c.tags[strings.ToLower(tagName)].(ConstraintPresent)
	return ok
}

// WithTag returns a new tagged URN with an added or replaced tag. The key must be
// an identifier and is lowered; the constraint must be one a URN can write.
func (c *TaggedUrn) WithTag(key string, constraint Constraint) (*TaggedUrn, error) {
	k, err := identifier("key", key)
	if err != nil {
		return nil, err
	}
	if err := checkConstraint(k, constraint); err != nil {
		return nil, err
	}
	newTags := make(map[string]Constraint, len(c.tags)+1)
	for tk, tv := range c.tags {
		newTags[tk] = tv
	}
	newTags[k] = constraint
	return assemble(c.prefix, newTags), nil
}

// WithoutTag returns a new tagged URN with a tag removed
// Key is normalized to lowercase for case-insensitive removal
func (c *TaggedUrn) WithoutTag(key string) *TaggedUrn {
	newTags := make(map[string]Constraint)
	key = strings.ToLower(key)
	for k, v := range c.tags {
		if k != key {
			newTags[k] = v
		}
	}
	return assemble(c.prefix, newTags)
}

// ConformsTo checks whether this URN (the instance) satisfies the pattern:
// everything this URN describes, the pattern describes.
//
// Decided by the proved model (formal.Refines): every tag form means the set of
// states it allows — on either side — and the instance satisfies the pattern
// when, key by key, its set lies inside the pattern's. A key the instance omits
// promises nothing, `?x` promises nothing, and `x` promises presence but no
// particular value. Both URNs must have the same prefix.
//
// This is a guarantee between two descriptions. Whether they COULD be about the
// same thing is Meets; whether a complete thing satisfies a description is
// Satisfies. Equivalent to pattern.Accepts(self).
func (c *TaggedUrn) ConformsTo(pattern *TaggedUrn) (bool, error) {
	if pattern == nil {
		return false, &TaggedUrnError{
			Code:    ErrorInvalidFormat,
			Message: "cannot match against nil pattern",
		}
	}
	return checkMatch(c, pattern)
}

// Accepts checks if this URN (pattern) accepts the given instance.
// Equivalent to instance.ConformsTo(self).
func (c *TaggedUrn) Accepts(instance *TaggedUrn) (bool, error) {
	if instance == nil {
		return false, &TaggedUrnError{
			Code:    ErrorInvalidFormat,
			Message: "cannot match against nil instance",
		}
	}
	return checkMatch(instance, c)
}

// checkMatch is the core matching: does instance satisfy pattern's constraints?
// Decided by the model (formal.Refines, proved to be refinement); only the
// prefix check, which the model folds into the answer, is here, because a
// comparison across prefixes is a caller's error and says so.
func checkMatch(instance, pattern *TaggedUrn) (bool, error) {
	if err := samePrefix(instance, pattern); err != nil {
		return false, err
	}
	return formal.Refines(instance.formal, pattern.formal)
}

// Meets reports whether this URN and other COULD be about the same thing: some
// thing is described by both. Symmetric — neither is the instance.
//
// `media:ext` (some ext) does not conform to `media:document;ext=pdf`, and is not
// excluded by it either: it meets it, and only the value that turns up says
// which. Whatever conforms meets; what meets need not conform, and meeting is
// not transitive (a pdf meets "some ext", which meets a png). Decided by the
// proved model (formal.Meets).
func (c *TaggedUrn) Meets(other *TaggedUrn) (bool, error) {
	if other == nil {
		return false, &TaggedUrnError{Code: ErrorInvalidFormat, Message: "cannot compare with a nil URN"}
	}
	if err := samePrefix(c, other); err != nil {
		return false, err
	}
	return formal.Meets(c.formal, other.formal)
}

// Satisfies reports whether this URN, read as a COMPLETE thing, satisfies the
// pattern.
//
// A description that omits a key says nothing about it, which is how ConformsTo
// reads both sides. A thing that exists — a value with these tags, a cap's own
// list of tags — omits a key because it does not have it. Read so, a thing that
// does not mention `x` satisfies `!x`. Use this where the receiver is what
// something IS; use ConformsTo where it is what something is declared to take
// or give. Decided by the proved model (formal.RefinesClosed).
func (c *TaggedUrn) Satisfies(pattern *TaggedUrn) (bool, error) {
	if pattern == nil {
		return false, &TaggedUrnError{Code: ErrorInvalidFormat, Message: "cannot match against nil pattern"}
	}
	if err := samePrefix(c, pattern); err != nil {
		return false, err
	}
	return formal.RefinesClosed(c.formal, pattern.formal)
}

// MaySatisfy reports whether this URN, read as a complete thing, COULD satisfy
// the pattern: Satisfies is to this as ConformsTo is to Meets.
func (c *TaggedUrn) MaySatisfy(pattern *TaggedUrn) (bool, error) {
	if pattern == nil {
		return false, &TaggedUrnError{Code: ErrorInvalidFormat, Message: "cannot match against nil pattern"}
	}
	if err := samePrefix(c, pattern); err != nil {
		return false, err
	}
	return formal.MeetsClosed(c.formal, pattern.formal)
}

// LeastUpperBound is the most specific URN every one of urns conforms to: their
// common denominator.
//
// Key by key, the bound allows every state some input allows. Where no single
// form says exactly that, there is no least bound at all — `ext=png` and
// `ext=jpeg` are both `ext`, both `ext!=pdf` and both `ext!=gif` — so the bound
// is the least among those that exclude only values the inputs name: here
// `ext`. When one input holds all the others, the bound is that input. The bound
// of no URNs is the URN with no tags.
//
// Made by the proved model (formal.LeastUpperBound), whose theorems say every
// input conforms to the bound and the bound conforms to every such bound of
// them, and no tag of it is ConstraintMissing. The prefix must be an
// identifier, and every URN must have it; another is a PrefixMismatch error.
func LeastUpperBound(prefix string, urns []*TaggedUrn) (*TaggedUrn, error) {
	prefix, err := identifier("prefix", prefix)
	if err != nil {
		return nil, err
	}
	handles := make([]formal.Wf, 0, len(urns))
	for _, u := range urns {
		if u == nil {
			return nil, &TaggedUrnError{Code: ErrorInvalidFormat, Message: "cannot bound a nil URN"}
		}
		handles = append(handles, u.formal)
	}
	tags, err := formal.LeastUpperBound(prefix, handles)
	if err != nil {
		return nil, err
	}
	if !tags.Valid {
		for _, u := range urns {
			if u.prefix != prefix {
				return nil, &TaggedUrnError{
					Code:    ErrorPrefixMismatch,
					Message: fmt.Sprintf("cannot bound URNs with different prefixes: '%s' vs '%s'", prefix, u.prefix),
				}
			}
		}
		panic(fmt.Sprintf("tagged-urn: the model made no bound of URNs that all have the prefix %s", prefix))
	}
	result := make(map[string]Constraint, len(tags.Value))
	for _, p := range tags.Value {
		result[p.First] = p.Second
	}
	return assemble(prefix, result), nil
}

// samePrefix refuses a comparison of URNs with different prefixes.
func samePrefix(a, b *TaggedUrn) error {
	if a.prefix != b.prefix {
		return &TaggedUrnError{
			Code:    ErrorPrefixMismatch,
			Message: fmt.Sprintf("cannot compare URNs with different prefixes: '%s' vs '%s'", a.prefix, b.prefix),
		}
	}
	return nil
}

// ValuesMatch decides one key: does the instance's constraint (ConstraintMissing
// for a key the URN omits, as TaggedUrn.Constraint gives it) satisfy the pattern's? Exposed for callers (e.g. CapUrn's
// y-axis matcher) that walk tag sets themselves and need the same per-key
// decision the tagged-URN matcher makes. Decided by the model
// (formal.ValuesMatch): every form has one meaning — the set of states the key
// may be in — and the instance satisfies the pattern when every state it
// allows, the pattern allows too (tagMatch_iff_allows in ../formal).
func ValuesMatch(inst, patt Constraint) bool {
	ok, err := formal.ValuesMatch(inst, patt)
	if err != nil {
		panic(fmt.Sprintf("tagged-urn: the model could not match one key: %v", err))
	}
	return ok
}

// ValuesMeet decides one key: do the two constraints allow a common state?
// (formal.ValuesMeet.)
func ValuesMeet(a, b Constraint) bool {
	ok, err := formal.ValuesMeet(a, b)
	if err != nil {
		panic(fmt.Sprintf("tagged-urn: the model could not compare one key: %v", err))
	}
	return ok
}

// ValuesMatchClosed decides one key of a complete thing against a pattern: an
// omitted key is absent. (formal.ValuesMatchClosed.)
func ValuesMatchClosed(inst, patt Constraint) bool {
	ok, err := formal.ValuesMatchClosed(inst, patt)
	if err != nil {
		panic(fmt.Sprintf("tagged-urn: the model could not match one key: %v", err))
	}
	return ok
}

// ConformsToStr checks if this URN (instance) satisfies a string pattern's constraints.
func (c *TaggedUrn) ConformsToStr(patternStr string) (bool, error) {
	pattern, err := NewTaggedUrnFromString(patternStr)
	if err != nil {
		return false, err
	}
	return c.ConformsTo(pattern)
}

// AcceptsStr checks if this URN (pattern) accepts a string instance.
func (c *TaggedUrn) AcceptsStr(instanceStr string) (bool, error) {
	instance, err := NewTaggedUrnFromString(instanceStr)
	if err != nil {
		return false, err
	}
	return c.Accepts(instance)
}

// Specificity returns the specificity score for URN matching: the sum of the
// per-tag truth-table scores across every tag in the URN, as the model computes
// it (formal.Specificity). More specific URNs have higher scores and are
// preferred.
func (c *TaggedUrn) Specificity() int {
	score, err := formal.Specificity(c.formal)
	if err != nil {
		panic(fmt.Sprintf("tagged-urn: the model could not score %s: %v", c.ToString(), err))
	}
	if !score.IsInt64() {
		panic(fmt.Sprintf("tagged-urn: the specificity of %s does not fit an int: %s", c.ToString(), score))
	}
	return int(score.Int64())
}

// SpecificityTuple counts how many tags fall into each non-zero form, from the
// highest score to the lowest: (absent, exact, presentNot, present, optionalNot).
func (c *TaggedUrn) SpecificityTuple() (int, int, int, int, int) {
	absent, exact, presentNot, present, optionalNot := 0, 0, 0, 0, 0
	for _, constraint := range c.tags {
		switch constraint.(type) {
		case ConstraintAbsent:
			absent++
		case ConstraintExact:
			exact++
		case ConstraintPresentNot:
			presentNot++
		case ConstraintPresent:
			present++
		case ConstraintOptionalNot:
			optionalNot++
		case ConstraintUnconstrained:
			// 0 points, not counted
		default:
			panic(fmt.Sprintf("tagged-urn: stored constraint %#v", constraint))
		}
	}
	return absent, exact, presentNot, present, optionalNot
}

// IsMoreSpecificThan checks if this URN is more specific than another
func (c *TaggedUrn) IsMoreSpecificThan(other *TaggedUrn) (bool, error) {
	if other == nil {
		return false, &TaggedUrnError{
			Code:    ErrorInvalidFormat,
			Message: "cannot compare against nil URN",
		}
	}

	// First check prefix
	if c.prefix != other.prefix {
		return false, &TaggedUrnError{
			Code:    ErrorPrefixMismatch,
			Message: fmt.Sprintf("cannot compare URNs with different prefixes: '%s' vs '%s'", c.prefix, other.prefix),
		}
	}

	return c.Specificity() > other.Specificity(), nil
}

// IsEquivalent checks if two URNs are equivalent (identical tag sets).
//
// From order theory: in the specialization partial order defined by
// Accepts/ConformsTo, two elements are **equivalent** when each
// accepts the other (antisymmetry: a ≤ b ∧ b ≤ a → a = b).
//
// This is stricter than IsComparable — it requires the tag sets to
// be identical, not just related by specialization.
//
//	a.IsEquivalent(b)  ≡  a.Accepts(b) && b.Accepts(a)
//
// Decided by the model (formal.Equivalent, proved equal to it). Returns a
// PrefixMismatch error if the prefixes differ.
func (c *TaggedUrn) IsEquivalent(other *TaggedUrn) (bool, error) {
	if other == nil {
		return false, &TaggedUrnError{
			Code:    ErrorInvalidFormat,
			Message: "cannot compare against nil URN",
		}
	}
	if err := samePrefix(c, other); err != nil {
		return false, err
	}
	return formal.Equivalent(c.formal, other.formal)
}

// IsComparable checks if two URNs are comparable (one is a specialization of the other).
//
// From order theory: in a partial order, two elements are **comparable**
// when one is ≤ the other. Elements that are NOT comparable are in
// different branches of the specialization lattice (e.g., `media:pdf;bytes`
// vs `media:enc=utf-8;txt` — neither accepts the other).
//
// This is the weakest relation: it finds all URNs on the same
// generalization/specialization chain. Use it when you want to discover
// all handlers that *could* service a request, whether they are more
// general (fallback) or more specific (exact match).
//
//	a.IsComparable(b)  ≡  a.Accepts(b) || b.Accepts(a)
//
// Decided by the model (formal.Comparable, proved equal to it). Returns a
// PrefixMismatch error if the prefixes differ.
func (c *TaggedUrn) IsComparable(other *TaggedUrn) (bool, error) {
	if other == nil {
		return false, &TaggedUrnError{
			Code:    ErrorInvalidFormat,
			Message: "cannot compare against nil URN",
		}
	}
	if err := samePrefix(c, other); err != nil {
		return false, err
	}
	return formal.Comparable(c.formal, other.formal)
}

// Compare returns -1, 0, or 1: by prefix, then by canonical text. Two URNs
// compare equal exactly when they are equal, since every URN prints so that it
// reads back as itself.
func (c *TaggedUrn) Compare(other *TaggedUrn) int {
	if c == nil && other == nil {
		return 0
	}
	if c == nil {
		return -1
	}
	if other == nil {
		return 1
	}
	if c.prefix != other.prefix {
		return strings.Compare(c.prefix, other.prefix)
	}
	return strings.Compare(c.ToString(), other.ToString())
}

// DeltaFrom computes the coordinate-space delta from base to c.
//
// Delta is defined over explicit canonical coordinates, not semantic
// equivalence classes. Equivalent URNs may still yield a non-empty delta if one
// side explicitly authors no-op coordinates.
func (c *TaggedUrn) DeltaFrom(base *TaggedUrn) (*TaggedUrnCoordinateDelta, error) {
	if base == nil {
		return nil, &TaggedUrnError{
			Code:    ErrorInvalidFormat,
			Message: "cannot derive delta from nil URN",
		}
	}
	if c.prefix != base.prefix {
		return nil, &TaggedUrnError{
			Code:    ErrorPrefixMismatch,
			Message: fmt.Sprintf("cannot compare URNs with different prefixes: '%s' vs '%s'", base.prefix, c.prefix),
		}
	}

	equivalent, err := c.IsEquivalent(base)
	if err != nil {
		return nil, err
	}
	comparable := false
	if !equivalent {
		comparable, err = c.IsComparable(base)
		if err != nil {
			return nil, err
		}
	}

	relationKind := TaggedUrnRelationIncomparable
	if equivalent {
		relationKind = TaggedUrnRelationEquivalent
	} else if comparable {
		relationKind = TaggedUrnRelationComparable
	}

	removed := make(map[string]Constraint)
	added := make(map[string]Constraint)
	allKeys := make(map[string]struct{}, len(base.tags)+len(c.tags))
	for key := range base.tags {
		allKeys[key] = struct{}{}
	}
	for key := range c.tags {
		allKeys[key] = struct{}{}
	}
	for key := range allKeys {
		baseValue, baseExists := base.tags[key]
		targetValue, targetExists := c.tags[key]
		if baseExists && targetExists && baseValue == targetValue {
			continue
		}
		if baseExists {
			removed[key] = baseValue
		}
		if targetExists {
			added[key] = targetValue
		}
	}

	return &TaggedUrnCoordinateDelta{
		Prefix:       c.prefix,
		Removed:      removed,
		Added:        added,
		RelationKind: relationKind,
	}, nil
}

// ApplyDelta applies a coordinate-space delta to this tagged URN.
//
// Keys named in Removed are deleted regardless of current value, then keys
// named in Added are inserted. Unrelated coordinates are preserved. A delta made
// by DeltaFrom holds only constraints a URN writes; one built by hand is checked
// as WithTag checks.
func (c *TaggedUrn) ApplyDelta(delta *TaggedUrnCoordinateDelta) (*TaggedUrn, error) {
	if delta == nil {
		return nil, &TaggedUrnError{
			Code:    ErrorInvalidFormat,
			Message: "cannot apply nil delta",
		}
	}
	if c.prefix != delta.Prefix {
		return nil, &TaggedUrnError{
			Code:    ErrorPrefixMismatch,
			Message: fmt.Sprintf("cannot apply delta with different prefix: '%s' vs '%s'", delta.Prefix, c.prefix),
		}
	}

	nextTags := make(map[string]Constraint, len(c.tags)+len(delta.Added))
	for key, value := range c.tags {
		nextTags[key] = value
	}
	for key := range delta.Removed {
		delete(nextTags, key)
	}
	for key, value := range delta.Added {
		k, err := identifier("key", key)
		if err != nil {
			return nil, err
		}
		if err := checkConstraint(k, value); err != nil {
			return nil, err
		}
		nextTags[k] = value
	}
	return assemble(c.prefix, nextTags), nil
}

// IsEquivalentStr is a string variant of IsEquivalent.
func (c *TaggedUrn) IsEquivalentStr(otherStr string) (bool, error) {
	other, err := NewTaggedUrnFromString(otherStr)
	if err != nil {
		return false, err
	}
	return c.IsEquivalent(other)
}

// IsComparableStr is a string variant of IsComparable.
func (c *TaggedUrn) IsComparableStr(otherStr string) (bool, error) {
	other, err := NewTaggedUrnFromString(otherStr)
	if err != nil {
		return false, err
	}
	return c.IsComparable(other)
}

// WithWildcardTag returns the same URN with key, where it has it, made a marker
// (ConstraintPresent).
func (c *TaggedUrn) WithWildcardTag(key string) *TaggedUrn {
	key = strings.ToLower(key)
	if _, exists := c.tags[key]; !exists {
		return c
	}
	newTags := make(map[string]Constraint, len(c.tags))
	for k, v := range c.tags {
		newTags[k] = v
	}
	newTags[key] = ConstraintPresent{}
	return assemble(c.prefix, newTags)
}

// Subset returns a new URN with only specified tags
func (c *TaggedUrn) Subset(keys []string) *TaggedUrn {
	newTags := make(map[string]Constraint)
	for _, key := range keys {
		key = strings.ToLower(key)
		if value, exists := c.tags[key]; exists {
			newTags[key] = value
		}
	}
	return assemble(c.prefix, newTags)
}

// Merge returns a new URN merged with another (other takes precedence for conflicts)
// Both must have the same prefix
func (c *TaggedUrn) Merge(other *TaggedUrn) (*TaggedUrn, error) {
	if other == nil {
		return nil, &TaggedUrnError{
			Code:    ErrorInvalidFormat,
			Message: "cannot merge with nil URN",
		}
	}

	if c.prefix != other.prefix {
		return nil, &TaggedUrnError{
			Code:    ErrorPrefixMismatch,
			Message: fmt.Sprintf("cannot merge URNs with different prefixes: '%s' vs '%s'", c.prefix, other.prefix),
		}
	}

	newTags := make(map[string]Constraint)
	for k, v := range c.tags {
		newTags[k] = v
	}
	for k, v := range other.tags {
		newTags[k] = v
	}
	return assemble(c.prefix, newTags), nil
}

// ToString returns the canonical string representation of this tagged URN: the
// prefix, `:`, and the tags in key order joined by `;`, with no trailing
// semicolon. Each constraint is written in its canonical form:
//
//	Present          -> k
//	Unconstrained    -> ?k
//	Absent           -> !k
//	OptionalNot v    -> k?=v
//	PresentNot v     -> k!=v
//	Exact v          -> k=v
//
// A value is quoted only when it must be to read back as itself (needsQuoting):
// Exact "*" prints k="*".
func (c *TaggedUrn) ToString() string {
	keys := make([]string, 0, len(c.tags))
	for key := range c.tags {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		switch v := c.tags[key].(type) {
		case ConstraintPresent:
			parts = append(parts, key)
		case ConstraintUnconstrained:
			parts = append(parts, "?"+key)
		case ConstraintAbsent:
			parts = append(parts, "!"+key)
		case ConstraintOptionalNot:
			parts = append(parts, key+"?="+writeValue(v.Value))
		case ConstraintPresentNot:
			parts = append(parts, key+"!="+writeValue(v.Value))
		case ConstraintExact:
			parts = append(parts, key+"="+writeValue(v.Value))
		default:
			panic(fmt.Sprintf("tagged-urn: stored constraint %#v", v))
		}
	}
	return c.prefix + ":" + strings.Join(parts, ";")
}

// String implements the Stringer interface
func (c *TaggedUrn) String() string {
	return c.ToString()
}

// Equals checks if this tagged URN is equal to another
func (c *TaggedUrn) Equals(other *TaggedUrn) bool {
	if other == nil {
		return false
	}

	if c.prefix != other.prefix {
		return false
	}

	if len(c.tags) != len(other.tags) {
		return false
	}

	for key, value := range c.tags {
		otherValue, exists := other.tags[key]
		if !exists || value != otherValue {
			return false
		}
	}

	return true
}

// Hash returns a hash of this tagged URN
// Two equivalent tagged URNs will have the same hash
func (c *TaggedUrn) Hash() string {
	// Use canonical string representation for consistent hashing
	canonical := c.ToString()
	h := sha256.Sum256([]byte(canonical))
	return fmt.Sprintf("%x", h)
}

// MarshalJSON implements the json.Marshaler interface
func (c *TaggedUrn) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.ToString())
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (c *TaggedUrn) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("failed to unmarshal TaggedUrn: expected string, got: %s", string(data))
	}

	taggedUrn, err := NewTaggedUrnFromString(s)
	if err != nil {
		return err
	}

	*c = *taggedUrn
	return nil
}

// UrnMatcher provides utility methods for matching URNs
type UrnMatcher struct{}

// FindBestMatch finds the most specific URN that conforms to a request's constraints.
// URNs are instances (capabilities), request is the pattern (requirement).
func (m *UrnMatcher) FindBestMatch(urns []*TaggedUrn, request *TaggedUrn) (*TaggedUrn, error) {
	var best *TaggedUrn
	bestSpecificity := -1

	for _, urn := range urns {
		ok, err := urn.ConformsTo(request)
		if err != nil {
			return nil, err
		}
		if ok {
			specificity := urn.Specificity()
			if specificity > bestSpecificity {
				best = urn
				bestSpecificity = specificity
			}
		}
	}

	return best, nil
}

// FindAllMatches finds all URNs that conform to a request's constraints, sorted by specificity.
// URNs are instances (capabilities), request is the pattern (requirement).
func (m *UrnMatcher) FindAllMatches(urns []*TaggedUrn, request *TaggedUrn) ([]*TaggedUrn, error) {
	var results []*TaggedUrn

	for _, urn := range urns {
		ok, err := urn.ConformsTo(request)
		if err != nil {
			return nil, err
		}
		if ok {
			results = append(results, urn)
		}
	}

	// Sort by specificity (most specific first)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Specificity() > results[j].Specificity()
	})

	return results, nil
}

// AreCompatible checks if two URN sets are compatible
// Two URNs are compatible if either accepts the other (bidirectional accepts)
func (m *UrnMatcher) AreCompatible(urns1, urns2 []*TaggedUrn) (bool, error) {
	for _, u1 := range urns1 {
		for _, u2 := range urns2 {
			fwd, err := u1.Accepts(u2)
			if err != nil {
				return false, err
			}
			if fwd {
				return true, nil
			}
			rev, err := u2.Accepts(u1)
			if err != nil {
				return false, err
			}
			if rev {
				return true, nil
			}
		}
	}
	return false, nil
}

// TaggedUrnBuilder provides a fluent builder interface for creating tagged URNs.
// The first error (a prefix or key that is no identifier, a constraint a URN
// cannot write, a key given twice) is kept and returned by Build.
type TaggedUrnBuilder struct {
	prefix string
	tags   map[string]Constraint
	err    error
}

// NewTaggedUrnBuilder creates a new builder with a specified prefix (required),
// which must be an identifier and is lowered.
func NewTaggedUrnBuilder(prefix string) *TaggedUrnBuilder {
	p, err := identifier("prefix", prefix)
	return &TaggedUrnBuilder{
		prefix: p,
		tags:   make(map[string]Constraint),
		err:    err,
	}
}

// Constraint adds a tag of any form.
func (b *TaggedUrnBuilder) Constraint(key string, constraint Constraint) *TaggedUrnBuilder {
	if b.err != nil {
		return b
	}
	k, err := identifier("key", key)
	if err != nil {
		b.err = err
		return b
	}
	if err := checkConstraint(k, constraint); err != nil {
		b.err = err
		return b
	}
	if _, exists := b.tags[k]; exists {
		b.err = &TaggedUrnError{Code: ErrorDuplicateKey, Message: fmt.Sprintf("duplicate tag key: %s", k)}
		return b
	}
	b.tags[k] = constraint
	return b
}

// Tag adds the exact value value for key (`key=value`), taken literally:
// Tag("x", "*") is the exact value `*`, not the marker.
func (b *TaggedUrnBuilder) Tag(key, value string) *TaggedUrnBuilder {
	return b.Constraint(key, ConstraintExact{Value: value})
}

// Marker adds key as a marker (ConstraintPresent, written as just the key).
func (b *TaggedUrnBuilder) Marker(key string) *TaggedUrnBuilder {
	return b.Constraint(key, ConstraintPresent{})
}

// Build creates the final TaggedUrn
func (b *TaggedUrnBuilder) Build() (*TaggedUrn, error) {
	if b.err != nil {
		return nil, b.err
	}

	if len(b.tags) == 0 {
		return nil, &TaggedUrnError{
			Code:    ErrorInvalidFormat,
			Message: "tagged URN cannot be empty",
		}
	}

	return assemble(b.prefix, b.tags), nil
}

// BuildAllowEmpty creates the final TaggedUrn, allowing empty tags
func (b *TaggedUrnBuilder) BuildAllowEmpty() (*TaggedUrn, error) {
	if b.err != nil {
		return nil, b.err
	}
	return assemble(b.prefix, b.tags), nil
}
