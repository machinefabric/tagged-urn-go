# Test catalogue — tagged-urn/tagged-urn-go

Generated from the test catalogue. Edit the tests, not this file.

121 tests: 121 numbered, 0 unnumbered.

## Numbered

| Number | Repository | Language | Test | Location | Description |
|---|---|---|---|---|---|
| TEST1 | tagged-urn/tagged-urn-go | go | `Test0001_TaggedUrnCreation` | tagged_urn_test.go:17 | TEST0001: Tagged urn creation |
| TEST2 | tagged-urn/tagged-urn-go | go | `Test0002_CustomPrefix` | tagged_urn_test.go:37 | TEST0002: Custom prefix |
| TEST3 | tagged-urn/tagged-urn-go | go | `Test0003_PrefixCaseInsensitive` | tagged_urn_test.go:47 | TEST0003: Prefix case insensitive |
| TEST4 | tagged-urn/tagged-urn-go | go | `Test0004_PrefixMismatchError` | tagged_urn_test.go:63 | TEST0004: Prefix mismatch error |
| TEST5 | tagged-urn/tagged-urn-go | go | `Test0005_BuilderWithPrefix` | tagged_urn_test.go:77 | TEST0005: Builder with prefix |
| TEST6 | tagged-urn/tagged-urn-go | go | `Test0006_CanonicalStringFormat` | tagged_urn_test.go:88 | TEST0006: Canonical string format |
| TEST7 | tagged-urn/tagged-urn-go | go | `Test0007_PrefixRequired` | tagged_urn_test.go:97 | TEST0007: Prefix required |
| TEST8 | tagged-urn/tagged-urn-go | go | `Test0008_TrailingSemicolonEquivalence` | tagged_urn_test.go:123 | TEST0008: Trailing semicolon equivalence |
| TEST9 | tagged-urn/tagged-urn-go | go | `Test0009_InvalidTaggedUrn` | tagged_urn_test.go:151 | TEST0009: Invalid tagged urn |
| TEST10 | tagged-urn/tagged-urn-go | go | `Test0010_ValuelessTagParsing` | tagged_urn_test.go:160 | TEST0010: Valueless tag parsing |
| TEST11 | tagged-urn/tagged-urn-go | go | `Test0011_InvalidCharacters` | tagged_urn_test.go:173 | TEST0011: Invalid characters |
| TEST12 | tagged-urn/tagged-urn-go | go | `Test0012_TagMatching` | tagged_urn_test.go:182 | TEST0012: Tag matching |
| TEST13 | tagged-urn/tagged-urn-go | go | `Test0013_MissingTagHandling` | tagged_urn_test.go:216 | TEST0013: Missing tag handling |
| TEST14 | tagged-urn/tagged-urn-go | go | `Test0014_Specificity` | tagged_urn_test.go:257 | TEST0014: Specificity |
| TEST15 | tagged-urn/tagged-urn-go | go | `Test0015_Compatibility` | tagged_urn_test.go:301 | TEST0015: Compatibility |
| TEST16 | tagged-urn/tagged-urn-go | go | `Test0016_ConvenienceMethods` | tagged_urn_test.go:334 | TEST0016: Convenience methods |
| TEST17 | tagged-urn/tagged-urn-go | go | `Test0017_Builder` | tagged_urn_test.go:354 | TEST0017: Builder |
| TEST18 | tagged-urn/tagged-urn-go | go | `Test0018_WithTag` | tagged_urn_test.go:371 | TEST0018: With tag |
| TEST19 | tagged-urn/tagged-urn-go | go | `Test0019_WithoutTag` | tagged_urn_test.go:385 | TEST0019: Without tag |
| TEST20 | tagged-urn/tagged-urn-go | go | `Test0020_WildcardTag` | tagged_urn_test.go:398 | TEST0020: Wildcard tag |
| TEST21 | tagged-urn/tagged-urn-go | go | `Test0021_Subset` | tagged_urn_test.go:422 | TEST0021: Subset |
| TEST22 | tagged-urn/tagged-urn-go | go | `Test0022_Merge` | tagged_urn_test.go:432 | TEST0022: Merge |
| TEST23 | tagged-urn/tagged-urn-go | go | `Test0023_MergePrefixMismatch` | tagged_urn_test.go:446 | TEST0023: Merge prefix mismatch |
| TEST24 | tagged-urn/tagged-urn-go | go | `Test0024_Equality` | tagged_urn_test.go:461 | TEST0024: Equality |
| TEST25 | tagged-urn/tagged-urn-go | go | `Test0025_EqualityDifferentPrefix` | tagged_urn_test.go:476 | TEST0025: Equality different prefix |
| TEST26 | tagged-urn/tagged-urn-go | go | `Test0026_UrnMatcher` | tagged_urn_test.go:487 | TEST0026: Urn matcher |
| TEST27 | tagged-urn/tagged-urn-go | go | `Test0027_UrnMatcherPrefixMismatch` | tagged_urn_test.go:516 | TEST0027: Urn matcher prefix mismatch |
| TEST28 | tagged-urn/tagged-urn-go | go | `Test0028_JSONSerialization` | tagged_urn_test.go:533 | TEST0028: J s o n serialization |
| TEST29 | tagged-urn/tagged-urn-go | go | `Test0029_JSONSerializationWithCustomPrefix` | tagged_urn_test.go:548 | TEST0029: J s o n serialization with custom prefix |
| TEST30 | tagged-urn/tagged-urn-go | go | `Test0030_UnquotedValuesLowercased` | tagged_urn_test.go:563 | TEST0030: Unquoted values lowercased |
| TEST31 | tagged-urn/tagged-urn-go | go | `Test0031_QuotedValuesPreserveCase` | tagged_urn_test.go:592 | TEST0031: Quoted values preserve case |
| TEST32 | tagged-urn/tagged-urn-go | go | `Test0032_QuotedValueSpecialChars` | tagged_urn_test.go:621 | TEST0032: Quoted value special chars |
| TEST33 | tagged-urn/tagged-urn-go | go | `Test0033_QuotedValueEscapeSequences` | tagged_urn_test.go:645 | TEST0033: Quoted value escape sequences |
| TEST34 | tagged-urn/tagged-urn-go | go | `Test0034_MixedQuotedUnquoted` | tagged_urn_test.go:669 | TEST0034: Mixed quoted unquoted |
| TEST35 | tagged-urn/tagged-urn-go | go | `Test0035_UnterminatedQuoteError` | tagged_urn_test.go:683 | TEST0035: Unterminated quote error |
| TEST36 | tagged-urn/tagged-urn-go | go | `Test0036_InvalidEscapeSequenceError` | tagged_urn_test.go:693 | TEST0036: Invalid escape sequence error |
| TEST37 | tagged-urn/tagged-urn-go | go | `Test0037_SerializationSmartQuoting` | tagged_urn_test.go:711 | TEST0037: Serialization smart quoting |
| TEST38 | tagged-urn/tagged-urn-go | go | `Test0038_RoundTripSimple` | tagged_urn_test.go:744 | TEST0038: Round trip simple |
| TEST39 | tagged-urn/tagged-urn-go | go | `Test0039_RoundTripQuoted` | tagged_urn_test.go:755 | TEST0039: Round trip quoted |
| TEST40 | tagged-urn/tagged-urn-go | go | `Test0040_RoundTripEscapes` | tagged_urn_test.go:769 | TEST0040: Round trip escapes |
| TEST41 | tagged-urn/tagged-urn-go | go | `Test0041_MatchingCaseSensitiveValues` | tagged_urn_test.go:783 | TEST0041: Matching case sensitive values |
| TEST42 | tagged-urn/tagged-urn-go | go | `Test0042_BuilderPreservesCase` | tagged_urn_test.go:807 | TEST0042: Builder preserves case |
| TEST43 | tagged-urn/tagged-urn-go | go | `Test0043_HasTagCaseSensitive` | tagged_urn_test.go:823 | TEST0043: Has tag case sensitive |
| TEST44 | tagged-urn/tagged-urn-go | go | `Test0044_WithTagPreservesValue` | tagged_urn_test.go:840 | TEST0044: With tag preserves value |
| TEST45 | tagged-urn/tagged-urn-go | go | `Test0045_SemanticEquivalence` | tagged_urn_test.go:852 | TEST0045: Semantic equivalence |
| TEST46 | tagged-urn/tagged-urn-go | go | `Test0046_EmptyTaggedUrn` | tagged_urn_test.go:866 | TEST0046: Empty tagged urn |
| TEST47 | tagged-urn/tagged-urn-go | go | `Test0047_EmptyWithCustomPrefix` | tagged_urn_test.go:904 | TEST0047: Empty with custom prefix |
| TEST48 | tagged-urn/tagged-urn-go | go | `Test0048_ExtendedCharacterSupport` | tagged_urn_test.go:912 | TEST0048: Extended character support |
| TEST49 | tagged-urn/tagged-urn-go | go | `Test0049_WildcardRestrictions` | tagged_urn_test.go:937 | TEST0049: Wildcard restrictions |
| TEST50 | tagged-urn/tagged-urn-go | go | `Test0050_DuplicateKeyRejection` | tagged_urn_test.go:957 | TEST0050: Duplicate key rejection |
| TEST51 | tagged-urn/tagged-urn-go | go | `Test0051_NumericKeyRestriction` | tagged_urn_test.go:968 | TEST0051: Numeric key restriction |
| TEST52 | tagged-urn/tagged-urn-go | go | `Test0052_EmptyValueError` | tagged_urn_test.go:995 | TEST0052: Empty value error |
| TEST53 | tagged-urn/tagged-urn-go | go | `Test0053_MatchingDifferentPrefixesError` | tagged_urn_test.go:1006 | TEST0053: Matching different prefixes error |
| TEST54 | tagged-urn/tagged-urn-go | go | `Test0054_MatchingSemantics_Test1_ExactMatch` | tagged_urn_test.go:1029 | MATCHING SEMANTICS SPECIFICATION TESTS These 9 tests verify the exact matching semantics from RULES.md Sections 12-17 All implementations (Rust, Go, JS, ObjC) must pass these identically |
| TEST55 | tagged-urn/tagged-urn-go | go | `Test0055_MatchingSemantics_Test2_InstanceMissingTag` | tagged_urn_test.go:1046 | TEST0055: Matching semantics  test2  instance missing tag |
| TEST56 | tagged-urn/tagged-urn-go | go | `Test0056_MatchingSemantics_Test3_UrnHasExtraTag` | tagged_urn_test.go:1073 | TEST0056: Matching semantics  test3  urn has extra tag |
| TEST57 | tagged-urn/tagged-urn-go | go | `Test0057_MatchingSemantics_Test4_RequestHasWildcard` | tagged_urn_test.go:1090 | TEST0057: Matching semantics  test4  request has wildcard |
| TEST58 | tagged-urn/tagged-urn-go | go | `Test0058_MatchingSemantics_Test5_UrnHasWildcard` | tagged_urn_test.go:1107 | TEST0058: Matching semantics  test5  urn has wildcard |
| TEST59 | tagged-urn/tagged-urn-go | go | `Test0059_MatchingSemantics_Test6_ValueMismatch` | tagged_urn_test.go:1125 | TEST0059: Matching semantics  test6  value mismatch |
| TEST60 | tagged-urn/tagged-urn-go | go | `Test0060_MatchingSemantics_Test7_PatternHasExtraTag` | tagged_urn_test.go:1142 | TEST0060: Matching semantics  test7  pattern has extra tag |
| TEST61 | tagged-urn/tagged-urn-go | go | `Test0061_MatchingSemantics_Test8_EmptyPatternMatchesAnything` | tagged_urn_test.go:1168 | TEST0061: Matching semantics  test8  empty pattern matches anything |
| TEST62 | tagged-urn/tagged-urn-go | go | `Test0062_MatchingSemantics_Test9_CrossDimensionConstraints` | tagged_urn_test.go:1197 | TEST0062: Matching semantics  test9  cross dimension constraints |
| TEST63 | tagged-urn/tagged-urn-go | go | `Test0063_ValuelessTagParsingSingle` | tagged_urn_test.go:1229 | VALUE-LESS TAG TESTS Value-less tags are equivalent to wildcard tags (key=*) |
| TEST64 | tagged-urn/tagged-urn-go | go | `Test0064_ValuelessTagParsingMultiple` | tagged_urn_test.go:1242 | TEST0064: Valueless tag parsing multiple |
| TEST65 | tagged-urn/tagged-urn-go | go | `Test0065_ValuelessTagMixedWithValued` | tagged_urn_test.go:1264 | TEST0065: Valueless tag mixed with valued |
| TEST66 | tagged-urn/tagged-urn-go | go | `Test0066_ValuelessTagAtEnd` | tagged_urn_test.go:1288 | TEST0066: Valueless tag at end |
| TEST67 | tagged-urn/tagged-urn-go | go | `Test0067_ValuelessTagEquivalenceToWildcard` | tagged_urn_test.go:1303 | TEST0067: Valueless tag equivalence to wildcard |
| TEST68 | tagged-urn/tagged-urn-go | go | `Test0068_ValuelessTagMatching` | tagged_urn_test.go:1318 | TEST0068: Valueless tag matching |
| TEST69 | tagged-urn/tagged-urn-go | go | `Test0069_ValuelessTagInPattern` | tagged_urn_test.go:1349 | TEST0069: Valueless tag in pattern |
| TEST70 | tagged-urn/tagged-urn-go | go | `Test0070_ValuelessTagSpecificity` | tagged_urn_test.go:1383 | TEST0070: Valueless tag specificity |
| TEST71 | tagged-urn/tagged-urn-go | go | `Test0071_ValuelessTagRoundtrip` | tagged_urn_test.go:1398 | TEST0071: Valueless tag roundtrip |
| TEST72 | tagged-urn/tagged-urn-go | go | `Test0072_ValuelessTagCaseNormalization` | tagged_urn_test.go:1411 | TEST0072: Valueless tag case normalization |
| TEST73 | tagged-urn/tagged-urn-go | go | `Test0073_EmptyValueStillError` | tagged_urn_test.go:1432 | TEST0073: Empty value still error |
| TEST74 | tagged-urn/tagged-urn-go | go | `Test0074_ValuelessTagCompatibility` | tagged_urn_test.go:1444 | TEST0074: Valueless tag compatibility |
| TEST75 | tagged-urn/tagged-urn-go | go | `Test0075_ValuelessNumericKeyStillRejected` | tagged_urn_test.go:1473 | TEST0075: Valueless numeric key still rejected |
| TEST76 | tagged-urn/tagged-urn-go | go | `Test0076_WhitespaceInInputRejected` | tagged_urn_test.go:1485 | TEST0076: Whitespace in input rejected |
| TEST77 | tagged-urn/tagged-urn-go | go | `Test0077_UnspecifiedQuestionMarkParsing` | tagged_urn_test.go:1530 | NEW SEMANTICS TESTS: ? (unspecified) and ! (must-not-have) |
| TEST78 | tagged-urn/tagged-urn-go | go | `Test0078_MustNotHaveExclamationParsing` | tagged_urn_test.go:1544 | TEST0078: Must not have exclamation parsing |
| TEST79 | tagged-urn/tagged-urn-go | go | `Test0079_QuestionMarkPatternMatchesAnything` | tagged_urn_test.go:1557 | TEST0079: Question mark pattern matches anything |
| TEST80 | tagged-urn/tagged-urn-go | go | `Test0080_QuestionMarkInInstance` | tagged_urn_test.go:1585 | TEST0080: Question mark in instance |
| TEST81 | tagged-urn/tagged-urn-go | go | `Test0081_MustNotHavePatternRequiresAbsent` | tagged_urn_test.go:1615 | TEST0081: Must not have pattern requires absent |
| TEST82 | tagged-urn/tagged-urn-go | go | `Test0082_MustNotHaveInInstance` | tagged_urn_test.go:1640 | TEST0082: Must not have in instance |
| TEST83 | tagged-urn/tagged-urn-go | go | `Test0083_FullCrossProductMatching` | tagged_urn_test.go:1668 | TEST0083: Full cross product matching |
| TEST84 | tagged-urn/tagged-urn-go | go | `Test0084_MixedSpecialValues` | tagged_urn_test.go:1718 | TEST0084: Mixed special values |
| TEST85 | tagged-urn/tagged-urn-go | go | `Test0085_SerializationRoundTripSpecialValues` | tagged_urn_test.go:1749 | TEST0085: Serialization round trip special values |
| TEST86 | tagged-urn/tagged-urn-go | go | `Test0086_CompatibilityWithSpecialValues` | tagged_urn_test.go:1769 | TEST0086: Compatibility with special values |
| TEST87 | tagged-urn/tagged-urn-go | go | `Test0087_SpecificityWithSpecialValues` | tagged_urn_test.go:1814 | TEST0087: Specificity with special values |
| TEST88 | tagged-urn/tagged-urn-go | go | `Test0088_BuilderRejectsEmptyValue` | tagged_urn_test.go:2198 | TEST: Builder rejects empty value |
| TEST392 | tagged-urn/tagged-urn-go | go | `Test0392_LeastUpperBound` | tagged_urn_test.go:2205 | TEST0392: the bound of several URNs: the values they share, a key every one sets to differing values as the set of those values, the coarser of two comparable URNs, nothing for none, and another prefix refused |
| TEST578 | tagged-urn/tagged-urn-go | go | `Test0578_EquivalentIdenticalTags` | tagged_urn_test.go:1872 | TEST0578: Equivalent URNs with identical tag sets |
| TEST579 | tagged-urn/tagged-urn-go | go | `Test0579_NotEquivalentWhenOneMoreSpecific` | tagged_urn_test.go:1884 | TEST0579: Non-equivalent URNs where one is more specific |
| TEST580 | tagged-urn/tagged-urn-go | go | `Test0580_ComparableSpecializationChain` | tagged_urn_test.go:1896 | TEST0580: Comparable URNs on the same specialization chain |
| TEST581 | tagged-urn/tagged-urn-go | go | `Test0581_IncomparableDifferentBranches` | tagged_urn_test.go:1908 | TEST0581: Incomparable URNs in different branches of the lattice |
| TEST582 | tagged-urn/tagged-urn-go | go | `Test0582_EquivalentImpliesComparable` | tagged_urn_test.go:1920 | TEST0582: Equivalent implies comparable but not vice versa |
| TEST583 | tagged-urn/tagged-urn-go | go | `Test0583_PrefixMismatchErrors` | tagged_urn_test.go:1942 | TEST0583: Prefix mismatch returns error for both relations |
| TEST584 | tagged-urn/tagged-urn-go | go | `Test0584_EmptyTagsComparableToAll` | tagged_urn_test.go:1952 | TEST0584: Empty tag set is comparable to everything with same prefix |
| TEST585 | tagged-urn/tagged-urn-go | go | `Test0585_StringVariants` | tagged_urn_test.go:1968 | TEST0585: String variants of IsEquivalent and IsComparable |
| TEST586 | tagged-urn/tagged-urn-go | go | `Test0586_SpecialValues` | tagged_urn_test.go:1985 | TEST0586: Special values (*, !, ?) with IsEquivalent and IsComparable |
| TEST587 | tagged-urn/tagged-urn-go | go | `Test0587_BuilderFluentAPI` | tagged_urn_test.go:2040 | TEST0587: Builder fluent API for tag manipulation |
| TEST588 | tagged-urn/tagged-urn-go | go | `Test0588_BuilderCustomTags` | tagged_urn_test.go:2062 | TEST0588: Builder with custom tags |
| TEST589 | tagged-urn/tagged-urn-go | go | `Test0589_BuilderTagOverrides` | tagged_urn_test.go:2079 | TEST0589: Builder tag overrides |
| TEST590 | tagged-urn/tagged-urn-go | go | `Test0590_BuilderEmptyBuild` | tagged_urn_test.go:2093 | TEST0590: Builder empty build returns error |
| TEST591 | tagged-urn/tagged-urn-go | go | `Test0591_BuilderSingleTag` | tagged_urn_test.go:2099 | TEST0591: Builder with single tag |
| TEST592 | tagged-urn/tagged-urn-go | go | `Test0592_BuilderComplex` | tagged_urn_test.go:2111 | TEST0592: Builder with complex multi-tag URN |
| TEST593 | tagged-urn/tagged-urn-go | go | `Test0593_BuilderWildcards` | tagged_urn_test.go:2137 | TEST0593: Builder with wildcards |
| TEST594 | tagged-urn/tagged-urn-go | go | `Test0594_BuilderCustomPrefix` | tagged_urn_test.go:2155 | TEST0594: Builder with custom prefix |
| TEST595 | tagged-urn/tagged-urn-go | go | `Test0595_BuilderMatchingWithBuiltUrn` | tagged_urn_test.go:2164 | TEST0595: Builder matching with built URN |
| TEST599 | tagged-urn/tagged-urn-go | go | `Test0599_EveryRowOfTheModelsTable` | formal_conformance_test.go:21 | TEST0599: every row of the proved model's table. The rules are proved in ../formal (Lean); this is what ties them to this mirror: every row of ../formal/conformance.json (written by the model, `lake exe conformance`) is parsed by this parser and must get the model's verdict — for the guarantee (ConformsTo), the possibility (Meets), and the complete reading of the instance (Satisfies, MaySatisfy) — and every set of URNs in it must get the model's bound (LeastUpperBound). The same table runs in every mirror. |
| TEST600 | tagged-urn/tagged-urn-go | go | `Test0600_EveryModelFunctionCarriesAProvedClaim` | formal_assurance_test.go:16 | TEST0600: every function of the proved model this mirror calls carries a proved claim. The model's package carries what is proved of each function it exports (its assurance document, generated from ../formal): each one decides, equals or keeps what its claim says, and none rests on an assumption about the host — the model needs none. |
| TEST6752 | tagged-urn/tagged-urn-go | go | `Test6752_QualifierValueAlias` | tagged_urn_test.go:2246 | TEST6752: An unquoted value that begins with ? or ! is that qualifier (x=?v is x?=v, x=!v is x!=v); a quoted one is literal |
| TEST6753 | tagged-urn/tagged-urn-go | go | `Test6753_PrefixAndKeyAreIdentifiers` | tagged_urn_test.go:2270 | TEST6753: The prefix and every key are identifiers (an ASCII letter, then ASCII letters, digits and -), lowered, in parsing and in every constructor |
| TEST6754 | tagged-urn/tagged-urn-go | go | `Test6754_EveryValueReadsBack` | tagged_urn_test.go:2317 | TEST6754: Every value a URN holds prints so that it reads back as itself |
| TEST6755 | tagged-urn/tagged-urn-go | go | `Test6755_UnquotedValueTakesAnyNonSeparator` | tagged_urn_test.go:2333 | TEST6755: An unquoted value takes any character but ; = \| " \ and whitespace |
| TEST6756 | tagged-urn/tagged-urn-go | go | `Test6756_QuotedValueIsLiteral` | tagged_urn_test.go:2353 | TEST6756: A quoted value is the literal value: a quoted *, ? or ! is never a sigil, and a quoted ?=v is never a qualifier |
| TEST7210 | tagged-urn/tagged-urn-go | go | `Test7210_AdmitsAbsenceByForm` | tagged_urn_test.go:2507 | TEST7210: A written tag lets its key be absent exactly when it is !k, ?k or k?=v; a key the URN does not mention says nothing of it, so it allows absence too. |
| TEST7223 | tagged-urn/tagged-urn-go | go | `Test7223_UrnReadBackFromItsHandle` | tagged_urn_test.go:2415 | TEST7223: A URN read back from its model handle is the URN: its prefix and every tag, each form and a quoted value included. |
| TEST7224 | tagged-urn/tagged-urn-go | go | `Test7224_SetFormsParseAndPrint` | tagged_urn_test.go:2428 | TEST7224: Each set form parses to the set it writes and prints canonically — values once, in increasing order, `!` last — reading back as itself; the qualifier aliases write the same sets. |
| TEST7225 | tagged-urn/tagged-urn-go | go | `Test7225_SetFormsRefused` | tagged_urn_test.go:2454 | TEST7225: A set of values refuses a value given twice, an empty member, a bare sigil or qualifier among its values, and `!` under a qualifier; a literal `\|` is quoted. |
| TEST7226 | tagged-urn/tagged-urn-go | go | `Test7226_BoundIsExactPerKey` | tagged_urn_test.go:2474 | TEST7226: The bound of several URNs is exact per key: png and jpeg are bound by ext=jpeg\|png, a value and its absence by ext=v\|!, and a value and any value by ext. |
| TEST7227 | tagged-urn/tagged-urn-go | go | `Test7227_SetScores` | tagged_urn_test.go:2493 | TEST7227: A set scores by what it allows — two or more values 3, absence and some values 1 — and narrowing never lowers a score, though two on one chain can tie. |
| TEST7241 | tagged-urn/tagged-urn-go | go | `Test7241_UrnsAreOrderedByTheirParts` | tagged_urn_test.go:2387 | TEST7241: URNs are ordered by their parts, never their text: the prefix, then each tag in key order by key, then by its constraint's form (unconstrained, optionalNot, present, presentNot, exact, absent), then by its value |

