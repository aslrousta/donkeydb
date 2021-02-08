package donkeydb

func newStorage() *storage {
	return &storage{
		Data: make(map[string]interface{}),
	}
}

type storage struct {
	Data map[string]interface{}
}

func (s *storage) Get(key string) (interface{}, error) {
	if value, exists := s.Data[key]; exists {
		return value, nil
	}
	return nil, ErrNothing
}

func (s *storage) Set(key string, value interface{}) error {
	s.Data[key] = value
	return nil
}
