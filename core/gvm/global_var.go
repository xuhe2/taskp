package gvm

import (
	"errors"
	"sync"
)

var globalVar = map[string]any{}
var globalVarMapMutex = sync.RWMutex{}

var ErrVarNotExist = errors.New("variable not exist")
var ErrVarTypeNotMatch = errors.New("variable type not match")

func SetGlobalVar[T any](key string, value T) {
	globalVarMapMutex.Lock()
	defer globalVarMapMutex.Unlock()

	globalVar[key] = value
}

func GetGlobalVar[T any](key string) (T, error) {
	globalVarMapMutex.Lock()

	value, ok := globalVar[key]
	if !ok {
		var zero T
		return zero, ErrVarNotExist
	}

	globalVarMapMutex.Unlock()

	valueAfterCast, ok := value.(T)
	if !ok {
		var zero T
		return zero, ErrVarTypeNotMatch
	}
	return valueAfterCast, nil
}
