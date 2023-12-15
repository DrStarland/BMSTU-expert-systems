package utils

func MySum(source [][]any) []any {
	result := []any{}
	for _, elem := range source {
		result = append(result, elem...)
	}
	return result
}
