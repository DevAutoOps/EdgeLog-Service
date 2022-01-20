package tools

//  duplicate removal
func RmRepeatList(list []int) []int {
	var x []int = []int{}
	for _, i := range list {
		if len(x) == 0 {
			x = append(x, i)
		} else {
			for k, v := range x {
				if i == v {
					break
				}
				if k == len(x)-1 {
					x = append(x, i)
				}
			}
		}
	}
	return x
}

//  Determine whether the element is in the array
func IsContain(items []interface{}, item interface{}) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}
