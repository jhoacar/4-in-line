package table

func GetTable[Type int | byte | string](rows int, columns int) [][]Type {
	table := make([][]Type, rows) // initialize a slice of rows slices
	for i := int(0); i < rows; i++ {
		table[i] = make([]Type, columns) // initialize a slice of columns unit8 in each of dy slices
	}
	return table
}
