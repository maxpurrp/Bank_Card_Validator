package lynh

import (
	"strconv"
)

func checksum(nums []int) bool {
	var sum int
	var num int
	var lst []int
	for i := 0; i < len(nums); i++ {
		if i%2 == 0 {
			num = nums[i] * 2
			if num > 9 {
				f, s := num/10, num%10
				sum += f + s
			} else {
				sum += num
			}
		} else {
			lst = append(lst, nums[i])
		}
	}
	for i := 0; i < len(lst); i++ {
		sum += lst[i]
	}
	return sum%10 == 0
}

func AlgLynh(number string) (bool, string) {
	var numslist []int
	for i := 0; i < len(number); i++ {
		num, err := strconv.Atoi(string(number[i]))
		if err != nil {
			return false, "Incorrect card number"
		}
		numslist = append(numslist, num)
	}
	res := checksum(numslist)
	if res {
		return true, "Valid success"
	} else {
		return false, "Bad card number"
	}
}
