package utils

func GetMultipleRelations(rel map[string]any) []map[string]any {
	if rel == nil || rel["data"] == nil {
		return []map[string]any{}
	}
	items := rel["data"].([]any)
	result := make([]map[string]any, len(items))
	for i, item := range items {
		result[i] = item.(map[string]any)
	}
	return result
}

func GetSingleRelation(rel map[string]any) map[string]any {
	if rel == nil || rel["data"] == nil {
		return nil
	}
	return rel["data"].(map[string]any)
}
