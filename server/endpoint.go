package server

import (
	"fmt"

	"github.com/kr328/snimark/config"
	"github.com/kr328/snimark/marker"
	"github.com/kr328/snimark/matcher"
	"github.com/kr328/snimark/sniffer"
)

type Endpoint struct {
	sniffer sniffer.Sniffer
	matcher matcher.Matcher
	marker  marker.Marker
}

func ParseEndpoints(cfg *config.Config) (map[uint16]*Endpoint, error) {
	r := map[string]marker.Marker{}
	e := map[uint16]*Endpoint{}

	for name, mark := range cfg.Markers {
		if len(mark.Marks) == 0 {
			return nil, fmt.Errorf("parse %s: empty marks", name)
		}

		switch mark.Mode {
		case "mono":
			r[name] = marker.NewMono(name, mark.Marks[0])
		case "round":
			r[name] = marker.NewRound(name, mark.Marks)
		}
	}

	for port, endpoint := range cfg.Endpoints {
		s, ok := sniffer.Sniffers[endpoint.Sniffer]
		if !ok {
			return nil, fmt.Errorf("%d: unsupported sniffer %s", port, endpoint.Sniffer)
		}
		cm, ok := matcher.Matchers[endpoint.Sniffer]
		if !ok {
			return nil, fmt.Errorf("%d: unsupported sniffer %s", port, endpoint.Sniffer)
		}

		m, err := cm(endpoint.Match)
		if err != nil {
			return nil, fmt.Errorf("%d: invalid match %s", port, err.Error())
		}

		r, ok := r[endpoint.Marker]
		if !ok {
			return nil, fmt.Errorf("%d: marker %s not found", port, endpoint.Marker)
		}

		e[port] = &Endpoint{
			sniffer: s,
			matcher: m,
			marker:  r,
		}
	}

	return e, nil
}
