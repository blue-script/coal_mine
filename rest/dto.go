package rest

import (
	"time"

	"github.com/blue-script/coal_mine/enterprise"
)

type MinerClassDTO struct {
	Class enterprise.MinerClass `json:"class"`
}

type EquipmentNameDTO struct {
	Name enterprise.EquipmentName `json:"name"`
}

type MinerInfoDTO struct {
	Class              enterprise.MinerClass `json:"class"`
	HireCost           enterprise.Coal       `json:"hire_cost"`
	CoalTotal          enterprise.Coal       `json:"coal_total"`
	CurrentCoalPerMine enterprise.Coal       `json:"current_coal_per_mine"`
	Working            bool                  `json:"working"`
	Energy             int                   `json:"energy"`
	RemainingEnergy    int                   `json:"remaining_energy"`
	Delay              int                   `json:"delay"`
}

type EquipmentPriceDTO struct {
	Name enterprise.EquipmentName `json:"name"`
	Cost enterprise.Coal          `json:"cost"`
}

type EquipmentItemDTO struct {
	Name  enterprise.EquipmentName `json:"name"`
	Cost  enterprise.Coal          `json:"cost"`
	Count int                      `json:"count"`
}

type EquipmentDTO struct {
	Pickaxe         EquipmentItemDTO `json:"pickaxe"`
	MineVentilation EquipmentItemDTO `json:"mine_ventilation"`
	OreCart         EquipmentItemDTO `json:"ore_cart"`
}

type EnterpriseDTO struct {
	Balance   enterprise.Coal               `json:"balance"`
	Hired     map[enterprise.MinerClass]int `json:"hired"`
	Miners    []MinerInfoDTO                `json:"miners"`
	Finished  bool                          `json:"finished"`
	Duration  time.Duration                 `json:"duration"`
	Equipment EquipmentDTO                  `json:"equipment"`
}

type EnterpriseStatisticDTO struct {
	Balance   enterprise.Coal               `json:"balance"`
	Hired     map[enterprise.MinerClass]int `json:"hired"`
	Miners    []MinerInfoDTO                `json:"miners"`
	Equipment EquipmentDTO                  `json:"equipment"`
	Duration  int                           `json:"duration"`
}
