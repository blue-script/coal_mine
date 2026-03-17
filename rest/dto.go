package rest

import (
	"github.com/blue-script/coal_mine/enterprise"
)

type HireMinerDTO struct {
	Class enterprise.MinerClass `json:"class"`
}

type EquipmentNameDTO struct {
	Name enterprise.EquipmentName `json:"name"`
}

type HireCostDTO struct {
	Class enterprise.MinerClass `json:"class"`
}

type MinerInfoDTO struct {
	Class              enterprise.MinerClass `json:"class"`
	HireCost           enterprise.Coal       `json:"hire_cost"`
	CoalTotal          enterprise.Coal       `json:"coal_total"`
	CurrentCoalPerMine enterprise.Coal       `json:"current_coal_per_mine"`
	Working            bool                  `json:"working"`
	RemainingEnergy    int                   `json:"remaining_energy"`
	Delay              int                   `json:"delay"`
}
