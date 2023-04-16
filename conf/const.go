package conf

// config paramaters for daemon and antisctl.
const (
	// Address of the API server. Not specifying the port will result in
	// an error.
	//
	// Eg: mycontrol.example.com:443
	Address = "localhost:2023"

	// true results in more verbose log messages.
	Debug = true

	// CommonSecret is used to encrypt sensitive data in transit perhaps
	// to avoid detection.
	CommonSecret = "5124e446534843af7eb935a9b8b369954d456849c5f27202300cf660cafd8952"
)
