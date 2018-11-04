package validator

import (
	"net/url"
	"strings"

	"github.com/dynamicgo/restrpc"
	"github.com/dynamicgo/xerrors"

	"github.com/Jeffail/gabs"
)

type queryReader struct {
	path   []string
	values url.Values
}

// NewQueryReader .
func NewQueryReader(values url.Values) restrpc.Reader {
	return &queryReader{
		values: values,
	}
}

func (reader *queryReader) Search(key string) ([]string, error) {

	path := strings.Join(append(reader.path, key), ".")

	values := reader.values[path]

	return values, nil
}
func (reader *queryReader) Get() ([]string, error) {

	path := strings.Join(reader.path, ".")

	values := reader.values[path]

	return values, nil
}

func (reader *queryReader) Range(f func(key string, reader restrpc.Reader) error) error {

	path := strings.Join(reader.path, ".")

	for key := range reader.values {
		if strings.HasPrefix(key, path) {
			err := f(key, &queryReader{
				path:   strings.Split(key, "."),
				values: reader.values,
			})

			if err != nil {
				return err
			}
		}
	}

	return nil
}
func (reader *queryReader) Path() string {
	return strings.Join(reader.path, ".")
}

func (reader *queryReader) Reader(key string) restrpc.Reader {
	return &queryReader{
		path:   append(reader.path, key),
		values: reader.values,
	}
}

type jsonReader struct {
	path      []string
	container *gabs.Container
}

// NewJSONReader .
func NewJSONReader(content []byte) (restrpc.Reader, error) {
	container, err := gabs.ParseJSON(content)

	if err != nil {
		return nil, xerrors.Wrapf(err, "parse input json content error: %s", string(content))
	}

	return &jsonReader{
		container: container,
	}, nil
}

func (reader *jsonReader) Search(key string) ([]string, error) {

	path := strings.Join(append(reader.path, key), ".")

	value := reader.container.Path(path).String()

	return []string{value}, nil
}
func (reader *jsonReader) Get() ([]string, error) {

	path := strings.Join(reader.path, ".")

	value := reader.container.Path(path).String()

	return []string{value}, nil
}
func (reader *jsonReader) Range(f func(key string, reader restrpc.Reader) error) error {

	path := strings.Join(reader.path, ".")

	children, err := reader.container.Path(path).ChildrenMap()

	if err != nil {
		return xerrors.Wrapf(err, "get path %s children map error", path)
	}

	for key := range children {
		err := f(key, &jsonReader{
			path:      append(reader.path, key),
			container: reader.container,
		})

		if err != nil {
			return err
		}
	}

	return nil
}
func (reader *jsonReader) Path() string {
	return strings.Join(reader.path, ".")
}

func (reader *jsonReader) Reader(key string) restrpc.Reader {
	return &jsonReader{
		path:      append(reader.path, key),
		container: reader.container,
	}
}
