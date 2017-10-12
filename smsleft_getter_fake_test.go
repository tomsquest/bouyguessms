package bouyguessms

type fakeSmsGetter struct {
	smsLeft SmsLeft
	err     error
}

func newFakeSmsGetter(smsLeft int, err error) *fakeSmsGetter {
	return &fakeSmsGetter{SmsLeft(smsLeft), err}
}

func (fake *fakeSmsGetter) Get() (SmsLeft, error) {
	return fake.smsLeft, fake.err
}
