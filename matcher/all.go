package matcher

type AllMatcher struct{}

func (a *AllMatcher) Match(string) bool {
	return true
}

func NewAll([]string) (Matcher, error) {
	return &AllMatcher{}, nil
}
