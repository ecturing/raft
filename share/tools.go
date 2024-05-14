package share

// VoteCheck : 检查投票资格
func VoteCheck(self uint, remote uint) bool {
	if self < remote {
		return true
	}
	return false //任期号小于等于自己的任期号，不具备当选leader的资格
}

func CountVote(voteNum []bool, f func(bool) bool) int {
	count := 0
	for _, v := range voteNum {
		if f(v) {
			count++
		}
	}
	return count
}
