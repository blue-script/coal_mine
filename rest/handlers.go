package rest

import (
	"net/http"

	"github.com/blue-script/coal_mine/enterprise"
	"github.com/gin-gonic/gin"
)

type HTTPHandlers struct {
	enterprise *enterprise.Enterprise
	onFinish   func()
}

func NewHTTPHandlers(enterprise *enterprise.Enterprise) *HTTPHandlers {
	return &HTTPHandlers{
		enterprise: enterprise,
	}
}

func (h *HTTPHandlers) SetOnFinish(fn func()) {
	h.onFinish = fn
}

func (h *HTTPHandlers) HireMiner(c *gin.Context) {
	var req MinerClassDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	minerInfo, err := h.enterprise.HireMiner(req.Class)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, minerInfo)
}

func (h *HTTPHandlers) BuyEquipment(c *gin.Context) {
	name := enterprise.EquipmentName(c.Param("name"))

	if err := h.enterprise.BuyEquipment(name); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *HTTPHandlers) MinerCost(c *gin.Context) {
	class := enterprise.MinerClass(c.Param("class"))

	info, err := h.enterprise.HireCost(class)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toMinerInfoDTO(info))
}

func (h *HTTPHandlers) ListMiners(c *gin.Context) {
	class := c.Query("class")

	var miners []enterprise.MinerInfo
	if class == "" {
		miners = h.enterprise.ListAllMiners()
	} else {
		minerClass := enterprise.MinerClass(class)
		miners = h.enterprise.ListMiners(&minerClass)
	}

	c.JSON(http.StatusOK, toMinerInfoDTOs(miners))
}

func (h *HTTPHandlers) EquipmentPrices(c *gin.Context) {
	costs := h.enterprise.EquipmentCosts()

	result := make([]EquipmentPriceDTO, 0, len(costs))
	for name, cost := range costs {
		result = append(result, EquipmentPriceDTO{
			Name: name,
			Cost: cost,
		})
	}

	c.JSON(http.StatusOK, result)
}

func (h *HTTPHandlers) ListEquipment(c *gin.Context) {
	equipment := h.enterprise.ListEquipment()

	c.JSON(http.StatusOK, toEquipmentDTO(equipment))
}

func (h *HTTPHandlers) EnterpriseStatistic(c *gin.Context) {
	statistic := h.enterprise.Statistic()

	c.JSON(http.StatusOK, toEnterpriseStatisticDTO(statistic))
}

func (h *HTTPHandlers) FinishGame(c *gin.Context) {
	statistic, err := h.enterprise.FinishGame()
	if err != nil {
		switch err {
		case enterprise.ErrGameNotCompleted:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case enterprise.ErrGameAlreadyFinished:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		return
	}

	c.JSON(http.StatusOK, toEnterpriseStatisticDTO(statistic))

	if h.onFinish != nil {
		go h.onFinish()
	}
}
