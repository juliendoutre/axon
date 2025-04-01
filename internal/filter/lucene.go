package filter

import (
	"fmt"

	"github.com/grindlemire/go-lucene"
	"github.com/grindlemire/go-lucene/pkg/driver"
)

func New() *Parser {
	return &Parser{driver: driver.NewPostgresDriver()}
}

type Parser struct {
	driver driver.PostgresDriver
}

func (p *Parser) Parse(query string) (string, []any, error) {
	ast, err := lucene.Parse(query)
	if err != nil {
		return "", nil, fmt.Errorf("failed parsing lucene query: %w", err)
	}

	// TODO: check filters are restricted to `attributes.*`, `observer_claims.*` and a few others.

	filter, params, err := p.driver.RenderParam(ast)
	if err != nil {
		return "", nil, fmt.Errorf("failed rendering lucene query: %w", err)
	}

	return filter, params, nil
}
