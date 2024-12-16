package discovery

type Broadcast interface {
	PublicAddrs() []string
	DeliverOnline(...string) error
	Reload()
}
