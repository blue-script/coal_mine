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
	mtx                  sync.RWMutex
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

	m.mtx.Lock()
	defer m.mtx.Unlock()
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
	m.mtx.RLock()
	defer m.mtx.RUnlock()

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

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				coal, ok := m.preview()
				if !ok {
					return
				}

				select {
				case <-ctx.Done():
					return
				case ch <- coal:
					m.commit(coal)
				}
			}
		}
	}()

	return ch
}

func (m *baseMiner) Info() MinerInfo {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

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

type MinerConfig struct {
	Class                MinerClass
	HireCost             int
	Energy               int
	CoalPerMine          Coal
	Delay                time.Duration
	PerformanceIncrement Coal
}

func newBaseMiner(cfg MinerConfig) *baseMiner {
	bm := &baseMiner{
		Class:                cfg.Class,
		HireCost:             cfg.HireCost,
		Energy:               cfg.Energy,
		CoalPerMine:          cfg.CoalPerMine,
		Delay:                cfg.Delay,
		PerformanceIncrement: cfg.PerformanceIncrement,
	}
	bm.RemainingEnergy = bm.Energy
	bm.CurrentCoalPerMine = bm.CoalPerMine

	bm.preview = bm.defaultPreview
	bm.commit = bm.defaultCommit

	return bm
}

func NewSmallMiner() *SmallMiner {
	return &SmallMiner{newBaseMiner(MinerConfig{
		Class:                MinerClassSmall,
		HireCost:             5,
		Energy:               30,
		CoalPerMine:          1,
		Delay:                3 * time.Second,
	})}
}

func NewNormalMiner() *NormalMiner {
	return &NormalMiner{newBaseMiner(MinerConfig{
		Class:                MinerClassNormal,
		HireCost:             50,
		Energy:               45,
		CoalPerMine:          3,
		Delay:                2 * time.Second,
	})}
}

func NewStrongMiner() *StrongMiner {
	sm := &StrongMiner{newBaseMiner(MinerConfig{
		Class:                MinerClassStrong,
		HireCost:             450,
		Energy:               60,
		CoalPerMine:          10,
		Delay:                1 * time.Second,
		PerformanceIncrement: 3,
	})}
	sm.commit = sm.strongCommit

	return sm
}
