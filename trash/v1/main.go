package main

// import "./flame"

func Hello() {
    for i := 0; i < 1e7; i++ {}
}

func main() {
    Record(Hello, "hello.svg")
}