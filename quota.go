package bouyguessms

type Quota int

const ExceededQuota = Quota(0)

func (quota Quota) IsExceeded() bool {
	return int(quota) <= 0
}

func (quota Quota) Remaining() int {
	return int(quota)
}
