package email

// SenderService is the package-level instance of the Sender.
var SenderService Sender

// NewSenderService sets the Sender for the package.
func NewSenderService(s Sender) {
	SenderService = s
}
