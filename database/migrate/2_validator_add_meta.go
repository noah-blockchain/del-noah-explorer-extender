package migrate

import "github.com/go-pg/migrations/v7"

const SqlCommand2 = `
ALTER TABLE public.validators ADD COLUMN name varchar(64);
ALTER TABLE public.validators ADD COLUMN site_url varchar(100);
ALTER TABLE public.validators ADD COLUMN icon_url varchar(100);
ALTER TABLE public.validators ADD COLUMN description text;
ALTER TABLE public.validators ADD COLUMN meta_updated_at_block_id integer;
`

func init() {
	_ = migrations.Register(func(db migrations.DB) error {
		_, err := db.Exec(SqlCommand2)
		if err != nil {
			return err
		}
		return nil
	})
}
