package game

import "uttt/pkg/board"

// whether or not a move is valid
func validateMove(b *board.Board, large, small *board.Coord) bool {
	if !validateCell(b, large) {
		return false
	}
	if b.Get(large).(*board.Cell).Get(small).Owner() != 0 {
		return false
	}
	return true
}

// whether or not a destination is fit for being the next
// curCell
func validateCell(b *board.Board, c *board.Coord) bool {
	if b.Get(c).Owner() != board.Owner_NONE || b.Get(c).(*board.Cell).Full() {
		return false
	}
	return true
}
