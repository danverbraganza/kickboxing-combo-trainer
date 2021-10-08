package combos

import (
	"strings"
	"time"

	"github.com/danverbraganza/go-mithril"
	"github.com/danverbraganza/go-mithril/moria"
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/dom"
)

var m = moria.M

type Move struct {
	ShortName, LongName string
	IsLeadSide          bool
}

var Moves = []Move{
	{"1", "Jab", true},
	{"2", "Cross", false},
	{"3", "Hook", true},
	{"4", "Rear Hook", false},
	{"5", "Lead Uppercut", true},
	{"6", "Rear Uppercut", false},
	{"1B", "Jab Body", true},
	{"2B", "Cross Body", false},
	{"3B", "Hook Body", true},
	{"4B", "Rear Hook Body", false},
	{"LK", "Lead Kick", true},
	{"RK", "Rear leg Kick", false},
	{"V", "Liver Shot", true},
	{"S", "Spleen Shot", false},
}

// Returns true if this move ends with a kick.
func (m Move) IsKick() bool {
	return strings.HasSuffix(m.ShortName, "K")
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
	return m("div.combo-picker", js.M{
		"onclick": func() {
			d := dom.GetWindow().Document()
			d.GetElementByID("combo-" + c.Name).(*dom.HTMLInputElement).Click()
		},
	},
		m("label[for='combo-"+c.Name+"']", nil,
			moria.S(c.Name)),
		m("input#combo-"+c.Name+"[type='checkbox']", js.M{
			"checked": selectedCombos[c.Name],
			"onchange": mithril.WithAttr("checked", func(checked bool) {
				selectedCombos[c.Name] = checked
			}),
			"onclick": func(event *js.Object) {
				// Don't propagate to the container, because that will just call click on us again.
				// TODO: js.Object and Call() are ugly.
				event.Call("stopPropagation")
			}},
		),
		m("br", nil),
		moria.F(func(children *[]moria.View) {
			for i, move := range c.Moves {
				*children = append(*children, moria.S(move.LongName))
				if i < len(c.Moves)-1 {
					*children = append(*children, m("br", nil))
				}
			}
			return
		}))
}

// Creates a new combo by joining
func Join(combos ...Combo) (retval Combo) {
	names := []string{}
	for _, combo := range combos {
		names = append(names, combo.Name)
		for _, move := range combo.Moves {
			retval.Moves = append(retval.Moves, move)
		}
	}
	retval.Name = strings.Join(names, "+")
	return
}

func (c Combo) Describe() (retval moria.VirtualElement) {
	return m(
		"div.combo-describe", nil,
		m("span.combo-name", nil, moria.S(c.Name)),
		m("span.combo-description", nil,
			moria.F(func(children *[]moria.View) {
				moveNames := []string{}
				for _, move := range c.Moves {
					moveNames = append(moveNames, move.LongName)
				}
				*children = append(*children, moria.S(strings.Join(moveNames, " | ")))
				return
			})))
}

type DisplayElement struct {
	Type   string
	Move   string
	Loaded bool
}

// Returns a channel that you can watch to get the current state
func (c Combo) NewChannel(beatTick chan time.Time) (displayChan chan moria.VirtualElement) {
	current := DisplayElement{Type: "combo", Move: c.Name}
	// Pause on the move intro for twice the beats.
	moveIndex := -2
	displayChan = make(chan moria.VirtualElement)
	go func() {
		for {
			select {
			case <-beatTick:
				print("Read on beatTick")
				moveIndex++
				if moveIndex < 0 {
					// Do nothing
				} else if moveIndex < len(c.Moves) {
					current = DisplayElement{Type: "move", Move: c.Moves[moveIndex].LongName}
				} else {
					print("closing displayChan")
					close(displayChan)
					return
				}
			default:
				if !current.Loaded {

					current.Loaded = true
					displayChan <- m("div#combo", nil, moria.S("-"))
				} else {
					displayChan <- m("div#"+current.Type, nil, moria.S(current.Move))
				}
			}
		}
	}()
	return displayChan
}

// WithKick finishes the combo off with the opposite side kick
func (c Combo) WithKick() Combo {
	lastMove := c.Moves[len(c.Moves)-1]
	if lastMove.IsKick() {
		return c
	}

	KickToThrow := Combo{"RK", FromNames("LK")}
	if lastMove.IsLeadSide {
		KickToThrow = Combo{"RK", FromNames("RK")}
	}
	return Join(c, KickToThrow)
}

// WithExtender finishes the combo off with X or Y.
func (c Combo) WithExtender() Combo {
	var extender Combo
	lastMove := c.Moves[len(c.Moves)-1]
	if !lastMove.IsKick() {
		if lastMove.IsLeadSide {
			extender = ByName("Y")
		} else {
			extender = ByName("X")
		}
	} else {
		if lastMove.IsLeadSide {
			extender = Combo{"Y", FromNames("2", "3", "RK")}
		} else {
			extender = Combo{"Y", FromNames("RK", "3", "2", "LK")}
		}
	}

	return Join(c, extender)
}

var List = []Combo{
	{"1", FromNames("1")},
	{"1-1", FromNames("1", "1")},
	{"2", FromNames("1", "2")},
	{"3", FromNames("1", "2", "3")},
	{"4", FromNames("1", "2", "3", "2")},
	{"5", FromNames("1", "2", "5", "2", "3")},
	{"3U", FromNames("1", "2", "3", "6")},
	{"4U", FromNames("1", "2", "3", "4", "5")},
	{"HCH", FromNames("3", "2", "3")},
	{"CHC", FromNames("2", "3", "2")},
	{"Liver go-around", FromNames("V", "4B", "3B", "S")},
	{"Spleen go-around", FromNames("S", "3B", "4B", "V")},
	{"A", FromNames("1B", "2", "3")},
	{"B", FromNames("2B", "3", "2")},
	{"X", FromNames("LK", "2", "3", "RK")},
	{"Y", FromNames("RK", "3", "2", "LK")},
}

type Round struct {
	Name   string
	Combos []Combo
}

func (r Round) NewRadioButton(selectedRound *string) (retval moria.VirtualElement) {
	return m(
		"div.round-picker",
		js.M{
			"onclick": func() {
				d := dom.GetWindow().Document()
				d.GetElementByID("round-" + r.Name).(*dom.HTMLInputElement).Click()
			}},
		m("label[for='round-"+r.Name+"']", nil, moria.S(r.Name)),
		m("input#round-"+r.Name+"[type='radio'][name='round']", js.M{
			"checked": *selectedRound == r.Name,
			"onchange": func() {
				*selectedRound = r.Name
			}},
		),
		moria.F(func(children *[]moria.View) {
			for _, combo := range r.Combos {
				*children = append(*children, combo.Describe())
			}
			return
		}))
}

// Returns a stringified list of the combos.
func (r Round) CombosAsString() string {
	combos := []string{}
	for _, combo := range r.Combos {
		combos = append(combos, combo.Name)
	}
	return strings.Join(combos, ",")

}

var RoundList = []Round{

	{"Getting Started", []Combo{ByName("1"), ByName("1-1"), ByName("2")}},
	{"Adding the hook", []Combo{ByName("1"), ByName("1-1"), ByName("2"), ByName("3")}},
	{"Adding the hook II", []Combo{ByName("1"), ByName("1-1"), ByName("2"), ByName("3"), ByName("HCH"), ByName("CHC")}},
	{"Adding the hook III", []Combo{ByName("1"), ByName("1-1"), ByName("2"), ByName("3"), ByName("HCH"), ByName("CHC"), Join(ByName("1"), ByName("CHC")), Join(ByName("2"), ByName("HCH")), Join(ByName("3"), ByName("CHC"))}},
	{"Uppercuts", []Combo{ByName("1"), ByName("1-1"), ByName("2"), ByName("3"), ByName("4"), ByName("3U"), ByName("4U")}},
	{"Intro to Body Shots", []Combo{ByName("1"), ByName("2"), ByName("A"), ByName("B"), ByName("Liver go-around"), ByName("Spleen go-around")}},

	{"Max hands", []Combo{ByName("1"), ByName("1-1"), ByName("2"), ByName("3"), ByName("HCH"), ByName("CHC"), Join(ByName("1"), ByName("CHC")), Join(ByName("2"), ByName("HCH")), Join(ByName("3"), ByName("CHC")), ByName("4"), ByName("Liver go-around"), ByName("Spleen go-around")}},

	{"Intro to Kicking", []Combo{ByName("1").WithKick(), ByName("1-1").WithKick(), ByName("2").WithKick()}},
	{"Multi-kicking combos", []Combo{ByName("1").WithKick(), ByName("1-1").WithKick(), ByName("2").WithKick(), ByName("3").WithKick(), ByName("4").WithKick(), ByName("5").WithKick(), ByName("X"), ByName("Y")}},
	{"Hooks and Kicks", []Combo{ByName("1").WithKick(), ByName("1-1").WithKick(), ByName("2").WithKick(), ByName("3").WithKick(), ByName("HCH").WithKick(), ByName("CHC").WithKick()}},
	{"Uppercuts and Kicks", []Combo{ByName("1").WithKick(), ByName("1-1").WithKick(), ByName("2").WithKick(), ByName("3").WithKick(), ByName("4").WithKick(), ByName("3U").WithKick(), ByName("4U").WithKick()}},

	{"Extended Combos", []Combo{ByName("1").WithExtender(), ByName("1-1").WithExtender(), ByName("2").WithExtender(), ByName("3").WithExtender(), ByName("HCH").WithExtender(), ByName("CHC").WithExtender(), Join(ByName("1"), ByName("CHC")).WithExtender(), Join(ByName("2"), ByName("HCH")).WithExtender(), Join(ByName("3"), ByName("CHC")).WithExtender(), ByName("4").WithExtender(), ByName("Liver go-around").WithExtender(), ByName("Spleen go-around").WithExtender()}},
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

// TODO: Dictionary lookup
func RoundByName(name string) Round {
	for _, round := range RoundList {
		if round.Name == name {
			return round
		}
	}
	return RoundList[0]
}
