package testit

import (
	"fmt"
	"log"
	"strings"

	"github.com/sonirico/vago/db"

	"gopkg.in/yaml.v2"
)

type Fixture []byte

type SQLFixture struct {
	db   db.Handler
	data Fixture
}

type SQLFixtures []SQLFixture

func (s SQLFixtures) Load() error {
	for _, f := range s {
		if err := f.Load(); err != nil {
			return err
		}
	}

	return nil
}

func NewSQLFixture(db db.Handler, data Fixture) SQLFixture {
	return SQLFixture{
		db:   db,
		data: data,
	}
}

func (s *SQLFixture) Load() error {

	data := make([]map[string][]map[string]interface{}, 0)

	if err := yaml.Unmarshal(s.data, &data); err != nil {
		log.Panicf("Error unmarshaling json %v", err)
	}

	for _, tableData := range data {
		for table, items := range tableData {
			for _, item := range items {

				keys := []string{}
				args := []string{}
				values := []any{}

				for key := range item {
					keys = append(keys, key)
					args = append(args, fmt.Sprintf("$%d", len(keys)))
					values = append(values, item[key])
				}

				sql := fmt.Sprintf(
					"insert into %s (%s) values (%s) on conflict do nothing",
					table,
					strings.Join(keys, ","),
					strings.Join(args, ","),
				)

				rows, err := s.db.Query(sql, values...)
				if err != nil {
					return err
				}
				rows.Close()
			}
		}
	}

	// Restart id sequences to a number that cannot generate conflicts with the fixtures
	s.db.Exec("ALTER SEQUENCE markets.datasource_assets_id_seq RESTART WITH 5000")
	s.db.Exec("ALTER SEQUENCE markets.datasource_pairs_id_seq RESTART WITH 5000")
	s.db.Exec("ALTER SEQUENCE markets.networks_id_seq RESTART WITH 5000")
	s.db.Exec("ALTER SEQUENCE markets.datasource_networks_id_seq RESTART WITH 5000")
	s.db.Exec("ALTER SEQUENCE markets.datasource_asset_networks_id_seq RESTART WITH 5000")
	s.db.Exec("ALTER SEQUENCE markets.networks_assets_id_seq RESTART WITH 5000")

	return nil
}
