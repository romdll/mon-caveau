package injector

import "moncaveau/database"

var AllToInject = []Injectable{
	NewGenericWrapper(
		regions,
		func(region database.WineRegion) interface{} {
			return region.Name + "_" + region.Country
		},
	),
	NewGenericWrapper(
		types,
		func(wineType database.WineType) interface{} {
			return wineType.Name
		},
	),
	NewGenericWrapper(
		bottleSizes,
		func(size database.WineBottleSize) interface{} {
			return size.Size
		},
	),
}

func SetupAndInjectAll() error {
	for _, injectable := range AllToInject {
		if err := injectable.Preload(); err != nil {
			return err
		}
		if err := injectable.Process(); err != nil {
			return err
		}
	}

	for _, injectable := range AllToInject {
		injectable.Display()
	}

	return nil
}
