package mapstore

import (
	"testing"

	"github.com/ineverbee/wbl0/internal/store"
	"github.com/stretchr/testify/require"
)

func TestSetGet(t *testing.T) {
	ms := NewMapStore(make(map[int]*store.Model, 1))
	id, model := 1, new(store.Model)
	require.NoError(t, ms.Set(&id, model))

	res, err := ms.Get(id)
	require.NoError(t, err)
	require.Equal(t, model, res)

	id = -1
	res, err = ms.Get(id)
	require.Error(t, err)
	require.Nil(t, res)
}
