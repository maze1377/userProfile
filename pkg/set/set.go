package set

type Set interface {
	Add(values ...string)
	Remove(value string) bool
	Has(value string) bool
	Size() int
	Clear()
	Items() []string
}
