# Server Game | Pong
Primarily server code but under all that is a janky version of Pong that is able to be played online.

## Building
1. Make sure you have the latest version of Golang installed along with RayLib dependencies.
2. Next clone the git repo `git clone https://github.com/Crumbed/server_game`
3. If you are on Windows you may need to place `raylib.dll` in the `server_game` directory.

## Running
- You cannot run the client without connecting to a server.
- To host a server, run the command `go run . server :PORT`. If no port is provided, it will default to `3000`.
- To connect to a server, run the command `go run . IP:PORT`. If no IP is provided, it will default to the last server you connected to.
- Whenever the server closes, all connected clients will close. And whenever a client closes, the server will close.
- First player to connect to the server will be player 1, and the second will be player 2.
- You move your paddle using `UP` & `DOWN` arrow keys, or `K` & `J` keys.
