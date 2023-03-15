package main

func Run() error {
	return nil
}

func main() {
	if err := Run(); err != nil {
		// oh no
		panic(err.Error())
	}
}
