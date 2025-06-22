package combination

func Combine(n int, k int) [][]int {
	result := [][]int{}
	currentCombination := []int{}

	backTrack(1, n, k, currentCombination, &result)

	return result
}

func backTrack(startNum int, n int, k int, currentCombination []int, results *[][]int) {
	if len(currentCombination) == k {
		temp := make([]int, k)
		copy(temp, currentCombination)
		*results = append(*results, temp)
		return
	}

	for i := startNum; i <= n-(k-len(currentCombination))+1; i++ {
		currentCombination = append(currentCombination, i)

		backTrack(i+1, n, k, currentCombination, results)

		currentCombination = currentCombination[:len(currentCombination)-1]
	}
}
