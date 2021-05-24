package game

import (
    "errors"

    u "GoChessgameServer/util"
)

var (
    // Errors
    CantPassError = errors.New("Cant pass with this figure on specified location")
)

// This type represent a chess table
type ChessTable [8][8]Figure

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

func (p *Pawn) CanFigurePass(altX int, altY int) bool {

    if p.isBlack {
        return p.canFigurePassBlack(-altX, -altY)
    } else {
        return p.canFigurePassWhite(altX, altY)
    }
}

func (p *Pawn) canFigurePassWhite(altX int, altY int) bool {
    table := (*p.cTable)

    // Check for forward going
    if p.isMoved && altY == 0 {
        if p.x <= 6 && altX == 1 && table[p.x + 1][p.y] == nil {
            return true
        }
    } else {
        if (p.x <= 6 && altX == 1 && table[p.x + 1][p.y] == nil) || (p.x <= 5 && (altX == 2 || altX == 1) && table[p.x + 1][p.y] == nil && table[p.x + 1][p.y] == nil) {
            return true
        }
    }

    // Check for figure beating
    if p.y >= 1 && p.x < 7 {
        figureOnTheLeft := table[p.x + 1][p.y - 1]

        switch fig := figureOnTheLeft.(type) {
        case *King:
        default:
            if fig != nil && fig.IsFigureBlack() && altY == -1 && altX == 1 {
                return true
            }
        }
    }

    if p.y <= 6 && p.x < 7 {
        figureOnTheRight := table[p.x + 1][p.y + 1]

        switch fig := figureOnTheRight.(type) {
        case *King:
        default:
            if fig != nil && fig.IsFigureBlack() && altY == 1 && altX == 1 {
                return true
            }
        }
    }

    return false
}

func (p *Pawn) canFigurePassBlack(altX int, altY int) bool {
    table := (*p.cTable)

    // Check for forward going
    if p.isMoved && altY == 0 {
        if p.x >= 1 && altX == -1 && table[p.x - 1][p.y] == nil {
            return true
        }
    } else {
        if (p.x >= 1 && altX == -1 && table[p.x - 1][p.y] == nil) || (p.x >= 2 && (altX == -2 || altX == -1) && table[p.x - 1][p.y] == nil && table[p.x - 2][p.y] == nil) {
            return true
        }
    }

    // Check for figure beating
    if p.y >= 1 && p.x > 0 {
        figureOnTheRight := table[p.x - 1][p.y - 1]

        switch fig := figureOnTheRight.(type) {
        case *King:
        default:
            if fig != nil && !fig.IsFigureBlack() && altY == -1 && altX == -1 {
                return true
            }
        }
    }

    if p.y <= 6 && p.x > 0 {
        figureOnTheLeft := table[p.x - 1][p.y + 1]

        switch fig := figureOnTheLeft.(type) {
        case *King:
        default:
            if fig != nil && !fig.IsFigureBlack() && altY == 1 && altX == -1 {
                return true
            }
        }
    }

    return false
}

func (p *Pawn) Pass(altX int, altY int) error {

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

        (*p.cTable)[p.x][p.y] = (*p.cTable)[p.x + altX][p.y + altY]
        (*p.cTable)[p.x + altX][p.y + altY] = nil
    } else {
        p.x += altX
        p.y += altY

        (*p.cTable)[p.x][p.y] = (*p.cTable)[p.x - altX][p.y - altY]
        (*p.cTable)[p.x - altX][p.y - altY] = nil
    }

    return nil
}

func (p *Pawn) IsFigureBlack() bool {
    return p.isBlack
}

type Rook struct {
    cTable *ChessTable
    x int
    y int
    isBlack bool
}

func (r *Rook) CanFigurePass(altX int, altY int) bool {

    if r.isBlack {
        return r.canFigurePassInner(-altX, -altY)
    } else {
        return r.canFigurePassInner(altX, altY)
    }
}

func (r *Rook) canFigurePassInner(altX int, altY int) bool {
    table := (*r.cTable)

    if altY == 0 && altX != 0 && r.x + altX <= 7 && r.x + altY >= 0 {
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

        if (table[r.x + altX][r.y] == nil) { return true; }
        return !(table[r.x + altX][r.y].IsFigureBlack() == r.isBlack)
    }

    if altX == 0 && altY != 0 && r.y + altY <= 7 && r.y + altY >= 0 {
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

        if (table[r.x][r.y + altY] == nil) { return true; }
        return !(table[r.x][r.y + altY].IsFigureBlack() == r.isBlack)
    }

    return false
}

func (r *Rook) Pass(altX int, altY int) error {

    if !r.CanFigurePass(altX, altY) {
        return CantPassError
    }

    if r.isBlack {
        r.x -= altX
        r.y -= altY

        (*r.cTable)[r.x][r.y] = (*r.cTable)[r.x + altX][r.y + altY]
        (*r.cTable)[r.x + altX][r.y + altY] = nil
    } else {
        r.x += altX
        r.y += altY

        (*r.cTable)[r.x][r.y] = (*r.cTable)[r.x - altX][r.y - altY]
        (*r.cTable)[r.x - altX][r.y - altY] = nil
    }

    return nil
}

func (r *Rook) IsFigureBlack() bool {
    return r.isBlack
}

type Knight struct {
    cTable *ChessTable
    x int
    y int
    isBlack bool
}

func (h *Knight) CanFigurePass(altX int, altY int) bool {

    if h.isBlack {
        return h.canFigurePassInner(-altX, -altY)
    } else {
        return h.canFigurePassInner(altX, altY)
    }
}

func (h *Knight) canFigurePassInner(altX int, altY int) bool {
    table := (*h.cTable)

    if altX != 0 && altY != 0 && ((u.Abs(altX) == 2 && u.Abs(altY) == 1) || (u.Abs(altX) == 1 && u.Abs(altY) == 2)) && h.x + altX <= 7 && h.x + altX >= 0 && h.y + altY <= 7 && h.y + altY >= 0 {
        destelem := table[h.x + altX][h.y + altY]
        if destelem != nil && (destelem.IsFigureBlack() == h.isBlack) {
            return false
        }
        return true
    }

    return false
}

func (h *Knight) Pass(altX int, altY int) error {

    if !h.CanFigurePass(altX, altY) {
        return CantPassError
    }

    if h.isBlack {
        h.x -= altX
        h.y -= altY

        (*h.cTable)[h.x][h.y] = (*h.cTable)[h.x + altX][h.y + altY]
        (*h.cTable)[h.x + altX][h.y + altY] = nil
    } else {
        h.x += altX
        h.y += altY

        (*h.cTable)[h.x][h.y] = (*h.cTable)[h.x - altX][h.y - altY]
        (*h.cTable)[h.x - altX][h.y - altY] = nil
    }

    return nil
}

func (h *Knight) IsFigureBlack() bool {
    return h.isBlack
}

type Bishop struct {
    cTable *ChessTable
    x int
    y int
    isBlack bool
}

func (e *Bishop) CanFigurePass(altX int, altY int) bool {

    if e.isBlack {
        return e.canFigurePassInner(-altX, -altY)
    } else {
        return e.canFigurePassInner(altX, altY)
    }
}

func (e *Bishop) canFigurePassInner(altX int, altY int) bool {
    table := (*e.cTable)

    if altY != 0 && altX != 0 && u.Abs(altX) == u.Abs(altY) && e.y + altY <= 7 && e.y + altY >= 0 && e.x + altX <= 7 && e.x + altX >= 0 {
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

        if (table[e.x + altX][e.y + altY] == nil) { return true; }
        return !(table[e.x + altX][e.y + altY].IsFigureBlack() == e.isBlack)
    }

    return false
}

func (e *Bishop) Pass(altX int, altY int) error {

    if !e.CanFigurePass(altX, altY) {
        return CantPassError
    }

    if e.isBlack {
        e.x -= altX
        e.y -= altY

        (*e.cTable)[e.x][e.y] = (*e.cTable)[e.x + altX][e.y + altY]
        (*e.cTable)[e.x + altX][e.y + altY] = nil
    } else {
        e.x += altX
        e.y += altY

        (*e.cTable)[e.x][e.y] = (*e.cTable)[e.x - altX][e.y - altY]
        (*e.cTable)[e.x - altX][e.y - altY] = nil
    }

    return nil
}

func (e *Bishop) IsFigureBlack() bool {
    return e.isBlack
}

type Queen struct {
    cTable *ChessTable
    x int
    y int
    isBlack bool
}

func (q *Queen) CanFigurePass(altX int, altY int) bool {

    if q.isBlack {
        return q.canFigurePassInner(-altX, -altY)
    } else {
        return q.canFigurePassInner(altX, altY)
    }
}

func (q *Queen) canFigurePassInner(altX int, altY int) bool {
    table := (*q.cTable)

    if q.y + altY <= 7 && q.y + altY >= 0 && q.x + altX <= 7 && q.x + altX >= 0 {
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

            if (table[q.x + altX][q.y] == nil) { return true; }
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

            if (table[q.x][q.y + altY] == nil) { return true; }
            return !(table[q.x][q.y + altY].IsFigureBlack() == q.isBlack)
        }

        if altY != 0 && altX != 0 && u.Abs(altX) == u.Abs(altY) {
            if altX > 0 && altY > 0 {
                countX := q.x + 1
                countY := q.y + 1

                for q.x + altX > countX && q.y + altY > countY {
                    if table[countX][countY] != nil {
                        return false
                    }
                    countX++
                    countY++
                }
            } else if altX > 0 && altY < 0 {
                countX := q.x + 1
                countY := q.y - 1

                for q.x + altX > countX && q.y + altY < countY {
                    if table[countX][countY] != nil {
                        return false
                    }
                    countX++
                    countY--
                }
            } else if altX < 0 && altY > 0 {
                countX := q.x - 1
                countY := q.y + 1

                for q.x + altX < countX && q.y + altY > countY {
                    if table[countX][countY] != nil {
                        return false
                    }
                    countX--
                    countY++
                }
            } else {
                countX := q.x - 1
                countY := q.y - 1

                for q.x + altX < countX && q.y + altY < countY {
                    if table[countX][countY] != nil {
                        return false
                    }
                    countX--
                    countY--
                }
            }

            if (table[q.x + altX][q.y + altY] == nil) { return true; }
            return !(table[q.x + altX][q.y + altY].IsFigureBlack() == q.isBlack)
        }
    }

    return false
}

func (q *Queen) Pass(altX int, altY int) error {

    if !q.CanFigurePass(altX, altY) {
        return CantPassError
    }

    if q.isBlack {
        q.x -= altX
        q.y -= altY

        (*q.cTable)[q.x][q.y] = (*q.cTable)[q.x + altX][q.y + altY]
        (*q.cTable)[q.x + altX][q.y + altY] = nil
    } else {
        q.x += altX
        q.y += altY

        (*q.cTable)[q.x][q.y] = (*q.cTable)[q.x - altX][q.y - altY]
        (*q.cTable)[q.x - altX][q.y - altY] = nil
    }

    return nil
}

func (q *Queen) IsFigureBlack() bool {
    return q.isBlack
}

type King struct {
    cTable *ChessTable
    x int
    y int
    isBlack bool
}

func (k *King) CanFigurePass(altX int, altY int) bool {

    if k.isBlack {
        return k.canFigurePassInner(-altX, -altY)
    } else {
        return k.canFigurePassInner(altX, altY)
    }
}

func (k *King) canFigurePassInner(altX int, altY int) bool {
    table := (*k.cTable)

    if altX == 0 && altY == 0 { return false; }
    if (altX == 0 || u.Abs(altX) == 1) && (altY == 0 || u.Abs(altY) == 1) && k.x + altX <= 7 && k.x + altX >= 0 && k.y + altY <= 7 && k.y + altY >= 0 {
        if table[k.x + altX][k.y + altY] != nil && table[k.x + altX][k.y + altY].IsFigureBlack() == k.isBlack {
            return false
        }
       return true
    }

    return false
}

func (k *King) Pass(altX int, altY int) error {

    if !k.CanFigurePass(altX, altY) {
        return CantPassError
    }

    if k.isBlack {
        k.x -= altX
        k.y -= altY

        (*k.cTable)[k.x][k.y] = (*k.cTable)[k.x + altX][k.y + altY]
        (*k.cTable)[k.x + altX][k.y + altY] = nil
    } else {
        k.x += altX
        k.y += altY

        (*k.cTable)[k.x][k.y] = (*k.cTable)[k.x - altX][k.y - altY]
        (*k.cTable)[k.x - altX][k.y - altY] = nil
    }

    return nil
}

func (k *King) IsFigureBlack() bool {
    return k.isBlack
}
