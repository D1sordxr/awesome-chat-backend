package apps

import "fmt"

type AppEnv int8

const (
	Unknown AppEnv = iota
	Local
	Development
	Staging
	Production
)

var envNames = map[AppEnv]string{
	Local:       "local",
	Development: "development",
	Staging:     "staging",
	Production:  "production",
}

func ParseAppEnv(s string) (AppEnv, error) {
	for env, name := range envNames {
		if name == s {
			return env, nil
		}
	}

	return Unknown, fmt.Errorf("unknown app env %q", s)
}

func (e AppEnv) String() string {
	if name, ok := envNames[e]; ok {
		return name
	}

	return "unknown"
}

func (e AppEnv) IsProduction() bool {
	return e == Production
}
