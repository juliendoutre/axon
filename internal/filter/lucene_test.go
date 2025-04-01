package filter_test

import (
	"testing"

	"github.com/juliendoutre/axon/internal/filter"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	t.Parallel()

	parser := filter.New()

	for testCaseName, testCase := range map[string]struct {
		query          string
		expectedFilter string
		expectedParams []any
	}{
		"basic": {
			query:          "attributes.a: test AND claims.issuer: me",
			expectedFilter: "(\"attributes.a\" = ?) AND (\"claims.issuer\" = ?)",
			expectedParams: []any{"test", "me"},
		},
	} {
		t.Run(testCaseName, func(t *testing.T) {
			t.Parallel()

			actualFilter, actualParams, err := parser.Parse(testCase.query)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedFilter, actualFilter)
			assert.ElementsMatch(t, testCase.expectedParams, actualParams)
		})
	}
}
