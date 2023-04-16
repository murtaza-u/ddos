package operate

import (
	"context"
	"fmt"

	"github.com/murtaza-u/ddos/store"
)

type ddosOpt struct {
	ctx      context.Context
	store    store.Storer
	resource string
	id       string
}

func NewDDoSOperator(
	ctx context.Context,
	s store.Storer,
	name string,
	id string) Operator {

	resource := fmt.Sprintf("/registry/%s/ddos", name)
	return NewDefaultOperator(ctx, s, resource, id)
}
