package bouyguessms

type fakeQuotaGetter struct {
	quota Quota
	err   error
}

func newFakeQuotaGetter(quota int, err error) *fakeQuotaGetter {
	return &fakeQuotaGetter{Quota(quota), err}
}

func (fake *fakeQuotaGetter) Get() (Quota, error) {
	return fake.quota, fake.err
}
