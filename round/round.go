package round

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/danverbraganza/go-mithril"
	"github.com/danverbraganza/go-mithril/moria"
	"github.com/gopherjs/gopherjs/js"

	"kickboxing-combo-trainer/combos"
)

func Init() {
	rand.Seed(time.Now().Unix())
}

type Round struct {
	sync.Mutex
	Duration, timeSpent time.Duration
	last                time.Time
	running             bool
	cleared             bool
	SelectedCombos      []combos.Combo
}

var m = moria.M

type s = moria.S

var (
	fps30           = time.Tick(time.Second / 30)
	DisplayChan     = make(chan moria.VirtualElement)
	beatTick        = time.Tick(800 * time.Millisecond)
	runningBeatTick = make(chan time.Time)
)

func ExtractCombos(rawComboString string) (sc []combos.Combo) {
	comboKeys := strings.Split(rawComboString, ",")
	for _, comboKey := range comboKeys {
		sc = append(sc, combos.ByName(comboKey))
	}
	return
}

func (r *Round) Controller() moria.Controller {
	r.Duration, _ = time.ParseDuration(
		mithril.RouteParam("duration").(string),
	)

	r.SelectedCombos = ExtractCombos(mithril.RouteParam("selectedCombos").(string))

	go func() {
		for !r.cleared {
			// Pick a combo
			comboTimer := r.RandomCombo().NewChannel(runningBeatTick)
			for innerElement := range comboTimer {
				DisplayChan <- innerElement
			}
		}
	}()

	r.timeSpent = 0 * time.Second
	r.Start()
	return r
}

func (r *Round) RandomCombo() combos.Combo {
	return r.SelectedCombos[rand.Intn(len(r.SelectedCombos))]
}

func (r *Round) Start() {
	r.Lock()
	defer r.Unlock()
	r.last = time.Now()
	r.running = true

	go func() {
		for r.running && r.timeSpent <= r.Duration {
			<-fps30
			now := time.Now()
			r.timeSpent += now.Sub(r.last)
			r.last = now
			mithril.Redraw(false)
			select {
			case x := <-beatTick:
				runningBeatTick <- x
			default:
			}
		}
		// Reroute back to main page
	}()
}

func (r *Round) Stop() {
	r.Lock()
	defer r.Unlock()
	r.running = false
	r.cleared = true
}

func FormatDuration(d time.Duration) string {
	return fmt.Sprintf("%02d:%02d",
		int(math.Abs(d.Minutes()))%60,
		int(math.Abs(d.Seconds()))%60,
	)
}

func (*Round) View(x moria.Controller) moria.View {
	r := x.(*Round)

	// Display pause button
	pauseSigil := s("\u23F8")
	if !r.running {
		// Display play button
		pauseSigil = s("\u25B6")
	}

	return m("div", nil,
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
					go r.Stop()
				} else {
					go r.Start()
				}
			},
		}, pauseSigil),
		m("button#stop.control", js.M{
			"config": mithril.RouteConfig,
			"onclick": func() {
				go r.Stop()
				mithril.RouteRedirect(
					"/",
					js.M{},
					false,
				)
			},
		},
			s("\u25a0")),
		m("div#move", nil, <-DisplayChan),
	)
}
