package main

//go:generate yarn --cwd ./sfts-web run build
func main() {
    StartServer(7464)
}
