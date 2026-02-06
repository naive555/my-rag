package classifier

func Resolve(domain Domain) Policy {
	switch domain {
	case DomainSMS:
		return Policy{
			UseRAG:       true,
			TopK:         5,
			MaxSentences: 3,
		}
	case DomainCampaign:
		return Policy{
			UseRAG:       true,
			TopK:         5,
			MaxSentences: 3,
		}
	case DomainContact:
		return Policy{
			UseRAG:       true,
			TopK:         4,
			MaxSentences: 3,
		}
	case DomainBilling:
		return Policy{
			UseRAG:       true,
			TopK:         1,
			MaxSentences: 2,
		}
	default:
		return Policy{
			UseRAG:       false,
			MaxSentences: 3,
		}
	}
}
