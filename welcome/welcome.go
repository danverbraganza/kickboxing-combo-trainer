package welcome

import (
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
		combos:   map[string]bool{},
		Duration: 180 * time.Second,
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
		// Add more components here
		m("div#disclaimer", nil, moria.S("Exercising is good for you! However, every individual is unique. By continuing to use this application, you recognize that you are taking full responsibility for the consequences. Make sure to check with your doctor before using this app if you need to. You agree that this is app is not responsible for any injuries.")),
		m("div", nil,
			m("select#select-duration", nil,
				moria.F(func(children *[]moria.View) {
					for _, duration := range []string{"1m", "2m", "3m"} {
						*children = append(*children, m("option[value='"+duration+"']", nil, moria.S(duration)))
					}
				},
				),
			),
			m("select#select-speed", nil,
				moria.F(func(children *[]moria.View) {
					for _, speed := range []string{"slow", "medium", "fast"} {
						*children = append(*children, m("option[value='"+speed+"']", nil, moria.S(speed)))
					}
				},
				),
			),
			m("button", js.M{
				"config": mithril.RouteConfig,
				"onclick": func() {
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

		m("div.combo-container", nil,
			moria.F(func(children *[]moria.View) {
				for _, combo := range combos.List {
					*children = append(*children, combo.NewCheckBox(w.combos))
				}
			},
			),
		),
	)

}
