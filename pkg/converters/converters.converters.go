package converters

import "encoding/json"

func ToJson(obj any) (map[string]any, error) {
	m := make(map[string]any)

	marshaledObj, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(marshaledObj, &m); err != nil {
		return nil, err
	}

	return m, nil
}
