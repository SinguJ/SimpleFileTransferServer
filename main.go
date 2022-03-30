package main

//go:generate yarn --cwd ./sfts-web build
func main() {
    StartServer(7464)
}
