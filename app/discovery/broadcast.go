// Define discovery broadcast
package discovery

type Broadcast interface {
	// Serve Public Addrs
	PublicAddrs() []string
	// DeliverOnline send a broadcast message to the online peers.
	//
	// The message will be sent to the peers whose addresses are in the parameter list.
	// If the parameter list is empty, the message will be sent to all online peers.
	DeliverOnline(...string) error
	// Reload reloads the broadcast module.
	//
	// It is called by the runtime to reload the broadcast module after the application has finished initializing.
	Reload()
}
