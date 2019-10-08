package migrate

import (
	"github.com/go-pg/migrations/v7"
)

const SqlCommand = `
	create schema if not exists public;

DROP TYPE IF EXISTS public.rewards_role;

create type public.rewards_role as enum ('Validator', 'Delegator', 'DAO', 'Developers');

alter type public.rewards_role owner to noah;

create table if not exists public.addresses
(
    id                  bigserial   not null
        constraint addresses_pkey
            primary key,
    address             varchar(64) not null,
    updated_at          timestamp with time zone,
    updated_at_block_id bigint
);

comment on column public.addresses.address is 'Address hex string without prefix(NOAHx****)';

comment on column public.addresses.updated_at is 'Last balance parsing time';

comment on column public.addresses.updated_at_block_id is 'Block id, that have transactions or events, that triggers address record to update from api-method GET /address';

alter table public.addresses
    owner to noah;

create unique index if not exists addresses_address_uindex
    on public.addresses (address);

create table if not exists public.blocks
(
    id                    integer                                not null
        constraint blocks_pkey
            primary key,
    total_txs             bigint                   default 0     not null,
    size                  bigint                                 not null,
    proposer_validator_id integer                                not null,
    num_txs               integer                  default 0     not null,
    block_time            bigint                                 not null,
    created_at            timestamp with time zone               not null,
    updated_at            timestamp with time zone default now() not null,
    block_reward          numeric(70)                            not null,
    hash                  varchar(64)                            not null
);

comment on table public.blocks is 'Address entity table';

comment on column public.blocks.total_txs is 'Total count of txs in blockchain';

comment on column public.blocks.proposer_validator_id is 'Proposer public key (Np***)';

comment on column public.blocks.num_txs is 'Count of txs in block';

comment on column public.blocks.block_time is 'Block operation time (???) in microseconds';

comment on column public.blocks.created_at is 'Datetime of block creation("time" field from api)';

comment on column public.blocks.updated_at is 'Time of record last update';

comment on column public.blocks.block_reward is 'Sum of all block rewards';

comment on column public.blocks.hash is 'Hex string';

alter table public.blocks
    owner to noah;

create index if not exists blocks_proposer_validator_id_index
    on public.blocks (proposer_validator_id);

create index if not exists blocks_created_at_index
    on public.blocks (created_at desc);

create table if not exists public.coins
(
    id                      serial                                 not null
        constraint coins_pkey
            primary key,
    creation_address_id     bigint
        constraint coins_addresses_id_fk
            references public.addresses,
    creation_transaction_id bigint,
    crr                     integer,
    volume                  numeric(70),
    reserve_balance         numeric(70),
    name                    varchar(255),
    symbol                  varchar(20)                            not null,
    updated_at              timestamp with time zone default now() not null,
    deleted_at              timestamp with time zone
);

comment on column public.coins.creation_address_id is 'Id of creator address in address table';

comment on column public.coins.reserve_balance is 'Reservation balance for coin creation
';

comment on column public.coins.name is 'Name of coin';

comment on column public.coins.symbol is 'Short symbol of coin';

comment on column public.coins.updated_at is 'Timestamp of coin balance/value updation(from api for example)';

alter table public.coins
    owner to noah;

create table if not exists public.balances
(
    id         bigserial   not null
        constraint balances_pkey
            primary key,
    address_id bigint      not null
        constraint balances_addresses_id_fk
            references public.addresses,
    coin_id    integer     not null
        constraint balances_coins_id_fk
            references public.coins,
    value      numeric(70) not null
);

alter table public.balances
    owner to noah;

create unique index if not exists balances_address_id_coind_id_uindex
    on public.balances (address_id, coin_id);

create index if not exists balances_address_id_index
    on public.balances (address_id);

create index if not exists balances_coind_id_index
    on public.balances (coin_id);

create unique index if not exists coins_creation_transaction_id_uindex
    on public.coins (creation_transaction_id);

create index if not exists coins_creator_address_id_index
    on public.coins (creation_address_id);

create unique index if not exists coins_symbol_uindex
    on public.coins (symbol);

create table if not exists public.invalid_transactions
(
    id              bigserial                not null
        constraint invalid_transactions_pkey
            primary key,
    from_address_id bigint                   not null
        constraint invalid_transactions_addresses_id_fk
            references public.addresses,
    block_id        integer                  not null
        constraint invalid_transactions_blocks_id_fk
            references public.blocks,
    created_at      timestamp with time zone not null,
    type            smallint                 not null,
    hash            varchar(64)              not null,
    tx_data         jsonb                    not null
);

comment on column public.invalid_transactions.created_at is 'Duplicate of block created_at for less joins listings';

alter table public.invalid_transactions
    owner to noah;

create index if not exists invalid_transactions_block_id_from_address_id_index
    on public.invalid_transactions (block_id desc, from_address_id asc);

create index if not exists invalid_transactions_from_address_id_index
    on public.invalid_transactions (from_address_id);

create index if not exists invalid_transactions_hash_index
    on public.invalid_transactions (hash);

create table if not exists public.transactions
(
    id              bigserial                not null
        constraint transactions_pkey
            primary key,
    from_address_id bigint                   not null
        constraint transactions_addresses_id_fk
            references public.addresses,
    nonce           bigint                   not null,
    gas_price       bigint                   not null,
    gas             bigint                   not null,
    block_id        integer                  not null
        constraint transactions_blocks_id_fk
            references public.blocks,
    gas_coin_id     integer                  not null
        constraint transactions_coins_id_fk
            references public.coins,
    created_at      timestamp with time zone not null,
    type            smallint                 not null,
    hash            varchar(64)              not null,
    service_data    text,
    data            jsonb                    not null,
    tags            jsonb                    not null,
    payload         bytea,
    raw_tx          bytea                    not null
);

comment on column public.transactions.from_address_id is 'Link to address, from that tx was signed';

comment on column public.transactions.block_id is 'Link to block';

comment on column public.transactions.created_at is 'Timestamp of tx = timestamp of block. Duplicate data for less joins on blocks';

comment on column public.transactions.type is 'Integer index of tx type';

comment on column public.transactions.hash is 'Tx hash 64 symbols hex string without prefix(Nt****). Because of key-value-only filtering uses hash index';

comment on column public.transactions.payload is 'transaction payload in bytes';

comment on column public.transactions.raw_tx is 'Raw tx data in bytes';

alter table public.transactions
    owner to noah;

alter table public.coins
    add constraint coins_transactions_id_fk
        foreign key (creation_transaction_id) references public.transactions;

create table if not exists public.transaction_outputs
(
    id             bigserial   not null
        constraint transaction_outputs_pkey
            primary key,
    transaction_id bigint      not null
        constraint transaction_outputs_transactions_id_fk
            references public.transactions,
    to_address_id  bigint      not null
        constraint transaction_outputs_addresses_id_fk
            references public.addresses,
    coin_id        integer     not null
        constraint transaction_outputs_coins_id_fk
            references public.coins,
    value          numeric(70) not null
);

comment on column public.transaction_outputs.value is 'Value of tx output';

alter table public.transaction_outputs
    owner to noah;

create index if not exists transaction_outputs_coin_id_index
    on public.transaction_outputs (coin_id);

create index if not exists transaction_outputs_transaction_id_index
    on public.transaction_outputs (transaction_id);

create index if not exists transaction_outputs_address_id_index
    on public.transaction_outputs (to_address_id);

create index if not exists transactions_block_id_from_address_id_index
    on public.transactions (block_id desc, from_address_id asc);

create index if not exists transactions_from_address_id_index
    on public.transactions (from_address_id);

create index if not exists transactions_hash_index
    on public.transactions (hash);

create table if not exists public.validators
(
    id                       serial                                 not null
        constraint validator_public_keys_pkey
            primary key,
    reward_address_id        bigint,
    owner_address_id         bigint,
    created_at_block_id      integer,
    status                   integer,
    commission               integer,
    total_stake              numeric(70),
    public_key               varchar(70)                            not null,
    update_at                timestamp with time zone default now() not null,
    name                     varchar(64),
    site_url                 varchar(100),
    icon_url                 varchar(100),
    description              text,
    meta_updated_at_block_id integer
);

comment on table public.validators is 'ATTENTION - only public _ey is not null field, other fields can be null';

alter table public.validators
    owner to noah;

create table if not exists public.block_validator
(
    block_id     bigint                not null
        constraint block_validator_blocks_id_fk
            references public.blocks,
    validator_id integer               not null
        constraint block_validator_validators_id_fk
            references public.validators,
    signed       boolean default false not null,
    constraint block_validator_pk
        primary key (block_id, validator_id)
);

alter table public.block_validator
    owner to noah;

create index if not exists block_validator_block_id_index
    on public.block_validator (block_id);

create index if not exists block_validator_validator_id_index
    on public.block_validator (validator_id);

create table if not exists public.rewards
(
    address_id   bigint       not null
        constraint rewards_addresses_id_fk
            references public.addresses,
    block_id     integer      not null
        constraint rewards_blocks_id_fk
            references public.blocks,
    validator_id integer      not null
        constraint rewards_validators_id_fk
            references public.validators,
    role         rewards_role not null,
    amount       numeric(70)  not null
);

alter table public.rewards
    owner to noah;

create index if not exists rewards_address_id_index
    on public.rewards (address_id);

create index if not exists rewards_block_id_index
    on public.rewards (block_id);

create index if not exists rewards_validator_id_index
    on public.rewards (validator_id);

create table if not exists public.aggregated_rewards
(
    time_id       timestamp with time zone not null,
    to_block_id   integer                  not null
        constraint aggregated_rewards_to_blocks_id_fk
            references public.blocks,
    from_block_id integer                  not null
        constraint aggregated_rewards_from_blocks_id_fk
            references public.blocks,
    address_id    bigint                   not null
        constraint aggregated_rewards_addresses_id_fk
            references public.addresses,
    validator_id  integer                  not null
        constraint aggregated_rewards_validators_id_fk
            references public.validators,
    role          rewards_role             not null,
    amount        numeric(70)              not null
);

alter table public.aggregated_rewards
    owner to noah;

create index if not exists aggregated_rewards_address_id_index
    on public.aggregated_rewards (address_id);

create index if not exists aggregated_rewards_validator_id_index
    on public.aggregated_rewards (validator_id);

create index if not exists aggregated_rewards_time_id_index
    on public.aggregated_rewards (time_id);

create unique index if not exists aggregated_rewards_unique_index
    on public.aggregated_rewards (time_id, address_id, validator_id, role);

create table if not exists public.slashes
(
    id           bigserial   not null
        constraint slashes_pkey
            primary key,
    address_id   bigint      not null
        constraint slashes_addresses_id_fk
            references public.addresses,
    block_id     integer     not null
        constraint slashes_blocks_id_fk
            references public.blocks,
    validator_id integer     not null
        constraint slashes_validators_id_fk
            references public.validators,
    coin_id      integer     not null
        constraint slashes_coins_id_fk
            references public.coins,
    amount       numeric(70) not null
);

alter table public.slashes
    owner to noah;

create index if not exists slashes_address_id_index
    on public.slashes (address_id);

create index if not exists slashes_block_id_index
    on public.slashes (block_id);

create index if not exists slashes_coin_id_index
    on public.slashes (coin_id);

create index if not exists slashes_validator_id_index
    on public.slashes (validator_id);

create table if not exists public.stakes
(
    id               serial      not null,
    owner_address_id bigint      not null
        constraint stakes_addresses_id_fk
            references public.addresses,
    validator_id     integer     not null
        constraint stakes_validators_id_fk
            references public.validators,
    coin_id          integer     not null
        constraint stakes_coins_id_fk
            references public.coins,
    value            numeric(70) not null,
    noah_value       numeric(70) not null,
    constraint stakes_pkey
        primary key (validator_id, owner_address_id, coin_id)
);

alter table public.stakes
    owner to noah;

create index if not exists stakes_coin_id_index
    on public.stakes (coin_id);

create index if not exists stakes_owner_address_id_index
    on public.stakes (owner_address_id);

create index if not exists stakes_validator_id_index
    on public.stakes (validator_id);

create table if not exists public.transaction_validator
(
    transaction_id bigint  not null
        constraint transaction_validator_transactions_id_fk
            references public.transactions,
    validator_id   integer not null
        constraint transaction_validator_validators_id_fk
            references public.validators,
    constraint transaction_validator_pk
        primary key (transaction_id, validator_id)
);

alter table public.transaction_validator
    owner to noah;

create index if not exists transaction_validator_validator_id_index
    on public.transaction_validator (validator_id);

create unique index if not exists validator_public_keys_public_key_uindex
    on public.validators (public_key);

create table if not exists public.index_transaction_by_address
(
    block_id       bigint not null
        constraint index_transaction_by_address_blocks_id_fk
            references public.blocks,
    address_id     bigint not null
        constraint index_transaction_by_address_addresses_id_fk
            references public.addresses,
    transaction_id bigint not null
        constraint index_transaction_by_address_transactions_id_fk
            references public.transactions,
    constraint index_transaction_by_address_pk
        primary key (block_id, address_id, transaction_id)
);

alter table public.index_transaction_by_address
    owner to noah;

create index if not exists index_transaction_by_address_address_id_index
    on public.index_transaction_by_address (address_id);

create index if not exists index_transaction_by_address_block_id_address_id_index
    on public.index_transaction_by_address (block_id, address_id);

create index if not exists index_transaction_by_address_transaction_id_index
    on public.index_transaction_by_address (transaction_id);

INSERT INTO public.coins (symbol) VALUES ('NOAH');
`

func init() {
	_ = migrations.Register(func(db migrations.DB) error {
		_, err := db.Exec(SqlCommand)
		if err != nil {
			return err
		}
		return nil
	})
}
