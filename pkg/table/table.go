package table

func GetTable[Type byte | string](rows uint8, columns uint8) [][]Type {
	table := make([][]Type, rows) // initialize a slice of rows slices
	for i := uint8(0); i < rows; i++ {
		table[i] = make([]Type, columns) // initialize a slice of columns unit8 in each of dy slices
	}
	return table
}
