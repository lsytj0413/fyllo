package compare

import "fmt"

var numbers = 100000
var findCount = 100

func treeLen(count int) int {
	// sum, i := 0, 1
	// for i < count {
	// 	sum += i
	// 	// fmt.Println(i)
	// 	i *= 2
	// }

	// return sum + count
	return 2*count - 1
}

func leftChild(i int) int {
	return (i+1)*2 - 1
}

func rightChild(i int) int {
	return (i + 1) * 2
}

func parentIndex(i int) int {
	return (i+1)/2 - 1
}

func topn(nums []int) ([]int, int) {
	ret, count, w := make([]int, findCount), 0, 0
	tree := make([]int, treeLen(numbers))

	for i, j := len(tree)-1, len(nums)-1; j >= 0; i, j = i-1, j-1 {
		tree[i] = nums[j]
	}

	// 初始化树结构
	for i := len(tree) - numbers - 1; i >= 0; i-- {
		left, right := 0, 0
		if leftChild(i) < len(tree) {
			left = tree[leftChild(i)]
		}
		if rightChild(i) < len(tree) {
			right = tree[rightChild(i)]
		}

		count++
		if left > right {
			tree[i] = left
		} else {
			tree[i] = right
		}
	}

	// 保存第一个最大值
	ret[w] = tree[0]
	w++

	// 循环查找剩下的
	for w < findCount {
		// 查找是哪个 index, 此处可以优化
		index := 0
		for j := len(nums) - 1; j >= 0; j = j - 1 {
			if nums[j] == ret[w-1] {
				index = j
				break
			}
		}

		index = len(tree) - len(nums) + index
		tree[index] = 0

		// 循环处理该 index 的 parent 并进行更新
		for {
			parent := parentIndex(index)
			if parent < 0 {
				break
			}

			left, right := 0, 0
			if leftChild(parent) < len(tree) {
				left = tree[leftChild(parent)]
			}
			if rightChild(parent) < len(tree) {
				right = tree[rightChild(parent)]
			}

			count++
			if left > right {
				tree[parent] = left
			} else {
				tree[parent] = right
			}

			index = parent
		}

		ret[w] = tree[0]
		w++
	}

	// fmt.Println(tree)
	fmt.Println(ret)
	fmt.Println(count)
	return ret, count
}
