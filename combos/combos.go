package combos

import (
	"time"

	"github.com/danverbraganza/go-mithril"
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

func (c Combo) NewCheckBox(selectedCombos map[string]bool) (retval moria.VirtualElement) {
	return m("li", nil,
		m("div", nil,
			m("label[for='combo-"+c.Name+"']", nil, moria.S(c.Name)),
			m("input#combo-"+c.Name+"[type='checkbox']", js.M{
				"checked": selectedCombos[c.Name],
				"onchange": mithril.WithAttr("checked", func(checked bool) {
					selectedCombos[c.Name] = checked
				})})),
		moria.F(func(children *[]moria.View) {
			for i, move := range c.Moves {
				*children = append(*children, moria.S(move.LongName))
				if i < len(c.Moves)-1 {
					*children = append(*children, moria.S(", "))
				}
			}
			return
		}),
	)
}

// Returns a channel that you can watch to get the current state
func (c Combo) NewChannel(beatTick chan time.Time) (retval chan moria.VirtualElement) {
	currentString := "Combo: " + c.Name
	// Pause on the move intro for twice the beats.
	moveIndex := -2
	retval = make(chan moria.VirtualElement)
	go func() {
		for {
			select {
			case <-beatTick:
				print("Read on beatTick")
				moveIndex++
				if moveIndex < 0 {
					// Do nothing
				} else if moveIndex < len(c.Moves) {
					currentString = "Move: " + c.Moves[moveIndex].LongName
				} else {
					close(retval)
					return
				}
			default:
				retval <- m("div#move", nil, moria.S(currentString))
			}
		}
	}()
	return retval
}

var List = []Combo{
	{"1", FromNames("1")},
	{"JabJab", FromNames("1", "1")},
	{"2", FromNames("1", "2")},
	{"3", FromNames("1", "2", "3")},
	{"4", FromNames("1", "2", "3", "2")},
	{"5", FromNames("1", "2", "5", "2", "3")},
}

// TODO: Dictionary lookup
func ByName(name string) Combo {
	for _, combo := range List {
		if combo.Name == name {
			return combo
		}
	}
	return List[0]
}
