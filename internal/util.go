package internal

func Reverse(x [][]interface{}) {
	n := len(x)
	for i := 0; i < n/2; i++ {
		x[i], x[n-i-1] = x[n-i-1], x[i]
	}
}
