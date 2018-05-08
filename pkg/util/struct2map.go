package util

import "encoding/json"
func ToMap(in2 interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	inrec, _ := json.Marshal(in2)
	json.Unmarshal(inrec, &m)
	return m
}
