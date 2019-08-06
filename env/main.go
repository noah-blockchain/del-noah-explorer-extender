package env

import (
	"flag"
	"github.com/noah-blockchain/noah-explorer-tools/models"
	"os"
)

func New() *models.ExtenderEnvironment {
	appName := flag.String("app_name", "Noah Extender", "App name")
	baseCoin := flag.String("base_coin", "MNT", "Base coin symbol")
	coinsUpdateTime := flag.Int("coins_upd_time", 3600, "Coins update time in minutes")
	debug := flag.Bool("debug", false, "Debug mode")
	dbName := flag.String("db_name", "", "DB name")
	dbUser := flag.String("db_user", "", "DB user")
	dbPassword := flag.String("db_password", "", "DB password")
	dbMinIdleConns := flag.Int("db_min_idle_conns", 10, "DB min idle connections")
	dbPoolSize := flag.Int("db_pool_size", 20, "DB pool size")
	nodeApi := flag.String("node_api", "", "DB password")
	txChunkSize := flag.Int("tx_chunk_size", 100, "Transactions chunk size")
	eventsChunkSize := flag.Int("event_chunk_size", 100, "Events chunk size")
	stakeChunkSize := flag.Int("stake_chunk_size", 100, "Stake chunk size")
	configFile := flag.String("config", "", "Env file")
	apiHost := flag.String("api_host", "", "API host")
	apiPort := flag.Int("api_port", 8000, "API port")
	wsLink := flag.String("ws_link", "", "WebSocket server link")
	wsKey := flag.String("ws_key", "", "WebSocket API key")
	wrkSaveTxsCount := flag.Int("wrk_save_txs_count", 3, "Count of workers that save transactions")
	wrkSaveTxsOutputCount := flag.Int("wrk_save_txs_output_count", 3, "Count of workers that save transactions output")
	wrkSaveInvalidTxsCount := flag.Int("wrk_save_invtxs_count", 3, "Count of workers that save invalid transactions")
	wrkSaveRewardsCount := flag.Int("wrk_save_rewards_count", 3, "Count of workers that save rewards")
	wrkSaveSlashesCount := flag.Int("wrk_save_slashes_count", 3, "Count of workers that save slashes")
	wrkSaveAddressesCount := flag.Int("wrk_save_addresses_count", 3, "Count of workers that save addresses")
	wrkSaveValidatorTxsCount := flag.Int("wrk_save_val_tx_count", 3, "Count of workers that save transaction-validator link")
	addrChunkSize := flag.Int("addr_chunk_size", 10, "Count of workers that save transaction-validator link")
	wrkUpdateBalanceCount := flag.Int("wrk_upd_balances_count", 1, "Count of workers that update balance")
	wrkGetBalancesFromNodeCount := flag.Int("wrk_node_balance_count", 1, "Count of workers that get balance from node")
	wrkUpdateTxsIndexNumBlocks := flag.Int("wrk_update_txs_index_num_blocks", 120, "Count of blocks that should be reindex")
	wrkUpdateTxsIndexTime := flag.Int("wrk_update_txs_index_time", 60, "Time in seconds which worker sleep before the next iteration")
	rewardAggregateEveryBlocksCount := flag.Int("reward_aggregate_every_blocks_count", 60, "Every X block will be launched reward aggregation")
	rewardAggregateTimeInterval := flag.String("reward_aggregate_time_interval", "hour", "Rewards aggregation time interval('hour' or 'day')")
	flag.Parse()

	envData := new(models.ExtenderEnvironment)

	if envData.DbUser == "" {
		dbUser := os.Getenv("EXPLORER_DB_USER")
		envData.DbUser = dbUser
	}
	if envData.DbName == "" {
		dbName := os.Getenv("EXPLORER_DB_NAME")
		envData.DbName = dbName
	}
	if envData.DbPassword == "" {
		dbPassword := os.Getenv("EXPLORER_DB_PASSWORD")
		envData.DbPassword = dbPassword
	}
	if envData.NodeApi == "" {
		nodeApi := os.Getenv("NOAH_NODE_API")
		envData.NodeApi = nodeApi
	}

	if *configFile != "" {
		config := NewViperConfig(*configFile)
		wsLink := `http://`
		if GetBool(`wsServer.isSecure`) {
			wsLink = `https://`
		}
		wsLink += GetString(`wsServer.link`)
		if GetString(`wsServer.port`) != `` {
			wsLink += `:` + GetString(`wsServer.port`)
		}
		nodeApi := "http://"
		if GetBool("noahApi.isSecure") {
			nodeApi = "https://"
		}
		nodeApi += GetString("noahApi.link") + ":" + GetString("noahApi.port")
		envData.Debug = GetBool("app.debug")
		envData.DbName = GetString("database.name")
		envData.DbUser = GetString("database.user")
		envData.DbPassword = GetString("database.password")
		envData.DbMinIdleConns = GetInt("database.minIdleConns")
		envData.DbPoolSize = GetInt("database.poolSize")
		envData.NodeApi = nodeApi
		envData.TxChunkSize = GetInt("app.txChunkSize")
		envData.AddrChunkSize = GetInt("app.addrChunkSize")
		envData.EventsChunkSize = GetInt("app.eventsChunkSize")
		envData.StakeChunkSize = GetInt("app.stakeChunkSize")
		envData.ApiHost = GetString("extenderApi.host")
		envData.ApiPort = GetInt("extenderApi.port")
		envData.WsLink = wsLink
		envData.WsKey = GetString(`wsServer.key`)
		envData.AppName = GetString("name")
		envData.WrkSaveTxsCount = GetInt("workers.saveTxs")
		envData.WrkSaveTxsOutputCount = GetInt("workers.saveTxsOutput")
		envData.WrkSaveInvTxsCount = GetInt("workers.saveInvalidTxs")
		envData.WrkSaveRewardsCount = GetInt("workers.saveRewards")
		envData.WrkSaveSlashesCount = GetInt("workers.saveSlashes")
		envData.WrkSaveAddressesCount = GetInt("workers.saveAddresses")
		envData.WrkSaveValidatorTxsCount = GetInt("workers.saveTxValidator")
		envData.WrkUpdateBalanceCount = GetInt("workers.updateBalance")
		envData.WrkGetBalancesFromNodeCount = GetInt("workers.balancesFromNode")
		envData.RewardAggregateEveryBlocksCount = GetInt("app.rewardsAggregateBlocksCount")
		envData.RewardAggregateTimeInterval = GetString("app.rewardsAggregateTimeInterval")
		envData.BaseCoin = GetString("app.baseCoin")
		envData.CoinsUpdateTime = GetInt("app.coinsUpdateTimeMinutes")
		envData.WrkUpdateTxsIndexNumBlocks = GetInt("workers.updateTxsIndexNumBlocks")
		envData.WrkUpdateTxsIndexTime = GetInt("workers.updateTxsIndexSleepSec")
	} else {
		envData.AppName = *appName
		envData.Debug = *debug
		envData.DbName = *dbName
		envData.DbUser = *dbUser
		envData.DbPassword = *dbPassword
		envData.DbMinIdleConns = *dbMinIdleConns
		envData.DbPoolSize = *dbPoolSize
		envData.NodeApi = *nodeApi
		envData.TxChunkSize = *txChunkSize
		envData.EventsChunkSize = *eventsChunkSize
		envData.StakeChunkSize = *stakeChunkSize
		envData.ApiHost = *apiHost
		envData.ApiPort = *apiPort
		envData.WsLink = *wsLink
		envData.WsKey = *wsKey
		envData.WrkSaveTxsCount = *wrkSaveTxsCount
		envData.WrkSaveTxsOutputCount = *wrkSaveTxsOutputCount
		envData.WrkSaveInvTxsCount = *wrkSaveInvalidTxsCount
		envData.WrkSaveRewardsCount = *wrkSaveRewardsCount
		envData.WrkSaveSlashesCount = *wrkSaveSlashesCount
		envData.WrkSaveAddressesCount = *wrkSaveAddressesCount
		envData.WrkSaveValidatorTxsCount = *wrkSaveValidatorTxsCount
		envData.AddrChunkSize = *addrChunkSize
		envData.WrkUpdateBalanceCount = *wrkUpdateBalanceCount
		envData.WrkGetBalancesFromNodeCount = *wrkGetBalancesFromNodeCount
		envData.BaseCoin = *baseCoin
		envData.CoinsUpdateTime = *coinsUpdateTime
		envData.WrkUpdateTxsIndexNumBlocks = *wrkUpdateTxsIndexNumBlocks
		envData.WrkUpdateTxsIndexTime = *wrkUpdateTxsIndexTime
		envData.RewardAggregateEveryBlocksCount = *rewardAggregateEveryBlocksCount
		envData.RewardAggregateTimeInterval = *rewardAggregateTimeInterval
	}
	return envData
}
