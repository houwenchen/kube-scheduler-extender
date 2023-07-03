package storage

import (
	"sync"
)

type Storage struct {
	// 定义一个存储 pod 与 node 对应关系的 map
	// eg: key for pod: nodename
	storage map[string]string
	m       sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		storage: make(map[string]string),
	}
}

func (s *Storage) DeletePodOfStorage(key string) error {
	s.m.Lock()
	defer s.m.Unlock()

	delete(s.storage, key)

	return nil
}

func (s *Storage) AddPodOfStorage(key string, nodeName string) error {
	s.m.Lock()
	defer s.m.Unlock()

	s.storage[key] = nodeName

	return nil
}

func (s *Storage) GetPodOfStorage(key string) (string, bool) {
	s.m.RLock()
	defer s.m.RUnlock()

	nodeName, exist := s.storage[key]

	return nodeName, exist
}
