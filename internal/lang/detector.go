package lang

const (
	EN  string = "en"
	TH  string = "th"
	UNK string = "unknown"
)

func Detect(text string) string {
	for _, r := range text {
		if r >= 0x0E00 && r <= 0x0E7F {
			return TH
		}
	}
	return EN
}
