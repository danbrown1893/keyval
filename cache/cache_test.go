package cache_test

import (
	"strconv"
	"testing"
	"time"

	"keyval/cache"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetGetEnsureKey(t *testing.T) {
	c := cache.NewCache()

	val1 := "v1"
	val2 := "v2"

	iVal1, notFound, err := c.EnsureKey("val1", val1, time.Now().Add(5*time.Minute))
	assert.True(t, notFound)
	assert.NoError(t, err)
	assert.Equal(t, val1, iVal1)

	iVal1, notFound, err = c.EnsureKey("val1", val1, time.Now().Add(5*time.Minute))
	assert.False(t, notFound)
	assert.NoError(t, err)
	assert.Equal(t, val1, iVal1)

	err = c.Set("val2", val2, time.Now().Add(5*time.Minute))
	assert.NoError(t, err)

	iVal2, err := c.Get("val2")
	assert.NoError(t, err)
	assert.Equal(t, val2, iVal2)
}

func TestSetDelete(t *testing.T) {
	c := cache.NewCache()

	val := "value"
	err := c.Set("val", val, time.Now().Add(5*time.Minute))
	assert.NoError(t, err)

	c.Delete("val")

	v, err := c.Get("val")
	assert.Nil(t, v)
	assert.Error(t, err)
}

func TestExpires(t *testing.T) {
	c := cache.NewCache()
	val := "value"
	err := c.Set("val", val, time.Now().Add(10*time.Millisecond))
	assert.NoError(t, err)

	time.Sleep(20 * time.Millisecond)

	v, err := c.Get("val")
	assert.Nil(t, v)
	assert.Error(t, err)
}

func TestLocalDataType(t *testing.T) {
	type localType struct {
		Data string `json:"my_data"`
	}

	obj := &localType{
		Data: "Some data",
	}

	c := cache.NewCache()

	require.NoError(t, c.Set("sampleValue", obj, time.Now().Add(time.Minute)))

	storedValue, err := c.Get("sampleValue")
	require.NoError(t, err)
	newObj, ok := storedValue.(*localType)
	require.True(t, ok)
	require.Equal(t, obj, newObj)
}

func BenchmarkCache_GetCustomData(b *testing.B) {
	c := cache.NewCache()
	type localType struct {
		Data string `json:"my_data"`
	}

	obj := &localType{
		Data: "Some data",
	}

	err := c.Set("val", obj, time.Now().Add(5*time.Minute))
	assert.NoError(b, err)

	for i := 0; i < b.N; i++ {
		c.Get("val")
	}
}

func BenchmarkCache_GetString(b *testing.B) {
	c := cache.NewCache()
	val := "val"

	err := c.Set("val", val, time.Now().Add(5*time.Minute))
	assert.NoError(b, err)

	for i := 0; i < b.N; i++ {
		c.Get("val")
	}
}

func BenchmarkCache_SetCustomData(b *testing.B) {
	c := cache.NewCache()
	type localType struct {
		Data string `json:"my_data"`
	}

	obj := &localType{
		Data: "Some data",
	}

	zeroTime := time.Time{}

	for i := 0; i < b.N; i++ {
		c.Set("val"+strconv.Itoa(i), obj, zeroTime)
	}
}
