package web

type Module interface {
	SetupToWeb(router Router)
}
