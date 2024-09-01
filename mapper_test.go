package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapper_Get(t *testing.T) {
	data1 := map[string]any{
		"key1": "value1",
		"key2": 2,
		"key3": 3.0,
		"key4": []any{
			1,
			2.0,
			"3",
			map[string]any{
				"key1": 1,
				"key2": []any{1, 2, 3},
				"key3": map[string]any{"key1": 1, "key2": "2", "key3": []any{"1", "2", "3", "4"}},
			},
		},
	}
	mapper := NewMapper(data1)

	// positive tests
	sub, err := mapper.Get("key1")
	assert.NoError(t, err)
	assert.Equal(t, sub.Val(), "value1")

	sub, err = mapper.Get("key2")
	assert.NoError(t, err)
	assert.Equal(t, sub.Val(), 2)

	sub, err = mapper.Get("key3")
	assert.NoError(t, err)
	assert.Equal(t, sub.Val(), 3.0)

	sub, err = mapper.Get("key4.0")
	assert.NoError(t, err)
	assert.Equal(t, sub.Val(), 1)

	sub, err = mapper.Get("key4.1")
	assert.NoError(t, err)
	assert.Equal(t, sub.Val(), 2.0)

	sub, err = mapper.Get("key4.2")
	assert.NoError(t, err)
	assert.Equal(t, sub.Val(), "3")

	sub, err = mapper.Get("key4.3.key1")
	assert.NoError(t, err)
	assert.Equal(t, sub.Val(), 1)

	sub, err = mapper.Get("key4.3.key2.0")
	assert.NoError(t, err)
	assert.Equal(t, sub.Val(), 1)

	sub, err = mapper.Get("key4.3.key2.1")
	assert.NoError(t, err)
	assert.Equal(t, sub.Val(), 2)

	sub, err = mapper.Get("key4.3.key2.2")
	assert.NoError(t, err)
	assert.Equal(t, sub.Val(), 3)

	sub, err = mapper.Get("key4.3.key3.key1")
	assert.NoError(t, err)
	assert.Equal(t, sub.Val(), 1)

	sub, err = mapper.Get("key4.3.key3.key2")
	assert.NoError(t, err)
	assert.Equal(t, sub.Val(), "2")

	sub, err = mapper.Get("key4.3.key3.key3.0")
	assert.NoError(t, err)
	assert.Equal(t, sub.Val(), "1")

	//negative tests
	_, err = mapper.Get("")
	assert.Error(t, err)

	_, err = mapper.Get("key0")
	assert.Error(t, err)

	_, err = mapper.Get("key4.4")
	assert.Error(t, err)

	_, err = mapper.Get("key4.3.key4")
	assert.Error(t, err)

	_, err = mapper.Get("key4.3.key3.key4")
	assert.Error(t, err)

	_, err = mapper.Get("key4.3.key3.key3.4")
	assert.Error(t, err)
}

func TestMapper_Set(t *testing.T) {
	data1 := map[string]any{
		"key1": "value1",
		"key2": 2,
		"key3": 3.0,
		"key4": []any{
			1,
			2.0,
			"3",
			map[string]any{
				"key1": 1,
				"key2": []any{1, 2, 3},
				"key3": map[string]any{"key1": 1, "key2": "2", "key3": []any{"1", "2", "3", "4"}},
			},
		},
	}
	mapper := NewMapper(data1)

	// positive tests
	err := mapper.Set("key1", "value2", false)
	assert.NoError(t, err)
	sub, _ := mapper.Get("key1")
	assert.Equal(t, sub.Val(), "value2")

	err = mapper.Set("key1", 1, false)
	assert.NoError(t, err)
	sub, _ = mapper.Get("key1")
	assert.Equal(t, sub.Val(), 1)

	err = mapper.Set("key2", "2", false)
	assert.NoError(t, err)
	err = mapper.Set("key4.0", 0, false)
	assert.NoError(t, err)
	sub, _ = mapper.Get("key4.0")
	assert.Equal(t, sub.Val(), 0)
	err = mapper.Set("key4.3.0", 0, true)
	assert.NoError(t, err)

	err = mapper.Set("key4.3.key2.2", "2", false)
	assert.NoError(t, err)
	sub, _ = mapper.Get("key4.3.key2.2")
	assert.Equal(t, sub.Val(), "2")

	err = mapper.Set("key4.3.key3.1", 0, true)
	assert.NoError(t, err)
	err = mapper.Set("key4.3.key3.key3.3", 0.3, true)
	assert.NoError(t, err)
	sub, _ = mapper.Get("key4.3.key3.key3.3")
	assert.Equal(t, sub.Val(), 0.3)

	//negative tests
	err = mapper.Set("key5", "value2", false)
	assert.Error(t, err)
	err = mapper.Set("key4.4", "value2", false)
	assert.Error(t, err)
	err = mapper.Set("key4.key1", "value2", true)
	assert.Error(t, err)
	err = mapper.Set("key4.3.key3.key1.1", 0, true)
	assert.Error(t, err)
	err = mapper.Set("key4.3.key3.key3.key1", 0, true)
	assert.Error(t, err)
	err = mapper.Set("key4.3.key3.key3.4", 0, false)
	assert.Error(t, err)
	err = mapper.Set("key4.3.key3.key3.4", 0, true)
	assert.Error(t, err)
}
