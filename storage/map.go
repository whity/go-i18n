package storage

type mapStorage struct {
	data Data
}

func NewMap(data Data) Storage {
	return &mapStorage{data: data}
}

func (s *mapStorage) Load(locale string) (Data, error) {
	localeData := s.data[locale]
	if localeData == nil {
		localeData = Data{}
	}

	data := localeData.(Data)

	return data, nil
}
