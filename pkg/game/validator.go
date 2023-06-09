package game

// whether or not a move is valid
func validateMove(b *board, large, small coord) bool {
	if !validateCell(b, large) {
		return false
	}
	if b.get(large).get(small).owner() != 0 {
		return false
	}
	return true
}

// whether or not a destination is fit for being the next
// curCell
func validateCell(b *board, c coord) bool {
	if b.get(c).owner() != 0 || b.get(c).full() {
		return false
	}
	return true
}
