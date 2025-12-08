package testit

import (
	"fmt"
	"log"
	"strings"

	"github.com/sonirico/vago/db"

	"gopkg.in/yaml.v2"
)

type CHFixture struct {
	db   db.Handler
	data []byte
}

func NewCHFixture(db db.Handler, data []byte) *CHFixture {
	return &CHFixture{
		db:   db,
		data: data,
	}
}

func (s *CHFixture) Exec(stmt string, args ...any) error {
	_, err := s.db.Exec(stmt, args...)
	if err != nil {
		return err
	}

	return nil
}

func (s *CHFixture) Load() error {

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
					"INSERT INTO %s (%s) VALUES (%s)",
					table,
					strings.Join(keys, ","),
					strings.Join(args, ","),
				)

				if err := s.Exec(sql, values...); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
