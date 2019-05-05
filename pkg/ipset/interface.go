package ipset

type Interface interface {
	Create(setname string, t SetType) error
	Add(setname string, entry string) error
	Delete(setname string, entry string) error
	List(setname string) error
	Swap(oldname string, newname string) error
	Destroy(setname string) error
}

type Set struct {
	Name     string
	Type     string
	Revision int
	Header   SetHeader
	Members  []string
}

type SetHeader struct {
	Family      string
	HashSize    int
	MaxElements int
	MemorySize  int
	References  int
}
