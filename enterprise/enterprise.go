package enterprise

import (
	"context"
	"maps"
	"sync"
	"time"
)

type Enterprise struct {
	mtx       sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	Equipment *Equipment
	Balance   Coal
	Miners    []Miner
	Hired     map[MinerClass]int
	Finished  bool
	StartedAt time.Time
}

func NewEnterprise(ctx context.Context, cancel context.CancelFunc) *Enterprise {
	equipment := NewEquipment()

	return &Enterprise{
		Equipment: equipment,
		ctx:       ctx,
		cancel:    cancel,
		StartedAt: time.Now(),
		Balance:   0,
		Miners:    make([]Miner, 0),
		Hired:     make(map[MinerClass]int),
	}
}

func (e *Enterprise) RunPassiveIncome() {
	e.wg.Add(1)
	go func() {
		defer e.wg.Done()

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-e.ctx.Done():
				return
			case <-ticker.C:
				e.mtx.Lock()
				e.Balance++
				e.mtx.Unlock()
			}
		}
	}()
}

func (e *Enterprise) newMinerByClass(class MinerClass) (Miner, error) {
	switch class {
	case MinerClassSmall:
		return NewSmallMiner(), nil
	case MinerClassNormal:
		return NewNormalMiner(), nil
	case MinerClassStrong:
		return NewStrongMiner(), nil
	default:
		return nil, ErrIncorrectMinerClass
	}
}

func (e *Enterprise) removeMiner(miner Miner) {
	e.mtx.Lock()
	defer e.mtx.Unlock()

	for i, value := range e.Miners {
		if value == miner {
			e.Miners = append(e.Miners[:i], e.Miners[i+1:]...)
			return
		}
	}
}

func (e *Enterprise) HireMiner(class MinerClass) (MinerInfo, error) {
	miner, err := e.newMinerByClass(class)
	if err != nil {
		return MinerInfo{}, err
	}
	info := miner.Info()

	e.mtx.Lock()
	if e.Finished {
		e.mtx.Unlock()
		return MinerInfo{}, ErrGameAlreadyFinished
	}

	if e.Balance < info.HireCost {
		e.mtx.Unlock()
		return MinerInfo{}, ErrNotEnoughMoney
	}

	e.Balance -= info.HireCost
	e.Miners = append(e.Miners, miner)
	e.Hired[class]++
	e.mtx.Unlock()

	e.wg.Add(1)
	ch := miner.Run(e.ctx)
	go func() {
		defer e.wg.Done()

		for coal := range ch {
			e.mtx.Lock()
			e.Balance += coal
			e.mtx.Unlock()
		}

		e.removeMiner(miner)
	}()

	return info, nil
}

func (e *Enterprise) BuyEquipment(name EquipmentName) error {
	e.mtx.Lock()
	defer e.mtx.Unlock()
	if e.Finished {
		return ErrGameAlreadyFinished
	}

	cost, err := e.Equipment.Cost(name)
	if err != nil {
		return err
	}

	if cost > e.Balance {
		return ErrNotEnoughMoney
	}

	err = e.Equipment.Buy(name)
	if err != nil {
		return err
	}
	e.Balance -= cost

	return nil
}

func (e *Enterprise) EquipmentCosts() map[EquipmentName]Coal {
	e.mtx.RLock()
	defer e.mtx.RUnlock()

	return e.Equipment.Costs()
}

func (e *Enterprise) ListEquipment() Equipment {
	e.mtx.RLock()
	defer e.mtx.RUnlock()

	return e.Equipment.Info()
}

type EnterpriseStats struct {
	Balance   Coal
	Hired     map[MinerClass]int
	Miners    []MinerInfo
	Equipment Equipment
	Duration  time.Duration
}

func (e *Enterprise) Statistic() EnterpriseStats {
	e.mtx.RLock()
	defer e.mtx.RUnlock()

	finishedAt := time.Now()

	hiredCopy := make(map[MinerClass]int, len(e.Hired))
	maps.Copy(hiredCopy, e.Hired)

	miners := make([]MinerInfo, 0, len(e.Miners))
	for _, val := range e.Miners {
		miners = append(miners, val.Info())
	}

	return EnterpriseStats{
		Balance:   e.Balance,
		Hired:     hiredCopy,
		Miners:    miners,
		Equipment: e.Equipment.Info(),
		Duration:  finishedAt.Sub(e.StartedAt),
	}
}

func (e *Enterprise) FinishGame() (EnterpriseStats, error) {
	e.mtx.Lock()
	if !e.Equipment.Completed() {
		e.mtx.Unlock()
		return EnterpriseStats{}, ErrGameNotCompleted
	}

	if e.Finished {
		e.mtx.Unlock()
		return EnterpriseStats{}, ErrGameAlreadyFinished
	}
	e.Finished = true
	finishedAt := time.Now()
	e.mtx.Unlock()

	e.cancel()
	e.wg.Wait()

	e.mtx.RLock()
	defer e.mtx.RUnlock()

	hiredCopy := make(map[MinerClass]int, len(e.Hired))
	maps.Copy(hiredCopy, e.Hired)

	return EnterpriseStats{
		Balance:   e.Balance,
		Hired:     hiredCopy,
		Miners:    []MinerInfo{},
		Equipment: e.Equipment.Info(),
		Duration:  finishedAt.Sub(e.StartedAt),
	}, nil
}

func (e *Enterprise) ListAllMiners() []MinerInfo {
	return e.ListMiners(nil)
}

func (e *Enterprise) ListMiners(class *MinerClass) []MinerInfo {
	e.mtx.RLock()
	defer e.mtx.RUnlock()

	miners := make([]MinerInfo, 0, len(e.Miners))
	for _, v := range e.Miners {
		info := v.Info()
		if class == nil || info.Class == *class {
			miners = append(miners, info)
		}
	}

	return miners
}

func (e *Enterprise) HireCost(class MinerClass) (MinerInfo, error) {
	miner, err := e.newMinerByClass(class)
	if err != nil {
		return MinerInfo{}, err
	}

	return miner.Info(), nil
}
