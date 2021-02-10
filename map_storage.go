package donkeydb

func newStorage() *mapStorage {
	return &mapStorage{
		Data: make(map[string]interface{}),
	}
}

type mapStorage struct {
	Data map[string]interface{}
}

func (s *mapStorage) Get(key string) (interface{}, error) {
	if value, exists := s.Data[key]; exists {
		return value, nil
	}
	return nil, ErrNothing
}

func (s *mapStorage) Set(key string, value interface{}) error {
	s.Data[key] = value
	return nil
}
