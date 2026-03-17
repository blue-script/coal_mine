package enterprise

import "errors"

var (
	ErrIncorrectEquipmentName = errors.New("Error: incorrect equipment name")
	ErrGameNotCompleted       = errors.New("Error: game is not completed")
	ErrIncorrectMinerClass    = errors.New("Error: incorrect miner class")
	ErrNotEnoughMoney         = errors.New("Error: not enough money")
	ErrGameAlreadyFinished         = errors.New("Error: game already finished")
)
