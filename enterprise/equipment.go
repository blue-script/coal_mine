package enterprise

type EquipmentName string

const (
	EquipmentPickaxe         EquipmentName = "pickaxe"
	EquipmentMineVentilation EquipmentName = "mine_ventilation"
	EquipmentOreCart         EquipmentName = "ore_cart"
)

type EquipmentItem struct {
	Name  EquipmentName
	Cost  Coal
	Count int
}

type Equipment struct {
	Pickaxe         EquipmentItem
	MineVentilation EquipmentItem
	OreCart         EquipmentItem
}

func NewEquipment() *Equipment {
	return &Equipment{
		Pickaxe: EquipmentItem{
			Name: EquipmentPickaxe,
			Cost: 3000,
		},
		MineVentilation: EquipmentItem{
			Name: EquipmentMineVentilation,
			Cost: 15000,
		},
		OreCart: EquipmentItem{
			Name: EquipmentOreCart,
			Cost: 50000,
		},
	}
}

func (e *Equipment) Info() Equipment {
	return Equipment{
		Pickaxe:         e.Pickaxe,
		MineVentilation: e.MineVentilation,
		OreCart:         e.OreCart,
	}
}

func (e *Equipment) Cost(name EquipmentName) (Coal, error) {
	switch name {
	case EquipmentPickaxe:
		return e.Pickaxe.Cost, nil
	case EquipmentMineVentilation:
		return e.MineVentilation.Cost, nil
	case EquipmentOreCart:
		return e.OreCart.Cost, nil
	default:
		return 0, ErrIncorrectEquipmentName
	}
}

func (e *Equipment) Costs() map[EquipmentName]Coal {
	return map[EquipmentName]Coal{
		EquipmentPickaxe:         e.Pickaxe.Cost,
		EquipmentMineVentilation: e.MineVentilation.Cost,
		EquipmentOreCart:         e.OreCart.Cost,
	}
}

func (e *Equipment) Buy(name EquipmentName) error {
	switch name {
	case EquipmentPickaxe:
		e.Pickaxe.Count++
	case EquipmentMineVentilation:
		e.MineVentilation.Count++
	case EquipmentOreCart:
		e.OreCart.Count++
	default:
		return ErrIncorrectEquipmentName
	}

	return nil
}

func (e *Equipment) Counts() map[EquipmentName]int {
	return map[EquipmentName]int{
		EquipmentPickaxe:         e.Pickaxe.Count,
		EquipmentMineVentilation: e.MineVentilation.Count,
		EquipmentOreCart:         e.OreCart.Count,
	}
}

// func (e *Equipment) Completed() bool {
// 	return e.Pickaxe.Count > 0 &&
// 		e.MineVentilation.Count > 0 &&
// 		e.OreCart.Count > 0
// }
func (e *Equipment) Completed() bool {
	return true
}
