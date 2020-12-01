package game

import (
    "errors"

    "GoChessgameServer/store"
    u "GoChessgameServer/util"
)

// Errors
var CantPassError = errors.New("Cant pass with this figure on specified location")

// This type represent a chess table
type ChessTable [8][8]Figure

// This type represents a game session
// that contains game states
type GameSession struct {
    gameStore *store.GameStore
    cTable *ChessTable
}

// Interface that represents a chess figure
type Figure interface {
    CanFigurePass(int, int) bool
    IsFigureBlack() bool
    Pass(int, int) error
}

// This types represent specific figures,
// that realizes figure interface
type Pawn struct {
    isMoved bool
    cTable *ChessTable
    x int
    y int
    isBlack bool
}

func (p Pawn) CanFigurePass(altX int, altY int) bool {

    if p.isBlack {
        return p.canFigurePassBlack(-altX, -altY)
    } else {
        return p.canFigurePassWhite(altX, altY)
    }
}

func (p Pawn) canFigurePassWhite(altX int, altY int) bool {
    table := (*p.cTable)

    // Check for forward going
    if p.isMoved && altX == 0 {
        if p.y <= 7 && altY == 1 && table[p.x][p.y + 1] == nil {
            return true
        }
    } else {
        if (p.y <= 7 && altY == 1 && table[p.x][p.y + 1] == nil) || (p.x <= 6 && (altY == 2 || altY == 1) && table[p.x][p.y + 1] == nil && table[p.x][p.y + 2] == nil) {
            return true
        }
    }

    // Check for figure beating
    figureOnTheRight := table[p.x + 1][p.y + 1]
    figureOnTheLeft := table[p.x - 1][p.y + 1]

    switch fig := figureOnTheRight.(type) {
    case King:
    default:
        if fig != nil && fig.IsFigureBlack() && altY == 1 && altX == 1 {
            return true
        }
    }

    switch fig := figureOnTheLeft.(type) {
    case King:
    default:
        if fig != nil && fig.IsFigureBlack() && altY == -1 && altX == 1 {
            return true
        }
    }

    return false
}

func (p Pawn) canFigurePassBlack(altX int, altY int) bool {
    table := (*p.cTable)

    // Check for forward going
    if p.isMoved && altX == 0 {
        if p.y >= 1 && altY == -1 && table[p.x][p.y - 1] == nil {
            return true
        }
    } else {
        if (p.y >= 1 && altY == -1 && table[p.x][p.y - 1] == nil) || (p.x >= 2 && (altY == -2 || altY == -1) && table[p.x][p.y - 1] == nil && table[p.x][p.y - 2] == nil) {
            return true
        }
    }

    // Check for figure beating
    figureOnTheLeft := table[p.x - 1][p.y - 1]
    figureOnTheRight := table[p.x + 1][p.y - 1]

    switch fig := figureOnTheRight.(type) {
    case King:
    default:
        if fig != nil && !fig.IsFigureBlack() && altY == -1 && altX == -1 {
            return true
        }
    }

    switch fig := figureOnTheLeft.(type) {
    case King:
    default:
        if fig != nil && !fig.IsFigureBlack() && altY == 1 && altX == -1 {
            return true
        }
    }

    return false
}

func (p Pawn) Pass(altX int, altY int) error {

    table := (*p.cTable)
    // Check is we can pass our figure to specified location
    if !p.CanFigurePass(altX, altY) {
        return CantPassError
    }

    if !p.isMoved {
        p.isMoved = true
    }

    if p.isBlack {
        p.x -= altX
        p.y -= altY

        table[p.x - altX][p.y - altY] = table[p.x][p.y]
        table[p.x + altX][p.y + altY] = nil
    } else {
        p.x += altX
        p.y += altY

        table[p.x + altX][p.y + altY] = table[p.x][p.y]
        table[p.x - altX][p.y - altY] = nil
    }

    return nil
}

func (p Pawn) IsFigureBlack() bool {
    return p.isBlack
}

type Rook struct {
    cTable *ChessTable
    x int
    y int
    isBlack bool
}

func (r Rook) CanFigurePass(altX int, altY int) bool {

    if r.isBlack {
        return r.canFigurePassInner(-altX, -altY)
    } else {
        return r.canFigurePassInner(altX, altY)
    }
}

func (r Rook) canFigurePassInner(altX int, altY int) bool {
    table := (*r.cTable)

    if altY == 0 && altX != 0 && r.x + altX <= 8 && r.x + altY >= 0 {
        if altX > 0 {
            count := r.x + 1
            for i := count; i < r.x + altX; i++ {
                if table[i][r.y] != nil {
                    return false
                }
            }
        } else {
            count := r.x - 1
            for i := count; i > r.x + altX; i-- {
                if table[i][r.y] != nil {
                    return false
                }
            }
        }

        return !(table[r.x + altX][r.y].IsFigureBlack() == r.isBlack)
    }

    if altX == 0 && altY != 0 && r.y + altY <= 8 && r.y + altY >= 0 {
        if altY > 0 {
            count := r.y + 1
            for i := count; i < r.y + altY; i++ {
                if table[r.x][i] != nil {
                    return false
                }
            }
        } else {
            count := r.y - 1
            for i := count; i > r.y + altY; i-- {
                if table[r.x][i] != nil {
                    return false
                }
            }
        }

        return !(table[r.x][r.y + altY].IsFigureBlack() == r.isBlack)
    }

    return false
}

func (r Rook) Pass(altX int, altY int) error {

    table := (*r.cTable)

    if !r.CanFigurePass(altX, altY) {
        return CantPassError
    }

    if r.isBlack {
        r.x -= altX
        r.y -= altY

        table[r.x - altX][r.y - altY] = table[r.x][r.y]
        table[r.x + altX][r.y + altY] = nil
    } else {
        r.x += altX
        r.y += altY

        table[r.x + altX][r.y + altY] = table[r.x][r.y]
        table[r.x - altX][r.y - altY] = nil
    }

    return nil
}

func (r Rook) IsFigureBlack() bool {
    return r.isBlack
}

type Horse struct {
    cTable *ChessTable
    x int
    y int
    isBlack bool
}

func (h Horse) CanFigurePass(altX int, altY int) bool {

    if h.isBlack {
        return h.canFigurePassInner(-altX, -altY)
    } else {
        return h.canFigurePassInner(altX, altY)
    }
}

func (h Horse) canFigurePassInner(altX int, altY int) bool {
    table := (*h.cTable)

    if altX != 0 && altY != 0 && ((u.Abs(altX) == 2 && u.Abs(altY) == 1) || (u.Abs(altX) == 1 && u.Abs(altY) == 2)) {
        destelem := table[h.x + altX][h.y + altY]
        if destelem != nil && (destelem.IsFigureBlack() == h.isBlack) {
            return false
        }
        return true
    }

    return false
}

func (h Horse) Pass(altX int, altY int) error {

    table := (*h.cTable)

    if !h.CanFigurePass(altX, altY) {
        return CantPassError
    }

    if h.isBlack {
        h.x -= altX
        h.y -= altY

        table[h.x - altX][h.y - altY] = table[h.x][h.y]
        table[h.x - altX][h.y + altY] = nil
    } else {
        h.x += altX
        h.y += altY

        table[h.x + altX][h.y + altY] = table[h.x][h.y]
        table[h.x - altX][h.y - altY] = nil
    }

    return nil
}

func (h Horse) IsFigureBlack() bool {
    return h.isBlack
}

type Elephant struct {
    cTable *ChessTable
    x int
    y int
    isBlack bool
}

func (e Elephant) CanFigurePass(altX int, altY int) bool {

    if e.isBlack {
        return e.canFigurePassInner(-altX, -altY)
    } else {
        return e.canFigurePassInner(altX, altY)
    }
}

func (e Elephant) canFigurePassInner(altX int, altY int) bool {
    table := (*e.cTable)

    if altY != 0 && altX != 0 && u.Abs(altX) == u.Abs(altY) && e.y + altY <= 8 && e.y + altY >= 0 && e.x + altX <= 8 && e.x + altX >= 0 {
        if altX > 0 && altY > 0 {
            countX := e.x + 1
            countY := e.y + 1

            for e.x + altX > countX && e.y + altX > countY {
                if table[countX][countY] != nil {
                    return false
                }
                countX++
                countY++
            }
        } else if altX > 0 && altY < 0 {
            countX := e.x + 1
            countY := e.y - 1

            for e.x + altX > countX && e.y + altX < countY {
                if table[countX][countY] != nil {
                    return false
                }
                countX++
                countY--
            }
        } else if altX < 0 && altY > 0 {
            countX := e.x - 1
            countY := e.y + 1

            for e.x + altX < countX && e.y + altX > countY {
                if table[countX][countY] != nil {
                    return false
                }
                countX--
                countY++
            }
        } else {
            countX := e.x - 1
            countY := e.y - 1

            for e.x + altX < countX && e.y + altX < countY {
                if table[countX][countY] != nil {
                    return false
                }
                countX--
                countY--
            }
        }

        return !(table[e.x + altX][e.y + altY].IsFigureBlack() == e.isBlack)
    }

    return false
}

func (e Elephant) Pass(altX int, altY int) error {

    table := (*e.cTable)

    if !e.CanFigurePass(altX, altY) {
        return CantPassError
    }

    if e.isBlack {
        e.x -= altX
        e.y -= altY

        table[e.x - altX][e.y - altY] = table[e.x][e.y]
        table[e.x - altX][e.y + altY] = nil
    } else {
        e.x += altX
        e.y += altY

        table[e.x + altX][e.y + altY] = table[e.x][e.y]
        table[e.x - altX][e.y - altY] = nil
    }

    return nil
}

func (e Elephant) IsFigureBlack() bool {
    return e.isBlack
}

type Queen struct {
    cTable *ChessTable
    x int
    y int
    isBlack bool
}

func (q Queen) CanFigurePass(altX int, altY int) bool {

    if q.isBlack {
        return q.canFigurePassInner(-altX, altY)
    } else {
        return q.canFigurePassInner(altX, altY)
    }
}

func (q Queen) canFigurePassInner(altX int, altY int) bool {
    table := (*q.cTable)

    if q.y + altY <= 8 && q.y + altY >= 0 && q.x + altX <= 8 && q.x + altX >= 0 {
        if altY == 0 && altX != 0 {
            if altX > 0 {
                count := q.x + 1
                for i := count; i < q.x + altX; i++ {
                    if table[i][q.y] != nil {
                        return false
                    }
                }
            } else {
                count := q.x - 1
                for i := count; i > q.x + altX; i-- {
                    if table[i][q.y] != nil {
                        return false
                    }
                }
            }

            return !(table[q.x + altX][q.y].IsFigureBlack() == q.isBlack)
        }

        if altX == 0 && altY != 0 {
            if altY > 0 {
                count := q.y + 1
                for i := count; i < q.y + altY; i++ {
                    if table[q.x][i] != nil {
                        return false
                    }
                }
            } else {
                count := q.y - 1
                for i := count; i > q.y + altY; i-- {
                    if table[q.x][i] != nil {
                        return false
                    }
                }
            }

            return !(table[q.x][q.y + altY].IsFigureBlack() == q.isBlack)
        }

        if altY != 0 && altX != 0 && u.Abs(altX) == u.Abs(altY) {
            if altX > 0 && altY > 0 {
                countX := q.x + 1
                countY := q.y + 1

                for q.x + altX > countX && q.y + altX > countY {
                    if table[countX][countY] != nil {
                        return false
                    }
                    countX++
                    countY++
                }
            } else if altX > 0 && altY < 0 {
                countX := q.x + 1
                countY := q.y - 1

                for q.x + altX > countX && q.y + altX < countY {
                    if table[countX][countY] != nil {
                        return false
                    }
                    countX++
                    countY--
                }
            } else if altX < 0 && altY > 0 {
                countX := q.x - 1
                countY := q.y + 1

                for q.x + altX < countX && q.y + altX > countY {
                    if table[countX][countY] != nil {
                        return false
                    }
                    countX--
                    countY++
                }
            } else {
                countX := q.x - 1
                countY := q.y - 1

                for q.x + altX < countX && q.y + altX < countY {
                    if table[countX][countY] != nil {
                        return false
                    }
                    countX--
                    countY--
                }
            }

            return !(table[q.x + altX][q.y + altY].IsFigureBlack() == q.isBlack)
        }
    }

    return false
}

func (q Queen) Pass(altX int, altY int) error {

    table := (*q.cTable)

    if !q.CanFigurePass(altX, altY) {
        return CantPassError
    }

    if q.isBlack {
        q.x -= altX
        q.y -= altY

        table[q.x - altX][q.y - altY] = table[q.x][q.y]
        table[q.x - altX][q.y + altY] = nil
    } else {
        q.x += altX
        q.y += altY

        table[q.x + altX][q.y + altY] = table[q.x][q.y]
        table[q.x - altX][q.y - altY] = nil
    }

    return nil
}

func (q Queen) IsFigureBlack() bool {
    return q.isBlack
}

type King struct {
    cTable *ChessTable
    x int
    y int
    isBlack bool
}

func (k King) CanFigurePass(altX int, altY int) bool {

    if k.isBlack {
        return k.canFigurePassInner(-altX, -altY)
    } else {
        return k.canFigurePassInner(altX, altY)
    }
}

func (k King) canFigurePassInner(altX int, altY int) bool {
    table := (*k.cTable)

    if altX == 0 && altY == 0 { return false; }
    if (altX == 0 || u.Abs(altX) == 1) && (altY == 0 || u.Abs(altY) == 1) {
        if table[k.x + altX][k.y + altY] != nil && table[k.x + altX][k.y + altY].IsFigureBlack() == k.isBlack {
            return false
        }
        return true
    }

    return false
}

func (k King) Pass(altX int, altY int) error {

    table := (*k.cTable)

    if !k.CanFigurePass(altX, altY) {
        return CantPassError
    }

    if k.isBlack {
        k.x -= altX
        k.y -= altY

        table[k.x - altX][k.y - altY] = table[k.x][k.y]
        table[k.x - altX][k.y + altY] = nil
    } else {
        k.x += altX
        k.y += altY

        table[k.x + altX][k.y + altY] = table[k.x][k.y]
        table[k.x - altX][k.y - altY] = nil
    }

    return nil
}

func (k King) IsFigureBlack() bool {
    return k.isBlack
}
