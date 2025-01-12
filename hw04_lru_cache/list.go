package hw04lrucache

import (
	"fmt"
)

const CacheCapacity = 5

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(data interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
	DisplayForward()
	DisplayBackward()
}

// Элемент двусвязного списка.
type ListItem struct {
	Value interface{}
	Prev  *ListItem
	Next  *ListItem
}

// Двусвязный список с его атрибутами.
type DLinkedList struct {
	head *ListItem
	tail *ListItem
	size int
}

// Добавление элемента в голову списка.
func (list *DLinkedList) PushFront(data interface{}) *ListItem {
	newNode := &ListItem{Value: data, Prev: nil, Next: nil}
	if list.head == nil {
		list.tail = newNode
	} else {
		newNode.Next = list.head
		list.head.Prev = newNode
	}
	list.head = newNode
	list.size++
	return newNode
}

func (list *DLinkedList) Front() *ListItem {
	return list.head
}

func (list *DLinkedList) Back() *ListItem {
	return list.tail
}

// Добавление элемента в хвост списка.
func (list *DLinkedList) PushBack(data interface{}) *ListItem {
	newNode := &ListItem{Value: data, Prev: nil, Next: nil}
	if list.head == nil {
		list.head = newNode
		list.tail = newNode
	} else {
		newNode.Prev = list.tail
		list.tail.Next = newNode
		list.tail = newNode
	}
	list.size++
	return newNode
}

// Перенос элемента в голову списка.
func (list *DLinkedList) MoveToFront(i *ListItem) {
	if i.Prev == nil {
		return // Голову можно не переносить в голову.
	}
	if i.Next == nil {
		// Переносим хвост в голову
		list.tail.Next, list.head.Prev = list.head, list.tail // Замыкаем список в кольцо.
		list.head, list.tail = list.tail, list.tail.Prev      // Смещаемся указателями головы и хвоста на 1 "назад"
		list.head.Prev, list.tail.Next = nil, nil             // Размыкаем список по новым голове и хвосту
		return
	}
	i.Prev.Next = i.Next // Замыкаем два соседних элемента...
	i.Next.Prev = i.Prev // ... между собой.
	list.head.Prev = i   // Пока ещё голове списка задаём предшественника.
	i.Next = list.head   // Будущей голове списка даём указатель на бывшую голову (становится вторым элементом).
	list.head = i        // Водружаем в голову списка новый текущий элемент.
	list.head.Prev = nil // Не забываем, что у головы нет предыдущего элемента.
}

// Удаление элемента списка.
func (list *DLinkedList) Remove(i *ListItem) {
	if i.Next == nil {
		if i.Prev == nil {
			//  У нас единственный элемент списка. Обнуляем всё на...
			i, list.head, list.tail = nil, nil, nil
			list.size = 0
			return
		}
		// Случай, когда у нас хвост. Откусываем последний элемент.
		list.tail = i.Prev
		list.tail.Next = nil
		list.size--
		return
	}
	// Удаление головы.
	if i.Prev == nil {
		list.head = list.head.Next
		list.head.Prev = nil
		list.size--
		return
	}
	//  Элемент между головой и хвостом. Просто замыкаем список "вокруг" элемента, который удаляем.
	i.Prev.Next, i.Next.Prev = i.Next, i.Prev
	list.size--
}

func (list *DLinkedList) Len() int {
	return list.size
}

// Вывод списка с начала. Для отладки. Если повлияет на оценку, удалю.
func (list *DLinkedList) DisplayForward() {
	fmt.Println("Вывод списка с головы")
	current := list.head
	for current != nil {
		fmt.Printf("%v prev[%p] self[%p] next[%p]-> \n", current.Value, current.Prev, current, current.Next)
		current = current.Next
	}
	fmt.Println("nil")
}

// Вывод списка с конца. Для отладки. Если повлияет на оценку, удалю.
func (list *DLinkedList) DisplayBackward() {
	current := list.tail
	for current != nil {
		fmt.Printf("%v prev[%p] self[%p] next[%p]-> \n", current.Value, current.Prev, current, current.Next)
		current = current.Prev
	}
	fmt.Println("nil")
}

func NewList() List {
	return new(DLinkedList)
}
