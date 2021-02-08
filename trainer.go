package main

import (
	"strings"
	"time"

	"github.com/danverbraganza/go-mithril/moria"
	"github.com/danverbraganza/go-mithril"
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/dom"

	"kickboxing-combo-trainer/combos"
	"kickboxing-combo-trainer/round"
)

var m = moria.M

type WelcomePage struct {
	combos []string
	Duration time.Duration
}

func (w *WelcomePage) Controller() moria.Controller {
	return w
}

func (*WelcomePage) View(x moria.Controller) moria.View {
	w := x.(*WelcomePage)
	return m("div#wrapper", nil,
		m("h1", nil, moria.S("Kickboxing Combo Trainer")),
		// Add more components here
		m("div#disclaimer", nil, moria.S("Exercising is good for you! However, every individual is unique. By continuing to use this application, you recognize that you are taking full responsibility for the consequences. Make sure to check with your doctor before using this app if you need to. You agree that you cannot hold the developer of this app responsible for any injuries.")),
		m("ul", nil,
		moria.F(func(children *[]moria.View) {
			for _, combo := range combos.List {
				*children = append(*children, combo.NewCheckBox())
			}
		},
		),
			),
		m("div", nil,
			m("button", js.M{
				"config": mithril.RouteConfig,
				"onclick": func() {
					mithril.RouteRedirect(
						strings.Join([]string{
							"",
							"round",
							w.Duration.String(),
						},
							"/"),
						js.M{},
						false,
					)
				}},
				moria.S("Go!"),
			),
		),
	)

}

func main() {
	welcomePage := &WelcomePage{}
	activeRound := &round.Round{}

	moria.Route(
		dom.GetWindow().Document().QuerySelector("body"), "/",
		map[string]moria.Component{
			"/":      welcomePage,
			"/round/:duration": activeRound,
		})

}
