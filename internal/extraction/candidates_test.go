package extraction_test

import (
	"testing"

	"github.com/juliendoutre/axon/internal/extraction"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"
)

//nolint:funlen
func TestCandidates(t *testing.T) {
	t.Parallel()

	for testCaseName, testCase := range map[string]struct {
		Value              *structpb.Value
		ExpectedCandidates []extraction.Candidate
	}{
		"empty": {
			Value:              &structpb.Value{},
			ExpectedCandidates: nil,
		},
		"null": {
			Value:              structpb.NewNullValue(),
			ExpectedCandidates: nil,
		},
		"number": {
			Value:              structpb.NewNumberValue(1.0),
			ExpectedCandidates: nil,
		},
		"boolean": {
			Value:              structpb.NewBoolValue(true),
			ExpectedCandidates: nil,
		},
		"string": {
			Value:              structpb.NewStringValue("test"),
			ExpectedCandidates: []extraction.Candidate{{Path: "$", Value: "test"}},
		},
		"nested": {
			Value: structpb.NewStructValue(&structpb.Struct{
				Fields: map[string]*structpb.Value{"a": structpb.NewStructValue(&structpb.Struct{
					Fields: map[string]*structpb.Value{"b": structpb.NewStringValue("test")},
				})},
			}),
			ExpectedCandidates: []extraction.Candidate{{Path: "$.a.b", Value: "test"}},
		},
		"list": {
			Value: structpb.NewStructValue(&structpb.Struct{
				Fields: map[string]*structpb.Value{"a": structpb.NewStructValue(&structpb.Struct{
					Fields: map[string]*structpb.Value{
						"b": structpb.NewStringValue("test"),
						"c": structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{
							structpb.NewStringValue("hello"),
							structpb.NewBoolValue(true),
							structpb.NewNumberValue(1.0),
							structpb.NewStringValue("world"),
						}}),
					},
				})},
			}),
			ExpectedCandidates: []extraction.Candidate{
				{Path: "$.a.b", Value: "test"},
				{Path: "$.a.c[0]", Value: "hello"},
				{Path: "$.a.c[3]", Value: "world"},
			},
		},
	} {
		t.Run(testCaseName, func(t *testing.T) {
			t.Parallel()

			actualCandidates := extraction.ExtractCandidatesFromValue(testCase.Value, "$")

			assert.ElementsMatch(t, testCase.ExpectedCandidates, actualCandidates)
		})
	}
}
