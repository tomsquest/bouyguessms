package bouyguessms

type Msg string

const maxLength = 160

func (m Msg) String() string {
	if len(m) > maxLength {
		return string(m)[0:maxLength]
	}
	return string(m)
}
