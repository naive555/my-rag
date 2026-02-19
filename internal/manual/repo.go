package manual

import "context"

type Repo interface {
	FindByDomain(ctx context.Context, domain string, limit int) ([]string, error)
}
