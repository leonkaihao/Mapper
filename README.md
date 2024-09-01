# Go Mapper
When you have a `map[string]interface{}` variable which may come from unmarshalling a json string, you wish to directly access a nested field with a known path. This package helps you to get the job done.

The package support single value and composite structure with map and array.

## Example
The path for accessing nested field has the format as below:  
To access a map field `xxx`: "xxx",  
To access an array index 2: "2",  
To access an array index 3 inside map field `userList`: "userList.3",  
To access an field `name` inside map field `user`: "user.name".  

A more complex example as below:
```Go
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

	err = mapper.Set("key4.3.key3.key3.3", 0.3, true)
	if err != nil {
        return err
    }
    ...
    ...
	sub, err = mapper.Get("key4.3.key3.key3.3")
	if sub.Val() != 0.3 {
        return fmt.Errorf("no match")
    }
```

Please refer to the test file for the detailed usage.