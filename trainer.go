package main

import (
	"github.com/danverbraganza/go-mithril/moria"
	"honnef.co/go/js/dom"

	"github.com/danverbraganza/kickboxing-combo-trainer/round"
	"github.com/danverbraganza/kickboxing-combo-trainer/welcome"
)

var m = moria.M

func main() {
	welcomePage := welcome.NewWelcomePage()
	activeRound := round.NewRound()

	moria.Route(
		dom.GetWindow().Document().QuerySelector("div#trainer"), "/",
		map[string]moria.Component{
			"/":                welcomePage,
			"/round/:duration/:speed/selectedCombos=:selectedCombos": activeRound,
		})

}
