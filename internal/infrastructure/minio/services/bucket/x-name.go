package bucket

type Name string

func (n Name) String() string {
	return string(n)
}

const (
	Voices Name = "voices"
)
