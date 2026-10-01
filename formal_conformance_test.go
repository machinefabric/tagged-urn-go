package taggedurn

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// TEST0599: every row of the proved model's table.
//
// The rules are proved in ../formal (Lean); this is what ties them to this
// mirror: every row of ../formal/conformance.json (written by the model,
// `lake exe conformance`) is parsed by this parser and must get the model's
// verdict — for the guarantee (ConformsTo), the possibility (Meets), and the
// complete reading of the instance (Satisfies, MaySatisfy). The same table runs
// in every mirror.
func Test0599_EveryRowOfTheModelsTable(t *testing.T) {
	raw, err := os.ReadFile("../formal/conformance.json")
	require.NoError(t, err, "the model's table")
	var table struct {
		Refines []struct {
			Instance   string `json:"instance"`
			Pattern    string `json:"pattern"`
			Refines    bool   `json:"refines"`
			Equivalent bool   `json:"equivalent"`
			Meets      bool   `json:"meets"`
			Satisfies  bool   `json:"satisfies"`
			MaySatisfy bool   `json:"may_satisfy"`
		} `json:"refines"`
		Scores []struct {
			Urn   string `json:"urn"`
			Score int    `json:"score"`
		} `json:"scores"`
	}
	require.NoError(t, json.Unmarshal(raw, &table))

	var wrong []string
	for _, row := range table.Refines {
		a, err := NewTaggedUrnFromString(row.Instance)
		require.NoError(t, err)
		b, err := NewTaggedUrnFromString(row.Pattern)
		require.NoError(t, err)
		got, err := a.ConformsTo(b)
		require.NoError(t, err)
		if got != row.Refines {
			wrong = append(wrong, fmt.Sprintf("%s ⪯ %s: model %v", row.Instance, row.Pattern, row.Refines))
		}
		got, err = a.IsEquivalent(b)
		require.NoError(t, err)
		if got != row.Equivalent {
			wrong = append(wrong, fmt.Sprintf("%s ≡ %s: model %v", row.Instance, row.Pattern, row.Equivalent))
		}
		for _, check := range []struct {
			name  string
			ask   func(*TaggedUrn) (bool, error)
			model bool
		}{
			{"meets", a.Meets, row.Meets},
			{"satisfies", a.Satisfies, row.Satisfies},
			{"may satisfy", a.MaySatisfy, row.MaySatisfy},
		} {
			got, err := check.ask(b)
			require.NoError(t, err)
			if got != check.model {
				wrong = append(wrong, fmt.Sprintf("%s %s %s: model %v", row.Instance, check.name, row.Pattern, check.model))
			}
		}
	}
	for _, row := range table.Scores {
		u, err := NewTaggedUrnFromString(row.Urn)
		require.NoError(t, err)
		if u.Specificity() != row.Score {
			wrong = append(wrong, fmt.Sprintf("specificity %s: model %d, got %d", row.Urn, row.Score, u.Specificity()))
		}
	}
	require.True(t, len(table.Refines) > 4000 && len(table.Scores) > 60, "the table is the full one")
	if len(wrong) > 0 {
		t.Fatalf("%d answer(s) differ from the model, e.g.\n  %v", len(wrong), wrong[:min(8, len(wrong))])
	}
}
