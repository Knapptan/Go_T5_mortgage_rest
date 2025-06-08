// Тесты пакета cache
package cache

import (
	"testing"

	"github.com/Knapptan/Go_T5_mortgage_rest/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	c := NewCache()
	assert.NotNil(t, c)

	item := &models.MortgageInfoResponse{ID: 1}
	c.set(*item)

	result, found := c.Get(1)
	assert.True(t, found)
	assert.Equal(t, 1, result.ID)

}
