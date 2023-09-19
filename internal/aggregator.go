package internal

func MergeConfigurations(configurations Configurations) Configuration {
	var config Configuration
	for _, c := range configurations {
		config.HTTP = deepMerge(config.HTTP, c.HTTP)
		config.TCP = deepMerge(config.TCP, c.TCP)
		config.UDP = deepMerge(config.UDP, c.UDP)
		config.TLS = deepMerge(config.TLS, c.TLS)
	}

	return config
}

// deepMerge recursively merges two YAML structures.
func deepMerge(a, b interface{}) interface{} {
	switch ta := a.(type) {
	case map[interface{}]interface{}:
		tb, ok := b.(map[interface{}]interface{})
		if !ok {
			return a
		}
		result := make(map[interface{}]interface{})
		for k, va := range ta {
			result[k] = va
		}
		for k, vb := range tb {
			if va, exists := result[k]; exists {
				result[k] = deepMerge(va, vb)
			} else {
				result[k] = vb
			}
		}
		return result
	case map[string]interface{}:
		tb, ok := b.(map[string]interface{})
		if !ok {
			return a
		}
		result := make(map[string]interface{})
		for k, va := range ta {
			result[k] = va
		}
		for k, vb := range tb {
			if va, exists := result[k]; exists {
				result[k] = deepMerge(va, vb)
			} else {
				result[k] = vb
			}
		}
		return result
	case []interface{}:
		tb, ok := b.([]interface{})
		if !ok {
			return a
		}
		result := append([]interface{}{}, ta...)
		result = append(result, tb...)
		return result

	default:
		return b
	}
}
