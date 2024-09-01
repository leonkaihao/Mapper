package main

import (
	"fmt"
	"strconv"
	"strings"
)

type MapperNodeType int

const (
	MAPPER_NODETYPE_NONE  MapperNodeType = 0
	MAPPER_NODETYPE_DICT  MapperNodeType = 1
	MAPPER_NODETYPE_ARRAY MapperNodeType = 2
	MAPPER_NODETYPE_LEAF  MapperNodeType = 3
)

type Mapper struct {
	nodeType MapperNodeType
	val      any
}

func NewMapper(src any) *Mapper {
	if src == nil {
		return &Mapper{MAPPER_NODETYPE_NONE, nil}
	}
	var nodeType MapperNodeType
	switch src.(type) {
	case map[string]any:
		nodeType = MAPPER_NODETYPE_DICT
	case []any:
		nodeType = MAPPER_NODETYPE_ARRAY
	default:
		nodeType = MAPPER_NODETYPE_LEAF
	}
	return &Mapper{nodeType, src}
}

func (mp *Mapper) Type() MapperNodeType {
	return mp.nodeType
}

func (mp *Mapper) Val() any {
	return mp.val
}

func (mp *Mapper) Get(path string) (*Mapper, error) {
	path = strings.TrimSpace(path)
	routes := strings.Split(path, ".")
	if len(routes) == 0 || routes[0] == "" {
		return nil, fmt.Errorf("MapperGet: Get should have at least 1 nested level in path param")
	}
	cur := mp.val
	if cur == nil {
		return nil, fmt.Errorf("MapperGet: can't get sub field from given path %v, parent node is empty", path)
	}
	var ok bool
	for _, route := range routes {
		switch item := cur.(type) {
		case map[string]any:
			if cur, ok = item[route]; !ok {
				return nil, fmt.Errorf("MapperGet: get dict field %v failed in path %v", route, path)
			}
		case []any:
			index, err := strconv.Atoi(route)
			if err != nil {
				return nil, fmt.Errorf("MapperGet: get array index %v failed in path %v", route, path)
			}
			if index >= len(item) {
				return nil, fmt.Errorf("MapperGet: get array index %v out of actual size %v in path %v", route, len(item), path)
			}
			cur = item[index]
		default:
			return nil, fmt.Errorf("MapperGet: get field %v from path %v is not supported", route, path)
		}
	}
	return NewMapper(cur), nil
}

func (mp *Mapper) Set(path string, val any, force bool) error {
	path = strings.TrimSpace(path)
	routes := strings.Split(path, ".")
	if len(routes) == 0 || routes[0] == "" {
		return fmt.Errorf("MapperSet: Set should have at least 1 nested level in path")
	}

	cur := mp.val
	if cur == nil {
		return fmt.Errorf("MapperSet: can't get sub field from given path %v, parent node is empty", path)
	}
	var ok bool
	for i, route := range routes {
		switch item := cur.(type) {
		case map[string]any:
			if cur, ok = item[route]; !ok {
				if !force {
					return fmt.Errorf("MapperSet: get dict field %v failed in path %v", route, path)
				}
				if i < len(routes)-1 {
					next := routes[i+1]
					if _, err := strconv.Atoi(next); err == nil {
						item[route] = []any{}
					} else {
						item[route] = map[string]any{}
					}
					cur = item[route]
				} else {
					item[route] = val
					return nil
				}
			} else {
				if i == len(routes)-1 {
					item[route] = val
					return nil
				}
			}
		case []any:
			index, err := strconv.Atoi(route)
			if err != nil {
				return fmt.Errorf("MapperSet: get array index %v failed in path %v", route, path)
			}
			if index >= len(item) {
				return fmt.Errorf("MapperSet: get array index %v out of actual size %v in path %v", route, len(item), path)
			}
			if i < len(routes)-1 {
				cur = item[index]
			} else {
				item[index] = val
				return nil
			}
		default:
			return fmt.Errorf("MapperSet: set field %v from path %v is not supported", route, path)
		}
	}
	return nil
}
