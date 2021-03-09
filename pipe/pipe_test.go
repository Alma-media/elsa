package pipe

import (
	"encoding/json"
	"testing"
)

func Test_ElementUnmarshalJSON(t *testing.T) {
	t.Run("given a single pipe element", func(t *testing.T) {
		var (
			element Element
			input   = `{"input":"/foo","output":"/bar","pipe":["print","reverse"]}`
		)

		if err := json.Unmarshal([]byte(input), &element); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
	})
}
