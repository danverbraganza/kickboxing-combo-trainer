package round

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/danverbraganza/go-mithril"
	"github.com/danverbraganza/go-mithril/moria"
	"github.com/gopherjs/gopherjs/js"

	"kickboxing-combo-trainer/combos"
)

type Round struct {
	sync.Mutex
	Duration, timeSpent time.Duration
	last                time.Time
	running             bool
	SelectedCombos      []combos.Combo
}

var m = moria.M

type s = moria.S

var (
	fps30 = time.Tick(time.Second / 30)
)

func (r *Round) Controller() moria.Controller {
	r.Duration, _ = time.ParseDuration(
		mithril.RouteParam("duration").(string),
	)

	r.timeSpent = 0 * time.Second
	r.Start()
	return r
}

func (r *Round) Start() {
	r.Lock()
	defer r.Unlock()
	r.last = time.Now()
	r.running = true

	go func() {
		for r.running {
			<-fps30
			now := time.Now()
			r.timeSpent += now.Sub(r.last)
			r.last = now
			mithril.Redraw(false)
		}
		// Reroute back to main page
	}()
}

func (r *Round) Stop() {
	r.Lock()
	defer r.Unlock()
	r.running = false
}

func FormatDuration(d time.Duration) string {
	return fmt.Sprintf("%02d:%02d:%03d",
		int(math.Abs(d.Minutes()))%60,
		int(math.Abs(d.Seconds()))%60,
		int(math.Abs(float64(d.Nanoseconds()/1e6)))%1000)
}

func (*Round) View(x moria.Controller) moria.View {
	r := x.(*Round)

	pauseSigil := s("\u23F8")
	if !r.running {
		pauseSigil = s("\u25B6")
	}

	return m("div#wrapper", nil,
		m("h1", nil, moria.S("Kickboxing Combo Trainer")),
		m("input#time-left", js.M{
			"value": FormatDuration(r.Duration - r.timeSpent),
		}),
		moria.F(func(children *[]moria.View) {
			// When a combo is introduced, we want two beats of
			// intro, and then one beat for each move.
		}),
		m("button#pause.control", js.M{
			"config": mithril.RouteConfig,
			"onclick": func() {
				if r.running {
					r.Stop()
				} else {
					r.Start()
				}
			},
		}, pauseSigil),
		m("button#stop.control", js.M{
			"config": mithril.RouteConfig,
			"onclick": func() {
				r.Stop()
				mithril.RouteRedirect(
					"/",
					js.M{},
					false,
				)
			},
		},
			s("\u25a0")),
	)
}
