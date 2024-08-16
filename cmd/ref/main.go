package main

//go:generate veil

// @d:service
type Foo struct {
}

func (f *Foo) Beep() error {
	return nil
}

// @d:service
type Bar struct {
}

func (f *Bar) Boop() {

}
