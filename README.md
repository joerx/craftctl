# Minecraft Server Control Panel

Run a Minecraft server:

```sh
java -Xmx1024M -Xms1024M -jar server.jar nogui
```

Hot reload using [nodemon](https://www.npmjs.com/package/nodemon):

```sh
nodemon -e go,html --exec go run main.go --signal SIGTERM
```
