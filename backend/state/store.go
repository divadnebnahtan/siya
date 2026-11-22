package state

type Store struct {
}

func (s *Store) Set(loc Location) {

}

func (s *Store) Get() Location {
	return Location{}
}

func (s *Store) GetAll() []Location {
	return []Location{}
}
