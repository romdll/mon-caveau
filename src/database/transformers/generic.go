package transformers

func PassListByTransformer[E, O any](entry []E, transformer func(E) O) []O {
	res := make([]O, len(entry))

	for i, val := range entry {
		res[i] = transformer(val)
	}

	return res
}
