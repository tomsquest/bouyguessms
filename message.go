package bouyguessms

type message string

const maxLength = 160

func (m message) String() string {
	if len(m) > maxLength {
		return string(m)[0:maxLength]
	}
	return string(m)
}
