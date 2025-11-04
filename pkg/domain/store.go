package domain


type Store interface {
 Atomic (func(Store) error) error
 ImageRepository() ImageRepository
}