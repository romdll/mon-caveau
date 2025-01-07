package injector

import (
	"moncaveau/database"
)

type EntityToMapKeyFunc[T any] func(entity T) interface{}

type Injector[T any] struct {
	entityToMapKeyFunc EntityToMapKeyFunc[T]
	data               []T

	// Acts as a set for quick existence checks but no real nead for a value
	existingSet map[interface{}]struct{}

	// Statistics tracking
	injectedCount  int
	existingCount  int
	totalProcessed int
}

func NewInjector[T any](data []T, keyFunc EntityToMapKeyFunc[T]) *Injector[T] {
	return &Injector[T]{
		data:               data,
		entityToMapKeyFunc: keyFunc,
		existingSet:        make(map[interface{}]struct{}),
	}
}

func (i *Injector[T]) Preload() error {
	existingEntities, err := database.GetAllEntities[T]()
	if err != nil {
		logger.Errorf("Failed to fetch existing entities (%T): %v", new(T), err)
		return err
	}

	for _, entity := range existingEntities {
		key := i.entityToMapKeyFunc(entity)
		i.existingSet[key] = struct{}{}
	}

	logger.Infof("Preloaded %d entities of type %T.", len(i.existingSet), new(T))
	return nil
}

func (i *Injector[T]) Process() error {
	for _, newItem := range i.data {
		i.totalProcessed++
		key := i.entityToMapKeyFunc(newItem)

		if _, exists := i.existingSet[key]; exists {
			i.existingCount++
			logger.Infof("Entity already exists (key: %v). Skipping injection.", key)
			continue
		}

		_, err := database.InsertEntityById(&newItem)
		if err != nil {
			logger.Errorf("Failed to insert entity (%T): %+v, Error: %v", new(T), newItem, err)
			return err
		}

		i.injectedCount++
		logger.Infof("Successfully injected entity: %+v", newItem)
	}

	return nil
}

func (i *Injector[T]) logInjectionStats() {
	logger.Infof("----------- Injection stats for type %T -----------", new(T))
	logger.Infof("Total processed entities: %d", i.totalProcessed)
	logger.Infof("Total already existing entities: %d", i.existingCount)
	logger.Infof("Total injected entities: %d", i.injectedCount)
}
