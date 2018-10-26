package validator

import (
	"net/url"
	"strings"
)

type queryReader struct {
	path   []string
	values url.Values
}

// NewQueryReader .
func NewQueryReader(values url.Values) Reader {
	return &queryReader{
		values: values,
	}
}

func (reader *queryReader) Search(key string) ([]string, error) {
	return nil, nil
}
func (reader *queryReader) Get() ([]string, error) {
	return nil, nil
}
func (reader *queryReader) Range(func(key string, reader Reader) error) error {
	return nil
}
func (reader *queryReader) Path() string {
	return strings.Join(reader.path, ".")
}

func (reader *queryReader) Reader(key string) Reader {
	return &queryReader{
		path:   append(reader.path, key),
		values: reader.values,
	}
}
