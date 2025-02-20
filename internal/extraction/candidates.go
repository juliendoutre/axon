package extraction

import (
	"fmt"

	"google.golang.org/protobuf/types/known/structpb"
)

type Candidate struct {
	Path  string
	Value string
}

func ExtractCandidatesFromValue(value *structpb.Value, path string) []Candidate {
	switch value := value.GetKind().(type) {
	case *structpb.Value_StringValue:
		return []Candidate{{Path: path, Value: value.StringValue}}
	case *structpb.Value_ListValue:
		return ExtractCandidatesFromList(value.ListValue, path)
	case *structpb.Value_StructValue:
		return ExtractCandidatesFromStruct(value.StructValue, path)
	default:
		return nil
	}
}

func ExtractCandidatesFromList(value *structpb.ListValue, path string) []Candidate {
	allValues := []Candidate{}

	for index, item := range value.GetValues() {
		allValues = append(allValues, ExtractCandidatesFromValue(item, fmt.Sprintf("%s[%d]", path, index))...)
	}

	return allValues
}

func ExtractCandidatesFromStruct(value *structpb.Struct, path string) []Candidate {
	allValues := []Candidate{}

	for key, value := range value.GetFields() {
		allValues = append(allValues, ExtractCandidatesFromValue(value, fmt.Sprintf("%s.%s", path, key))...)
	}

	return allValues
}
