# ws-game
Number guessing game created with websocket in golang

Setup Local Environment

For the first time setup
1. Install Go (https://golang.org/dl/)
2. Use your preferred IDE (e.g. Visual Studio Code) and install the required plugin [Go golang.go]
3. Create a folder for your golang project
4. Use Terminal -> Run "go mod init example.com/ws-game" to init golang project
5. Create a golang file (e.g. "test.go") to trigger the installation of required plugin for developing in golang

Setup and run ws-game
1. Download the project ws-game to local
2. Run "go get github.com/gorilla/websocket" to install the websocket plugin with terminal in Visual Studio Code
3. To start the websocket game server, run "go run main.go"
4. Use websocket test client to connect the game server (https://chrome.google.com/webstore/detail/websocket-test-client/fgponpodhbmadfljofbimhhlengambbn/related)
5. Insert "ws://localhost:8089/ws" as the URL
6. Submit the following request for starting the game
- Registration
{"message":"29c9c30e0604515ced98b3d14fd88751a8f8e4b9bc69d483a67a257c14ab79fb","playerName":"Geane","timestamp":1615140324}

- Guess
{"message":"f1abe1b083d12d181ae136cfc75b8d18a8ecb43ac4e9d1a36d6a9c75b6016b61","guess":123,"timestamp":1615140124,"gameId":1}
