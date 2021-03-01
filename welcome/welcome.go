package welcome

import (
	"fmt"
	"strings"
	"time"

	"github.com/danverbraganza/go-mithril"
	"github.com/danverbraganza/go-mithril/moria"
	"github.com/gopherjs/gopherjs/js"

	"kickboxing-combo-trainer/combos"
)

var m = moria.M

type WelcomePage struct {
	combos   map[string]bool
	Duration time.Duration
}

func NewWelcomePage() *WelcomePage {
	return &WelcomePage{
		combos: map[string]bool{},
	}
}

func (w *WelcomePage) Controller() moria.Controller {
	return w
}

func (w *WelcomePage) SelectedCombosAsString() string {
	selectedCombos := []string{}
	for combo, selected := range w.combos {
		if selected {
			selectedCombos = append(selectedCombos, combo)
		}

	}
	return strings.Join(selectedCombos, ",")
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
					*children = append(*children, combo.NewCheckBox(w.combos))
				}
			},
			),
		),
		m("div", nil,
			m("button", js.M{
				"config": mithril.RouteConfig,
				"onclick": func() {
					fmt.Println(
						strings.Join([]string{
							"",
							"round",
							w.Duration.String(),
						},
							"/") + "/selectedCombos=" + w.SelectedCombosAsString())

					mithril.RouteRedirect(
						strings.Join([]string{
							"",
							"round",
							w.Duration.String(),
						},
							"/")+"/selectedCombos="+w.SelectedCombosAsString(),

						js.M{},
						false,
					)
				}},
				moria.S("Go!"),
			),
		),
	)

}
