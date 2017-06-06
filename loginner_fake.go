package bouyguessms

type fakeLoginner struct {
	err            error
	hasCalledLogin bool
}

func (l *fakeLoginner) Login(login, pass string) error {
	l.hasCalledLogin = true
	return l.err
}
