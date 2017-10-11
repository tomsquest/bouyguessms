package bouyguessms

// Quota is the amount of remaining SMS
type Quota int

// ExceededQuota is the quota when there is no remaining SMS
const ExceededQuota = Quota(0)

// IsExceeded returns true when there is no remaining SMS
func (quota Quota) IsExceeded() bool {
	return int(quota) <= 0
}

// Remaining returns the amount of remaining SMS as an integer
func (quota Quota) Remaining() int {
	return int(quota)
}
