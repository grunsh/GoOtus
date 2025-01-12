package hw04lrucache

import (
	"fmt"
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	s        sync.Mutex
}

// Струтур для закладки ключа в интерфейс значения. Используется при удалении из кэша.
type DataWithKey struct {
	data interface{}
	Key  Key
}

// Упаковываем в интерфейс данные с ключом. Ключ нам понадобится для удаления элемента индекса кэша.
func PackKeyToData(key Key, data interface{}) interface{} {
	return DataWithKey{
		data: data,
		Key:  key,
	}
}

// Выковыриваем ключ из упаковки интерфейса.
func GetKeyGromData(data interface{}) Key {
	if d, norm := data.(DataWithKey); norm {
		return d.Key
	}
	panic("Чот сломалось :(")
}

func GetValueFromData(data interface{}) interface{} {
	if d, norm := data.(DataWithKey); norm {
		return d.data
	}
	panic("Чот сломалось :(")
}

func (cache *lruCache) Set(key Key, value interface{}) bool {
	cache.s.Lock()
	defer cache.s.Unlock()
	element, ok := cache.items[key]
	if !ok {
		// В индексе элемента нет, значит он добавляется в голову списка.
		cache.queue.PushFront(PackKeyToData(key, value))
		cache.items[key] = cache.queue.Front()
		// Проверяем длину на превышение ёмкости. Если превысили, отрубаем хвост.
		if cache.queue.Len() > cache.capacity {
			delete(cache.items, GetKeyGromData(cache.queue.Back().Value))
			cache.queue.Remove(cache.queue.Back())
		}
	} else {
		// Элемент уже есть. Двигаем его в голову. Следом перезапишем значение, если оно вдруг новое.
		cache.queue.MoveToFront(element)
		cache.queue.Front().Value = PackKeyToData(key, value)
	}
	return ok
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	cache.s.Lock()
	defer cache.s.Unlock()
	element, ok := cache.items[key]
	if ok {
		cache.queue.MoveToFront(element)
		return GetValueFromData(cache.queue.Front().Value), true
	}
	return nil, false
}

func (cache *lruCache) Clear() {
	cache.s.Lock()
	defer cache.s.Unlock()
	cache.queue.DisplayForward()
	for k, v := range cache.items {
		fmt.Println("-", k, v)
		cache.queue.Remove(v)
	}
	cache.items = make(map[Key]*ListItem, cache.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
