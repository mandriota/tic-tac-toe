# Tic-Tac-Toe
 tic-tac-toe multiplayer game

https://user-images.githubusercontent.com/62650188/228930784-f7605f7b-3f61-4a06-8a1c-91b62b0399a3.mp4


### Run locally
First, start the server with:
```bash
go run cmd/server/main.go
```

Then start two clients, each with:
```bash
go run cmd/client/main.go
```

### Structure
* cmd
    * client: minimal client implementation
    * server: minimal server implementation
* internal
    * client/ui: UI model
    * server: handle clients requests and distributes into game rooms (each room for 2 players)
* pkg/world: base logic
