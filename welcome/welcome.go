package welcome

import (
	"strings"
	"time"

	"github.com/danverbraganza/go-mithril"
	"github.com/danverbraganza/go-mithril/moria"
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/dom"

	"github.com/danverbraganza/kickboxing-combo-trainer/combos"
)

var m = moria.M

type WelcomePage struct {
	combos        map[string]bool
	selectedRound string
	Duration      time.Duration
	Speed         string
}

func NewWelcomePage() *WelcomePage {
	return &WelcomePage{
		combos:        map[string]bool{},
		selectedRound: "",
		Duration:      60 * time.Second,
		Speed:         "slow",
	}
}

func (w *WelcomePage) Controller() moria.Controller {
	return w
}

func (w *WelcomePage) RoundCombosAsString() string {
	if w.selectedRound == "custom" || w.selectedRound == "" {
		// If there's no selectedRound, this is custom
		selectedCombos := []string{}
		for combo, selected := range w.combos {
			if selected {
				selectedCombos = append(selectedCombos, combo)
			}

		}
		return strings.Join(selectedCombos, ",")
	} else {
		round := combos.RoundByName(w.selectedRound)
		return round.CombosAsString()
	}
}

func (*WelcomePage) View(x moria.Controller) moria.View {
	w := x.(*WelcomePage)

	return m("div#wrapper", nil,
		m("div#subtitle", nil, moria.S("Drill rounds of randomized combos to improve your timing and movement. Rounds are categorised by the types of moves they include so you can customize the level of complexity, intensity and duration.")),
		// Add more components here
		m("div", nil,
			m("div#options", nil,
				m("select#select-duration",
					js.M{
						"onchange": mithril.WithAttr("value", func(value string) {
							w.Duration, _ = time.ParseDuration(value)
						}),
					},
					moria.F(func(children *[]moria.View) {
						for _, duration := range []string{"30s", "1m", "2m", "3m"} {
							*children = append(*children, m("option[value='"+duration+"']", js.M{"selected": strings.HasPrefix(w.Duration.String(), duration)}, moria.S(duration)))
						}
					})),
				m("select#select-speed",
					js.M{
						"onchange": mithril.WithAttr("value", func(value string) {
							w.Speed = value
						}),
					},
					moria.F(func(children *[]moria.View) {
						for _, speed := range []string{"slow", "medium", "fast"} {
							*children = append(*children, m("option[value='"+speed+"']",
								js.M{"selected": speed == w.Speed},
								moria.S(speed)))
						}
					}),
				),
				m("button#go", js.M{
					"config":   mithril.RouteConfig,
					"disabled": w.selectedRound == "",
					"onclick": func() {
						mithril.RouteRedirect(
							strings.Join([]string{
								"",
								"round",
								w.Duration.String(),
								w.Speed,
							},
								"/")+"/selectedCombos="+w.RoundCombosAsString(),

							js.M{},
							false,
						)
					}},
					moria.S("Begin Round"),
				),
			),
		),

		m("div.round-container", nil,
			// TODO: Make collapsible
			m("div.container-title", nil, moria.S("Rounds")),
			moria.F(func(children *[]moria.View) {
				for _, round := range combos.RoundList {
					*children = append(*children, round.NewRadioButton(&w.selectedRound))
				}
			},
			),

			m("div.round-picker", js.M{
				"onclick": func() {
					d := dom.GetWindow().Document()
					d.GetElementByID("round-custom").(*dom.HTMLInputElement).Click()
				}},
				m("label[for='round-custom']", nil, moria.S("Customize your own Round")),
				m("input#round-custom[type='radio'][name='round']", js.M{
					"checked": w.selectedRound == "custom",
					"onchange": func() {
						w.selectedRound = "custom"
					}}),
				m("div.combo-container", nil,

					moria.F(func(children *[]moria.View) {
						for _, combo := range combos.List {
							*children = append(*children, combo.NewCheckBox(w.combos))
						}
					},
					)),
			),
		),
		m("div#disclaimer", nil,
			m("input#disclaimer-collapse.toggle[type='checkbox']", nil),
			m("label#label-disclaimer[for='disclaimer-collapse']", nil, moria.S("Disclaimer")),
			m("div.disclaimer-content", nil,
				moria.S("Exercising is good for you! However, every individual is unique. By continuing to use this application, you recognize that you are taking full responsibility for the consequences. Make sure to check with your doctor before using this app if you need to. You agree that this is app is not responsible for any injuries.")),
		),
		// End disclaimer
)
}
