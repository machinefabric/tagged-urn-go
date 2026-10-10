package taggedurn

import (
	"slices"
	"testing"

	"github.com/machinefabric/tagged-urn-go/formal"
	"github.com/stretchr/testify/require"
)

// TEST0600: every function of the proved model this mirror calls carries a proved claim.
//
// The model's package carries what is proved of each function it exports (its assurance
// document, generated from ../formal): each one decides, equals or keeps what its claim says,
// and none rests on an assumption about the host — the model needs none.
func Test0600_EveryModelFunctionCarriesAProvedClaim(t *testing.T) {
	a := formal.Assurance()
	require.Empty(t, a.Facilities, "the model assumes nothing of the host")
	require.Empty(t, a.Assumptions)
	require.NotEmpty(t, a.Exports)
	for _, e := range a.Exports {
		require.NotEmpty(t, e.Claims, "%s carries no claim", e.Name)
		require.Empty(t, e.Assumptions, "%s rests on an assumption", e.Name)
		for _, name := range e.Claims {
			c := a.Claim(name)
			require.NotNil(t, c, name)
			require.Equal(t, "proved", c.Status, name)
			require.True(t, slices.Contains(c.Subjects, e.Name), "%s is about %s", name, e.Name)
		}
	}
	refines := a.Claim("TaggedUrn.Exec.refines_decides")
	require.Equal(t, "lungo.decides", refines.Relation)
	require.Equal(t, []string{"TaggedUrn.refines"}, refines.Specifications)
}
