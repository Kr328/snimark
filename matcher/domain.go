package matcher

import (
	"net"

	"github.com/Dreamacro/clash/component/trie"
)

type DomainMatcher struct {
	domains *trie.DomainTrie
}

func (d *DomainMatcher) Match(host string) bool {
	if net.ParseIP(host) != nil {
		return false
	}

	return d.domains.Search(host) != nil
}

func NewDomain(domains []string) (Matcher, error) {
	t := trie.New()

	for _, domain := range domains {
		if err := t.Insert(domain, struct{}{}); err != nil {
			return nil, err
		}
	}

	return &DomainMatcher{domains: t}, nil
}
