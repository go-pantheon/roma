package align

func Align(slice []string, size int) []string {
	if len(slice) > size {
		return slice[:size]
	}

	if len(slice) < size {
		slice = append(slice, make([]string, size-len(slice))...)
	}

	return slice
}
