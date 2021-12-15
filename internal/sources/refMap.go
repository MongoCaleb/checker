package sources

type RefMap map[string]Ref

type Ref struct {
	Target string
	Type   string
}

func (r *RefMap) Get(key string) (Ref, bool) {
	val, ok := (*r)[key]
	return val, ok
}
