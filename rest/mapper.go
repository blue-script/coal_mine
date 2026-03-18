package rest

import (
	"github.com/blue-script/coal_mine/enterprise"
)

func toMinerInfoDTO(info enterprise.MinerInfo) MinerInfoDTO {
	return MinerInfoDTO{
		Class:              info.Class,
		HireCost:           info.HireCost,
		CoalTotal:          info.CoalTotal,
		CurrentCoalPerMine: info.CurrentCoalPerMine,
		Working:            info.Working,
		Energy:             info.Energy,
		RemainingEnergy:    info.RemainingEnergy,
		Delay:              int(info.Delay.Seconds()),
	}
}

func toMinerInfoDTOs(miners []enterprise.MinerInfo) []MinerInfoDTO {
	dtos := make([]MinerInfoDTO, 0, len(miners))
	for _, m := range miners {
		dtos = append(dtos, toMinerInfoDTO(m))
	}
	return dtos
}

func toEquipmentDTO(e enterprise.Equipment) EquipmentDTO {
	return EquipmentDTO{
		Pickaxe: EquipmentItemDTO{
			Name:  e.Pickaxe.Name,
			Cost:  e.Pickaxe.Cost,
			Count: e.Pickaxe.Count,
		},
		MineVentilation: EquipmentItemDTO{
			Name:  e.MineVentilation.Name,
			Cost:  e.MineVentilation.Cost,
			Count: e.MineVentilation.Count,
		},
		OreCart: EquipmentItemDTO{
			Name:  e.OreCart.Name,
			Cost:  e.OreCart.Cost,
			Count: e.OreCart.Count,
		},
	}
}

func toEnterpriseStatisticDTO(e enterprise.EnterpriseStats) EnterpriseStatisticDTO {
	miners := make([]enterprise.MinerInfo, 0, len(e.Miners))
	for _, val := range e.Miners {
		miners = append(miners, val)
	}

	return EnterpriseStatisticDTO{
		Balance:   e.Balance,
		Hired:     e.Hired,
		Miners:    toMinerInfoDTOs(miners),
		Equipment: toEquipmentDTO(e.Equipment),
		Duration:  int(e.Duration.Seconds()),
	}
}
