# Go Test Catalog

**Total Tests:** 106

**Numbered Tests:** 106

**Unnumbered Tests:** 0

**Numbered Tests Missing Descriptions:** 0

**Numbering Mismatches:** 0

All numbered test numbers are unique.

This catalog lists all tests in the Go codebase.

| Test # | Function Name | Description | File |
|--------|---------------|-------------|------|
| test0001 | `Test0001_TaggedUrnCreation` | TEST0001: Tagged urn creation | tagged_urn_test.go:12 |
| test0002 | `Test0002_CustomPrefix` | TEST0002: Custom prefix | tagged_urn_test.go:32 |
| test0003 | `Test0003_PrefixCaseInsensitive` | TEST0003: Prefix case insensitive | tagged_urn_test.go:42 |
| test0004 | `Test0004_PrefixMismatchError` | TEST0004: Prefix mismatch error | tagged_urn_test.go:58 |
| test0005 | `Test0005_BuilderWithPrefix` | TEST0005: Builder with prefix | tagged_urn_test.go:72 |
| test0006 | `Test0006_CanonicalStringFormat` | TEST0006: Canonical string format | tagged_urn_test.go:83 |
| test0007 | `Test0007_PrefixRequired` | TEST0007: Prefix required | tagged_urn_test.go:92 |
| test0008 | `Test0008_TrailingSemicolonEquivalence` | TEST0008: Trailing semicolon equivalence | tagged_urn_test.go:118 |
| test0009 | `Test0009_InvalidTaggedUrn` | TEST0009: Invalid tagged urn | tagged_urn_test.go:146 |
| test0010 | `Test0010_ValuelessTagParsing` | TEST0010: Valueless tag parsing | tagged_urn_test.go:155 |
| test0011 | `Test0011_InvalidCharacters` | TEST0011: Invalid characters | tagged_urn_test.go:168 |
| test0012 | `Test0012_TagMatching` | TEST0012: Tag matching | tagged_urn_test.go:177 |
| test0013 | `Test0013_MissingTagHandling` | TEST0013: Missing tag handling | tagged_urn_test.go:211 |
| test0014 | `Test0014_Specificity` | TEST0014: Specificity | tagged_urn_test.go:252 |
| test0015 | `Test0015_Compatibility` | TEST0015: Compatibility | tagged_urn_test.go:296 |
| test0016 | `Test0016_ConvenienceMethods` | TEST0016: Convenience methods | tagged_urn_test.go:329 |
| test0017 | `Test0017_Builder` | TEST0017: Builder | tagged_urn_test.go:349 |
| test0018 | `Test0018_WithTag` | TEST0018: With tag | tagged_urn_test.go:366 |
| test0019 | `Test0019_WithoutTag` | TEST0019: Without tag | tagged_urn_test.go:379 |
| test0020 | `Test0020_WildcardTag` | TEST0020: Wildcard tag | tagged_urn_test.go:392 |
| test0021 | `Test0021_Subset` | TEST0021: Subset | tagged_urn_test.go:416 |
| test0022 | `Test0022_Merge` | TEST0022: Merge | tagged_urn_test.go:426 |
| test0023 | `Test0023_MergePrefixMismatch` | TEST0023: Merge prefix mismatch | tagged_urn_test.go:440 |
| test0024 | `Test0024_Equality` | TEST0024: Equality | tagged_urn_test.go:455 |
| test0025 | `Test0025_EqualityDifferentPrefix` | TEST0025: Equality different prefix | tagged_urn_test.go:470 |
| test0026 | `Test0026_UrnMatcher` | TEST0026: Urn matcher | tagged_urn_test.go:481 |
| test0027 | `Test0027_UrnMatcherPrefixMismatch` | TEST0027: Urn matcher prefix mismatch | tagged_urn_test.go:510 |
| test0028 | `Test0028_JSONSerialization` | TEST0028: J s o n serialization | tagged_urn_test.go:527 |
| test0029 | `Test0029_JSONSerializationWithCustomPrefix` | TEST0029: J s o n serialization with custom prefix | tagged_urn_test.go:542 |
| test0030 | `Test0030_UnquotedValuesLowercased` | TEST0030: Unquoted values lowercased | tagged_urn_test.go:557 |
| test0031 | `Test0031_QuotedValuesPreserveCase` | TEST0031: Quoted values preserve case | tagged_urn_test.go:586 |
| test0032 | `Test0032_QuotedValueSpecialChars` | TEST0032: Quoted value special chars | tagged_urn_test.go:615 |
| test0033 | `Test0033_QuotedValueEscapeSequences` | TEST0033: Quoted value escape sequences | tagged_urn_test.go:639 |
| test0034 | `Test0034_MixedQuotedUnquoted` | TEST0034: Mixed quoted unquoted | tagged_urn_test.go:663 |
| test0035 | `Test0035_UnterminatedQuoteError` | TEST0035: Unterminated quote error | tagged_urn_test.go:677 |
| test0036 | `Test0036_InvalidEscapeSequenceError` | TEST0036: Invalid escape sequence error | tagged_urn_test.go:687 |
| test0037 | `Test0037_SerializationSmartQuoting` | TEST0037: Serialization smart quoting | tagged_urn_test.go:705 |
| test0038 | `Test0038_RoundTripSimple` | TEST0038: Round trip simple | tagged_urn_test.go:738 |
| test0039 | `Test0039_RoundTripQuoted` | TEST0039: Round trip quoted | tagged_urn_test.go:749 |
| test0040 | `Test0040_RoundTripEscapes` | TEST0040: Round trip escapes | tagged_urn_test.go:763 |
| test0041 | `Test0041_MatchingCaseSensitiveValues` | TEST0041: Matching case sensitive values | tagged_urn_test.go:777 |
| test0042 | `Test0042_BuilderPreservesCase` | TEST0042: Builder preserves case | tagged_urn_test.go:801 |
| test0043 | `Test0043_HasTagCaseSensitive` | TEST0043: Has tag case sensitive | tagged_urn_test.go:817 |
| test0044 | `Test0044_WithTagPreservesValue` | TEST0044: With tag preserves value | tagged_urn_test.go:834 |
| test0045 | `Test0045_SemanticEquivalence` | TEST0045: Semantic equivalence | tagged_urn_test.go:844 |
| test0046 | `Test0046_EmptyTaggedUrn` | TEST0046: Empty tagged urn | tagged_urn_test.go:858 |
| test0047 | `Test0047_EmptyWithCustomPrefix` | TEST0047: Empty with custom prefix | tagged_urn_test.go:896 |
| test0048 | `Test0048_ExtendedCharacterSupport` | TEST0048: Extended character support | tagged_urn_test.go:904 |
| test0049 | `Test0049_WildcardRestrictions` | TEST0049: Wildcard restrictions | tagged_urn_test.go:920 |
| test0050 | `Test0050_DuplicateKeyRejection` | TEST0050: Duplicate key rejection | tagged_urn_test.go:940 |
| test0051 | `Test0051_NumericKeyRestriction` | TEST0051: Numeric key restriction | tagged_urn_test.go:951 |
| test0052 | `Test0052_EmptyValueError` | TEST0052: Empty value error | tagged_urn_test.go:980 |
| test0053 | `Test0053_MatchingDifferentPrefixesError` | TEST0053: Matching different prefixes error | tagged_urn_test.go:991 |
| test0054 | `Test0054_MatchingSemantics_Test1_ExactMatch` | MATCHING SEMANTICS SPECIFICATION TESTS These 9 tests verify the exact matching semantics from RULES.md Sections 12-17 All implementations (Rust, Go, JS, ObjC) must pass these identically | tagged_urn_test.go:1014 |
| test0055 | `Test0055_MatchingSemantics_Test2_InstanceMissingTag` | TEST0055: Matching semantics  test2  instance missing tag | tagged_urn_test.go:1031 |
| test0056 | `Test0056_MatchingSemantics_Test3_UrnHasExtraTag` | TEST0056: Matching semantics  test3  urn has extra tag | tagged_urn_test.go:1058 |
| test0057 | `Test0057_MatchingSemantics_Test4_RequestHasWildcard` | TEST0057: Matching semantics  test4  request has wildcard | tagged_urn_test.go:1075 |
| test0058 | `Test0058_MatchingSemantics_Test5_UrnHasWildcard` | TEST0058: Matching semantics  test5  urn has wildcard | tagged_urn_test.go:1092 |
| test0059 | `Test0059_MatchingSemantics_Test6_ValueMismatch` | TEST0059: Matching semantics  test6  value mismatch | tagged_urn_test.go:1109 |
| test0060 | `Test0060_MatchingSemantics_Test7_PatternHasExtraTag` | TEST0060: Matching semantics  test7  pattern has extra tag | tagged_urn_test.go:1126 |
| test0061 | `Test0061_MatchingSemantics_Test8_EmptyPatternMatchesAnything` | TEST0061: Matching semantics  test8  empty pattern matches anything | tagged_urn_test.go:1152 |
| test0062 | `Test0062_MatchingSemantics_Test9_CrossDimensionConstraints` | TEST0062: Matching semantics  test9  cross dimension constraints | tagged_urn_test.go:1181 |
| test0063 | `Test0063_ValuelessTagParsingSingle` | VALUE-LESS TAG TESTS Value-less tags are equivalent to wildcard tags (key=*) | tagged_urn_test.go:1213 |
| test0064 | `Test0064_ValuelessTagParsingMultiple` | TEST0064: Valueless tag parsing multiple | tagged_urn_test.go:1226 |
| test0065 | `Test0065_ValuelessTagMixedWithValued` | TEST0065: Valueless tag mixed with valued | tagged_urn_test.go:1248 |
| test0066 | `Test0066_ValuelessTagAtEnd` | TEST0066: Valueless tag at end | tagged_urn_test.go:1272 |
| test0067 | `Test0067_ValuelessTagEquivalenceToWildcard` | TEST0067: Valueless tag equivalence to wildcard | tagged_urn_test.go:1287 |
| test0068 | `Test0068_ValuelessTagMatching` | TEST0068: Valueless tag matching | tagged_urn_test.go:1302 |
| test0069 | `Test0069_ValuelessTagInPattern` | TEST0069: Valueless tag in pattern | tagged_urn_test.go:1328 |
| test0070 | `Test0070_ValuelessTagSpecificity` | TEST0070: Valueless tag specificity | tagged_urn_test.go:1362 |
| test0071 | `Test0071_ValuelessTagRoundtrip` | TEST0071: Valueless tag roundtrip | tagged_urn_test.go:1377 |
| test0072 | `Test0072_ValuelessTagCaseNormalization` | TEST0072: Valueless tag case normalization | tagged_urn_test.go:1390 |
| test0073 | `Test0073_EmptyValueStillError` | TEST0073: Empty value still error | tagged_urn_test.go:1411 |
| test0074 | `Test0074_ValuelessTagCompatibility` | TEST0074: Valueless tag compatibility | tagged_urn_test.go:1423 |
| test0075 | `Test0075_ValuelessNumericKeyStillRejected` | TEST0075: Valueless numeric key still rejected | tagged_urn_test.go:1452 |
| test0076 | `Test0076_WhitespaceInInputRejected` | TEST0076: Whitespace in input rejected | tagged_urn_test.go:1464 |
| test0077 | `Test0077_UnspecifiedQuestionMarkParsing` | NEW SEMANTICS TESTS: ? (unspecified) and ! (must-not-have) | tagged_urn_test.go:1509 |
| test0078 | `Test0078_MustNotHaveExclamationParsing` | TEST0078: Must not have exclamation parsing | tagged_urn_test.go:1523 |
| test0079 | `Test0079_QuestionMarkPatternMatchesAnything` | TEST0079: Question mark pattern matches anything | tagged_urn_test.go:1536 |
| test0080 | `Test0080_QuestionMarkInInstance` | TEST0080: Question mark in instance | tagged_urn_test.go:1564 |
| test0081 | `Test0081_MustNotHavePatternRequiresAbsent` | TEST0081: Must not have pattern requires absent | tagged_urn_test.go:1592 |
| test0082 | `Test0082_MustNotHaveInInstance` | TEST0082: Must not have in instance | tagged_urn_test.go:1616 |
| test0083 | `Test0083_FullCrossProductMatching` | TEST0083: Full cross product matching | tagged_urn_test.go:1644 |
| test0084 | `Test0084_MixedSpecialValues` | TEST0084: Mixed special values | tagged_urn_test.go:1694 |
| test0085 | `Test0085_SerializationRoundTripSpecialValues` | TEST0085: Serialization round trip special values | tagged_urn_test.go:1721 |
| test0086 | `Test0086_CompatibilityWithSpecialValues` | TEST0086: Compatibility with special values | tagged_urn_test.go:1741 |
| test0087 | `Test0087_SpecificityWithSpecialValues` | TEST0087: Specificity with special values | tagged_urn_test.go:1786 |
| test0088 | `Test0088_BuilderRejectsEmptyValue` | TEST: Builder rejects empty value | tagged_urn_test.go:2158 |
| test0578 | `Test0578_EquivalentIdenticalTags` | TEST0578: Equivalent URNs with identical tag sets | tagged_urn_test.go:1844 |
| test0579 | `Test0579_NotEquivalentWhenOneMoreSpecific` | TEST0579: Non-equivalent URNs where one is more specific | tagged_urn_test.go:1856 |
| test0580 | `Test0580_ComparableSpecializationChain` | TEST0580: Comparable URNs on the same specialization chain | tagged_urn_test.go:1868 |
| test0581 | `Test0581_IncomparableDifferentBranches` | TEST0581: Incomparable URNs in different branches of the lattice | tagged_urn_test.go:1880 |
| test0582 | `Test0582_EquivalentImpliesComparable` | TEST0582: Equivalent implies comparable but not vice versa | tagged_urn_test.go:1892 |
| test0583 | `Test0583_PrefixMismatchErrors` | TEST0583: Prefix mismatch returns error for both relations | tagged_urn_test.go:1914 |
| test0584 | `Test0584_EmptyTagsComparableToAll` | TEST0584: Empty tag set is comparable to everything with same prefix | tagged_urn_test.go:1924 |
| test0585 | `Test0585_StringVariants` | TEST0585: String variants of IsEquivalent and IsComparable | tagged_urn_test.go:1940 |
| test0586 | `Test0586_SpecialValues` | TEST0586: Special values (*, !, ?) with IsEquivalent and IsComparable | tagged_urn_test.go:1957 |
| test0587 | `Test0587_BuilderFluentAPI` | TEST0587: Builder fluent API for tag manipulation | tagged_urn_test.go:2000 |
| test0588 | `Test0588_BuilderCustomTags` | TEST0588: Builder with custom tags | tagged_urn_test.go:2022 |
| test0589 | `Test0589_BuilderTagOverrides` | TEST0589: Builder tag overrides | tagged_urn_test.go:2039 |
| test0590 | `Test0590_BuilderEmptyBuild` | TEST0590: Builder empty build returns error | tagged_urn_test.go:2053 |
| test0591 | `Test0591_BuilderSingleTag` | TEST0591: Builder with single tag | tagged_urn_test.go:2059 |
| test0592 | `Test0592_BuilderComplex` | TEST0592: Builder with complex multi-tag URN | tagged_urn_test.go:2071 |
| test0593 | `Test0593_BuilderWildcards` | TEST0593: Builder with wildcards | tagged_urn_test.go:2097 |
| test0594 | `Test0594_BuilderCustomPrefix` | TEST0594: Builder with custom prefix | tagged_urn_test.go:2115 |
| test0595 | `Test0595_BuilderMatchingWithBuiltUrn` | TEST0595: Builder matching with built URN | tagged_urn_test.go:2124 |
---

*Generated from Go source tree*
*Total tests: 106*
*Total numbered tests: 106*
*Total unnumbered tests: 0*
*Total numbered tests missing descriptions: 0*
*Total numbering mismatches: 0*
