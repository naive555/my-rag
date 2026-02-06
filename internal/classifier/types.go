package classifier

type Domain string

const (
	DomainSMS      Domain = "sms"
	DomainCampaign Domain = "campaign"
	DomainContact  Domain = "contact"
	DomainBilling  Domain = "billing"
	DomainGeneral  Domain = "general"
)

type Policy struct {
	UseRAG       bool
	TopK         int
	MaxSentences int
}

type Result struct {
	Domain Domain
	TopK   int
}
