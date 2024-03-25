package game

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"text/template"

	"github.com/gorilla/websocket"
	"github.com/sessionsdev/blue-octopus/internal/auth"
)

type HtmxWebSocketMsg struct {
	Prompt  string   `json:"prompt"`
	Headers struct{} `json:"HEADERS"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeGamePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/game-websocket.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandleGameWebSocket(w http.ResponseWriter, r *http.Request) {
	// Get the user from the context
	user := r.Context().Value("user")
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Upgrade initial GET request to a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Make sure we close the connection when the function returns
	defer ws.Close()

	for {
		// Read the message from the browser
		messageType, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// unmarshal the message into a socket response
		var socketResponse HtmxWebSocketMsg
		err = json.Unmarshal(msg, &socketResponse)
		if err != nil {
			log.Println(err)
			return
		}

		// handle the message
		err = handleMessage(r.Context(), socketResponse, ws, messageType)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func handleMessage(ctx context.Context, msg HtmxWebSocketMsg, ws *websocket.Conn, messageType int) error {
	log.Println("Received message:", msg)

	// Get user form context
	user := ctx.Value("user")
	if user == nil {
		return fmt.Errorf("Unauthorized")
	}

	userEmail := user.(*auth.User).Email
	var g *Game

	switch msg.Prompt {
	case "RESET GAME":
		g = InitializeNewGame()
		html := wrapGameOutputDiv(`
			<p>[GAME MASTER]: New game created!</p>
		`)
		writeMessageAndHandleError(ws, messageType, []byte(html))
		SaveGameToRedis(ctx, g, userEmail)
		return nil

	case "":
		html := wrapGameOutputDiv(`
				<p>[ERROR]: Please enter a prompt</p>
			`)

		writeMessageAndHandleError(ws, messageType, []byte(html))
		return nil

	default:
		// echo back a formatted message
		// Reconcile the game state
		// Save the game state
		// Write message back to browser
		return handlePlayerPrompt(ctx, msg, ws, messageType, g, userEmail)
	}
}

func handlePlayerPrompt(ctx context.Context, msg HtmxWebSocketMsg, ws *websocket.Conn, messageType int, g *Game, userEmail string) error {
	// Get the game from redis
	g, err := LoadGameFromRedis(ctx, userEmail)
	if err != nil {
		return err
	}

	if g.Processing {
		html := wrapGameOutputDiv(`
			<p>[ERROR]: game command processing is already in progress. Please wait a moment and try again.</p>
		`)
		writeMessageAndHandleError(ws, messageType, []byte(html))
		return nil
	}

	g.Processing = true

	html := wrapGameOutputDiv(fmt.Sprintf(`
			<p>[PLAYER]: %s </p>
		`, msg.Prompt))

	writeMessageAndHandleError(ws, messageType, []byte(html))

	narrativeResponse, err := g.ProcessPlayerPrompt(string(msg.Prompt))
	if err != nil {
		html := wrapGameOutputDiv(fmt.Sprintf(`
				<p>[ERROR]: %s</p>
			`, err.Error()))

		writeMessageAndHandleError(ws, messageType, []byte(html))
		return nil
	}

	done := make(chan bool)

	go func() {

		writeMessageAndHandleError(ws, messageType, []byte(formatStatsPanelHtml(`<img src="static/svg-loaders/puff.svg" alt="">`)))
		g.reconcileGameState()
		g.populatePreparedStatsCache()
		html, err := getTemplateHtml("templates/stats-panel.html", "stats-panel", PreparedStatsCache)
		if err != nil {
			log.Println(err)
		}

		writeMessageAndHandleError(ws, messageType, []byte(formatStatsPanelHtml(html)))
		done <- true
	}()

	go func() {

		g.progressStoryThreads()
		done <- true
	}()

	go func() {
		<-done
		<-done
		g.Processing = false
		SaveGameToRedis(ctx, g, userEmail)
	}()

	html = wrapGameOutputDiv(fmt.Sprintf(`<p>[GAME MASTER]: %s</p>`, narrativeResponse))

	writeMessageAndHandleError(ws, messageType, []byte(html))

	return nil
}

var writelock sync.Mutex

func writeMessageAndHandleError(ws *websocket.Conn, messageType int, data []byte) {
	writelock.Lock()
	defer writelock.Unlock()

	err := ws.WriteMessage(messageType, data)
	if err != nil {
		log.Println("Failed to write message:", err)

		// Close the connection if it's still open
		if closeErr := ws.Close(); closeErr != nil {
			log.Println("Failed to close WebSocket:", closeErr)
		}
		return
	}
}

func wrapGameOutputDiv(html string) string {
	return fmt.Sprintf(`
	<div id="game-output" hx-swap-oob="beforeend">
		%s
	</div>
	`, html)
}

func formatStatsPanelHtml(html string) string {
	return fmt.Sprintf(`
	<div id="game-state-panel" hx-swap-oob="innerHTML">
		%s
	</div>
	`, html)
}

func getTemplateHtml(templateFile string, templateName string, data interface{}) (string, error) {
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err := tmpl.ExecuteTemplate(&tpl, templateName, data); err != nil {
		return "", err
	}

	return tpl.String(), nil
}
