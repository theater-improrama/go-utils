package slice

func PartitionSplit[V any](vs []V, size int) [][]V {
	if len(vs) == 0 {
		return make([][]V, 0)
	}

	groupCount := 1 + (len(vs)-1)/size
	r := make([][]V, groupCount)

	for i := range groupCount {
		groupSize := 0
		if len(vs) >= (i+1)*size {
			groupSize = size
		} else {
			groupSize = len(vs) - i*size
		}

		r[i] = vs[i*size : i*size+groupSize]
	}

	return r
}
