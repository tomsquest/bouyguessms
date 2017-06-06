package bouyguessms

type fakeSmsSender struct {
	err    error
	called bool
}

func (fake *fakeSmsSender) SendSms(msg string, to string) error {
	fake.called = true
	return fake.err
}
