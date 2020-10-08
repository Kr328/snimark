package matcher

type Matcher interface {
	Match(host string) bool
}

type ConstructMatcher func(payload []string) (Matcher, error)

var Matchers = map[string]ConstructMatcher{
	"http": NewDomain,
	"tls":  NewDomain,
	"tcp":  NewAll,
}
