package taggedurn

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST0001: Tagged urn creation
func Test0001_TaggedUrnCreation(t *testing.T) {
	taggedUrn, err := NewTaggedUrnFromString("cap:transform;format=json;data_processing")

	assert.NoError(t, err)
	assert.NotNil(t, taggedUrn)
	assert.Equal(t, "cap", taggedUrn.GetPrefix())

	// data_processing is a valueless tag, stored as * (must-have-any)
	dataProcessing, exists := taggedUrn.GetTag("data_processing")
	assert.True(t, exists)
	assert.Equal(t, "*", dataProcessing)

	assert.True(t, taggedUrn.HasMarkerTag("transform"))

	format, exists := taggedUrn.GetTag("format")
	assert.True(t, exists)
	assert.Equal(t, "json", format)
}

// TEST0002: Custom prefix
func Test0002_CustomPrefix(t *testing.T) {
	urn, err := NewTaggedUrnFromString("myapp:generate;ext=pdf")
	require.NoError(t, err)

	assert.Equal(t, "myapp", urn.GetPrefix())
	assert.True(t, urn.HasMarkerTag("generate"))
	assert.Equal(t, "myapp:ext=pdf;generate", urn.ToString())
}

// TEST0003: Prefix case insensitive
func Test0003_PrefixCaseInsensitive(t *testing.T) {
	urn1, err := NewTaggedUrnFromString("CAP:test")
	require.NoError(t, err)
	urn2, err := NewTaggedUrnFromString("cap:test")
	require.NoError(t, err)
	urn3, err := NewTaggedUrnFromString("Cap:test")
	require.NoError(t, err)

	assert.Equal(t, "cap", urn1.GetPrefix())
	assert.Equal(t, "cap", urn2.GetPrefix())
	assert.Equal(t, "cap", urn3.GetPrefix())
	assert.True(t, urn1.Equals(urn2))
	assert.True(t, urn2.Equals(urn3))
}

// TEST0004: Prefix mismatch error
func Test0004_PrefixMismatchError(t *testing.T) {
	urn1, err := NewTaggedUrnFromString("cap:test")
	require.NoError(t, err)
	urn2, err := NewTaggedUrnFromString("myapp:test")
	require.NoError(t, err)

	_, err = urn1.ConformsTo(urn2)
	assert.Error(t, err)
	capError, ok := err.(*TaggedUrnError)
	assert.True(t, ok)
	assert.Equal(t, ErrorPrefixMismatch, capError.Code)
}

// TEST0005: Builder with prefix
func Test0005_BuilderWithPrefix(t *testing.T) {
	urn, err := NewTaggedUrnBuilder("custom").
		Tag("key", "value").
		Build()
	require.NoError(t, err)

	assert.Equal(t, "custom", urn.GetPrefix())
	assert.Equal(t, "custom:key=value", urn.ToString())
}

// TEST0006: Canonical string format
func Test0006_CanonicalStringFormat(t *testing.T) {
	taggedUrn, err := NewTaggedUrnFromString("cap:generate;target=thumbnail;ext=pdf")
	require.NoError(t, err)

	// Should be sorted alphabetically and have no trailing semicolon in canonical form
	assert.Equal(t, "cap:ext=pdf;generate;target=thumbnail", taggedUrn.ToString())
}

// TEST0007: Prefix required
func Test0007_PrefixRequired(t *testing.T) {
	// Missing prefix should fail
	taggedUrn, err := NewTaggedUrnFromString("generate;ext=pdf")
	assert.Nil(t, taggedUrn)
	assert.Error(t, err)
	assert.Equal(t, ErrorMissingPrefix, err.(*TaggedUrnError).Code)

	// Empty prefix should fail
	taggedUrn, err = NewTaggedUrnFromString(":generate")
	assert.Nil(t, taggedUrn)
	assert.Error(t, err)
	assert.Equal(t, ErrorEmptyPrefix, err.(*TaggedUrnError).Code)

	// Valid prefix should work
	taggedUrn, err = NewTaggedUrnFromString("cap:generate;ext=pdf")
	assert.NoError(t, err)
	assert.NotNil(t, taggedUrn)
	assert.True(t, taggedUrn.HasMarkerTag("generate"))

	// Case-insensitive prefix
	taggedUrn, err = NewTaggedUrnFromString("CAP:generate")
	assert.NoError(t, err)
	assert.True(t, taggedUrn.HasMarkerTag("generate"))
}

// TEST0008: Trailing semicolon equivalence
func Test0008_TrailingSemicolonEquivalence(t *testing.T) {
	// Both with and without trailing semicolon should be equivalent
	urn1, err := NewTaggedUrnFromString("cap:generate;ext=pdf")
	require.NoError(t, err)

	urn2, err := NewTaggedUrnFromString("cap:generate;ext=pdf;")
	require.NoError(t, err)

	// They should be equal
	assert.True(t, urn1.Equals(urn2))

	// They should have same hash
	assert.Equal(t, urn1.Hash(), urn2.Hash())

	// They should have same string representation (canonical form)
	assert.Equal(t, urn1.ToString(), urn2.ToString())

	// They should match each other
	matches1, err := urn1.ConformsTo(urn2)
	require.NoError(t, err)
	assert.True(t, matches1)

	matches2, err := urn2.ConformsTo(urn1)
	require.NoError(t, err)
	assert.True(t, matches2)
}

// TEST0009: Invalid tagged urn
func Test0009_InvalidTaggedUrn(t *testing.T) {
	taggedUrn, err := NewTaggedUrnFromString("")

	assert.Nil(t, taggedUrn)
	assert.Error(t, err)
	assert.Equal(t, ErrorInvalidFormat, err.(*TaggedUrnError).Code)
}

// TEST0010: Valueless tag parsing
func Test0010_ValuelessTagParsing(t *testing.T) {
	// Value-less tag is now valid and treated as wildcard
	taggedUrn, err := NewTaggedUrnFromString("cap:optimize")

	assert.NotNil(t, taggedUrn)
	assert.NoError(t, err)
	value, exists := taggedUrn.GetTag("optimize")
	assert.True(t, exists)
	assert.Equal(t, "*", value)
	assert.Equal(t, "cap:optimize", taggedUrn.ToString())
}

// TEST0011: Invalid characters
func Test0011_InvalidCharacters(t *testing.T) {
	taggedUrn, err := NewTaggedUrnFromString("cap:type@invalid=value")

	assert.Nil(t, taggedUrn)
	assert.Error(t, err)
	assert.Equal(t, ErrorInvalidCharacter, err.(*TaggedUrnError).Code)
}

// TEST0012: Tag matching
func Test0012_TagMatching(t *testing.T) {
	urn, err := NewTaggedUrnFromString("cap:generate;ext=pdf;target=thumbnail")
	require.NoError(t, err)

	// Exact match
	request1, err := NewTaggedUrnFromString("cap:generate;ext=pdf;target=thumbnail")
	require.NoError(t, err)
	matches, err := urn.ConformsTo(request1)
	require.NoError(t, err)
	assert.True(t, matches)

	// Subset match
	request2, err := NewTaggedUrnFromString("cap:generate")
	require.NoError(t, err)
	matches, err = urn.ConformsTo(request2)
	require.NoError(t, err)
	assert.True(t, matches)

	// Wildcard match
	request3, err := NewTaggedUrnFromString("cap:ext=*")
	require.NoError(t, err)
	matches, err = urn.ConformsTo(request3)
	require.NoError(t, err)
	assert.True(t, matches)

	// No match - conflicting value
	request4, err := NewTaggedUrnFromString("cap:extract")
	require.NoError(t, err)
	matches, err = urn.ConformsTo(request4)
	require.NoError(t, err)
	assert.False(t, matches)
}

// TEST0013: Missing tag handling
func Test0013_MissingTagHandling(t *testing.T) {
	// NEW SEMANTICS: Missing tag in instance means the tag doesn't exist.
	// Pattern constraints must be satisfied by instance.

	instance, err := NewTaggedUrnFromString("cap:generate")
	require.NoError(t, err)

	// Pattern with tag that instance doesn't have: NO MATCH
	// Pattern ext=pdf requires instance to have ext=pdf, but instance doesn't have ext
	pattern1, err := NewTaggedUrnFromString("cap:ext=pdf")
	require.NoError(t, err)
	matches, err := instance.ConformsTo(pattern1)
	require.NoError(t, err)
	assert.False(t, matches) // Instance missing ext, pattern wants ext=pdf

	// Pattern missing tag = no constraint: MATCH
	// Instance has generate, pattern has no constraint on op
	instance2, err := NewTaggedUrnFromString("cap:generate;ext=pdf")
	require.NoError(t, err)
	pattern2, err := NewTaggedUrnFromString("cap:generate")
	require.NoError(t, err)
	matches, err = instance2.ConformsTo(pattern2)
	require.NoError(t, err)
	assert.True(t, matches) // Instance has ext=pdf, pattern doesn't constrain ext

	// To match any value of a tag, use explicit ? or *
	pattern3, err := NewTaggedUrnFromString("cap:ext=?")
	require.NoError(t, err)
	matches, err = instance.ConformsTo(pattern3)
	require.NoError(t, err)
	assert.True(t, matches) // Instance missing ext, pattern doesn't care

	// * means must-have-any - instance must have the tag
	pattern4, err := NewTaggedUrnFromString("cap:ext=*")
	require.NoError(t, err)
	matches, err = instance.ConformsTo(pattern4)
	require.NoError(t, err)
	assert.False(t, matches) // Instance missing ext, pattern requires ext to be present
}

// TEST0014: Specificity
func Test0014_Specificity(t *testing.T) {
	// Six-form per-tag specificity ladder:
	//   ?x        : 0  (no constraint)
	//   x?=v      : 1  (absent OR not v)
	//   x (=x=*)  : 2  (must-have-any)
	//   x!=v      : 3  (present and not v)
	//   x=v       : 4  (must-have-this-value)
	//   !x        : 5  (must-not-have)

	urn1, err := NewTaggedUrnFromString("cap:op") // bare marker x=* -> 2
	require.NoError(t, err)

	urn2, err := NewTaggedUrnFromString("cap:op=generate") // exact -> 4
	require.NoError(t, err)

	urn3, err := NewTaggedUrnFromString("cap:op;ext=pdf") // marker(2) + exact(4) = 6
	require.NoError(t, err)

	urn4, err := NewTaggedUrnFromString("cap:?op") // ?x -> 0
	require.NoError(t, err)

	urn5, err := NewTaggedUrnFromString("cap:!op") // !x -> 5
	require.NoError(t, err)

	urn6, err := NewTaggedUrnFromString("cap:op?=generate") // x?=v -> 1
	require.NoError(t, err)

	urn7, err := NewTaggedUrnFromString("cap:op!=generate") // x!=v -> 3
	require.NoError(t, err)

	assert.Equal(t, 2, urn1.Specificity()) // x=* = 2
	assert.Equal(t, 4, urn2.Specificity()) // exact = 4
	assert.Equal(t, 6, urn3.Specificity()) // 2 + 4
	assert.Equal(t, 0, urn4.Specificity()) // ?x = 0
	assert.Equal(t, 5, urn5.Specificity()) // !x = 5
	assert.Equal(t, 1, urn6.Specificity()) // x?=v = 1
	assert.Equal(t, 3, urn7.Specificity()) // x!=v = 3

	moreSpecific, err := urn2.IsMoreSpecificThan(urn1)
	require.NoError(t, err)
	assert.True(t, moreSpecific) // exact(4) > marker(2)
}

// TEST0015: Compatibility
func Test0015_Compatibility(t *testing.T) {
	// TEST526: Compatibility now uses directional Accepts
	// General pattern (fewer tags) accepts specific instance (more tags)
	general, err := NewTaggedUrnFromString("cap:generate")
	require.NoError(t, err)

	specific, err := NewTaggedUrnFromString("cap:generate;ext=pdf")
	require.NoError(t, err)

	// General pattern accepts specific instance (missing tag = no constraint)
	accepts, err := general.Accepts(specific)
	require.NoError(t, err)
	assert.True(t, accepts)

	// Specific does NOT accept general (missing tag in instance = fails pattern requirement)
	accepts, err = specific.Accepts(general)
	require.NoError(t, err)
	assert.False(t, accepts)

	// Different values: neither direction accepts
	different, err := NewTaggedUrnFromString("cap:extract;ext=pdf")
	require.NoError(t, err)

	accepts, err = specific.Accepts(different)
	require.NoError(t, err)
	assert.False(t, accepts)

	accepts, err = different.Accepts(specific)
	require.NoError(t, err)
	assert.False(t, accepts)
}

// TEST0016: Convenience methods
func Test0016_ConvenienceMethods(t *testing.T) {
	urn, err := NewTaggedUrnFromString("cap:generate;ext=pdf;output=binary;target=thumbnail")
	require.NoError(t, err)

	assert.True(t, urn.HasMarkerTag("generate"))

	target, exists := urn.GetTag("target")
	assert.True(t, exists)
	assert.Equal(t, "thumbnail", target)

	format, exists := urn.GetTag("ext")
	assert.True(t, exists)
	assert.Equal(t, "pdf", format)

	output, exists := urn.GetTag("output")
	assert.True(t, exists)
	assert.Equal(t, "binary", output)
}

// TEST0017: Builder
func Test0017_Builder(t *testing.T) {
	urn, err := NewTaggedUrnBuilder("cap").
		Marker("generate").
		Tag("target", "thumbnail").
		Tag("ext", "pdf").
		Tag("output", "binary").
		Build()
	require.NoError(t, err)

	assert.True(t, urn.HasMarkerTag("generate"))

	output, exists := urn.GetTag("output")
	assert.True(t, exists)
	assert.Equal(t, "binary", output)
}

// TEST0018: With tag
func Test0018_WithTag(t *testing.T) {
	original, err := NewTaggedUrnFromString("cap:generate")
	require.NoError(t, err)

	modified := original.WithTag("ext", "pdf")

	assert.Equal(t, "cap:ext=pdf;generate", modified.ToString())

	// Original should be unchanged
	assert.Equal(t, "cap:generate", original.ToString())
}

// TEST0019: Without tag
func Test0019_WithoutTag(t *testing.T) {
	original, err := NewTaggedUrnFromString("cap:generate;ext=pdf")
	require.NoError(t, err)

	modified := original.WithoutTag("ext")

	assert.Equal(t, "cap:generate", modified.ToString())

	// Original should be unchanged
	assert.Equal(t, "cap:ext=pdf;generate", original.ToString())
}

// TEST0020: Wildcard tag
func Test0020_WildcardTag(t *testing.T) {
	urn, err := NewTaggedUrnFromString("cap:ext=pdf")
	require.NoError(t, err)

	wildcarded := urn.WithWildcardTag("ext")

	// Wildcard serializes as value-less tag
	assert.Equal(t, "cap:ext", wildcarded.ToString())

	// Test that wildcarded URN can match more requests
	request, err := NewTaggedUrnFromString("cap:ext=jpg")
	require.NoError(t, err)
	matches, err := urn.ConformsTo(request)
	require.NoError(t, err)
	assert.False(t, matches)

	wildcardRequest, err := NewTaggedUrnFromString("cap:ext")
	require.NoError(t, err)
	matches, err = wildcarded.ConformsTo(wildcardRequest)
	require.NoError(t, err)
	assert.True(t, matches)
}

// TEST0021: Subset
func Test0021_Subset(t *testing.T) {
	urn, err := NewTaggedUrnFromString("cap:generate;ext=pdf;output=binary;target=thumbnail;")
	require.NoError(t, err)

	subset := urn.Subset([]string{"type", "ext"})

	assert.Equal(t, "cap:ext=pdf", subset.ToString())
}

// TEST0022: Merge
func Test0022_Merge(t *testing.T) {
	urn1, err := NewTaggedUrnFromString("cap:generate")
	require.NoError(t, err)

	urn2, err := NewTaggedUrnFromString("cap:ext=pdf;output=binary")
	require.NoError(t, err)

	merged, err := urn1.Merge(urn2)
	require.NoError(t, err)

	assert.Equal(t, "cap:ext=pdf;generate;output=binary", merged.ToString())
}

// TEST0023: Merge prefix mismatch
func Test0023_MergePrefixMismatch(t *testing.T) {
	urn1, err := NewTaggedUrnFromString("cap:generate")
	require.NoError(t, err)

	urn2, err := NewTaggedUrnFromString("myapp:ext=pdf")
	require.NoError(t, err)

	_, err = urn1.Merge(urn2)
	assert.Error(t, err)
	capError, ok := err.(*TaggedUrnError)
	assert.True(t, ok)
	assert.Equal(t, ErrorPrefixMismatch, capError.Code)
}

// TEST0024: Equality
func Test0024_Equality(t *testing.T) {
	urn1, err := NewTaggedUrnFromString("cap:generate")
	require.NoError(t, err)

	urn2, err := NewTaggedUrnFromString("cap:generate") // different order
	require.NoError(t, err)

	urn3, err := NewTaggedUrnFromString("cap:generate;image")
	require.NoError(t, err)

	assert.True(t, urn1.Equals(urn2)) // order doesn't matter
	assert.False(t, urn1.Equals(urn3))
}

// TEST0025: Equality different prefix
func Test0025_EqualityDifferentPrefix(t *testing.T) {
	urn1, err := NewTaggedUrnFromString("cap:generate")
	require.NoError(t, err)

	urn2, err := NewTaggedUrnFromString("myapp:generate")
	require.NoError(t, err)

	assert.False(t, urn1.Equals(urn2))
}

// TEST0026: Urn matcher
func Test0026_UrnMatcher(t *testing.T) {
	matcher := &UrnMatcher{}

	urns := []*TaggedUrn{}

	urn1, err := NewTaggedUrnFromString("cap:op")
	require.NoError(t, err)
	urns = append(urns, urn1)

	urn2, err := NewTaggedUrnFromString("cap:generate")
	require.NoError(t, err)
	urns = append(urns, urn2)

	urn3, err := NewTaggedUrnFromString("cap:generate;ext=pdf")
	require.NoError(t, err)
	urns = append(urns, urn3)

	request, err := NewTaggedUrnFromString("cap:generate")
	require.NoError(t, err)

	best, err := matcher.FindBestMatch(urns, request)
	require.NoError(t, err)
	require.NotNil(t, best)

	// Most specific URN that can handle the request
	assert.Equal(t, "cap:ext=pdf;generate", best.ToString())
}

// TEST0027: Urn matcher prefix mismatch
func Test0027_UrnMatcherPrefixMismatch(t *testing.T) {
	matcher := &UrnMatcher{}

	urns := []*TaggedUrn{}

	urn1, err := NewTaggedUrnFromString("cap:generate")
	require.NoError(t, err)
	urns = append(urns, urn1)

	request, err := NewTaggedUrnFromString("myapp:generate")
	require.NoError(t, err)

	_, err = matcher.FindBestMatch(urns, request)
	assert.Error(t, err)
}

// TEST0028: J s o n serialization
func Test0028_JSONSerialization(t *testing.T) {
	original, err := NewTaggedUrnFromString("cap:generate")
	require.NoError(t, err)

	data, err := json.Marshal(original)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	var decoded TaggedUrn
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.True(t, original.Equals(&decoded))
}

// TEST0029: J s o n serialization with custom prefix
func Test0029_JSONSerializationWithCustomPrefix(t *testing.T) {
	original, err := NewTaggedUrnFromString("myapp:key=value")
	require.NoError(t, err)

	data, err := json.Marshal(original)
	assert.NoError(t, err)

	var decoded TaggedUrn
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.True(t, original.Equals(&decoded))
	assert.Equal(t, "myapp", decoded.GetPrefix())
}

// TEST0030: Unquoted values lowercased
func Test0030_UnquotedValuesLowercased(t *testing.T) {
	// Unquoted keys AND values are normalized to lowercase.
	urn, err := NewTaggedUrnFromString("cap:ext=pdf;generate;target=thumbnail;")
	require.NoError(t, err)

	assert.True(t, urn.HasMarkerTag("generate"))

	ext, exists := urn.GetTag("ext")
	assert.True(t, exists)
	assert.Equal(t, "pdf", ext)

	target, exists := urn.GetTag("target")
	assert.True(t, exists)
	assert.Equal(t, "thumbnail", target)

	// Key lookup is case-insensitive: uppercase variants of `ext`
	// resolve to the same keyed tag.
	extUpper, exists := urn.GetTag("EXT")
	assert.True(t, exists)
	assert.Equal(t, "pdf", extUpper)

	// Tag order in the source string does not affect canonical form.
	urn2, err := NewTaggedUrnFromString("cap:target=thumbnail;ext=pdf;generate;")
	require.NoError(t, err)
	assert.Equal(t, urn.ToString(), urn2.ToString())
	assert.True(t, urn.Equals(urn2))
}

// TEST0031: Quoted values preserve case
func Test0031_QuotedValuesPreserveCase(t *testing.T) {
	// Quoted values preserve their case
	urn, err := NewTaggedUrnFromString(`cap:key="Value With Spaces"`)
	require.NoError(t, err)
	value, exists := urn.GetTag("key")
	assert.True(t, exists)
	assert.Equal(t, "Value With Spaces", value)

	// Key is still lowercase
	urn2, err := NewTaggedUrnFromString(`cap:KEY="Value With Spaces"`)
	require.NoError(t, err)
	value2, exists := urn2.GetTag("key")
	assert.True(t, exists)
	assert.Equal(t, "Value With Spaces", value2)

	// Unquoted vs quoted case difference
	unquoted, err := NewTaggedUrnFromString("cap:key=UPPERCASE")
	require.NoError(t, err)
	quoted, err := NewTaggedUrnFromString(`cap:key="UPPERCASE"`)
	require.NoError(t, err)

	unquotedVal, _ := unquoted.GetTag("key")
	quotedVal, _ := quoted.GetTag("key")
	assert.Equal(t, "uppercase", unquotedVal) // lowercase
	assert.Equal(t, "UPPERCASE", quotedVal)   // preserved
	assert.False(t, unquoted.Equals(quoted))  // NOT equal
}

// TEST0032: Quoted value special chars
func Test0032_QuotedValueSpecialChars(t *testing.T) {
	// Semicolons in quoted values
	urn, err := NewTaggedUrnFromString(`cap:key="value;with;semicolons"`)
	require.NoError(t, err)
	value, exists := urn.GetTag("key")
	assert.True(t, exists)
	assert.Equal(t, "value;with;semicolons", value)

	// Equals in quoted values
	urn2, err := NewTaggedUrnFromString(`cap:key="value=with=equals"`)
	require.NoError(t, err)
	value2, exists := urn2.GetTag("key")
	assert.True(t, exists)
	assert.Equal(t, "value=with=equals", value2)

	// Spaces in quoted values
	urn3, err := NewTaggedUrnFromString(`cap:key="hello world"`)
	require.NoError(t, err)
	value3, exists := urn3.GetTag("key")
	assert.True(t, exists)
	assert.Equal(t, "hello world", value3)
}

// TEST0033: Quoted value escape sequences
func Test0033_QuotedValueEscapeSequences(t *testing.T) {
	// Escaped quotes
	urn, err := NewTaggedUrnFromString(`cap:key="value\"quoted\""`)
	require.NoError(t, err)
	value, exists := urn.GetTag("key")
	assert.True(t, exists)
	assert.Equal(t, `value"quoted"`, value)

	// Escaped backslashes
	urn2, err := NewTaggedUrnFromString(`cap:key="path\\file"`)
	require.NoError(t, err)
	value2, exists := urn2.GetTag("key")
	assert.True(t, exists)
	assert.Equal(t, `path\file`, value2)

	// Mixed escapes
	urn3, err := NewTaggedUrnFromString(`cap:key="say \"hello\\world\""`)
	require.NoError(t, err)
	value3, exists := urn3.GetTag("key")
	assert.True(t, exists)
	assert.Equal(t, `say "hello\world"`, value3)
}

// TEST0034: Mixed quoted unquoted
func Test0034_MixedQuotedUnquoted(t *testing.T) {
	urn, err := NewTaggedUrnFromString(`cap:a="Quoted";b=simple`)
	require.NoError(t, err)

	a, exists := urn.GetTag("a")
	assert.True(t, exists)
	assert.Equal(t, "Quoted", a)

	b, exists := urn.GetTag("b")
	assert.True(t, exists)
	assert.Equal(t, "simple", b)
}

// TEST0035: Unterminated quote error
func Test0035_UnterminatedQuoteError(t *testing.T) {
	urn, err := NewTaggedUrnFromString(`cap:key="unterminated`)
	assert.Nil(t, urn)
	assert.Error(t, err)
	urnError, ok := err.(*TaggedUrnError)
	assert.True(t, ok)
	assert.Equal(t, ErrorUnterminatedQuote, urnError.Code)
}

// TEST0036: Invalid escape sequence error
func Test0036_InvalidEscapeSequenceError(t *testing.T) {
	urn, err := NewTaggedUrnFromString(`cap:key="bad\n"`)
	assert.Nil(t, urn)
	assert.Error(t, err)
	urnError, ok := err.(*TaggedUrnError)
	assert.True(t, ok)
	assert.Equal(t, ErrorInvalidEscapeSequence, urnError.Code)

	// Invalid escape at end
	urn2, err := NewTaggedUrnFromString(`cap:key="bad\x"`)
	assert.Nil(t, urn2)
	assert.Error(t, err)
	urnError2, ok := err.(*TaggedUrnError)
	assert.True(t, ok)
	assert.Equal(t, ErrorInvalidEscapeSequence, urnError2.Code)
}

// TEST0037: Serialization smart quoting
func Test0037_SerializationSmartQuoting(t *testing.T) {
	// Simple lowercase value - no quoting needed
	urn, err := NewTaggedUrnBuilder("cap").Tag("key", "simple").Build()
	require.NoError(t, err)
	assert.Equal(t, "cap:key=simple", urn.ToString())

	// Value with spaces - needs quoting
	urn2, err := NewTaggedUrnBuilder("cap").Tag("key", "has spaces").Build()
	require.NoError(t, err)
	assert.Equal(t, `cap:key="has spaces"`, urn2.ToString())

	// Value with semicolons - needs quoting
	urn3, err := NewTaggedUrnBuilder("cap").Tag("key", "has;semi").Build()
	require.NoError(t, err)
	assert.Equal(t, `cap:key="has;semi"`, urn3.ToString())

	// Value with uppercase - needs quoting to preserve
	urn4, err := NewTaggedUrnBuilder("cap").Tag("key", "HasUpper").Build()
	require.NoError(t, err)
	assert.Equal(t, `cap:key="HasUpper"`, urn4.ToString())

	// Value with quotes - needs quoting and escaping
	urn5, err := NewTaggedUrnBuilder("cap").Tag("key", `has"quote`).Build()
	require.NoError(t, err)
	assert.Equal(t, `cap:key="has\"quote"`, urn5.ToString())

	// Value with backslashes - needs quoting and escaping
	urn6, err := NewTaggedUrnBuilder("cap").Tag("key", `path\file`).Build()
	require.NoError(t, err)
	assert.Equal(t, `cap:key="path\\file"`, urn6.ToString())
}

// TEST0038: Round trip simple
func Test0038_RoundTripSimple(t *testing.T) {
	original := "cap:generate;ext=pdf"
	urn, err := NewTaggedUrnFromString(original)
	require.NoError(t, err)
	serialized := urn.ToString()
	reparsed, err := NewTaggedUrnFromString(serialized)
	require.NoError(t, err)
	assert.True(t, urn.Equals(reparsed))
}

// TEST0039: Round trip quoted
func Test0039_RoundTripQuoted(t *testing.T) {
	original := `cap:key="Value With Spaces"`
	urn, err := NewTaggedUrnFromString(original)
	require.NoError(t, err)
	serialized := urn.ToString()
	reparsed, err := NewTaggedUrnFromString(serialized)
	require.NoError(t, err)
	assert.True(t, urn.Equals(reparsed))
	value, exists := reparsed.GetTag("key")
	assert.True(t, exists)
	assert.Equal(t, "Value With Spaces", value)
}

// TEST0040: Round trip escapes
func Test0040_RoundTripEscapes(t *testing.T) {
	original := `cap:key="value\"with\\escapes"`
	urn, err := NewTaggedUrnFromString(original)
	require.NoError(t, err)
	value, exists := urn.GetTag("key")
	assert.True(t, exists)
	assert.Equal(t, `value"with\escapes`, value)
	serialized := urn.ToString()
	reparsed, err := NewTaggedUrnFromString(serialized)
	require.NoError(t, err)
	assert.True(t, urn.Equals(reparsed))
}

// TEST0041: Matching case sensitive values
func Test0041_MatchingCaseSensitiveValues(t *testing.T) {
	// Values with different case should NOT match
	urn1, err := NewTaggedUrnFromString(`cap:key="Value"`)
	require.NoError(t, err)
	urn2, err := NewTaggedUrnFromString(`cap:key="value"`)
	require.NoError(t, err)

	matches1, err := urn1.ConformsTo(urn2)
	require.NoError(t, err)
	assert.False(t, matches1)

	matches2, err := urn2.ConformsTo(urn1)
	require.NoError(t, err)
	assert.False(t, matches2)

	// Same case should match
	urn3, err := NewTaggedUrnFromString(`cap:key="Value"`)
	require.NoError(t, err)
	matches3, err := urn1.ConformsTo(urn3)
	require.NoError(t, err)
	assert.True(t, matches3)
}

// TEST0042: Builder preserves case
func Test0042_BuilderPreservesCase(t *testing.T) {
	urn, err := NewTaggedUrnBuilder("cap").
		Tag("KEY", "ValueWithCase").
		Build()
	require.NoError(t, err)

	// Key is lowercase
	value, exists := urn.GetTag("key")
	assert.True(t, exists)
	assert.Equal(t, "ValueWithCase", value)

	// Value case preserved, so needs quoting
	assert.Equal(t, `cap:key="ValueWithCase"`, urn.ToString())
}

// TEST0043: Has tag case sensitive
func Test0043_HasTagCaseSensitive(t *testing.T) {
	urn, err := NewTaggedUrnFromString(`cap:key="Value"`)
	require.NoError(t, err)

	// Exact case match works
	assert.True(t, urn.HasTag("key", "Value"))

	// Different case does not match
	assert.False(t, urn.HasTag("key", "value"))
	assert.False(t, urn.HasTag("key", "VALUE"))

	// Key lookup is case-insensitive
	assert.True(t, urn.HasTag("KEY", "Value"))
	assert.True(t, urn.HasTag("Key", "Value"))
}

// TEST0044: With tag preserves value
func Test0044_WithTagPreservesValue(t *testing.T) {
	urn := NewTaggedUrnFromTags("cap", map[string]string{})
	modified := urn.WithTag("key", "ValueWithCase")

	value, exists := modified.GetTag("key")
	assert.True(t, exists)
	assert.Equal(t, "ValueWithCase", value)
}

// TEST0045: Semantic equivalence
func Test0045_SemanticEquivalence(t *testing.T) {
	// Unquoted and quoted simple lowercase values are equivalent
	unquoted, err := NewTaggedUrnFromString("cap:key=simple")
	require.NoError(t, err)
	quoted, err := NewTaggedUrnFromString(`cap:key="simple"`)
	require.NoError(t, err)
	assert.True(t, unquoted.Equals(quoted))

	// Both serialize the same way (unquoted)
	assert.Equal(t, "cap:key=simple", unquoted.ToString())
	assert.Equal(t, "cap:key=simple", quoted.ToString())
}

// TEST0046: Empty tagged urn
func Test0046_EmptyTaggedUrn(t *testing.T) {
	// Empty tagged URN is valid
	empty, err := NewTaggedUrnFromString("cap:")
	assert.NoError(t, err)
	assert.NotNil(t, empty)
	assert.Equal(t, "cap:", empty.ToString())

	// NEW SEMANTICS:
	// Empty PATTERN matches any INSTANCE (pattern has no constraints)
	// Empty INSTANCE only matches patterns that have no required tags

	specific, err := NewTaggedUrnFromString("cap:generate;ext=pdf")
	assert.NoError(t, err)

	// Empty instance vs specific pattern: NO MATCH
	// Pattern requires generate and ext=pdf, instance doesn't have them
	matches, err := empty.ConformsTo(specific)
	require.NoError(t, err)
	assert.False(t, matches)

	// Specific instance vs empty pattern: MATCH
	// Pattern has no constraints, instance can have anything
	matches, err = specific.ConformsTo(empty)
	require.NoError(t, err)
	assert.True(t, matches)

	// Empty instance vs empty pattern: MATCH
	matches, err = empty.ConformsTo(empty)
	require.NoError(t, err)
	assert.True(t, matches)

	// With trailing semicolon
	empty2, err := NewTaggedUrnFromString("cap:;")
	assert.NoError(t, err)
	assert.Equal(t, "cap:", empty2.ToString())
}

// TEST0047: Empty with custom prefix
func Test0047_EmptyWithCustomPrefix(t *testing.T) {
	empty, err := NewTaggedUrnFromString("myapp:")
	require.NoError(t, err)
	assert.Equal(t, "myapp", empty.GetPrefix())
	assert.Equal(t, "myapp:", empty.ToString())
}

// TEST0048: Extended character support
func Test0048_ExtendedCharacterSupport(t *testing.T) {
	// Test forward slashes and colons in tag components
	urn, err := NewTaggedUrnFromString("cap:url=https://example_org/api;path=/some/file")
	assert.NoError(t, err)
	assert.NotNil(t, urn)

	url, exists := urn.GetTag("url")
	assert.True(t, exists)
	assert.Equal(t, "https://example_org/api", url)

	path, exists := urn.GetTag("path")
	assert.True(t, exists)
	assert.Equal(t, "/some/file", path)
}

// TEST0049: Wildcard restrictions
func Test0049_WildcardRestrictions(t *testing.T) {
	// Wildcard should be rejected in keys
	invalidKey, err := NewTaggedUrnFromString("cap:*=value")
	assert.Error(t, err)
	assert.Nil(t, invalidKey)
	urnError, ok := err.(*TaggedUrnError)
	assert.True(t, ok)
	assert.Equal(t, ErrorInvalidCharacter, urnError.Code)

	// Wildcard should be accepted in values
	validValue, err := NewTaggedUrnFromString("cap:key=*")
	assert.NoError(t, err)
	assert.NotNil(t, validValue)

	value, exists := validValue.GetTag("key")
	assert.True(t, exists)
	assert.Equal(t, "*", value)
}

// TEST0050: Duplicate key rejection
func Test0050_DuplicateKeyRejection(t *testing.T) {
	// Duplicate keys should be rejected
	duplicate, err := NewTaggedUrnFromString("cap:key=value1;key=value2")
	assert.Error(t, err)
	assert.Nil(t, duplicate)
	urnError, ok := err.(*TaggedUrnError)
	assert.True(t, ok)
	assert.Equal(t, ErrorDuplicateKey, urnError.Code)
}

// TEST0051: Numeric key restriction
func Test0051_NumericKeyRestriction(t *testing.T) {
	// Pure numeric keys should be rejected
	numericKey, err := NewTaggedUrnFromString("cap:123=value")
	assert.Error(t, err)
	assert.Nil(t, numericKey)
	urnError, ok := err.(*TaggedUrnError)
	assert.True(t, ok)
	assert.Equal(t, ErrorNumericKey, urnError.Code)

	// Mixed alphanumeric keys should be allowed
	mixedKey1, err := NewTaggedUrnFromString("cap:key123=value")
	assert.NoError(t, err)
	assert.NotNil(t, mixedKey1)

	mixedKey2, err := NewTaggedUrnFromString("cap:123key=value")
	assert.NoError(t, err)
	assert.NotNil(t, mixedKey2)

	// Pure numeric values should be allowed
	numericValue, err := NewTaggedUrnFromString("cap:key=123")
	assert.NoError(t, err)
	assert.NotNil(t, numericValue)

	value, exists := numericValue.GetTag("key")
	assert.True(t, exists)
	assert.Equal(t, "123", value)
}

// TEST0052: Empty value error
func Test0052_EmptyValueError(t *testing.T) {
	urn, err := NewTaggedUrnFromString("cap:key=")
	assert.Nil(t, urn)
	assert.Error(t, err)

	urn2, err := NewTaggedUrnFromString("cap:key=;other=value")
	assert.Nil(t, urn2)
	assert.Error(t, err)
}

// TEST0053: Matching different prefixes error
func Test0053_MatchingDifferentPrefixesError(t *testing.T) {
	// URNs with different prefixes should cause an error, not just return false
	urn1, err := NewTaggedUrnFromString("cap:test")
	require.NoError(t, err)
	urn2, err := NewTaggedUrnFromString("other:test")
	require.NoError(t, err)

	_, err = urn1.ConformsTo(urn2)
	assert.Error(t, err)

	_, err = urn1.Accepts(urn2)
	assert.Error(t, err)

	_, err = urn1.IsMoreSpecificThan(urn2)
	assert.Error(t, err)
}

// ============================================================================
// MATCHING SEMANTICS SPECIFICATION TESTS
// These 9 tests verify the exact matching semantics from RULES.md Sections 12-17
// All implementations (Rust, Go, JS, ObjC) must pass these identically
// ============================================================================

func Test0054_MatchingSemantics_Test1_ExactMatch(t *testing.T) {
	// Test 1: Exact match
	// URN:     cap:generate;ext=pdf
	// Request: cap:generate;ext=pdf
	// Result:  MATCH
	urn, err := NewTaggedUrnFromString("cap:generate;ext=pdf")
	require.NoError(t, err)

	request, err := NewTaggedUrnFromString("cap:generate;ext=pdf")
	require.NoError(t, err)

	matches, err := urn.ConformsTo(request)
	require.NoError(t, err)
	assert.True(t, matches, "Test 1: Exact match should succeed")
}

// TEST0055: Matching semantics  test2  instance missing tag
func Test0055_MatchingSemantics_Test2_InstanceMissingTag(t *testing.T) {
	// Test 2: Instance missing tag
	// Instance: cap:generate;in=media:;out=media:
	// Pattern:  cap:generate;ext=pdf
	// Result:   NO MATCH (pattern requires ext=pdf, instance doesn't have ext)
	//
	// NEW SEMANTICS: Missing tag in instance means it doesn't exist.
	// Pattern K=v requires instance to have K=v.
	instance, err := NewTaggedUrnFromString("cap:generate")
	require.NoError(t, err)

	pattern, err := NewTaggedUrnFromString("cap:generate;ext=pdf")
	require.NoError(t, err)

	matches, err := instance.ConformsTo(pattern)
	require.NoError(t, err)
	assert.False(t, matches, "Test 2: Instance missing tag should NOT match when pattern requires it")

	// To accept any ext (or missing), use pattern with ext=?
	patternOptional, err := NewTaggedUrnFromString("cap:generate;ext=?")
	require.NoError(t, err)
	matches, err = instance.ConformsTo(patternOptional)
	require.NoError(t, err)
	assert.True(t, matches, "Pattern with ext=? should match instance without ext")
}

// TEST0056: Matching semantics  test3  urn has extra tag
func Test0056_MatchingSemantics_Test3_UrnHasExtraTag(t *testing.T) {
	// Test 3: URN has extra tag
	// URN:     cap:generate;ext=pdf;version=2
	// Request: cap:generate;ext=pdf
	// Result:  MATCH (request doesn't constrain version)
	urn, err := NewTaggedUrnFromString("cap:generate;ext=pdf;version=2")
	require.NoError(t, err)

	request, err := NewTaggedUrnFromString("cap:generate;ext=pdf")
	require.NoError(t, err)

	matches, err := urn.ConformsTo(request)
	require.NoError(t, err)
	assert.True(t, matches, "Test 3: URN with extra tag should match")
}

// TEST0057: Matching semantics  test4  request has wildcard
func Test0057_MatchingSemantics_Test4_RequestHasWildcard(t *testing.T) {
	// Test 4: Request has wildcard
	// URN:     cap:generate;ext=pdf
	// Request: cap:generate;ext=*
	// Result:  MATCH (request accepts any ext)
	urn, err := NewTaggedUrnFromString("cap:generate;ext=pdf")
	require.NoError(t, err)

	request, err := NewTaggedUrnFromString("cap:generate;ext=*")
	require.NoError(t, err)

	matches, err := urn.ConformsTo(request)
	require.NoError(t, err)
	assert.True(t, matches, "Test 4: Request wildcard should match")
}

// TEST0058: Matching semantics  test5  urn has wildcard
func Test0058_MatchingSemantics_Test5_UrnHasWildcard(t *testing.T) {
	// Test 5: URN has wildcard
	// URN:     cap:generate;ext=*
	// Request: cap:generate;ext=pdf
	// Result:  MATCH (URN handles any ext)
	urn, err := NewTaggedUrnFromString("cap:generate;ext=*")
	require.NoError(t, err)

	request, err := NewTaggedUrnFromString("cap:generate;ext=pdf")
	require.NoError(t, err)

	matches, err := urn.ConformsTo(request)
	require.NoError(t, err)
	assert.True(t, matches, "Test 5: URN wildcard should match")
}

// TEST0059: Matching semantics  test6  value mismatch
func Test0059_MatchingSemantics_Test6_ValueMismatch(t *testing.T) {
	// Test 6: Value mismatch
	// URN:     cap:generate;ext=pdf
	// Request: cap:generate;ext=docx
	// Result:  NO MATCH
	urn, err := NewTaggedUrnFromString("cap:generate;ext=pdf")
	require.NoError(t, err)

	request, err := NewTaggedUrnFromString("cap:generate;ext=docx")
	require.NoError(t, err)

	matches, err := urn.ConformsTo(request)
	require.NoError(t, err)
	assert.False(t, matches, "Test 6: Value mismatch should not match")
}

// TEST0060: Matching semantics  test7  pattern has extra tag
func Test0060_MatchingSemantics_Test7_PatternHasExtraTag(t *testing.T) {
	// Test 7: Pattern has extra tag that instance doesn't have
	// Instance: cap:generate-thumbnail;out="media:binary"
	// Pattern:  cap:generate-thumbnail;out="media:binary";ext=wav
	// Result:   NO MATCH (pattern requires ext=wav, instance doesn't have ext)
	//
	// NEW SEMANTICS: Pattern K=v requires instance to have K=v
	instance, err := NewTaggedUrnFromString(`cap:generate-thumbnail;out="media:binary"`)
	require.NoError(t, err)

	pattern, err := NewTaggedUrnFromString(`cap:generate-thumbnail;out="media:binary";ext=wav`)
	require.NoError(t, err)

	matches, err := instance.ConformsTo(pattern)
	require.NoError(t, err)
	assert.False(t, matches, "Test 7: Instance missing ext should NOT match when pattern requires ext=wav")

	// Instance vs pattern that doesn't constrain ext: MATCH
	patternNoExt, err := NewTaggedUrnFromString(`cap:generate-thumbnail;out="media:binary"`)
	require.NoError(t, err)
	matches, err = instance.ConformsTo(patternNoExt)
	require.NoError(t, err)
	assert.True(t, matches)
}

// TEST0061: Matching semantics  test8  empty pattern matches anything
func Test0061_MatchingSemantics_Test8_EmptyPatternMatchesAnything(t *testing.T) {
	// Test 8: Empty PATTERN matches any INSTANCE
	// Instance: cap:generate;ext=pdf
	// Pattern:  cap:
	// Result:   MATCH (pattern has no constraints)
	//
	// NEW SEMANTICS: Empty pattern = no constraints = matches any instance
	// But empty instance only matches patterns that don't require tags
	instance, err := NewTaggedUrnFromString("cap:generate;ext=pdf")
	require.NoError(t, err)

	emptyPattern, err := NewTaggedUrnFromString("cap:")
	require.NoError(t, err)

	matches, err := instance.ConformsTo(emptyPattern)
	require.NoError(t, err)
	assert.True(t, matches, "Test 8: Any instance should match empty pattern")

	// Empty instance vs pattern with requirements: NO MATCH
	emptyInstance, err := NewTaggedUrnFromString("cap:")
	require.NoError(t, err)
	pattern, err := NewTaggedUrnFromString("cap:generate;ext=pdf")
	require.NoError(t, err)
	matches, err = emptyInstance.ConformsTo(pattern)
	require.NoError(t, err)
	assert.False(t, matches, "Empty instance should NOT match pattern with requirements")
}

// TEST0062: Matching semantics  test9  cross dimension constraints
func Test0062_MatchingSemantics_Test9_CrossDimensionConstraints(t *testing.T) {
	// Test 9: Cross-dimension constraints
	// Instance: cap:generate;in=media:;out=media:
	// Pattern:  cap:ext=pdf
	// Result:   NO MATCH (pattern requires ext=pdf, instance doesn't have ext)
	//
	// NEW SEMANTICS: Pattern K=v requires instance to have K=v
	instance, err := NewTaggedUrnFromString("cap:generate")
	require.NoError(t, err)

	pattern, err := NewTaggedUrnFromString("cap:ext=pdf")
	require.NoError(t, err)

	matches, err := instance.ConformsTo(pattern)
	require.NoError(t, err)
	assert.False(t, matches, "Test 9: Instance without ext should NOT match pattern requiring ext")

	// Instance with ext vs pattern with different tag only: MATCH
	instance2, err := NewTaggedUrnFromString("cap:generate;ext=pdf")
	require.NoError(t, err)
	pattern2, err := NewTaggedUrnFromString("cap:ext=pdf")
	require.NoError(t, err)
	matches, err = instance2.ConformsTo(pattern2)
	require.NoError(t, err)
	assert.True(t, matches, "Instance with ext=pdf should match pattern requiring ext=pdf")
}

// ============================================================================
// VALUE-LESS TAG TESTS
// Value-less tags are equivalent to wildcard tags (key=*)
// ============================================================================

func Test0063_ValuelessTagParsingSingle(t *testing.T) {
	// Single value-less tag
	urn, err := NewTaggedUrnFromString("cap:optimize")
	require.NoError(t, err)

	value, exists := urn.GetTag("optimize")
	assert.True(t, exists)
	assert.Equal(t, "*", value)
	// Serializes as value-less (no =*)
	assert.Equal(t, "cap:optimize", urn.ToString())
}

// TEST0064: Valueless tag parsing multiple
func Test0064_ValuelessTagParsingMultiple(t *testing.T) {
	// Multiple value-less tags
	urn, err := NewTaggedUrnFromString("cap:fast;optimize;secure")
	require.NoError(t, err)

	fast, exists := urn.GetTag("fast")
	assert.True(t, exists)
	assert.Equal(t, "*", fast)

	optimize, exists := urn.GetTag("optimize")
	assert.True(t, exists)
	assert.Equal(t, "*", optimize)

	secure, exists := urn.GetTag("secure")
	assert.True(t, exists)
	assert.Equal(t, "*", secure)

	// Serializes alphabetically as value-less
	assert.Equal(t, "cap:fast;optimize;secure", urn.ToString())
}

// TEST0065: Valueless tag mixed with valued
func Test0065_ValuelessTagMixedWithValued(t *testing.T) {
	// Mix of value-less and valued tags
	urn, err := NewTaggedUrnFromString("cap:generate;optimize;ext=pdf;secure")
	require.NoError(t, err)

	assert.True(t, urn.HasMarkerTag("generate"))

	optimize, exists := urn.GetTag("optimize")
	assert.True(t, exists)
	assert.Equal(t, "*", optimize)

	ext, exists := urn.GetTag("ext")
	assert.True(t, exists)
	assert.Equal(t, "pdf", ext)

	secure, exists := urn.GetTag("secure")
	assert.True(t, exists)
	assert.Equal(t, "*", secure)

	// Serializes alphabetically
	assert.Equal(t, "cap:ext=pdf;generate;optimize;secure", urn.ToString())
}

// TEST0066: Valueless tag at end
func Test0066_ValuelessTagAtEnd(t *testing.T) {
	// Value-less tag at the end (no trailing semicolon)
	urn, err := NewTaggedUrnFromString("cap:generate;optimize")
	require.NoError(t, err)

	assert.True(t, urn.HasMarkerTag("generate"))

	optimize, exists := urn.GetTag("optimize")
	assert.True(t, exists)
	assert.Equal(t, "*", optimize)

	assert.Equal(t, "cap:generate;optimize", urn.ToString())
}

// TEST0067: Valueless tag equivalence to wildcard
func Test0067_ValuelessTagEquivalenceToWildcard(t *testing.T) {
	// Value-less tag is equivalent to explicit wildcard
	valueless, err := NewTaggedUrnFromString("cap:ext")
	require.NoError(t, err)

	wildcard, err := NewTaggedUrnFromString("cap:ext=*")
	require.NoError(t, err)

	assert.True(t, valueless.Equals(wildcard))
	// Both serialize to value-less form
	assert.Equal(t, "cap:ext", valueless.ToString())
	assert.Equal(t, "cap:ext", wildcard.ToString())
}

// TEST0068: Valueless tag matching
func Test0068_ValuelessTagMatching(t *testing.T) {
	// Value-less tag (wildcard) matches any value
	urn, err := NewTaggedUrnFromString("cap:generate;ext")
	require.NoError(t, err)

	requestPdf, err := NewTaggedUrnFromString("cap:generate;ext=pdf")
	require.NoError(t, err)
	requestDocx, err := NewTaggedUrnFromString("cap:generate;ext=docx")
	require.NoError(t, err)
	requestAny, err := NewTaggedUrnFromString("cap:generate;ext=anything")
	require.NoError(t, err)

	matches, err := urn.ConformsTo(requestPdf)
	require.NoError(t, err)
	assert.True(t, matches)

	matches, err = urn.ConformsTo(requestDocx)
	require.NoError(t, err)
	assert.True(t, matches)

	matches, err = urn.ConformsTo(requestAny)
	require.NoError(t, err)
	assert.True(t, matches)
}

// TEST0069: Valueless tag in pattern
func Test0069_ValuelessTagInPattern(t *testing.T) {
	// Pattern with value-less tag (K=*) requires instance to have the tag
	pattern, err := NewTaggedUrnFromString("cap:generate;ext")
	require.NoError(t, err)

	instancePdf, err := NewTaggedUrnFromString("cap:generate;ext=pdf")
	require.NoError(t, err)
	instanceDocx, err := NewTaggedUrnFromString("cap:generate;ext=docx")
	require.NoError(t, err)
	instanceMissing, err := NewTaggedUrnFromString("cap:generate")
	require.NoError(t, err)

	// NEW SEMANTICS: K=* (valueless tag) means must-have-any
	matches, err := instancePdf.ConformsTo(pattern)
	require.NoError(t, err)
	assert.True(t, matches) // Has ext=pdf

	matches, err = instanceDocx.ConformsTo(pattern)
	require.NoError(t, err)
	assert.True(t, matches) // Has ext=docx

	matches, err = instanceMissing.ConformsTo(pattern)
	require.NoError(t, err)
	assert.False(t, matches) // Missing ext, pattern requires it

	// To accept missing ext, use ? instead
	patternOptional, err := NewTaggedUrnFromString("cap:generate;ext=?")
	require.NoError(t, err)
	matches, err = instanceMissing.ConformsTo(patternOptional)
	require.NoError(t, err)
	assert.True(t, matches)
}

// TEST0070: Valueless tag specificity
func Test0070_ValuelessTagSpecificity(t *testing.T) {
	// Six-form ladder: ?x=0, x?=v=1, x=*=2, x!=v=3, x=v=4, !x=5.
	urn1, err := NewTaggedUrnFromString("cap:generate") // 1 marker (=*)
	require.NoError(t, err)
	urn2, err := NewTaggedUrnFromString("cap:generate;optimize") // 2 markers
	require.NoError(t, err)
	urn3, err := NewTaggedUrnFromString("cap:generate;ext=pdf") // marker + exact
	require.NoError(t, err)

	assert.Equal(t, 2, urn1.Specificity()) // 1 marker = 2
	assert.Equal(t, 4, urn2.Specificity()) // 2 markers = 2 + 2 = 4
	assert.Equal(t, 6, urn3.Specificity()) // 1 marker + 1 exact = 2 + 4 = 6
}

// TEST0071: Valueless tag roundtrip
func Test0071_ValuelessTagRoundtrip(t *testing.T) {
	// Round-trip parsing and serialization
	original := "cap:ext=pdf;generate;optimize;secure"
	urn, err := NewTaggedUrnFromString(original)
	require.NoError(t, err)
	serialized := urn.ToString()
	reparsed, err := NewTaggedUrnFromString(serialized)
	require.NoError(t, err)
	assert.True(t, urn.Equals(reparsed))
	assert.Equal(t, original, serialized)
}

// TEST0072: Valueless tag case normalization
func Test0072_ValuelessTagCaseNormalization(t *testing.T) {
	// Value-less tags are normalized to lowercase like other keys
	urn, err := NewTaggedUrnFromString("cap:OPTIMIZE;Fast;SECURE")
	require.NoError(t, err)

	optimize, exists := urn.GetTag("optimize")
	assert.True(t, exists)
	assert.Equal(t, "*", optimize)

	fast, exists := urn.GetTag("fast")
	assert.True(t, exists)
	assert.Equal(t, "*", fast)

	secure, exists := urn.GetTag("secure")
	assert.True(t, exists)
	assert.Equal(t, "*", secure)

	assert.Equal(t, "cap:fast;optimize;secure", urn.ToString())
}

// TEST0073: Empty value still error
func Test0073_EmptyValueStillError(t *testing.T) {
	// Empty value with = is still an error (different from value-less)
	urn, err := NewTaggedUrnFromString("cap:key=")
	assert.Nil(t, urn)
	assert.Error(t, err)

	urn2, err := NewTaggedUrnFromString("cap:key=;other=value")
	assert.Nil(t, urn2)
	assert.Error(t, err)
}

// TEST0074: Valueless tag compatibility
func Test0074_ValuelessTagCompatibility(t *testing.T) {
	// TEST564: Value-less tag (ext=*) accepts any specific value
	pattern, err := NewTaggedUrnFromString("cap:generate;ext")
	require.NoError(t, err)
	instance1, err := NewTaggedUrnFromString("cap:generate;ext=pdf")
	require.NoError(t, err)
	instance2, err := NewTaggedUrnFromString("cap:generate;ext=docx")
	require.NoError(t, err)

	// Pattern with wildcard tag accepts specific instances
	accepts, err := pattern.Accepts(instance1)
	require.NoError(t, err)
	assert.True(t, accepts)

	accepts, err = pattern.Accepts(instance2)
	require.NoError(t, err)
	assert.True(t, accepts)

	// Two different specific values: neither accepts the other
	accepts, err = instance1.Accepts(instance2)
	require.NoError(t, err)
	assert.False(t, accepts)

	accepts, err = instance2.Accepts(instance1)
	require.NoError(t, err)
	assert.False(t, accepts)
}

// TEST0075: Valueless numeric key still rejected
func Test0075_ValuelessNumericKeyStillRejected(t *testing.T) {
	// Purely numeric keys are still rejected for value-less tags
	urn, err := NewTaggedUrnFromString("cap:123")
	assert.Nil(t, urn)
	assert.Error(t, err)

	urn2, err := NewTaggedUrnFromString("cap:generate;456")
	assert.Nil(t, urn2)
	assert.Error(t, err)
}

// TEST0076: Whitespace in input rejected
func Test0076_WhitespaceInInputRejected(t *testing.T) {
	// Leading whitespace fails hard
	urn, err := NewTaggedUrnFromString(" cap:test")
	assert.Nil(t, urn)
	assert.Error(t, err)
	urnError, ok := err.(*TaggedUrnError)
	assert.True(t, ok)
	assert.Equal(t, ErrorWhitespaceInInput, urnError.Code)

	// Trailing whitespace fails hard
	urn, err = NewTaggedUrnFromString("cap:in=media:;out=media:;test ")
	assert.Nil(t, urn)
	assert.Error(t, err)
	urnError, ok = err.(*TaggedUrnError)
	assert.True(t, ok)
	assert.Equal(t, ErrorWhitespaceInInput, urnError.Code)

	// Both leading and trailing whitespace fails hard
	urn, err = NewTaggedUrnFromString(" cap:in=media:;out=media:;test ")
	assert.Nil(t, urn)
	assert.Error(t, err)
	urnError, ok = err.(*TaggedUrnError)
	assert.True(t, ok)
	assert.Equal(t, ErrorWhitespaceInInput, urnError.Code)

	// Tab and newline also count as whitespace
	urn, err = NewTaggedUrnFromString("\tcap:test")
	assert.Nil(t, urn)
	assert.Error(t, err)

	urn, err = NewTaggedUrnFromString("cap:in=media:;out=media:;test\n")
	assert.Nil(t, urn)
	assert.Error(t, err)

	// Clean input works
	urn, err = NewTaggedUrnFromString("cap:test")
	assert.NoError(t, err)
	assert.NotNil(t, urn)
	assert.True(t, urn.HasMarkerTag("test"))
}

// ============================================================================
// NEW SEMANTICS TESTS: ? (unspecified) and ! (must-not-have)
// ============================================================================

func Test0077_UnspecifiedQuestionMarkParsing(t *testing.T) {
	// ? parses as unspecified
	urn, err := NewTaggedUrnFromString("cap:ext=?")
	require.NoError(t, err)

	value, exists := urn.GetTag("ext")
	assert.True(t, exists)
	assert.Equal(t, "?", value)
	// All three input aliases (?x, x?, x=?) parse to stored value
	// "?" and serialize as the canonical prefix form `?x`.
	assert.Equal(t, "cap:?ext", urn.ToString())
}

// TEST0078: Must not have exclamation parsing
func Test0078_MustNotHaveExclamationParsing(t *testing.T) {
	urn, err := NewTaggedUrnFromString("cap:ext=!")
	require.NoError(t, err)

	value, exists := urn.GetTag("ext")
	assert.True(t, exists)
	assert.Equal(t, "!", value)
	// All three input aliases (!x, x!, x=!) parse to stored value
	// "!" and serialize as the canonical prefix form `!x`.
	assert.Equal(t, "cap:!ext", urn.ToString())
}

// TEST0079: Question mark pattern matches anything
func Test0079_QuestionMarkPatternMatchesAnything(t *testing.T) {
	// Pattern with K=? matches any instance (with or without K)
	pattern, err := NewTaggedUrnFromString("cap:ext=?")
	require.NoError(t, err)

	instancePdf, _ := NewTaggedUrnFromString("cap:ext=pdf")
	instanceDocx, _ := NewTaggedUrnFromString("cap:ext=docx")
	instanceMissing, _ := NewTaggedUrnFromString("cap:")
	instanceWildcard, _ := NewTaggedUrnFromString("cap:ext=*")
	instanceMustNot, _ := NewTaggedUrnFromString("cap:ext=!")

	matches, _ := instancePdf.ConformsTo(pattern)
	assert.True(t, matches, "ext=pdf should match ext=?")

	matches, _ = instanceDocx.ConformsTo(pattern)
	assert.True(t, matches, "ext=docx should match ext=?")

	matches, _ = instanceMissing.ConformsTo(pattern)
	assert.True(t, matches, "(no ext) should match ext=?")

	matches, _ = instanceWildcard.ConformsTo(pattern)
	assert.True(t, matches, "ext=* should match ext=?")

	matches, _ = instanceMustNot.ConformsTo(pattern)
	assert.True(t, matches, "ext=! should match ext=?")
}

// TEST0080: Question mark in instance
func Test0080_QuestionMarkInInstance(t *testing.T) {
	// Instance with K=? matches any pattern constraint
	instance, err := NewTaggedUrnFromString("cap:ext=?")
	require.NoError(t, err)

	patternPdf, _ := NewTaggedUrnFromString("cap:ext=pdf")
	patternWildcard, _ := NewTaggedUrnFromString("cap:ext=*")
	patternMustNot, _ := NewTaggedUrnFromString("cap:ext=!")
	patternQuestion, _ := NewTaggedUrnFromString("cap:ext=?")
	patternMissing, _ := NewTaggedUrnFromString("cap:")

	matches, _ := instance.ConformsTo(patternPdf)
	assert.True(t, matches, "ext=? should match ext=pdf")

	matches, _ = instance.ConformsTo(patternWildcard)
	assert.True(t, matches, "ext=? should match ext=*")

	matches, _ = instance.ConformsTo(patternMustNot)
	assert.True(t, matches, "ext=? should match ext=!")

	matches, _ = instance.ConformsTo(patternQuestion)
	assert.True(t, matches, "ext=? should match ext=?")

	matches, _ = instance.ConformsTo(patternMissing)
	assert.True(t, matches, "ext=? should match (no ext)")
}

// TEST0081: Must not have pattern requires absent
func Test0081_MustNotHavePatternRequiresAbsent(t *testing.T) {
	// Pattern with K=! requires instance to NOT have K
	pattern, err := NewTaggedUrnFromString("cap:ext=!")
	require.NoError(t, err)

	instanceMissing, _ := NewTaggedUrnFromString("cap:")
	instancePdf, _ := NewTaggedUrnFromString("cap:ext=pdf")
	instanceWildcard, _ := NewTaggedUrnFromString("cap:ext=*")
	instanceMustNot, _ := NewTaggedUrnFromString("cap:ext=!")

	matches, _ := instanceMissing.ConformsTo(pattern)
	assert.True(t, matches, "(no ext) should match ext=!")

	matches, _ = instancePdf.ConformsTo(pattern)
	assert.False(t, matches, "ext=pdf should NOT match ext=!")

	matches, _ = instanceWildcard.ConformsTo(pattern)
	assert.False(t, matches, "ext=* should NOT match ext=!")

	matches, _ = instanceMustNot.ConformsTo(pattern)
	assert.True(t, matches, "ext=! should match ext=!")
}

// TEST0082: Must not have in instance
func Test0082_MustNotHaveInInstance(t *testing.T) {
	// Instance with K=! conflicts with patterns requiring K
	instance, err := NewTaggedUrnFromString("cap:ext=!")
	require.NoError(t, err)

	patternPdf, _ := NewTaggedUrnFromString("cap:ext=pdf")
	patternWildcard, _ := NewTaggedUrnFromString("cap:ext=*")
	patternMustNot, _ := NewTaggedUrnFromString("cap:ext=!")
	patternQuestion, _ := NewTaggedUrnFromString("cap:ext=?")
	patternMissing, _ := NewTaggedUrnFromString("cap:")

	matches, _ := instance.ConformsTo(patternPdf)
	assert.False(t, matches, "ext=! should NOT match ext=pdf")

	matches, _ = instance.ConformsTo(patternWildcard)
	assert.False(t, matches, "ext=! should NOT match ext=*")

	matches, _ = instance.ConformsTo(patternMustNot)
	assert.True(t, matches, "ext=! should match ext=!")

	matches, _ = instance.ConformsTo(patternQuestion)
	assert.True(t, matches, "ext=! should match ext=?")

	matches, _ = instance.ConformsTo(patternMissing)
	assert.True(t, matches, "ext=! should match (no ext)")
}

// TEST0083: Full cross product matching
func Test0083_FullCrossProductMatching(t *testing.T) {
	// Comprehensive test of all instance/pattern combinations
	check := func(instance, pattern string, expected bool, msg string) {
		inst, err := NewTaggedUrnFromString(instance)
		require.NoError(t, err)
		patt, err := NewTaggedUrnFromString(pattern)
		require.NoError(t, err)
		matches, err := inst.ConformsTo(patt)
		require.NoError(t, err)
		assert.Equal(t, expected, matches, "%s: instance=%s, pattern=%s", msg, instance, pattern)
	}

	// Instance missing, Pattern variations
	check("cap:", "cap:", true, "(none)/(none)")
	check("cap:", "cap:k=?", true, "(none)/K=?")
	check("cap:", "cap:k=!", true, "(none)/K=!")
	check("cap:", "cap:k", false, "(none)/K=*")
	check("cap:", "cap:k=v", false, "(none)/K=v")

	// Instance K=?, Pattern variations
	check("cap:k=?", "cap:", true, "K=?/(none)")
	check("cap:k=?", "cap:k=?", true, "K=?/K=?")
	check("cap:k=?", "cap:k=!", true, "K=?/K=!")
	check("cap:k=?", "cap:k", true, "K=?/K=*")
	check("cap:k=?", "cap:k=v", true, "K=?/K=v")

	// Instance K=!, Pattern variations
	check("cap:k=!", "cap:", true, "K=!/(none)")
	check("cap:k=!", "cap:k=?", true, "K=!/K=?")
	check("cap:k=!", "cap:k=!", true, "K=!/K=!")
	check("cap:k=!", "cap:k", false, "K=!/K=*")
	check("cap:k=!", "cap:k=v", false, "K=!/K=v")

	// Instance K=*, Pattern variations
	check("cap:k", "cap:", true, "K=*/(none)")
	check("cap:k", "cap:k=?", true, "K=*/K=?")
	check("cap:k", "cap:k=!", false, "K=*/K=!")
	check("cap:k", "cap:k", true, "K=*/K=*")
	check("cap:k", "cap:k=v", true, "K=*/K=v")

	// Instance K=v, Pattern variations
	check("cap:k=v", "cap:", true, "K=v/(none)")
	check("cap:k=v", "cap:k=?", true, "K=v/K=?")
	check("cap:k=v", "cap:k=!", false, "K=v/K=!")
	check("cap:k=v", "cap:k", true, "K=v/K=*")
	check("cap:k=v", "cap:k=v", true, "K=v/K=v")
	check("cap:k=v", "cap:k=w", false, "K=v/K=w")
}

// TEST0084: Mixed special values
func Test0084_MixedSpecialValues(t *testing.T) {
	// Test URNs with multiple special values
	pattern, err := NewTaggedUrnFromString("cap:required;optional=?;forbidden=!;exact=pdf")
	require.NoError(t, err)

	// Instance that satisfies all constraints
	goodInstance, _ := NewTaggedUrnFromString("cap:required=yes;optional=maybe;exact=pdf")
	matches, _ := goodInstance.ConformsTo(pattern)
	assert.True(t, matches)

	// Instance missing required tag
	missingRequired, _ := NewTaggedUrnFromString("cap:optional=maybe;exact=pdf")
	matches, _ = missingRequired.ConformsTo(pattern)
	assert.False(t, matches)

	// Instance has forbidden tag
	hasForbidden, _ := NewTaggedUrnFromString("cap:required=yes;forbidden=oops;exact=pdf")
	matches, _ = hasForbidden.ConformsTo(pattern)
	assert.False(t, matches)

	// Instance with wrong exact value
	wrongExact, _ := NewTaggedUrnFromString("cap:required=yes;exact=doc")
	matches, _ = wrongExact.ConformsTo(pattern)
	assert.False(t, matches)
}

// TEST0085: Serialization round trip special values
func Test0085_SerializationRoundTripSpecialValues(t *testing.T) {
	// All special values round-trip correctly
	originals := []string{
		"cap:ext=?",
		"cap:ext=!",
		"cap:ext", // * serializes as valueless
		"cap:a=?;b=!;c;d=exact",
	}

	for _, original := range originals {
		urn, err := NewTaggedUrnFromString(original)
		require.NoError(t, err, "Failed to parse: %s", original)
		serialized := urn.ToString()
		reparsed, err := NewTaggedUrnFromString(serialized)
		require.NoError(t, err, "Failed to reparse: %s", serialized)
		assert.True(t, urn.Equals(reparsed), "Round-trip failed for: %s", original)
	}
}

// TEST0086: Compatibility with special values
func Test0086_CompatibilityWithSpecialValues(t *testing.T) {
	// TEST576: Compatibility via bidirectional accepts
	mustNot, _ := NewTaggedUrnFromString("cap:ext=!")
	mustHave, _ := NewTaggedUrnFromString("cap:ext=*")
	specific, _ := NewTaggedUrnFromString("cap:ext=pdf")
	unspecified, _ := NewTaggedUrnFromString("cap:ext=?")
	missing, _ := NewTaggedUrnFromString("cap:")

	// Helper: bidirectional accepts = compatible
	biAccepts := func(a, b *TaggedUrn) bool {
		fwd, _ := a.Accepts(b)
		rev, _ := b.Accepts(a)
		return fwd || rev
	}

	// ! vs *: neither direction accepts (absent vs present)
	assert.False(t, biAccepts(mustNot, mustHave))

	// ! vs specific: neither direction accepts
	assert.False(t, biAccepts(mustNot, specific))

	// ! vs ?: ? accepts anything in either direction
	assert.True(t, biAccepts(mustNot, unspecified))

	// ! vs missing: ! accepts missing (absent)
	assert.True(t, biAccepts(mustNot, missing))

	// ! vs !: both want absent, accepts each other
	assert.True(t, biAccepts(mustNot, mustNot))

	// * vs specific: * accepts specific
	assert.True(t, biAccepts(mustHave, specific))

	// * vs *: both accept each other
	assert.True(t, biAccepts(mustHave, mustHave))

	// ? vs everything: ? accepts anything
	assert.True(t, biAccepts(unspecified, mustNot))
	assert.True(t, biAccepts(unspecified, mustHave))
	assert.True(t, biAccepts(unspecified, specific))
	assert.True(t, biAccepts(unspecified, unspecified))
	assert.True(t, biAccepts(unspecified, missing))
}

// TEST0087: Specificity with special values
func Test0087_SpecificityWithSpecialValues(t *testing.T) {
	// Six-form ladder: ?x=0, x?=v=1, x=*=2, x!=v=3, x=v=4, !x=5.
	exact, _ := NewTaggedUrnFromString("cap:a=x;b=y;c=z")    // 3 * 4 = 12
	mustHave, _ := NewTaggedUrnFromString("cap:a;b;c")       // 3 * 2 = 6
	mustNotUrn, _ := NewTaggedUrnFromString("cap:!a;!b;!c")  // 3 * 5 = 15
	unspecified, _ := NewTaggedUrnFromString("cap:?a;?b;?c") // 3 * 0 = 0
	// mixed: a=x (4) + b (2) + !c (5) + ?d (0) = 11
	mixed, _ := NewTaggedUrnFromString("cap:!c;?d;a=x;b")

	assert.Equal(t, 12, exact.Specificity())
	assert.Equal(t, 6, mustHave.Specificity())
	assert.Equal(t, 15, mustNotUrn.Specificity())
	assert.Equal(t, 0, unspecified.Specificity())
	assert.Equal(t, 11, mixed.Specificity())

	// Five-tuple counts: (must_not_have, exact, present_not_value,
	// must_have_any, absent_or_not_value).
	mn, e, pnv, mha, anv := exact.SpecificityTuple()
	assert.Equal(t, 0, mn)
	assert.Equal(t, 3, e)
	assert.Equal(t, 0, pnv)
	assert.Equal(t, 0, mha)
	assert.Equal(t, 0, anv)

	mn, e, pnv, mha, anv = mustHave.SpecificityTuple()
	assert.Equal(t, 0, mn)
	assert.Equal(t, 0, e)
	assert.Equal(t, 0, pnv)
	assert.Equal(t, 3, mha)
	assert.Equal(t, 0, anv)

	mn, e, pnv, mha, anv = mustNotUrn.SpecificityTuple()
	assert.Equal(t, 3, mn)
	assert.Equal(t, 0, e)
	assert.Equal(t, 0, pnv)
	assert.Equal(t, 0, mha)
	assert.Equal(t, 0, anv)

	mn, e, pnv, mha, anv = unspecified.SpecificityTuple()
	assert.Equal(t, 0, mn)
	assert.Equal(t, 0, e)
	assert.Equal(t, 0, pnv)
	assert.Equal(t, 0, mha)
	assert.Equal(t, 0, anv)

	mn, e, pnv, mha, anv = mixed.SpecificityTuple()
	assert.Equal(t, 1, mn)
	assert.Equal(t, 1, e)
	assert.Equal(t, 0, pnv)
	assert.Equal(t, 1, mha)
	assert.Equal(t, 0, anv)
}

// =========================================================================
// ORDER-THEORETIC RELATIONS: IsEquivalent, IsComparable
// =========================================================================

// TEST0578: Equivalent URNs with identical tag sets
func Test0578_EquivalentIdenticalTags(t *testing.T) {
	a, _ := NewTaggedUrnFromString("cap:generate;ext=pdf")
	b, _ := NewTaggedUrnFromString("cap:ext=pdf;generate") // same tags, different order
	equiv, err := a.IsEquivalent(b)
	require.NoError(t, err)
	assert.True(t, equiv)
	equiv, err = b.IsEquivalent(a) // symmetric
	require.NoError(t, err)
	assert.True(t, equiv)
}

// TEST0579: Non-equivalent URNs where one is more specific
func Test0579_NotEquivalentWhenOneMoreSpecific(t *testing.T) {
	general, _ := NewTaggedUrnFromString("media:")
	specific, _ := NewTaggedUrnFromString("media:pdf")
	equiv, err := general.IsEquivalent(specific)
	require.NoError(t, err)
	assert.False(t, equiv)
	equiv, err = specific.IsEquivalent(general)
	require.NoError(t, err)
	assert.False(t, equiv)
}

// TEST0580: Comparable URNs on the same specialization chain
func Test0580_ComparableSpecializationChain(t *testing.T) {
	general, _ := NewTaggedUrnFromString("media:")
	specific, _ := NewTaggedUrnFromString("media:pdf")
	comp, err := general.IsComparable(specific)
	require.NoError(t, err)
	assert.True(t, comp)
	comp, err = specific.IsComparable(general) // symmetric
	require.NoError(t, err)
	assert.True(t, comp)
}

// TEST0581: Incomparable URNs in different branches of the lattice
func Test0581_IncomparableDifferentBranches(t *testing.T) {
	pdf, _ := NewTaggedUrnFromString("media:pdf")
	txt, _ := NewTaggedUrnFromString("media:enc=utf-8;txt")
	comp, err := pdf.IsComparable(txt)
	require.NoError(t, err)
	assert.False(t, comp)
	comp, err = txt.IsComparable(pdf)
	require.NoError(t, err)
	assert.False(t, comp)
}

// TEST0582: Equivalent implies comparable but not vice versa
func Test0582_EquivalentImpliesComparable(t *testing.T) {
	a, _ := NewTaggedUrnFromString("cap:test;ext=pdf")
	b, _ := NewTaggedUrnFromString("cap:test;ext=pdf")
	equiv, err := a.IsEquivalent(b)
	require.NoError(t, err)
	assert.True(t, equiv)
	comp, err := a.IsComparable(b)
	require.NoError(t, err)
	assert.True(t, comp)

	// comparable but NOT equivalent
	general, _ := NewTaggedUrnFromString("cap:test")
	specific, _ := NewTaggedUrnFromString("cap:test;ext=pdf")
	equiv, err = general.IsEquivalent(specific)
	require.NoError(t, err)
	assert.False(t, equiv)
	comp, err = general.IsComparable(specific)
	require.NoError(t, err)
	assert.True(t, comp)
}

// TEST0583: Prefix mismatch returns error for both relations
func Test0583_PrefixMismatchErrors(t *testing.T) {
	cap, _ := NewTaggedUrnFromString("cap:test")
	media, _ := NewTaggedUrnFromString("media:")
	_, err := cap.IsEquivalent(media)
	assert.Error(t, err)
	_, err = cap.IsComparable(media)
	assert.Error(t, err)
}

// TEST0584: Empty tag set is comparable to everything with same prefix
func Test0584_EmptyTagsComparableToAll(t *testing.T) {
	empty, _ := NewTaggedUrnFromString("media:")
	specific, _ := NewTaggedUrnFromString("media:pdf;thumbnail")
	comp, err := empty.IsComparable(specific)
	require.NoError(t, err)
	assert.True(t, comp)
	equiv, err := empty.IsEquivalent(specific)
	require.NoError(t, err)
	assert.False(t, equiv)
	empty2, _ := NewTaggedUrnFromString("media:")
	equiv, err = empty.IsEquivalent(empty2)
	require.NoError(t, err)
	assert.True(t, equiv)
}

// TEST0585: String variants of IsEquivalent and IsComparable
func Test0585_StringVariants(t *testing.T) {
	urn, _ := NewTaggedUrnFromString("media:pdf")
	equiv, err := urn.IsEquivalentStr("media:pdf")
	require.NoError(t, err)
	assert.True(t, equiv)
	equiv, err = urn.IsEquivalentStr("media:")
	require.NoError(t, err)
	assert.False(t, equiv)
	comp, err := urn.IsComparableStr("media:")
	require.NoError(t, err)
	assert.True(t, comp)
	comp, err = urn.IsComparableStr("media:enc=utf-8;txt")
	require.NoError(t, err)
	assert.False(t, comp)
}

// TEST0586: Special values (*, !, ?) with IsEquivalent and IsComparable
func Test0586_SpecialValues(t *testing.T) {
	mustHave, _ := NewTaggedUrnFromString("cap:ext")
	exact, _ := NewTaggedUrnFromString("cap:ext=pdf")
	mustNot, _ := NewTaggedUrnFromString("cap:ext=!")
	unspecified, _ := NewTaggedUrnFromString("cap:ext=?")

	equiv, err := mustHave.IsEquivalent(exact)
	require.NoError(t, err)
	assert.True(t, equiv)
	comp, err := mustHave.IsComparable(exact)
	require.NoError(t, err)
	assert.True(t, comp)

	comp, err = mustNot.IsComparable(exact)
	require.NoError(t, err)
	assert.False(t, comp)
	equiv, err = mustNot.IsEquivalent(exact)
	require.NoError(t, err)
	assert.False(t, equiv)

	comp, err = mustNot.IsComparable(mustHave)
	require.NoError(t, err)
	assert.False(t, comp)
	equiv, err = mustNot.IsEquivalent(mustHave)
	require.NoError(t, err)
	assert.False(t, equiv)

	equiv, err = unspecified.IsEquivalent(exact)
	require.NoError(t, err)
	assert.True(t, equiv)
	equiv, err = unspecified.IsEquivalent(mustHave)
	require.NoError(t, err)
	assert.True(t, equiv)
	equiv, err = unspecified.IsEquivalent(mustNot)
	require.NoError(t, err)
	assert.True(t, equiv)
}

// =========================================================================
// BUILDER TESTS
// =========================================================================

// TEST0587: Builder fluent API for tag manipulation
func Test0587_BuilderFluentAPI(t *testing.T) {
	urn, err := NewTaggedUrnBuilder("cap").
		Marker("generate").
		Tag("target", "thumbnail").
		Tag("format", "pdf").
		Tag("output", "binary").
		Build()
	require.NoError(t, err)

	assert.True(t, urn.HasMarkerTag("generate"))
	val, exists := urn.GetTag("target")
	assert.True(t, exists)
	assert.Equal(t, "thumbnail", val)
	val, exists = urn.GetTag("format")
	assert.True(t, exists)
	assert.Equal(t, "pdf", val)
	val, exists = urn.GetTag("output")
	assert.True(t, exists)
	assert.Equal(t, "binary", val)
}

// TEST0588: Builder with custom tags
func Test0588_BuilderCustomTags(t *testing.T) {
	urn, err := NewTaggedUrnBuilder("cap").
		Tag("engine", "v2").
		Tag("quality", "high").
		Tag("op", "compress").
		Build()
	require.NoError(t, err)

	val, _ := urn.GetTag("engine")
	assert.Equal(t, "v2", val)
	val, _ = urn.GetTag("quality")
	assert.Equal(t, "high", val)
	val, _ = urn.GetTag("op")
	assert.Equal(t, "compress", val)
}

// TEST0589: Builder tag overrides
func Test0589_BuilderTagOverrides(t *testing.T) {
	urn, err := NewTaggedUrnBuilder("cap").
		Tag("op", "convert").
		Tag("format", "jpg").
		Build()
	require.NoError(t, err)

	val, _ := urn.GetTag("op")
	assert.Equal(t, "convert", val)
	val, _ = urn.GetTag("format")
	assert.Equal(t, "jpg", val)
}

// TEST0590: Builder empty build returns error
func Test0590_BuilderEmptyBuild(t *testing.T) {
	_, err := NewTaggedUrnBuilder("cap").Build()
	assert.Error(t, err)
}

// TEST0591: Builder with single tag
func Test0591_BuilderSingleTag(t *testing.T) {
	urn, err := NewTaggedUrnBuilder("cap").Tag("type", "utility").Build()
	require.NoError(t, err)

	assert.Equal(t, "cap:type=utility", urn.ToString())
	val, _ := urn.GetTag("type")
	assert.Equal(t, "utility", val)
	// Six-form ladder: exact value = 4 points.
	assert.Equal(t, 4, urn.Specificity())
}

// TEST0592: Builder with complex multi-tag URN
func Test0592_BuilderComplex(t *testing.T) {
	urn, err := NewTaggedUrnBuilder("cap").
		Tag("type", "media").
		Marker("transcode").
		Tag("target", "video").
		Tag("format", "mp4").
		Tag("codec", "h264").
		Tag("quality", "1080p").
		Tag("framerate", "30fps").
		Tag("output", "binary").
		Build()
	require.NoError(t, err)

	assert.Equal(t, "media", urn.AllTags()["type"])
	assert.True(t, urn.HasMarkerTag("transcode"))
	assert.Equal(t, "video", urn.AllTags()["target"])
	assert.Equal(t, "mp4", urn.AllTags()["format"])
	assert.Equal(t, "h264", urn.AllTags()["codec"])
	assert.Equal(t, "1080p", urn.AllTags()["quality"])
	assert.Equal(t, "30fps", urn.AllTags()["framerate"])
	assert.Equal(t, "binary", urn.AllTags()["output"])
	// Six-form ladder: 7 exact tags × 4 + 1 marker × 2 = 28 + 2 = 30.
	assert.Equal(t, 30, urn.Specificity())
}

// TEST0593: Builder with wildcards
func Test0593_BuilderWildcards(t *testing.T) {
	urn, err := NewTaggedUrnBuilder("cap").
		Marker("convert").
		Marker("ext").
		Marker("quality").
		Build()
	require.NoError(t, err)

	assert.Equal(t, "cap:convert;ext;quality", urn.ToString())
	// Six-form ladder: 3 markers × 2 = 6.
	assert.Equal(t, 6, urn.Specificity())
	val, _ := urn.GetTag("ext")
	assert.Equal(t, "*", val)
	val, _ = urn.GetTag("quality")
	assert.Equal(t, "*", val)
}

// TEST0594: Builder with custom prefix
func Test0594_BuilderCustomPrefix(t *testing.T) {
	urn, err := NewTaggedUrnBuilder("myapp").Tag("key", "value").Build()
	require.NoError(t, err)

	assert.Equal(t, "myapp", urn.GetPrefix())
	assert.Equal(t, "myapp:key=value", urn.ToString())
}

// TEST0595: Builder matching with built URN
func Test0595_BuilderMatchingWithBuiltUrn(t *testing.T) {
	specificInstance, _ := NewTaggedUrnBuilder("cap").
		Tag("op", "generate").
		Tag("target", "thumbnail").
		Tag("format", "pdf").
		Build()

	generalPattern, _ := NewTaggedUrnBuilder("cap").Tag("op", "generate").Build()

	wildcardPattern, _ := NewTaggedUrnBuilder("cap").
		Tag("op", "generate").
		Tag("target", "thumbnail").
		Marker("ext").
		Build()

	conforms, err := specificInstance.ConformsTo(generalPattern)
	require.NoError(t, err)
	assert.True(t, conforms)

	conforms, err = specificInstance.ConformsTo(wildcardPattern)
	require.NoError(t, err)
	assert.False(t, conforms)

	moreSpec, err := specificInstance.IsMoreSpecificThan(generalPattern)
	require.NoError(t, err)
	assert.True(t, moreSpec)

	// Six-form ladder: exact = 4, marker (=*) = 2.
	assert.Equal(t, 12, specificInstance.Specificity()) // 3 exact × 4 = 12
	assert.Equal(t, 4, generalPattern.Specificity())    // 1 exact × 4 = 4
	assert.Equal(t, 10, wildcardPattern.Specificity())  // 2 exact × 4 + 1 * × 2 = 10
}

// TEST: Builder rejects empty value
func Test0088_BuilderRejectsEmptyValue(t *testing.T) {
	_, err := NewTaggedUrnBuilder("cap").Tag("key", "").Build()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty value")
}
