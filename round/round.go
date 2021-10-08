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
	counter             int
	SelectedCombos      []combos.Combo
	Speed               string
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
		comboNames := strings.Split(comboKey, "-")
		individualCombos := []combos.Combo{}
		for _, comboName := range comboNames {
			individualCombos = append(individualCombos, combos.ByName(comboName))
		}
		sc = append(sc, combos.Join(individualCombos...))
	}
	return
}

func (r *Round) Controller() moria.Controller {
	r.Duration, _ = time.ParseDuration(
		mithril.RouteParam("duration").(string),
	)

	r.SelectedCombos = ExtractCombos(mithril.RouteParam("selectedCombos").(string))
	r.Speed = (mithril.RouteParam("speed").(string))

	if r.Speed == "slow" {
		beatTick = time.Tick(800 * time.Millisecond)
	} else if r.Speed == "medium" {
		beatTick = time.Tick(600 * time.Millisecond)
	} else {
		beatTick = time.Tick(400 * time.Millisecond)
	}

	r.timeSpent = 0 * time.Second
	r.Start()
	return r
}

func NewRound() *Round {
	r := Round{
		SelectedCombos: []combos.Combo{},
	}
	go func() {
		counter := r.counter
		// block here for the beattick to prevent calling randomCombo
		// with the wrong combo.
		<-runningBeatTick
		for {
			// Wait for a beat
			// Pick a combo
			comboTimer := r.RandomCombo().NewChannel(runningBeatTick)
			for innerElement := range comboTimer {
				print("Writing to display", counter)
				DisplayChan <- innerElement
			}
			print("next move")
		}
		print("R cleared")
	}()
	return &r
}

func (r *Round) RandomCombo() combos.Combo {
	if len(r.SelectedCombos) > 0 {
		return r.SelectedCombos[rand.Intn(len(r.SelectedCombos))]
	} else {
		return combos.ByName("1")
	}
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
		r.counter++
		mithril.RouteRedirect(
			"/",
			js.M{},
			false,
		)
	}()
}

func (r *Round) Pause() {
	r.Lock()
	defer r.Unlock()
	r.running = false
}

func (r *Round) Stop() {
	r.Lock()
	defer r.Unlock()
	r.running = false
	r.counter++
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

	var innerCard moria.VirtualElement
	select {
	case innerCard = <-DisplayChan:
	default:
		innerCard = m("div", nil)
	}

	return m("div", nil,
		m("input#time-left", js.M{
			"value": FormatDuration(r.Duration - r.timeSpent),
		}),
		m("button#pause.control", js.M{
			"config": mithril.RouteConfig,
			"onclick": func() {
				if r.running {
					go r.Pause()
				} else {
					go r.Start()
				}
			},
		}, pauseSigil),
		m("button#stop.control", js.M{
			"config": mithril.RouteConfig,
			"onclick": func() {
				print("Stopping")
				go r.Stop()
			},
		},
			s("\u25a0")),
		m("div#move", nil, innerCard),
	)
}
