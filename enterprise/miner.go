package enterprise

import (
	"context"
	"sync"
	"time"
)

type Coal int

type MinerClass string

const (
	MinerClassSmall  MinerClass = "small"
	MinerClassNormal MinerClass = "normal"
	MinerClassStrong MinerClass = "strong"
)

type MinerInfo struct {
	Class              MinerClass
	HireCost           int
	CoalTotal          Coal
	CurrentCoalPerMine Coal
	Working            bool
	RemainingEnergy    int
	Delay              time.Duration
}

type Miner interface {
	Run(ctx context.Context) <-chan Coal
	Info() MinerInfo
}

type previewFunc func() (Coal, bool)
type commitFunc func(coal Coal)

type baseMiner struct {
	mtx                  sync.Mutex
	Class                MinerClass
	HireCost             int
	Energy               int
	CoalPerMine          Coal
	Delay                time.Duration
	PerformanceIncrement Coal

	RemainingEnergy    int
	CurrentCoalPerMine Coal
	CoalTotal          Coal
	Working            bool

	preview previewFunc
	commit  commitFunc
}

type SmallMiner struct {
	*baseMiner
}

type NormalMiner struct {
	*baseMiner
}

type StrongMiner struct {
	*baseMiner
}

func (m *baseMiner) startRun() (chan Coal, bool) {
	ch := make(chan Coal)

	defer m.mtx.Unlock()
	m.mtx.Lock()
	if m.Working {
		close(ch)
		return ch, false
	}
	m.Working = true
	return ch, true
}

func (m *baseMiner) finishRun(ch chan Coal) {
	m.mtx.Lock()
	m.Working = false
	m.mtx.Unlock()

	close(ch)
}

func (m *baseMiner) defaultPreview() (Coal, bool) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if m.RemainingEnergy == 0 {
		return 0, false
	}

	return m.CurrentCoalPerMine, true
}

func (m *baseMiner) defaultCommit(coal Coal) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.CoalTotal += coal
	m.RemainingEnergy--
}

func (m *baseMiner) strongCommit(coal Coal) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.CoalTotal += coal
	m.RemainingEnergy--
	m.CurrentCoalPerMine += m.PerformanceIncrement
}

func (m *baseMiner) Run(ctx context.Context) <-chan Coal {
	ch, ok := m.startRun()
	if !ok {
		return ch
	}

	go func() {
		defer m.finishRun(ch)

		ticker := time.NewTicker(m.Delay)
		defer ticker.Stop()

	Loop:
		for {
			select {
			case <-ctx.Done():
				break Loop
			case <-ticker.C:
				coal, ok := m.preview()
				if !ok {
					break Loop
				}

				select {
				case <-ctx.Done():
					break Loop
				case ch <- Coal(coal):
					m.commit(coal)
				}
			}
		}
	}()

	return ch
}

func (m *baseMiner) Info() MinerInfo {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	return MinerInfo{
		Class:              m.Class,
		HireCost:           m.HireCost,
		CoalTotal:          m.CoalTotal,
		CurrentCoalPerMine: m.CurrentCoalPerMine,
		Working:            m.Working,
		RemainingEnergy:    m.RemainingEnergy,
		Delay:              m.Delay,
	}
}

