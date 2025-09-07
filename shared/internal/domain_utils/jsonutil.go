package domain_tools

import (
	"encoding/json"
)

// MarshalIndentWithoutNulls serializa uma struct para JSON identado omitindo campos nulos (nulls) do resultado.
// Para funcionar, os campos da struct devem ter a tag `json:",omitempty"`.
func MarshalIndentWithoutNulls(v interface{}, prefix, indent string) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	m = cleanNulls(m)
	return json.MarshalIndent(m, prefix, indent)
}

// cleanNulls remove todos os campos nulos de um map[string]interface{} recursivamente.
func cleanNulls(m map[string]interface{}) map[string]interface{} {
	for k, v := range m {
		switch vv := v.(type) {
		case nil:
			delete(m, k)
		case map[string]interface{}:
			m[k] = cleanNulls(vv)
		case []interface{}:
			for i := range vv {
				if sub, ok := vv[i].(map[string]interface{}); ok {
					vv[i] = cleanNulls(sub)
				}
			}
		}
	}
	return m
}
