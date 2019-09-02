package main

import (
	"github.com/noah-blockchain/noah-explorer-extender/api"
	"github.com/noah-blockchain/noah-explorer-extender/core"
	"github.com/noah-blockchain/noah-explorer-extender/database/migrate"
	"github.com/noah-blockchain/noah-explorer-extender/env"
)

func main() {
	migrate.Migrate()

	envData := env.New()
	extenderApi := api.New(envData.ApiHost, envData.ApiPort)
	go extenderApi.Run()
	ext := core.NewExtender(envData)
	ext.Run()
}
