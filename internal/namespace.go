package internal

type Namespace string

func (n Namespace) String() string {
	return string(n)
}
