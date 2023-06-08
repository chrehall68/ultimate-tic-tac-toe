package game

import "uttt/pkg/color"

// ========== abstractions ==========

type abstractSpace interface {
	owner() uint8
}

type abstractCell interface {
	cells() [3][3]abstractSpace
	owner() uint8
}

// ========== struct and interface definitions ==========

// internal coordinate representation
type coord struct {
	row uint8
	col uint8
}

// a certain space
type space struct {
	val uint8
}

// a cell (one tic tac toe game board) on the board
type cell struct {
	spaces [3][3]space
}

// the full ultimate tic tac toe game board
type board struct {
	spaces  [3][3]cell
	curCell coord
}

// ========== owner method for boards and cells ==========
func getOwner(ac abstractCell) uint8 {
	// check if any of the rows are claimed
	for row := 0; row < 3; row++ {
		allSame := true
		for col := 1; col < 3; col++ {
			if ac.cells()[row][col].owner() != ac.cells()[row][col-1].owner() {
				allSame = false
				break
			}
		}
		if allSame && ac.cells()[row][0].owner() != 0 {
			return ac.cells()[row][0].owner()
		}
	}

	// check if any of the columns have been claimed
	for col := 0; col < 3; col++ {
		allSame := true
		for row := 1; row < 3; row++ {
			if ac.cells()[row][col].owner() != ac.cells()[row-1][col].owner() {
				allSame = false
				break
			}
		}
		if allSame && ac.cells()[0][col].owner() != 0 {
			return ac.cells()[0][col].owner()
		}
	}

	// check the left diagonal
	allSameLeft := true
	for i := 1; i < 3; i++ {
		if ac.cells()[i][i].owner() != ac.cells()[i-1][i-1].owner() {
			allSameLeft = false
			break
		}
	}
	if allSameLeft && ac.cells()[0][0].owner() != 0 {
		return ac.cells()[0][0].owner()
	}

	// check right diagonal
	allSameRight := true
	for row, col := 1, 1; row < 3; row, col = row+1, col-1 {
		if ac.cells()[row][col].owner() != ac.cells()[row-1][col+1].owner() {
			allSameRight = false
			break
		}
	}
	if allSameRight && ac.cells()[0][2].owner() != 0 {
		return ac.cells()[0][0].owner()
	}

	// no one is the owner
	return 0
}

// ========== space methods ==========
func (s *space) owner() uint8 {
	return s.val
}

// ========== cell methods ==========
func (c *cell) cells() [3][3]abstractSpace {
	var abstractSpaces [3][3]abstractSpace
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			abstractSpaces[row][col] = &c.spaces[row][col]
		}
	}
	return abstractSpaces
}
func (c *cell) owner() uint8 {
	return getOwner(c)
}

// ========== board methods ==========
func (b *board) cells() [3][3]abstractSpace {
	var abstractSpaces [3][3]abstractSpace
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			abstractSpaces[row][col] = &b.spaces[row][col]
		}
	}
	return abstractSpaces
}
func (b *board) owner() uint8 {
	return getOwner(b)
}
func (b *board) String() string {
	dims := 3
	ret := ""
	for row := 0; row < dims; row++ {
		for innerRow := 0; innerRow < dims; innerRow++ {
			for col := 0; col < dims; col++ {
				if col == int(b.curCell.col) && row == int(b.curCell.row) {
					ret += color.Red
				}
				for innerCol := 0; innerCol < dims; innerCol++ {
					if b.spaces[row][col].spaces[innerRow][innerCol].owner() == 1 {
						ret += "X "
					} else if b.spaces[row][col].spaces[innerRow][innerCol].owner() == 2 {
						ret += "O "
					} else {
						ret += "_ "
					}
				}
				ret += color.Reset
				ret += "| "
			}
			ret += "\n"
		}
		ret += "------------------------\n"
	}

	ret += "\n"
	return ret
}

// ========== New Methods ==========

func newSpace() space {
	s := space{val: 0}
	return s
}

func newCell() cell {
	// initialize cell spaces
	var spaces [3][3]space
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			spaces[row][col] = newSpace()
		}
	}

	c := cell{spaces}
	return c
}

func newBoard() board {
	// initialize board spaces
	var spaces [3][3]cell
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			spaces[row][col] = newCell()
		}
	}
	// board starting cell
	startCell := coord{row: 1, col: 1}

	b := board{curCell: startCell, spaces: spaces}
	return b
}
