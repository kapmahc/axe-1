package axe

type route struct {
	method   string
	path     string
	handlers []HandlerFunc
}

type routes []route

func (p routes) Len() int {
	return len(p)
}

func (p routes) Less(i, j int) bool {
	if p[i].method < p[j].method {
		return true
	}
	if p[i].method == p[j].method {
		return p[i].path < p[j].path
	}
	return false
}

func (p routes) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
