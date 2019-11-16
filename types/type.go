package types

const (
	// connection stat
	StatConnected    = "connected"
	StatDisconnected = "disconnected"
	// connection use type
	// connection only for message
	UseTypeMessage string = "msg"
	// connection only for stream
	UseTypeStream string = "str"
	// connection only can be used for message and stream
	UseTypeShare string = "shr"	
)
