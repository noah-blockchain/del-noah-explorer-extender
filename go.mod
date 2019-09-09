module github.com/noah-blockchain/noah-explorer-extender

go 1.12

replace mellium.im/sasl v0.2.1 => github.com/mellium/sasl v0.2.1

require (
	github.com/centrifugal/gocent v2.0.2+incompatible
	github.com/go-pg/migrations/v7 v7.1.3
	github.com/go-pg/pg v8.0.5+incompatible
	github.com/go-pg/pg/v9 v9.0.0-beta.7
	github.com/noah-blockchain/api v0.1.1
	github.com/noah-blockchain/noah-explorer-tools v0.1.1
	github.com/noah-blockchain/noah-go-node v0.1.2
	github.com/noah-blockchain/noah-node-go-api v0.1.1
	github.com/prometheus/client_golang v1.1.0
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/viper v1.4.0
)
