package apps

type AppEnv int8

const (
	Local AppEnv = iota
	Development
	Staging
	Production
)

func New(s string) AppEnv {
	switch s {
	case "local":
		return Local
	case "development":
		return Development
	case "staging":
		return Staging
	case "production":
		return Production
	}
}
