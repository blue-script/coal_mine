package rest

import (
	"net/http"

	"github.com/blue-script/coal_mine/enterprise"
	"github.com/gin-gonic/gin"
)

type HTTPHandlers struct {
	enterprise *enterprise.Enterprise
}

func NewHTTPHandlers(enterprise *enterprise.Enterprise) *HTTPHandlers {
	return &HTTPHandlers{
		enterprise: enterprise,
	}
}

func (h *HTTPHandlers) HireMiner(c *gin.Context) {
	var req HireMinerDTO

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

func (h *HTTPHandlers) HireCost(c *gin.Context) {
	class := enterprise.MinerClass(c.Param("class"))

	info, err := h.enterprise.HireCost(class)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toMinerInfoDTO(info))
}
