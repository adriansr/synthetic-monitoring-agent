package synthetic_monitoring

import (
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestVoidVsKnownChecks(t *testing.T) {
	vEmpty := marshal(t, &Void{})
	kEmpty := marshal(t, &ProbeState{})
	t.Logf("For Void got %+v", vEmpty)
	t.Logf("For KnCh got %+v", kEmpty)
	require.Equal(t, len(vEmpty), len(kEmpty))
	var v Void
	unmarshal(t, &v, kEmpty)
	kZero := marshal(t, &ProbeState{
		Checks: make([]Entity, 0),
	})
	t.Logf("For KnCh got %+v", kZero)
	require.Equal(t, len(vEmpty), len(kZero))
	kData := marshal(t, &ProbeState{
		Checks: []Entity{
			{
				Id:       123,
				Modified: 789.0,
			},
			{
				Id:       456,
				Modified: 1234567.0,
			},
		},
	})
	t.Logf("For KnCh got %+v", kData)
	var k ProbeState
	unmarshal(t, &k, kEmpty)
	require.Nil(t, k.Checks)
}

func TestCheckChangesV2(t *testing.T) {
	// This is an original CheckChanges message without the optional deltaFirstBatch field.
	original := []byte{
		10, 4, 18, 2, 66, 0, 10, 6, 8, 1, 18, 2, 66, 0, 10, 6, 8, 2, 18, 2,
		66, 0, 18, 5, 8, 123, 16, 200, 3, 18, 6, 8, 149, 6, 16, 243, 7,
	}
	cc := &Changes{
		Checks: []CheckChange{
			{
				Operation: CheckOperation_CHECK_ADD,
			},
			{
				Operation: CheckOperation_CHECK_UPDATE,
			},
			{
				Operation: CheckOperation_CHECK_DELETE,
			},
		},
		Tenants: []Tenant{
			{
				Id:    123,
				OrgId: 456,
			},
			{
				Id:    789,
				OrgId: 1011,
			},
		},
	}
	var decoded Changes
	unmarshal(t, &decoded, original)
	require.Nil(t, decoded.DeltaFirstBatch)
	b1 := marshal(t, cc)
	t.Logf("For B1 got %d bytes %+v", len(b1), b1)
	require.Equal(t, original, b1)
	cc.DeltaFirstBatch = []bool{true}
	b2 := marshal(t, cc)
	t.Logf("For B2 got %d bytes %+v", len(b2), b2)
	require.True(t, len(b1) < len(b2))
}

func marshal(t testing.TB, m proto.Marshaler) []byte {
	t.Helper()
	data, err := m.Marshal()
	require.NoError(t, err)
	return data
}

func unmarshal(t testing.TB, m proto.Unmarshaler, d []byte) {
	t.Helper()
	require.NoError(t, m.Unmarshal(d))
}
