package main

//Post blah blh
type Post struct {
	Name  string
	Views int
	Date  int
}

//Posts blah blh
type Posts []Post

func (slice Posts) Len() int {
	return len(slice)
}

func (slice Posts) Less(i, j int) bool {
	return slice[i].Date > slice[j].Date
}

func (slice Posts) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
