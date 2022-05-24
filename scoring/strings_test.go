package scoring

import "testing"

func TestRemoveNonAlphaNumericWithSpace(t *testing.T) {
	newVal := removeNonAlphaNumerics(`/FR7630003034950005005419318
CHARLES DUPONT COMPANY
RUE GENERAL DE GAULLE, 21
75013 PARIS`)
	t.Log("#" + newVal + "#")
}
