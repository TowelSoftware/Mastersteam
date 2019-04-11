// Licensed under the GNU General Public License, version 3 or higher.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	batch "github.com/TowelSoftware/Mastersteam/batch"
	valve "github.com/TowelSoftware/Mastersteam/valve"
)

var (
	sOutputBuffer bytes.Buffer
	sNumServers   int64
	master        *valve.MasterServerQuerier
)

/*
ErrorObject ...
*/
type ErrorObject struct {
	IP    string `json:"ip"`
	Error string `json:"error"`
}

/*
ServerObject ...
*/
type ServerObject struct {
	Address     string      `json:"ip"`
	Protocol    uint8       `json:"protocol"`
	Name        string      `json:"name"`
	MapName     string      `json:"map"`
	Folder      string      `json:"folder"`
	Game        string      `json:"game"`
	Players     uint8       `json:"players"`
	MaxPlayers  uint8       `json:"max_players"`
	Bots        uint8       `json:"bots"`
	Type        string      `json:"type"`
	Os          string      `json:"os"`
	Visibility  string      `json:"visibility"`
	Vac         bool        `json:"vac"`
	AppID       valve.AppId `json:"appid,omitempty"`
	GameVersion string      `json:"game_version,omitempty"`
	Port        uint16      `json:"port,omitempty"`
	SteamID     string      `json:"steamid,omitempty"`
	GameMode    string      `json:"game_mode,omitempty"`
	GameID      string      `json:"gameid,omitempty"`

	PlayersOnline []*valve.Player `json:"players_online,omitempty"`
}

func addJSON(hostAndPort string, obj interface{}) {
	buf, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}

	var indented bytes.Buffer
	json.Indent(&indented, buf, "\t", "\t")

	if sNumServers != 0 {
		sOutputBuffer.WriteString(",")
	}

	sOutputBuffer.WriteString(fmt.Sprintf("\n\t\"%s\": ", hostAndPort))
	sOutputBuffer.WriteString(indented.String())

	sNumServers++
}

func addError(hostAndPort string, err error) {
	addJSON(hostAndPort, &ErrorObject{
		IP:    hostAndPort,
		Error: err.Error(),
	})
}

/*
Log ...
*/
func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("access: %s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func httpMasterSearch(w http.ResponseWriter, r *http.Request) {
	uriSegments := strings.Split(r.URL.String(), "/")
	appID, _ := strconv.Atoi(uriSegments[2])
	hostname, _ := url.QueryUnescape(uriSegments[3])

	newMasterServerQuerier()

	// Set up the filter list.
	master.FilterAppId(valve.AppId(appID))
	master.FilterName(hostname)

	newServerQuerier()

	//defer master.Close()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	fmt.Fprintf(w, "%s", sOutputBuffer.String())
}

func httpServer(w http.ResponseWriter, r *http.Request) {
	uriSegments := strings.Split(r.URL.String(), "/")
	host, _ := url.QueryUnescape(uriSegments[2])

	newMasterServerQuerier()

	master.FilterGameaddr(host)

	newServerQuerier()

	//defer master.Close()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	fmt.Fprintf(w, "%s", sOutputBuffer.String())
}

func newMasterServerQuerier() {
	m, err := valve.NewMasterServerQuerier(valve.MasterServer)
	if err != nil {
		log.Printf("Could not query master: %s", err.Error())
	}
	master = m
	//defer m.Close()
}

func newServerQuerier() {
	flagTimeout := time.Second * 3
	flagJ := 20
	sNumServers = 0

	sOutputBuffer.Reset()

	bp := batch.NewBatchProcessor(func(item interface{}) {
		addr := item.(*net.TCPAddr)
		query, err := valve.NewServerQuerier(addr.String(), flagTimeout)
		if err != nil {
			addError(addr.String(), err)
			return
		}
		defer query.Close()

		info, err := query.QueryInfo()
		if err != nil {
			addError(addr.String(), err)
			return
		}

		log.Printf("%s - %s\n", addr.String(), info.Name)

		out := &ServerObject{
			Address:    addr.String(),
			Protocol:   info.Protocol,
			Name:       info.Name,
			MapName:    info.MapName,
			Folder:     info.Folder,
			Game:       info.Game,
			Players:    info.Players,
			MaxPlayers: info.MaxPlayers,
			Bots:       info.Bots,
			Type:       info.Type.String(),
			Os:         info.OS.String(),
		}
		if info.Vac == 1 {
			out.Vac = true
		}
		if info.Visibility == 0 {
			out.Visibility = "public"
		} else {
			out.Visibility = "private"
		}
		if info.Ext != nil {
			out.AppID = info.Ext.AppId
			out.GameVersion = info.Ext.GameVersion
			out.Port = info.Ext.Port
			out.SteamID = fmt.Sprintf("%d", info.Ext.SteamId)
			out.GameMode = info.Ext.GameModeDescription
			out.GameID = fmt.Sprintf("%d", info.Ext.GameId)
		}

		if info.Players > 0 {
			players, err := query.QueryPlayers()
			if err != nil {
				out.PlayersOnline = nil
			} else {
				out.PlayersOnline = players
			}
		}

		addJSON(addr.String(), out)
	}, flagJ)

	defer bp.Terminate()

	// TOP OF JSON FILE
	sOutputBuffer.WriteString("{\n")
	sOutputBuffer.WriteString("\t\"data\" : [{")

	// Query the master.
	err := master.Query(func(servers valve.ServerList) error {
		bp.AddBatch(servers)
		return nil
	})

	if err != nil {
		log.Printf("Could not query the master: %s\n", err.Error())
		os.Exit(1)
	}

	// Wait for batch processing to complete.
	bp.Finish()

	if sNumServers != 0 {
		//sOutputBuffer.WriteString("\n")
	}

	sOutputBuffer.WriteString("}],\n")
	sOutputBuffer.WriteString(fmt.Sprintf("\t\"total\":%d\n", sNumServers))
	sOutputBuffer.WriteString("}\n")
	//BOTTOM OF JSON FILE

}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	http.HandleFunc("/search/", httpMasterSearch)
	http.HandleFunc("/server/", httpServer)
	log.Fatal(http.ListenAndServe(":8080", Log(http.DefaultServeMux)))
}
