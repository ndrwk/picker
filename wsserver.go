package picker

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Command struct {
	CommandName string
}

var upgrader = websocket.Upgrader{}
var logger *log.Logger

func runWSServer(ip string, pickerLog *log.Logger) error {
	logger = pickerLog
	http.HandleFunc("/picker", echo)
	http.HandleFunc("/", home)
	logger.Println("Serve ", config.Device.Ip)
	http.ListenAndServe(config.Device.Ip, nil)
	return nil
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		// mt, message, err := c.ReadMessage()
		var cmd Command
		err := c.ReadJSON(&cmd)
		if err != nil {
			logger.Println("read:", err)
			// break
		}
		logger.Printf("recv: %s", cmd.CommandName)
		logger.Println(c.RemoteAddr())
		message := Command{
			CommandName: "123",
		}
		err = c.WriteJSON(message)
		if err != nil {
			logger.Println("write:", err)
			// break
		}
		logger.Printf("sent: %s", message)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/picker")
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>
window.addEventListener("load", function(evt) {

    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;

    var print = function(message) {
        var d = document.createElement("div");
        d.innerHTML = message;
        output.appendChild(d);
    };

    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };

    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };

    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };

});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="{}">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))
