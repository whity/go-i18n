package storage

type mapStorage struct {
	data Data
}

func NewMap(data Data) (Storage, error) {
	return &mapStorage{data: data}, nil
}

func (s *mapStorage) Load(locale string) (Data, error) {
	localeData := s.data[locale]
	if localeData == nil {
		localeData = Data{}
	}

	data := localeData.(Data)

	return data, nil
}
