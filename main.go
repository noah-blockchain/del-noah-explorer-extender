package noah_explorer_extender

import (
	"github.com/noah-blockchain/noah-explorer-extender/api"
	"github.com/noah-blockchain/noah-explorer-extender/core"
	"github.com/noah-blockchain/noah-explorer-extender/env"
)

func main() {
	envData := env.New()
	extenderApi := api.New(envData.ApiHost, envData.ApiPort)
	go extenderApi.Run()
	ext := core.NewExtender(envData)
	ext.Run()
}
