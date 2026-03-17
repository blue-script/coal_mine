package rest

import "github.com/blue-script/coal_mine/enterprise"

func toMinerInfoDTO(info enterprise.MinerInfo) MinerInfoDTO {
	return MinerInfoDTO{
		Class:              info.Class,
		HireCost:           info.HireCost,
		CoalTotal:          info.CoalTotal,
		CurrentCoalPerMine: info.CurrentCoalPerMine,
		Working:            info.Working,
		RemainingEnergy:    info.RemainingEnergy,
		Delay:              int(info.Delay.Seconds()),
	}
}
