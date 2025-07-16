package types

func ToInterfaceSlice[T OneToMany](slice []T) []OneToMany {
	result := make([]OneToMany, len(slice))
	for i, v := range slice {
		result[i] = v
	}
	return result
}
