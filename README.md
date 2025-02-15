# Server Game | Pong
Primarily server code but under all that is a version of Pong that is able to be played online.

---

## Building
1. Make sure you have the latest version of Golang installed along with RayLib dependencies.
2. Next clone the git repo `git clone https://github.com/Crumbed/server_game`.
3. If you are on Windows you may need to place `raylib.dll` in the `server_game` directory.
4. Finally, while in the `server_game` directory, run the command `go build -o pong`. This will generate an executable called `pong`.
---

## Running & Playing
- On **Linux** run `./pong`.
- On **Windows** run `start pong.exe`.
- First player to connect to the server will be player 1, and the second will be player 2.
- You move your paddle using `UP` & `DOWN` , `K` & `J`, or `W` & `S`.
---

## Server Hosting
I don't have any public game servers since this is just a toy project so you will have to host your own. 
- Find an open port on your network. If you do not know what this means or how to do it, watch [this](https://www.youtube.com/watch?v=WOZQppVNGvA) video.
- If you've ever hosted a Minecraft server on your own computer you've likely had to set up port forwarding for `:25565`, this is the port I used for testing since I already had it open.
- To host a server, use the same run command from before but add `server :OPEN_PORT`. If no port is provided, it will default to `3000`.
- When you want to stop the server, use the keybind `ctrl+c`.

### Server Settings
After you've run the server once a new file, `server_rules.json` will be generated with the default values.
- If you edit this file make sure that the JSON formatting is correct. I have no error handling for invalid JSON. If you are unsure of the formatting rules, you can use [this](https://jsonlint.com/) website.
- `initial_velocity` is the initial "speed" of the puck. This is `300` by default.
- `increment_velocity` is the amount of "speed" the puck gains after hitting a paddle. This is `100` by default.
