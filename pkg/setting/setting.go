package setting

import (
	"game_slots_vsn/pkg/app"
	"game_slots_vsn/pkg/ggeoip"
	"game_slots_vsn/pkg/logger"
	"game_slots_vsn/pkg/redis"
	"game_slots_vsn/pkg/server"
)

func Setup() {
	app.SetUp()
	logger.SetUp()
	redis.SetUp()
	ggeoip.SetUp()
	server.SetUp()
}
