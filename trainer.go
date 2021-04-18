package main

import (
	"github.com/danverbraganza/go-mithril/moria"
	"honnef.co/go/js/dom"

	"kickboxing-combo-trainer/round"
	"kickboxing-combo-trainer/welcome"
)

var m = moria.M

func main() {
	welcomePage := welcome.NewWelcomePage()
	activeRound := round.NewRound()

	moria.Route(
		dom.GetWindow().Document().QuerySelector("div#trainer"), "/",
		map[string]moria.Component{
			"/":                welcomePage,
			"/round/:duration/selectedCombos=:selectedCombos": activeRound,
		})

}
