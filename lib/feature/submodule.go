package feature

type SubModuleNewFunc[T any] func(parentModule T) interface{}

func NewSubModules[T any](module T, newFuncArr ...SubModuleNewFunc[T]) []interface{} {
	return AppendSubModules([]interface{}{}, module, newFuncArr...)
}

func AppendSubModules[T any](modules []interface{}, parentModule T, newFuncArr ...SubModuleNewFunc[T]) []interface{} {
	if len(newFuncArr) > 0 {
		for _, newFunc := range newFuncArr {
			modules = append(modules, newFunc(parentModule))
		}
	}
	return modules
}
