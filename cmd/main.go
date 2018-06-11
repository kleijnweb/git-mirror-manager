package main

func main() {
  container := &Container{}
  server := container.Server()
  server.Start()
}
