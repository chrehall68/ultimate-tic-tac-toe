package board

import "uttt/pkg/color"

// ========== abstractions ==========

type abstractProtoSpace interface {
	Owner() Owner
}

type abstractProtoCell interface {
	Get(c *Coord) abstractProtoSpace
	Owner() Owner
	Full() bool
}

// ========== abstract functions ==========

// This returns who the Owner of the  abstractProtoCell is.
// It checks diagonals as well as rows and columns
func getOwner(ac abstractProtoCell) Owner {
	c := Coord{}
	prev := Coord{}

	// check if any of the rows are claimed
	for row := 0; row < ROWS; row++ {
		allSame := true
		for col := 1; col < COLS; col++ {
			c.Row, c.Col = int32(row), int32(col)
			prev.Row, prev.Col = int32(row), int32(col-1)
			if ac.Get(&c).Owner() != ac.Get(&prev).Owner() {
				allSame = false
				break
			}
		}
		if allSame && ac.Get(&c).Owner() != Owner_NONE {
			return ac.Get(&c).Owner()
		}
	}

	// check if any of the columns have been claimed
	for col := 0; col < COLS; col++ {
		allSame := true
		for row := 1; row < ROWS; row++ {
			c.Row, c.Col = int32(row), int32(col)
			prev.Row, prev.Col = int32(row-1), int32(col)
			if ac.Get(&c).Owner() != ac.Get(&prev).Owner() {
				allSame = false
				break
			}
		}
		if allSame && ac.Get(&c).Owner() != Owner_NONE {
			return ac.Get(&c).Owner()
		}
	}

	// check the left diagonal
	allSameLeft := true
	for i := 1; i < 3; i++ {
		c.Row, c.Col = int32(i), int32(i)
		prev.Row, prev.Col = int32(i-1), int32(i-1)
		if ac.Get(&c).Owner() != ac.Get(&prev).Owner() {
			allSameLeft = false
			break
		}
	}
	if allSameLeft && ac.Get(&c).Owner() != Owner_NONE {
		return ac.Get(&c).Owner()
	}

	// check right diagonal
	allSameRight := true
	for row, col := 1, 1; row < 3; row, col = row+1, col-1 {
		c.Row, c.Col = int32(row), int32(col)
		prev.Row, prev.Col = int32(row-1), int32(col+1)
		if ac.Get(&c).Owner() != ac.Get(&prev).Owner() {
			allSameRight = false
			break
		}
	}
	if allSameRight && ac.Get(&c).Owner() != Owner_NONE {
		return ac.Get(&c).Owner()
	}

	// no one is the Owner
	return 0
}

// An abstractProtoCell is full if all of the spaces
// within it have an Owner
func isFull(ac abstractProtoCell) bool {
	c := Coord{}
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			c.Row, c.Col = int32(row), int32(col)
			if ac.Get(&c).Owner() == Owner_NONE {
				return false
			}
		}
	}
	return true
}

// ========== Coord Methods ==========

// converts the coordinate to a 1D array index
func (c *Coord) Index() (idx uint32, valid bool) {
	if !c.Valid() {
		valid = false
		return
	}
	idx = uint32(c.Col + c.Row*COLS)
	return
}

// converts a 1D array index to a coordinate
func ToCoord(idx uint32) *Coord {
	return &Coord{Row: int32(idx / COLS), Col: int32(idx % COLS)}
}

// whether or not the coordinate is valid. It is invalid if either Col or Row is negative
// or if either Col or Row is greater than the max number of columns or rows, respectively
func (c *Coord) Valid() bool {
	return (c.Col >= 0 && c.Row >= 0) && (c.Col < COLS && c.Row < ROWS)
}

// ========== Space Methods ==========
func NewProtoSpace() *Space {
	return &Space{Val: Owner_NONE}
}
func (s *Space) Owner() Owner {
	return s.Val
}

// ========== Cell Methods ==========
func NewProtoCell() *Cell {
	spaces := make([]*Space, CELLS)
	for i := 0; i < CELLS; i++ {
		spaces[i] = NewProtoSpace()
	}
	return &Cell{Spaces: spaces}
}

// Returns the corresponding *Space
func (ce *Cell) Get(co *Coord) abstractProtoSpace {
	idx, _ := co.Index()
	return ce.Spaces[idx]
}
func (c *Cell) Full() bool {
	return isFull(c)
}
func (c *Cell) Owner() Owner {
	return getOwner(c)
}

// ========== Board Methods ==========
func NewProtoBoard() *Board {
	cells := make([]*Cell, CELLS)
	for i := 0; i < CELLS; i++ {
		cells[i] = NewProtoCell()
	}
	return &Board{Cells: cells, CurCell: &Coord{Row: ROWS / 2, Col: COLS / 2}}
}

// returns the corresponding *Coord
func (b *Board) Get(c *Coord) abstractProtoSpace {
	idx, _ := c.Index()
	return b.Cells[idx]
}
func (b *Board) Full() bool {
	return isFull(b)
}
func (b *Board) Owner() Owner {
	return getOwner(b)
}

// Returns a string printable to color-supporting terminals
func (b *Board) TerminalString() string {
	ret := ""
	for row := 0; row < ROWS; row++ {
		for innerRow := 0; innerRow < ROWS; innerRow++ {
			for col := 0; col < COLS; col++ {
				if col == int(b.CurCell.Col) && row == int(b.CurCell.Row) {
					ret += color.Red
				}
				for innerCol := 0; innerCol < COLS; innerCol++ {
					switch b.Cells[row*ROWS+col].Spaces[innerRow*ROWS+innerCol].Owner() {
					case 1:
						ret += "X "
					case 2:
						ret += "O "
					default:
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
