package bouyguessms

// SmsLeft is the amount of remaining SMS
type SmsLeft int

// NoSmsLeft represents an exceeded count of SMS
const NoSmsLeft = SmsLeft(0)

// IsExceeded returns true when there is no remaining SMS
func (smsLeft SmsLeft) IsExceeded() bool {
	return int(smsLeft) <= 0
}

// Remaining returns the amount of remaining SMS as an integer
func (smsLeft SmsLeft) Remaining() int {
	return int(smsLeft)
}
