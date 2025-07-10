package main

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

/*
	Идея подхода copy-on-write заключается в том, что при чтении данных используется общая копия данных буффера,
	но в случае изменения данных — создается новая копия данных буффера.

	Для реализации такого подхода можно использовать разделяемый счетчик ссылок (refs *int):
	- если при изменении данных буффера кто-то еще ссылается на этот буффер,
	  	то нужно будет сначала произвести копию данных буффера,
	  	изменить счетчик ссылок и только затем произвести изменение
	- если никто не ссылается на буффер, то копировать данные буффера не нужно при изменении данных

	Дополнительно еще нужно реализовать метод конвертации данных буффера в строку без копирования
	и дополнительного выделения памяти.
*/

type COWBuffer struct {
	data []byte
	refs *int
}

func NewCOWBuffer(data []byte) COWBuffer {
	refs := 1
	return COWBuffer{
		data: data,
		refs: &refs,
	}
}

func (b *COWBuffer) Clone() COWBuffer {
	*b.refs++
	return *b
}

func (b *COWBuffer) Close() {
	*b.refs--
}

func (b *COWBuffer) Update(index int, value byte) bool {
	if index < 0 || index > len(b.data)-1 {
		return false
	}

	if *b.refs > 1 {
		newData := make([]byte, len(b.data))
		copy(newData, b.data)
		b.Close()
		*b = NewCOWBuffer(newData)
	}

	b.data[index] = value
	return true
}

func (b *COWBuffer) String() string {
	if len(b.data) == 0 {
		return ""
	}

	return unsafe.String(unsafe.SliceData(b.data), len(b.data))
}

func TestCOWBuffer(t *testing.T) {
	data := []byte{'a', 'b', 'c', 'd'}
	buffer := NewCOWBuffer(data)
	defer buffer.Close()

	copy1 := buffer.Clone()
	copy2 := buffer.Clone()

	assert.Equal(t, unsafe.SliceData(data), unsafe.SliceData(buffer.data))
	assert.Equal(t, unsafe.SliceData(buffer.data), unsafe.SliceData(copy1.data))
	assert.Equal(t, unsafe.SliceData(copy1.data), unsafe.SliceData(copy2.data))

	assert.True(t, (*byte)(unsafe.SliceData(data)) == unsafe.StringData(buffer.String()))
	assert.True(t, (*byte)(unsafe.StringData(buffer.String())) == unsafe.StringData(copy1.String()))
	assert.True(t, (*byte)(unsafe.StringData(copy1.String())) == unsafe.StringData(copy2.String()))

	assert.True(t, buffer.Update(0, 'g'))
	assert.False(t, buffer.Update(-1, 'g'))
	assert.False(t, buffer.Update(4, 'g'))

	assert.True(t, reflect.DeepEqual([]byte{'g', 'b', 'c', 'd'}, buffer.data))
	assert.True(t, reflect.DeepEqual([]byte{'a', 'b', 'c', 'd'}, copy1.data))
	assert.True(t, reflect.DeepEqual([]byte{'a', 'b', 'c', 'd'}, copy2.data))

	assert.NotEqual(t, unsafe.SliceData(buffer.data), unsafe.SliceData(copy1.data))
	assert.Equal(t, unsafe.SliceData(copy1.data), unsafe.SliceData(copy2.data))

	copy1.Close()

	previous := copy2.data
	copy2.Update(0, 'f')
	current := copy2.data

	// 1 reference - don't need to copy buffer during update
	assert.Equal(t, unsafe.SliceData(previous), unsafe.SliceData(current))

	copy2.Close()
}

func TestCOWBufferEmptyData(t *testing.T) {
	// Тест с пустыми данными
	buffer := NewCOWBuffer([]byte{})
	defer buffer.Close()

	assert.Equal(t, "", buffer.String())
	assert.False(t, buffer.Update(0, 'a'))
	assert.False(t, buffer.Update(-1, 'a'))

	copy1 := buffer.Clone()
	defer copy1.Close()
	assert.Equal(t, "", copy1.String())
}

func TestCOWBufferNilData(t *testing.T) {
	// Тест с nil данными
	buffer := NewCOWBuffer(nil)
	defer buffer.Close()

	assert.Equal(t, "", buffer.String())
	assert.False(t, buffer.Update(0, 'a'))

	copy1 := buffer.Clone()
	defer copy1.Close()
	assert.Equal(t, "", copy1.String())
}

func TestCOWBufferMultipleUpdates(t *testing.T) {
	// Тест множественных обновлений
	data := []byte{'a', 'b', 'c'}
	buffer := NewCOWBuffer(data)
	defer buffer.Close()

	copy1 := buffer.Clone()
	copy2 := buffer.Clone()
	defer copy1.Close()
	defer copy2.Close()

	// Первое обновление создает новый буфер
	assert.True(t, buffer.Update(0, 'x'))
	assert.Equal(t, "xbc", buffer.String())
	assert.Equal(t, "abc", copy1.String())

	// Второе обновление не должно создавать новый буфер
	oldData := unsafe.SliceData(buffer.data)
	assert.True(t, buffer.Update(1, 'y'))
	newData := unsafe.SliceData(buffer.data)
	assert.Equal(t, oldData, newData)
	assert.Equal(t, "xyc", buffer.String())
}

func TestCOWBufferChainedClones(t *testing.T) {
	// Тест цепочки клонов
	original := NewCOWBuffer([]byte{'a', 'b', 'c'})
	defer original.Close()

	// Создаем цепочку клонов
	clone1 := original.Clone()
	clone2 := clone1.Clone()
	clone3 := clone2.Clone()
	defer clone1.Close()
	defer clone2.Close()
	defer clone3.Close()

	// Все должны указывать на одни данные
	assert.Equal(t, unsafe.SliceData(original.data), unsafe.SliceData(clone1.data))
	assert.Equal(t, unsafe.SliceData(clone1.data), unsafe.SliceData(clone2.data))
	assert.Equal(t, unsafe.SliceData(clone2.data), unsafe.SliceData(clone3.data))

	// Счетчик ссылок должен быть 4
	assert.Equal(t, 4, *original.refs)
}

func TestCOWBufferUpdateBoundaryIndices(t *testing.T) {
	// Тест граничных индексов
	buffer := NewCOWBuffer([]byte{'a', 'b', 'c', 'd', 'e'})
	defer buffer.Close()

	// Обновление первого элемента
	assert.True(t, buffer.Update(0, 'x'))
	assert.Equal(t, byte('x'), buffer.data[0])

	// Обновление последнего элемента
	assert.True(t, buffer.Update(4, 'z'))
	assert.Equal(t, byte('z'), buffer.data[4])

	// Выход за границы
	assert.False(t, buffer.Update(5, 'w'))
	assert.False(t, buffer.Update(-1, 'w'))
	assert.False(t, buffer.Update(100, 'w'))
}

func TestCOWBufferReferenceCountAfterClose(t *testing.T) {
	// Тест счетчика ссылок после закрытия
	buffer := NewCOWBuffer([]byte{'a', 'b', 'c'})

	copy1 := buffer.Clone()
	copy2 := buffer.Clone()

	assert.Equal(t, 3, *buffer.refs)

	buffer.Close()
	assert.Equal(t, 2, *copy1.refs)

	copy1.Close()
	assert.Equal(t, 1, *copy2.refs)

	// После обновления с refs=1 не должно быть копирования
	oldData := unsafe.SliceData(copy2.data)
	copy2.Update(0, 'x')
	assert.Equal(t, oldData, unsafe.SliceData(copy2.data))

	copy2.Close()
	assert.Equal(t, 0, *copy2.refs)
}

func TestCOWBufferLargeData(t *testing.T) {
	// Тест с большими данными
	largeData := make([]byte, 10000)
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	buffer := NewCOWBuffer(largeData)
	defer buffer.Close()

	copy1 := buffer.Clone()
	defer copy1.Close()

	// Проверяем, что данные совпадают
	assert.Equal(t, unsafe.SliceData(buffer.data), unsafe.SliceData(copy1.data))

	// Обновляем и проверяем, что создался новый буфер
	buffer.Update(5000, 255)
	assert.Equal(t, byte(255), buffer.data[5000])
	assert.Equal(t, byte(5000%256), copy1.data[5000])
}

func TestCOWBufferStringConversion(t *testing.T) {
	// Тест конвертации в строку без копирования
	data := []byte("Hello, World!")
	buffer := NewCOWBuffer(data)
	defer buffer.Close()

	str := buffer.String()

	// Проверяем, что строка указывает на те же данные
	strData := unsafe.StringData(str)
	sliceData := (*byte)(unsafe.SliceData(buffer.data))
	assert.Equal(t, sliceData, strData)

	// Проверяем корректность строки
	assert.Equal(t, "Hello, World!", str)
}

func TestCOWBufferSpecialCharacters(t *testing.T) {
	// Тест со специальными символами и UTF-8
	data := []byte("Привет, мир! 🌍")
	buffer := NewCOWBuffer(data)
	defer buffer.Close()

	copy1 := buffer.Clone()
	defer copy1.Close()

	assert.Equal(t, "Привет, мир! 🌍", buffer.String())
	assert.Equal(t, buffer.String(), copy1.String())

	// Обновляем байт в середине UTF-8 символа (может привести к некорректной строке)
	buffer.Update(0, 'H')
	assert.NotEqual(t, buffer.String(), copy1.String())
}

func TestCOWBufferConcurrentClones(t *testing.T) {
	// Тест создания множества клонов
	original := NewCOWBuffer([]byte{'a', 'b', 'c'})
	defer original.Close()

	clones := make([]COWBuffer, 100)
	for i := range clones {
		clones[i] = original.Clone()
	}

	// Проверяем счетчик ссылок
	assert.Equal(t, 101, *original.refs)

	// Закрываем все клоны
	for i := range clones {
		clones[i].Close()
	}

	assert.Equal(t, 1, *original.refs)
}

func TestCOWBufferUpdateAfterAllClonesGone(t *testing.T) {
	// Тест обновления после удаления всех клонов
	buffer := NewCOWBuffer([]byte{'a', 'b', 'c'})

	// Создаем и сразу закрываем клоны
	for i := 0; i < 5; i++ {
		clone := buffer.Clone()
		clone.Close()
	}

	// refs должен быть 1
	assert.Equal(t, 1, *buffer.refs)

	// Обновление не должно создавать копию
	oldData := unsafe.SliceData(buffer.data)
	buffer.Update(1, 'x')
	assert.Equal(t, oldData, unsafe.SliceData(buffer.data))

	buffer.Close()
}
