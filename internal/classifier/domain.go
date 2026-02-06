package classifier

import "strings"

func Classify(q string) Result {
	q = strings.ToLower(q)

	switch {
	case containsAny(q, "sms", "ส่ง", "ข้อความ", "otp", "delivery", "dr", "failed"):
		return Result{DomainSMS, 5}
	case containsAny(q, "campaign", "แคมเปน", "send", "broadcast", "marketing", "การตลาด", "schedule", "ตั้งเวลา"):
		return Result{DomainCampaign, 3}
	case containsAny(q, "contact", "group", "กลุ่ม", "number", "เบอร์", "โทรศัพท์"):
		return Result{DomainContact, 3}
	case containsAny(q, "plan", "package", "quota", "limit", "แพ็ค", "เงิน", "ตังค์", "โควต้า"):
		return Result{DomainBilling, 1}
	default:
		return Result{DomainGeneral, 0}
	}
}

func containsAny(q string, keywords ...string) bool {
	for _, k := range keywords {
		if strings.Contains(q, k) {
			return true
		}
	}
	return false
}
