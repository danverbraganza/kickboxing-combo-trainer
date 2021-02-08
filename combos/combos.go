package combos

import (
	"github.com/danverbraganza/go-mithril/moria"
	"github.com/gopherjs/gopherjs/js"
)

var m = moria.M

type Move struct {
	ShortName, LongName string
}

var Moves = []Move{
	{"1", "Jab"},
	{"2", "Cross"},
	{"3", "Lead Hook"},
	{"4", "Rear Hook"},
	{"5", "Lead Uppercut"},
	{"6", "Rear Uppercut"},
	{"7", "Lead Hook Body"},
	{"8", "Read Hook Body"},
	{"9", "Jab Body"},
	{"10", "Cross Body"},
}

func FromNames(names ...string) (moves []Move) {
	for _, name := range names {
		// Jab is the default if this doesn't match.
		pickedMove := Moves[0]
		for _, move := range Moves {
			if move.ShortName == name {
				pickedMove = move
				break
			}
		}
		moves = append(moves, pickedMove)
	}
	return
}

type Combo struct {
	Name  string
	Moves []Move
}

func (c Combo) NewCheckBox() (retval moria.VirtualElement) {
	return m("li", nil,
		m("div", nil,
			m("label[for='combo-"+c.Name+"']", nil, moria.S(c.Name)),
			m("input#combo-"+c.Name, js.M{"type": "checkbox"}),
		),
		moria.F(func(children *[]moria.View) {
			for i, move := range c.Moves {
				*children = append(*children, moria.S(move.LongName))
				if i < len(c.Moves) - 1 {
					*children = append(*children, moria.S(", "))
}
			}
			return
		}),
	)
}

var List = []Combo{
	{"1", FromNames("1")},
	{"2", FromNames("1", "2")},
	{"3", FromNames("1", "2", "3")},
	{"4", FromNames("1", "2", "3", "2")},
	{"5", FromNames("1", "2", "5", "2", "3")},
}
