package settings

type SettingsChangedTriggerBase[T any] interface {
	OnChanged(settings T)
}
