Mastersteam
=======

Mastersteam is a web api for querying servers from the Valve Master Server List. 
it's lightweight and fast

Valve's master server has a rate limit of about 15 queries per minute, and returns a batch of ~220 servers for each query. For a popular game, it can take a long time (around ten minutes) to retrieve its entire server list. Mastersteam will query individual game servers in the background to lessen the overall waiting time. At the moment it will process 20 servers in the background, concurrently. Go's scheduler is still weak so it's not recommended to use more.

Usage
-----
You can use Mastersteam either across all of Half-Life 1 or Half-Life 2, or with a specific list of Application IDs. For a full list of Application IDs, see: https://developer.valvesoftware.com/wiki/Steam_Application_IDs

To start the webservice on port 8080
```
$ go get github.com/TowelSoftware/Mastersteam
$ go build
$ Mastersteam
```

You can also use this it in docker.
https://github.com/TowelSoftware/Mastersteam-docker

You can now search through the browser, Curl or an from an web app

#### Search on server name limited by appid
`http://localhost:8080/search/[APP_ID]/[NAME]`

#### Search by serer ip
`http://localhost:8080/server/[IP]`

```
curl "http://localhost:8080/search/107410/*Bohemia%20Interactive*"

Would give you something like this.
{
	"data" : [{
	"85.190.155.160:2403": {
		"ip": "85.190.155.160:2403",
		"protocol": 17,
		"name": "\ufffd [ OFFICIAL ] Arma 3 Vanguard by Bohemia Interactive (EU) #02",
		"map": "Tanoa",
		"folder": "Arma3",
		"game": "Vanguard 50 Power Plant",
		"players": 2,
		"max_players": 50,
		"bots": 0,
		"type": "dedicated",
		"os": "windows",
		"visibility": "public",
		"vac": false,
		"appid": 107410,
		"game_version": "1.90.145471",
		"port": 2402,
		"steamid": "90124885451686921",
		"game_mode": "bt,r190,n145381,s3,i1,mf,lf,vt,dt,tvanguar,g65545,h87a3a791,f0,c-2147483648--2147483648,pw,e0,j0,k0,",
		"gameid": "107410",
		"players_online": [
			{
				"Name": "Smith",
				"Score": 4294967291,
				"Duration": 2401.2808
			},
			{
				"Name": "jonas",
				"Score": 87,
				"Duration": 2358.7737
			}
		]
	},
	"85.190.155.59:2303": {
		"ip": "85.190.155.59:2303",
		"protocol": 17,
		"name": "\ufffd [ OFFICIAL ] Arma 3 EndGame by Bohemia Interactive (EU) #01",
		"map": "Tanoa",
		"folder": "Arma3",
		"game": "End Game 24 Balavu",
		"players": 0,
		"max_players": 28,
		"bots": 0,
		"type": "dedicated",
		"os": "windows",
		"visibility": "public",
		"vac": false,
		"appid": 107410,
		"game_version": "1.90.145471",
		"port": 2302,
		"steamid": "90124884470531080",
		"game_mode": "bt,r190,n145381,s3,i0,mf,lf,vt,dt,tendgame,g65545,h87a3a791,f0,c14-50,pw,e0,j0,k0,",
		"gameid": "107410"
	}
	...
	}],
	"total":72
}
```

```
curl "http://localhost:8080/server/85.190.158.12"

{
	"data" : [{
	"85.190.158.12:2303": {
		"ip": "85.190.158.12:2303",
		"protocol": 17,
		"name": "\ufffd [ OFFICIAL ] Arma 3 CP by Bohemia Interactive (USA) #01",
		"map": "Malden",
		"folder": "Arma3",
		"game": "Escape 10 Malden",
		"players": 0,
		"max_players": 10,
		"bots": 0,
		"type": "dedicated",
		"os": "windows",
		"visibility": "public",
		"vac": false,
		"appid": 107410,
		"game_version": "1.90.145471",
		"port": 2302,
		"steamid": "90124960057373704",
		"game_mode": "bt,r190,n145381,s7,i2,mf,lf,vt,dt,tescape,g65545,h87a3a791,f0,c-2147483648--2147483648,pw,e15,j0,k0,",
		"gameid": "107410"
	},
	"85.190.158.12:2403": {
		"ip": "85.190.158.12:2403",
		"protocol": 17,
		"name": "\ufffd [ OFFICIAL ] Arma 3 CP by Bohemia Interactive (USA) #02",
		"map": "Malden",
		"folder": "Arma3",
		"game": "Combat Patrol",
		"players": 0,
		"max_players": 12,
		"bots": 0,
		"type": "dedicated",
		"os": "windows",
		"visibility": "public",
		"vac": false,
		"appid": 107410,
		"game_version": "1.90.145471",
		"port": 2402,
		"steamid": "90124961543184388",
		"game_mode": "bt,r190,n145381,s7,i1,mf,lf,vt,dt,tpatrol,g65545,h87a3a791,f0,c-2147483648--2147483648,pw,e15,j0,k0,",
		"gameid": "107410"
	}}],
	"total":2
}

```

Building
--------

1. Make sure you have Golang installed, (see: http://golang.org/)
2. Make sure your Go environment is set up. Example:

        export GOROOT=~/tools/go
        export GOPATH=~/go
        export PATH="$PATH:$GOROOT/bin:$GOPATH/bin"

3. Get the source code and its dependencies:

        go get https://github.com/TowelSoftware/Mastersteam

4. Build:

        go install

5. The `Mastersteam` binary wll be in `$GOPATH/bin/`.

Resources
---------
https://developer.valvesoftware.com/wiki/Master_Server_Query_Protocol \
https://developer.valvesoftware.com/wiki/Server_queries

Init code and insparation
---------
https://github.com/alliedmodders/blaster \
https://github.com/rumblefrog/go-a2s/
