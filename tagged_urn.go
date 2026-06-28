// Package taggedurn provides the fundamental tagged URN system with flat tag-based
// naming, pattern matching, and graded specificity comparison.
//
// Special pattern values:
//   - K=v: Must have key K with exact value v
//   - K=*: Must have key K with any value (presence required)
//   - K=!: Must NOT have key K (absence required)
//   - K=?: No constraint on key K (explicit don't-care)
//   - (missing): Same as K=? - no constraint
//
// Graded specificity scoring:
//   - Exact value (K=v): 3 points
//   - Must-have-any (K=*): 2 points
//   - Must-not-have (K=!): 1 point
//   - Unspecified (K=?) or missing: 0 points
package taggedurn

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

// TaggedUrn represents a tagged URN using flat, ordered tags with a configurable prefix.
//
// Examples:
//   - cap:generate;ext=pdf;out=binary;target=thumbnail
//   - cap:format=*;debug=!  (format required, debug forbidden)
//   - myapp:key="Value With Spaces"
type TaggedUrn struct {
	prefix string
	tags   map[string]string
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
	Removed      map[string]string
	Added        map[string]string
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

// Error codes for tagged URN operations
const (
	ErrorInvalidFormat         = 1
	ErrorEmptyTag              = 2
	ErrorInvalidCharacter      = 3
	ErrorInvalidTagFormat      = 4
	ErrorMissingPrefix         = 5
	ErrorDuplicateKey          = 6
	ErrorNumericKey            = 7
	ErrorUnterminatedQuote     = 8
	ErrorInvalidEscapeSequence = 9
	ErrorEmptyPrefix           = 10
	ErrorPrefixMismatch        = 11
	ErrorWhitespaceInInput     = 12
)

// Parser states for state machine.
//
// The parser handles six tag forms — the canonical alphabet of the
// constraint truth table:
//
//	| Authored                | Canonical | Stored value | Score | Reading                                  |
//	|-------------------------|-----------|--------------|------:|------------------------------------------|
//	| `?x` ≡ `x?`             | `?x`      | "?"          |     0 | no constraint                            |
//	| `?x=v` ≡ `x?=v`         | `x?=v`    | "?=v"        |     1 | absent OR (present and not v)            |
//	| `x` ≡ `x=*`             | `x`       | "*"          |     2 | present with any value                   |
//	| `!x=v` ≡ `x!=v`         | `x!=v`    | "!=v"        |     3 | present and not v                        |
//	| `x=v`                   | `x=v`     | "v"          |     4 | present and exactly v (`v ∉ {?, !, *}`)  |
//	| `!x` ≡ `x!`             | `!x`      | "!"          |     5 | absent (must-not-have)                   |
//
// Qualifier `?` or `!` may appear EITHER as a key prefix (`?x`, `!x`,
// `?x=v`, `!x=v`) OR as an infix immediately before `=` (`x?`, `x!`,
// `x?=v`, `x!=v`). The two notations are exact aliases.
//
// Disallowed (hard parse errors): `?x?`, `?x?=v`, `!x!=v`, `?!x`,
// `!?x`, `?x=*`, `!x=*`, `?x=?`, `!x=?`, mixed prefix+infix.
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

var numericPattern = regexp.MustCompile(`^[0-9]+$`)

// isValidKeyChar checks if a character is valid for a key
func isValidKeyChar(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_' || c == '-' || c == '/' || c == ':' || c == '.'
}

// isValidUnquotedValueChar checks if a character is valid for an unquoted value
func isValidUnquotedValueChar(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_' || c == '-' || c == '/' || c == ':' || c == '.' || c == '*' || c == '?' || c == '!'
}

// needsQuoting checks if a value needs quoting for serialization
func needsQuoting(value string) bool {
	for _, c := range value {
		if c == ';' || c == '=' || c == '"' || c == '\\' || c == ' ' || unicode.IsUpper(c) {
			return true
		}
	}
	return false
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

	prefix := strings.ToLower(s[:colonPos])
	tagsPart := s[colonPos+1:]
	tags := make(map[string]string)

	// Handle empty tagged URN (prefix: with no tags or just semicolon)
	if tagsPart == "" || tagsPart == ";" {
		return &TaggedUrn{prefix: prefix, tags: tags}, nil
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

	canonicalNoValue := func() string {
		switch qualifier {
		case 0:
			return "*"
		case '?':
			return "?"
		case '!':
			return "!"
		}
		panic(fmt.Sprintf("invalid qualifier rune: %q", qualifier))
	}

	canonicalizeValue := func() error {
		if qualifier == 0 {
			return nil
		}
		v := currentValue.String()
		if v == "*" || v == "?" || v == "!" {
			return &TaggedUrnError{
				Code: ErrorInvalidCharacter,
				Message: fmt.Sprintf(
					"qualifier '%c' on key '%s' cannot combine with sigil value '%s': "+
						"use a real value or drop the qualifier",
					qualifier, currentKey.String(), v),
			}
		}
		var b strings.Builder
		b.WriteRune(qualifier)
		b.WriteRune('=')
		b.WriteString(v)
		currentValue.Reset()
		currentValue.WriteString(b.String())
		return nil
	}

	finishTag := func() error {
		key := currentKey.String()
		value := currentValue.String()

		if key == "" {
			return &TaggedUrnError{
				Code:    ErrorEmptyTag,
				Message: "empty key",
			}
		}
		if value == "" {
			return &TaggedUrnError{
				Code:    ErrorEmptyTag,
				Message: fmt.Sprintf("empty value for key '%s'", key),
			}
		}

		// Check for duplicate keys
		if _, exists := tags[key]; exists {
			return &TaggedUrnError{
				Code:    ErrorDuplicateKey,
				Message: fmt.Sprintf("duplicate tag key: %s", key),
			}
		}

		// Validate key cannot be purely numeric
		if numericPattern.MatchString(key) {
			return &TaggedUrnError{
				Code:    ErrorNumericKey,
				Message: fmt.Sprintf("tag key cannot be purely numeric: %s", key),
			}
		}

		tags[key] = value
		currentKey.Reset()
		currentValue.Reset()
		qualifier = 0
		return nil
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
			} else if isValidKeyChar(c) {
				currentKey.WriteRune(unicode.ToLower(c))
				state = stateInKey
			} else {
				return nil, &TaggedUrnError{
					Code:    ErrorInvalidCharacter,
					Message: fmt.Sprintf("invalid character '%c' at position %d", c, pos),
				}
			}

		case stateAfterPrefixQuestion, stateAfterPrefixBang:
			if isValidKeyChar(c) {
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
				if currentKey.Len() == 0 {
					return nil, &TaggedUrnError{
						Code:    ErrorEmptyTag,
						Message: "empty key",
					}
				}
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
				if currentKey.Len() == 0 {
					return nil, &TaggedUrnError{
						Code:    ErrorEmptyTag,
						Message: "empty key",
					}
				}
				currentValue.WriteString(canonicalNoValue())
				if err := finishTag(); err != nil {
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
				currentValue.WriteString(canonicalNoValue())
				if err := finishTag(); err != nil {
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
				return nil, &TaggedUrnError{
					Code:    ErrorEmptyTag,
					Message: fmt.Sprintf("empty value for key '%s'", currentKey.String()),
				}
			} else if isValidUnquotedValueChar(c) {
				currentValue.WriteRune(unicode.ToLower(c))
				state = stateInUnquotedValue
			} else {
				return nil, &TaggedUrnError{
					Code:    ErrorInvalidCharacter,
					Message: fmt.Sprintf("invalid character '%c' in value at position %d", c, pos),
				}
			}

		case stateInUnquotedValue:
			if c == ';' {
				if err := canonicalizeValue(); err != nil {
					return nil, err
				}
				if err := finishTag(); err != nil {
					return nil, err
				}
				state = stateExpectingKey
			} else if isValidUnquotedValueChar(c) {
				currentValue.WriteRune(unicode.ToLower(c))
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
				if err := canonicalizeValue(); err != nil {
					return nil, err
				}
				if err := finishTag(); err != nil {
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
	case stateInUnquotedValue, stateExpectingSemiOrEnd:
		if err := canonicalizeValue(); err != nil {
			return nil, err
		}
		if err := finishTag(); err != nil {
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
	case stateInKey:
		if currentKey.Len() == 0 {
			return nil, &TaggedUrnError{
				Code:    ErrorEmptyTag,
				Message: "empty key",
			}
		}
		currentValue.WriteString(canonicalNoValue())
		if err := finishTag(); err != nil {
			return nil, err
		}
	case stateInKeyAfterQuestion, stateInKeyAfterBang:
		currentValue.WriteString(canonicalNoValue())
		if err := finishTag(); err != nil {
			return nil, err
		}
	case stateExpectingValue:
		return nil, &TaggedUrnError{
			Code:    ErrorEmptyTag,
			Message: fmt.Sprintf("empty value for key '%s'", currentKey.String()),
		}
	}

	return &TaggedUrn{prefix: prefix, tags: tags}, nil
}

// NewTaggedUrnFromTags creates a tagged URN from tags with a specified prefix (required)
// Keys are normalized to lowercase; values are preserved as-is
func NewTaggedUrnFromTags(prefix string, tags map[string]string) *TaggedUrn {
	result := make(map[string]string)
	for k, v := range tags {
		result[strings.ToLower(k)] = v
	}
	return &TaggedUrn{prefix: strings.ToLower(prefix), tags: result}
}

// Empty creates an empty tagged URN with the specified prefix (required)
func Empty(prefix string) *TaggedUrn {
	return &TaggedUrn{prefix: strings.ToLower(prefix), tags: make(map[string]string)}
}

// GetPrefix returns the prefix of this tagged URN
func (c *TaggedUrn) GetPrefix() string {
	return c.prefix
}

// GetTag returns the value of a specific tag
// Key is normalized to lowercase for lookup
func (c *TaggedUrn) GetTag(key string) (string, bool) {
	value, exists := c.tags[strings.ToLower(key)]
	return value, exists
}

// AllTags returns a copy of all tags in this URN
func (c *TaggedUrn) AllTags() map[string]string {
	result := make(map[string]string, len(c.tags))
	for k, v := range c.tags {
		result[k] = v
	}
	return result
}

// HasTag checks if this URN has a specific tag with a specific value
// Key is normalized to lowercase; value comparison is case-sensitive
func (c *TaggedUrn) HasTag(key, value string) bool {
	tagValue, exists := c.tags[strings.ToLower(key)]
	return exists && tagValue == value
}

// HasMarkerTag checks whether a marker tag (a tag whose value is "*") is
// present at the given key. Equivalent to HasTag(tagName, "*") but
// expresses authorial intent: this tag is present as a marker (a
// wildcard-valued tag that serializes as just the key), not as a
// key=value pair. Example: "cap:constrained;..." has marker tag
// "constrained".
func (c *TaggedUrn) HasMarkerTag(tagName string) bool {
	tagValue, exists := c.tags[strings.ToLower(tagName)]
	return exists && tagValue == "*"
}

// WithTag returns a new tagged URN with an added or updated tag
// Key is normalized to lowercase; value is preserved as-is
func (c *TaggedUrn) WithTag(key, value string) *TaggedUrn {
	newTags := make(map[string]string)
	for k, v := range c.tags {
		newTags[k] = v
	}
	newTags[strings.ToLower(key)] = value
	return &TaggedUrn{prefix: c.prefix, tags: newTags}
}

// WithoutTag returns a new tagged URN with a tag removed
// Key is normalized to lowercase for case-insensitive removal
func (c *TaggedUrn) WithoutTag(key string) *TaggedUrn {
	newTags := make(map[string]string)
	key = strings.ToLower(key)
	for k, v := range c.tags {
		if k != key {
			newTags[k] = v
		}
	}
	return &TaggedUrn{prefix: c.prefix, tags: newTags}
}

// Matches checks if this URN (instance) matches a pattern based on tag compatibility
//
// IMPORTANT: Both URNs must have the same prefix. Comparing URNs with
// different prefixes is a programming error and will return an error.
//
// Per-tag matching semantics:
// | Pattern Form | Interpretation              | Instance Missing | Instance = v | Instance = x≠v |
// |--------------|-----------------------------|--------------------|--------------|----------------|
// | (no entry)   | no constraint               | OK match           | OK match     | OK match       |
// | K=?          | no constraint (explicit)    | OK                 | OK           | OK             |
// | K=!          | must-not-have               | OK                 | NO           | NO             |
// | K=*          | must-have, any value        | NO                 | OK           | OK             |
// | K=v          | must-have, exact value      | NO                 | OK           | NO             |
//
// Special values work symmetrically on both instance and pattern sides.
//
// ConformsTo checks if this URN (instance) satisfies the pattern's constraints.
// Equivalent to pattern.Accepts(self).
func (c *TaggedUrn) ConformsTo(pattern *TaggedUrn) (bool, error) {
	if pattern == nil {
		return false, &TaggedUrnError{
			Code:    ErrorInvalidFormat,
			Message: "cannot match against nil pattern",
		}
	}
	return checkMatch(c.tags, c.prefix, pattern.tags, pattern.prefix)
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
	return checkMatch(instance.tags, instance.prefix, c.tags, c.prefix)
}

// checkMatch is the core matching: does instance satisfy pattern's constraints?
func checkMatch(instanceTags map[string]string, instancePrefix string, patternTags map[string]string, patternPrefix string) (bool, error) {
	if instancePrefix != patternPrefix {
		return false, &TaggedUrnError{
			Code:    ErrorPrefixMismatch,
			Message: fmt.Sprintf("cannot compare URNs with different prefixes: '%s' vs '%s'", instancePrefix, patternPrefix),
		}
	}

	allKeys := make(map[string]bool)
	for key := range instanceTags {
		allKeys[key] = true
	}
	for key := range patternTags {
		allKeys[key] = true
	}

	for key := range allKeys {
		inst, instExists := instanceTags[key]
		patt, pattExists := patternTags[key]

		var instVal, pattVal *string
		if instExists {
			instVal = &inst
		}
		if pattExists {
			pattVal = &patt
		}

		if !valuesMatch(instVal, pattVal) {
			return false, nil
		}
	}
	return true, nil
}

// valuesMatch checks if instance value matches pattern constraint
//
// Full cross-product truth table (instance = cap, pattern = request):
// | Instance | Pattern | Match? | Reason |
// |----------|---------|--------|--------|
// | (none)   | (none)  | OK     | No constraint either side |
// | (none)   | K=?     | OK     | Pattern doesn't care |
// | (none)   | K=!     | OK     | Pattern wants absent, it is |
// | (none)   | K=*     | NO     | Pattern wants present |
// | (none)   | K=v     | OK     | Instance missing = wildcard (Rust semantics) |
// | K=?      | (any)   | OK     | Instance doesn't care |
// | K=!      | (none)  | OK     | Symmetric: absent |
// | K=!      | K=?     | OK     | Pattern doesn't care |
// | K=!      | K=!     | OK     | Both want absent |
// | K=!      | K=*     | NO     | Conflict: absent vs present |
// | K=!      | K=v     | NO     | Conflict: absent vs value |
// | K=*      | (none)  | OK     | Pattern has no constraint |
// | K=*      | K=?     | OK     | Pattern doesn't care |
// | K=*      | K=!     | NO     | Conflict: present vs absent |
// | K=*      | K=*     | OK     | Both accept any presence |
// | K=*      | K=v     | OK     | Instance accepts any, v is fine |
// | K=v      | (none)  | OK     | Pattern has no constraint |
// | K=v      | K=?     | OK     | Pattern doesn't care |
// | K=v      | K=!     | NO     | Conflict: value vs absent |
// | K=v      | K=*     | OK     | Pattern wants any, v satisfies |
// | K=v      | K=v     | OK     | Exact match |
// | K=v      | K=w     | NO     | Value mismatch (v≠w) |

// formKind classifies a stored tag value into one of the six
// canonical constraint forms (plus "missing" for nil). The remaining
// matcher logic operates on (kind, optional value) pairs uniformly.
type formKind int

const (
	formMissing             formKind = iota // key absent from tag map
	formNoConstraint                        // "?" — no constraint
	formAbsentOrNotValue                    // "?=v" — absent OR (present and not v)
	formMustHaveAny                         // "*" — present with any value
	formPresentNotValue                     // "!=v" — present and not v
	formExact                               // exact value
	formMustNotHave                         // "!" — must not have
)

// classifyForm parses the stored value into (kind, raw value).
// The returned value is the inner v for ?=v and !=v, the literal
// value for exact, and "" for the sigil-only forms.
func classifyForm(value *string) (formKind, string) {
	if value == nil {
		return formMissing, ""
	}
	v := *value
	switch v {
	case "?":
		return formNoConstraint, ""
	case "*":
		return formMustHaveAny, ""
	case "!":
		return formMustNotHave, ""
	}
	if strings.HasPrefix(v, "?=") {
		return formAbsentOrNotValue, v[2:]
	}
	if strings.HasPrefix(v, "!=") {
		return formPresentNotValue, v[2:]
	}
	return formExact, v
}

// ValuesMatch evaluates the truth-table cell for (instance, pattern)
// over the six canonical forms. Both arguments are stored tag-value
// pointers (or nil to mean "key absent"). Exposed for callers (e.g.
// CapUrn's y-axis matcher) that walk tag sets themselves and need
// the same per-cell decision the tagged-URN matcher uses internally.
func ValuesMatch(inst, patt *string) bool {
	return valuesMatch(inst, patt)
}

// valuesMatch is the package-private worker. Prefer the exported
// wrapper above for callers outside this package.
func valuesMatch(inst, patt *string) bool {
	iKind, iVal := classifyForm(inst)
	pKind, pVal := classifyForm(patt)

	// Pattern unconditionally permissive.
	if pKind == formMissing || pKind == formNoConstraint {
		return true
	}

	// Instance unconditionally permissive — defers to pattern.
	if iKind == formNoConstraint {
		return true
	}

	switch pKind {
	case formMustNotHave:
		// Pattern requires absent. Only absent-side instances pass.
		switch iKind {
		case formMissing, formMustNotHave, formAbsentOrNotValue:
			return true
		default:
			return false
		}

	case formMustHaveAny:
		// Pattern requires present (any value).
		switch iKind {
		case formMissing, formAbsentOrNotValue, formMustNotHave:
			return false
		default:
			return true
		}

	case formPresentNotValue:
		// Pattern requires present-and-not-pVal.
		switch iKind {
		case formMissing, formAbsentOrNotValue, formMustNotHave:
			return false
		case formMustHaveAny, formPresentNotValue:
			return true // defer on actual value identity
		case formExact:
			return iVal != pVal
		}

	case formAbsentOrNotValue:
		// Pattern allows absent OR (present and not pVal).
		switch iKind {
		case formMissing, formAbsentOrNotValue, formMustNotHave:
			return true
		case formMustHaveAny, formPresentNotValue:
			return true // defer
		case formExact:
			return iVal != pVal
		}

	case formExact:
		// Pattern requires exact pVal.
		switch iKind {
		case formMissing, formAbsentOrNotValue, formMustNotHave:
			return false
		case formMustHaveAny:
			return true // defer
		case formPresentNotValue:
			return iVal != pVal
		case formExact:
			return iVal == pVal
		}
	}

	return false
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

// ScoreTagValue returns the per-tag truth-table specificity score.
// Applied uniformly to any stored tag value — media-URN tags, cap-tag
// y-axis, any other Tagged URN dimension. Missing keys score 0; the
// caller filters them out before calling.
//
//	"?"          -> 0   (no constraint)
//	starts "?="  -> 1   (absent or not v)
//	"*"          -> 2   (must-have-any)
//	starts "!="  -> 3   (present and not v)
//	"!"          -> 5   (must-not-have)
//	otherwise    -> 4   (exact value)
func ScoreTagValue(value string) int {
	switch value {
	case "?":
		return 0
	case "*":
		return 2
	case "!":
		return 5
	}
	if strings.HasPrefix(value, "?=") {
		return 1
	}
	if strings.HasPrefix(value, "!=") {
		return 3
	}
	return 4
}

// Specificity returns the specificity score for URN matching.
// More specific URNs have higher scores and are preferred. Sum of
// the per-tag truth-table score across every tag in the URN.
func (c *TaggedUrn) Specificity() int {
	score := 0
	for _, value := range c.tags {
		score += ScoreTagValue(value)
	}
	return score
}

// SpecificityTuple returns specificity as a tuple for tie-breaking.
// Counts how many tags fall into each non-zero form bucket, ordered
// from highest score to lowest:
//
//	(must_not_have, exact, present_not_value, must_have_any, absent_or_not_value)
//
// Compare tuples lexicographically when sum scores are equal.
func (c *TaggedUrn) SpecificityTuple() (int, int, int, int, int) {
	mustNotHave := 0
	exact := 0
	presentNotValue := 0
	mustHaveAny := 0
	absentOrNotValue := 0
	for _, value := range c.tags {
		v := value
		kind, _ := classifyForm(&v)
		switch kind {
		case formMustNotHave:
			mustNotHave++
		case formExact:
			exact++
		case formPresentNotValue:
			presentNotValue++
		case formMustHaveAny:
			mustHaveAny++
		case formAbsentOrNotValue:
			absentOrNotValue++
		case formMissing, formNoConstraint:
			// 0 points, not counted
		}
	}
	return mustNotHave, exact, presentNotValue, mustHaveAny, absentOrNotValue
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
// Returns PrefixMismatch error if prefixes differ (inherited from
// Accepts/ConformsTo — both sides return false on mismatch, but
// since we AND them, the error propagates).
func (c *TaggedUrn) IsEquivalent(other *TaggedUrn) (bool, error) {
	if other == nil {
		return false, &TaggedUrnError{
			Code:    ErrorInvalidFormat,
			Message: "cannot compare against nil URN",
		}
	}

	aAcceptsB, err := c.Accepts(other)
	if err != nil {
		return false, err
	}

	bAcceptsA, err := other.Accepts(c)
	if err != nil {
		return false, err
	}

	return aAcceptsB && bAcceptsA, nil
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
// Returns PrefixMismatch error if prefixes differ (inherited from
// Accepts/ConformsTo).
func (c *TaggedUrn) IsComparable(other *TaggedUrn) (bool, error) {
	if other == nil {
		return false, &TaggedUrnError{
			Code:    ErrorInvalidFormat,
			Message: "cannot compare against nil URN",
		}
	}

	aAcceptsB, err := c.Accepts(other)
	if err != nil {
		return false, err
	}

	bAcceptsA, err := other.Accepts(c)
	if err != nil {
		return false, err
	}

	return aAcceptsB || bAcceptsA, nil
}

// Compare returns -1, 0, or 1 using structural tagged-URN ordering:
// prefix first, then lexicographic order over sorted tag key/value pairs.
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
	if c.prefix < other.prefix {
		return -1
	}
	if c.prefix > other.prefix {
		return 1
	}
	keysA := make([]string, 0, len(c.tags))
	keysB := make([]string, 0, len(other.tags))
	for k := range c.tags {
		keysA = append(keysA, k)
	}
	for k := range other.tags {
		keysB = append(keysB, k)
	}
	sort.Strings(keysA)
	sort.Strings(keysB)
	limit := len(keysA)
	if len(keysB) < limit {
		limit = len(keysB)
	}
	for i := 0; i < limit; i++ {
		if keysA[i] < keysB[i] {
			return -1
		}
		if keysA[i] > keysB[i] {
			return 1
		}
		va := c.tags[keysA[i]]
		vb := other.tags[keysB[i]]
		if va < vb {
			return -1
		}
		if va > vb {
			return 1
		}
	}
	if len(keysA) < len(keysB) {
		return -1
	}
	if len(keysA) > len(keysB) {
		return 1
	}
	return 0
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

	removed := make(map[string]string)
	added := make(map[string]string)
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
// named in Added are inserted. Unrelated coordinates are preserved.
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

	nextTags := make(map[string]string, len(c.tags)+len(delta.Added))
	for key, value := range c.tags {
		nextTags[key] = value
	}
	for key := range delta.Removed {
		delete(nextTags, key)
	}
	for key, value := range delta.Added {
		nextTags[key] = value
	}
	return &TaggedUrn{prefix: c.prefix, tags: nextTags}, nil
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

// WithWildcardTag returns a new URN with a specific tag set to wildcard
func (c *TaggedUrn) WithWildcardTag(key string) *TaggedUrn {
	if _, exists := c.tags[key]; exists {
		return c.WithTag(key, "*")
	}
	return c
}

// Subset returns a new URN with only specified tags
func (c *TaggedUrn) Subset(keys []string) *TaggedUrn {
	newTags := make(map[string]string)
	for _, key := range keys {
		if value, exists := c.tags[key]; exists {
			newTags[key] = value
		}
	}
	return &TaggedUrn{prefix: c.prefix, tags: newTags}
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

	newTags := make(map[string]string)
	for k, v := range c.tags {
		newTags[k] = v
	}
	for k, v := range other.tags {
		newTags[k] = v
	}
	return &TaggedUrn{prefix: c.prefix, tags: newTags}, nil
}

// ToString returns the canonical string representation of this tagged URN
// Uses the stored prefix
// Tags are sorted alphabetically for consistent representation
// No trailing semicolon in canonical form
// Values are quoted only when necessary (smart quoting)
// Special value serialization:
// - * (must-have-any): serialized as value-less tag (just the key)
// - ? (unspecified): serialized as key=?
// - ! (must-not-have): serialized as key=!
func (c *TaggedUrn) ToString() string {
	if len(c.tags) == 0 {
		return fmt.Sprintf("%s:", c.prefix)
	}

	// Sort keys for canonical representation
	keys := make([]string, 0, len(c.tags))
	for key := range c.tags {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Build tag string in canonical form. Stored values map to
	// emitted forms as follows:
	//
	//   "*"           -> "k"          (bare key, must-have-any)
	//   "?"           -> "?k"         (prefix qualifier, no constraint)
	//   "!"           -> "!k"         (prefix qualifier, must-not-have)
	//   "?=v"         -> "k?=v"       (infix qualifier, absent or not v)
	//   "!=v"         -> "k!=v"       (infix qualifier, present and not v)
	//   exact "v"     -> "k=v" / "k=\"v\"" (exact value, with quoting if needed)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		value := c.tags[key]
		switch {
		case value == "*":
			parts = append(parts, key)
		case value == "?":
			parts = append(parts, fmt.Sprintf("?%s", key))
		case value == "!":
			parts = append(parts, fmt.Sprintf("!%s", key))
		case strings.HasPrefix(value, "?="):
			raw := value[2:]
			if needsQuoting(raw) {
				parts = append(parts, fmt.Sprintf("%s?=%s", key, quoteValue(raw)))
			} else {
				parts = append(parts, fmt.Sprintf("%s?=%s", key, raw))
			}
		case strings.HasPrefix(value, "!="):
			raw := value[2:]
			if needsQuoting(raw) {
				parts = append(parts, fmt.Sprintf("%s!=%s", key, quoteValue(raw)))
			} else {
				parts = append(parts, fmt.Sprintf("%s!=%s", key, raw))
			}
		default:
			if needsQuoting(value) {
				parts = append(parts, fmt.Sprintf("%s=%s", key, quoteValue(value)))
			} else {
				parts = append(parts, fmt.Sprintf("%s=%s", key, value))
			}
		}
	}

	tagsStr := strings.Join(parts, ";")
	return fmt.Sprintf("%s:%s", c.prefix, tagsStr)
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

	c.prefix = taggedUrn.prefix
	c.tags = taggedUrn.tags
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

// TaggedUrnBuilder provides a fluent builder interface for creating tagged URNs
type TaggedUrnBuilder struct {
	prefix string
	tags   map[string]string
	err    error
}

// NewTaggedUrnBuilder creates a new builder with a specified prefix (required)
func NewTaggedUrnBuilder(prefix string) *TaggedUrnBuilder {
	return &TaggedUrnBuilder{
		prefix: strings.ToLower(prefix),
		tags:   make(map[string]string),
	}
}

// Tag adds or updates a tag
// Key is normalized to lowercase; value is preserved as-is
// Tracks error if value is empty (use Marker for wildcard)
// Error is returned at Build() time
func (b *TaggedUrnBuilder) Tag(key, value string) *TaggedUrnBuilder {
	if b.err != nil {
		return b // Already have an error, don't process further
	}
	if value == "" {
		b.err = &TaggedUrnError{
			Code:    ErrorEmptyTag,
			Message: fmt.Sprintf("empty value for key '%s' (use '*' for wildcard)", key),
		}
		return b
	}
	b.tags[strings.ToLower(key)] = value
	return b
}

// Marker adds a tag with wildcard value (*)
// Key is normalized to lowercase
func (b *TaggedUrnBuilder) Marker(key string) *TaggedUrnBuilder {
	b.tags[strings.ToLower(key)] = "*"
	return b
}

// Build creates the final TaggedUrn
func (b *TaggedUrnBuilder) Build() (*TaggedUrn, error) {
	// Check for errors accumulated during building
	if b.err != nil {
		return nil, b.err
	}

	if len(b.tags) == 0 {
		return nil, &TaggedUrnError{
			Code:    ErrorInvalidFormat,
			Message: "tagged URN cannot be empty",
		}
	}

	return &TaggedUrn{prefix: b.prefix, tags: b.tags}, nil
}

// BuildAllowEmpty creates the final TaggedUrn, allowing empty tags
func (b *TaggedUrnBuilder) BuildAllowEmpty() *TaggedUrn {
	return &TaggedUrn{prefix: b.prefix, tags: b.tags}
}
