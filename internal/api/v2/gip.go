package v2

import (
	"game_slots_vsn/pkg/e"
	"game_slots_vsn/pkg/ggeoip"
	"github.com/gin-gonic/gin"
)

func GetCountry(c *gin.Context) {
	appG := e.Gin{C: c}
	ipStr := c.Query("ip")
	country := ggeoip.Gip.GetCountryByIP(ipStr)
	appG.Success(e.SUCCESS, map[string]interface{}{
		"country": country,
	})
	return
}
func GetCountryAndCity(c *gin.Context) {
	appG := e.Gin{C: c}
	ipStr := c.Query("ip")
	country, cities := ggeoip.Gip.GetCountryAndCityByIP(ipStr)
	appG.Success(e.SUCCESS, map[string]interface{}{
		"country": country,
		"cities":  cities,
	})
	return
}
