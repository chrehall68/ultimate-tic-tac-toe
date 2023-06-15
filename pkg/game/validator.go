package game

import "uttt/pkg/board"

// whether or not a move is valid
func validateMove(b *board.Board, m *board.Move) bool {
	// if the cell isn't valid for the current move
	if !validateMoveCell(b, m.Large) {
		return false
	}

	// if the destination space is taken
	if b.Get(m.Large).(*board.Cell).Get(m.Small).Owner() != 0 {
		return false
	}

	return true
}

// whether or not a cell is fit for being the site of the current move
func validateMoveCell(b *board.Board, c *board.Coord) bool {
	// if it isn't valid for moves in general, it's not valid
	if !validateCell(b, c) {
		return false
	}

	// if the CurCell is valid, check if it isn't
	if b.CurCell.Valid() {
		return c.Row == b.CurCell.Row && c.Col == b.CurCell.Col
	}

	return true
}

// whether or not a cell is fit for moves in general
func validateCell(b *board.Board, c *board.Coord) bool {
	// the cell isn't owned and the cell isn't full
	return b.Get(c).Owner() == board.Owner_NONE && !b.Get(c).(*board.Cell).Full()
}
