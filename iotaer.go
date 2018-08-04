package main

const (
	ewfw = iota
	a
	b
	c
	d
	e
	f
	g
	h
	i
	j
	k
	l
	m
	n
	o
	p
	q
	r
	s
	t
	u
	v
	w
	x
	y
	z
)

var joiner string = "\n"

func main() {
	print(a, joiner, b, joiner, c, joiner, d, joiner, e, joiner, f, joiner)
	print(g, joiner, h, joiner, i, joiner, j, joiner, k, joiner)
	print(l, joiner, m, joiner, n, joiner, o, joiner, p, joiner, q, joiner)
	print(r, joiner, s, joiner, t, joiner, u, joiner, v, joiner)
	print(w, joiner, x, joiner, y, joiner, z, joiner)
}
