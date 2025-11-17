package types

type HashMap[K comparable, V any] map[K]V

func (m HashMap[K, V]) Has(key K) bool {
	_, ok := m[key]
	return ok
}

func (m HashMap[K, V]) Get(key K) V {
	val := m[key]
	return val
}

func (m HashMap[K, V]) Set(key K, val V) {
	m[key] = val
}

func (m HashMap[K, V]) Del(key K) {
	delete(m, key)
}
