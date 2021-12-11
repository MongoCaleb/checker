package sources

type RefMap map[string]string

func (r *RefMap) Get(key string) (string, bool) {
	val, ok := (*r)[key]
	return val, ok
}
