package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	tableWidth   = 1200
	tableHeight  = 600
	paddleWidth  = 20
	paddleHeight = 120
	ballRadius   = 10
	paddleSpeed  = 10
	maxScore     = 11
)

type Vec2 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Ball struct {
	Pos    Vec2    `json:"pos"`
	Vel    Vec2    `json:"vel"`
	Radius float64 `json:"radius"`
}

type Paddle struct {
	Y      float64 `json:"y"`
	Height float64 `json:"height"`
	Width  float64 `json:"width"`
}

type GameState struct {
	Ball        Ball   `json:"ball"`
	LeftPaddle  Paddle `json:"leftPaddle"`
	RightPaddle Paddle `json:"rightPaddle"`
	LeftScore   int    `json:"leftScore"`
	RightScore  int    `json:"rightScore"`
	Paused      bool   `json:"paused"`
	GameOver    bool   `json:"gameOver"`
	Winner      string `json:"winner"`
	GameMode    string `json:"gameMode"`
	Difficulty  string `json:"difficulty"`
	InMenu      bool   `json:"inMenu"`
	mu          sync.Mutex
}

var game *GameState

func newGame() *GameState {
	rand.Seed(time.Now().UnixNano())
	return &GameState{
		Ball: Ball{
			Pos:    Vec2{X: tableWidth / 2, Y: tableHeight / 2},
			Vel:    Vec2{X: 6, Y: 4},
			Radius: ballRadius,
		},
		LeftPaddle: Paddle{
			Y:      tableHeight/2 - paddleHeight/2,
			Height: paddleHeight,
			Width:  paddleWidth,
		},
		RightPaddle: Paddle{
			Y:      tableHeight/2 - paddleHeight/2,
			Height: paddleHeight,
			Width:  paddleWidth,
		},
		LeftScore:  0,
		RightScore: 0,
		Paused:     false,
		GameOver:   false,
		GameMode:   "ai",
		Difficulty: "medium",
		InMenu:     true,
	}
}

func (g *GameState) reset() {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.Ball.Pos = Vec2{X: tableWidth / 2, Y: tableHeight / 2}
	direction := 1.0
	if g.RightScore > g.LeftScore {
		direction = -1.0
	}
	g.Ball.Vel = Vec2{X: 6 * direction, Y: 4}
	g.LeftPaddle.Y = tableHeight/2 - paddleHeight/2
	g.RightPaddle.Y = tableHeight/2 - paddleHeight/2
}

func (g *GameState) resetGame() {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.LeftScore = 0
	g.RightScore = 0
	g.GameOver = false
	g.Winner = ""
	g.Paused = false
	g.Ball.Pos = Vec2{X: tableWidth / 2, Y: tableHeight / 2}
	g.Ball.Vel = Vec2{X: 6, Y: 4}
	g.LeftPaddle.Y = tableHeight/2 - paddleHeight/2
	g.RightPaddle.Y = tableHeight/2 - paddleHeight/2
}

func (g *GameState) updateAI() {
	if g.GameMode != "ai" || g.Paused || g.GameOver {
		return
	}

	targetY := g.Ball.Pos.Y
	paddleCenter := g.RightPaddle.Y + g.RightPaddle.Height/2

	var aiSpeed float64
	var reactionDelay float64

	switch g.Difficulty {
	case "easy":
		aiSpeed = 4
		reactionDelay = 50
	case "medium":
		aiSpeed = 7
		reactionDelay = 20
	case "hard":
		aiSpeed = 10
		reactionDelay = 5
	}

	if math.Abs(targetY-paddleCenter) > reactionDelay {
		if targetY > paddleCenter {
			g.RightPaddle.Y = math.Min(tableHeight-g.RightPaddle.Height, g.RightPaddle.Y+aiSpeed)
		} else {
			g.RightPaddle.Y = math.Max(0, g.RightPaddle.Y-aiSpeed)
		}
	}
}

func (g *GameState) update() {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.Paused || g.GameOver || g.InMenu {
		return
	}

	g.Ball.Pos.X += g.Ball.Vel.X
	g.Ball.Pos.Y += g.Ball.Vel.Y

	if g.Ball.Pos.Y-g.Ball.Radius <= 0 || g.Ball.Pos.Y+g.Ball.Radius >= tableHeight {
		g.Ball.Vel.Y = -g.Ball.Vel.Y
		g.Ball.Pos.Y = math.Max(g.Ball.Radius, math.Min(tableHeight-g.Ball.Radius, g.Ball.Pos.Y))
	}

	if g.Ball.Pos.X-g.Ball.Radius <= paddleWidth {
		if g.Ball.Pos.Y >= g.LeftPaddle.Y && g.Ball.Pos.Y <= g.LeftPaddle.Y+g.LeftPaddle.Height {
			relativeY := (g.Ball.Pos.Y - (g.LeftPaddle.Y + g.LeftPaddle.Height/2)) / (g.LeftPaddle.Height / 2)
			angle := relativeY * math.Pi / 3
			speed := math.Sqrt(g.Ball.Vel.X*g.Ball.Vel.X + g.Ball.Vel.Y*g.Ball.Vel.Y)
			speed *= 1.08
			g.Ball.Vel.X = speed * math.Cos(angle)
			g.Ball.Vel.Y = speed * math.Sin(angle)
			g.Ball.Pos.X = paddleWidth + g.Ball.Radius
		}
	}

	if g.Ball.Pos.X+g.Ball.Radius >= tableWidth-paddleWidth {
		if g.Ball.Pos.Y >= g.RightPaddle.Y && g.Ball.Pos.Y <= g.RightPaddle.Y+g.RightPaddle.Height {
			relativeY := (g.Ball.Pos.Y - (g.RightPaddle.Y + g.RightPaddle.Height/2)) / (g.RightPaddle.Height / 2)
			angle := relativeY * math.Pi / 3
			speed := math.Sqrt(g.Ball.Vel.X*g.Ball.Vel.X + g.Ball.Vel.Y*g.Ball.Vel.Y)
			speed *= 1.08
			g.Ball.Vel.X = -speed * math.Cos(angle)
			g.Ball.Vel.Y = speed * math.Sin(angle)
			g.Ball.Pos.X = tableWidth - paddleWidth - g.Ball.Radius
		}
	}

	if g.Ball.Pos.X < 0 {
		g.RightScore++
		if g.RightScore >= maxScore {
			g.GameOver = true
			if g.GameMode == "ai" {
				g.Winner = "Computer Wins!"
			} else {
				g.Winner = "Right Player Wins!"
			}
		} else {
			g.mu.Unlock()
			g.reset()
			g.mu.Lock()
		}
	} else if g.Ball.Pos.X > tableWidth {
		g.LeftScore++
		if g.LeftScore >= maxScore {
			g.GameOver = true
			if g.GameMode == "ai" {
				g.Winner = "You Win!"
			} else {
				g.Winner = "Left Player Wins!"
			}
		} else {
			g.mu.Unlock()
			g.reset()
			g.mu.Lock()
		}
	}
}

func (g *GameState) movePaddle(paddle string, direction string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if paddle == "left" {
		if direction == "up" {
			g.LeftPaddle.Y = math.Max(0, g.LeftPaddle.Y-paddleSpeed)
		} else {
			g.LeftPaddle.Y = math.Min(tableHeight-g.LeftPaddle.Height, g.LeftPaddle.Y+paddleSpeed)
		}
	} else if paddle == "right" {
		if direction == "up" {
			g.RightPaddle.Y = math.Max(0, g.RightPaddle.Y-paddleSpeed)
		} else {
			g.RightPaddle.Y = math.Min(tableHeight-g.RightPaddle.Height, g.RightPaddle.Y+paddleSpeed)
		}
	}
}

func handleState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	game.mu.Lock()
	defer game.mu.Unlock()
	json.NewEncoder(w).Encode(game)
}

func handleMove(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Paddle    string `json:"paddle"`
		Direction string `json:"direction"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	game.movePaddle(req.Paddle, req.Direction)
	w.WriteHeader(http.StatusOK)
}

func handlePause(w http.ResponseWriter, r *http.Request) {
	game.mu.Lock()
	game.Paused = !game.Paused
	game.mu.Unlock()
	w.WriteHeader(http.StatusOK)
}

func handleReset(w http.ResponseWriter, r *http.Request) {
	game.resetGame()
	w.WriteHeader(http.StatusOK)
}

func handleStartGame(w http.ResponseWriter, r *http.Request) {
	var req struct {
		GameMode   string `json:"gameMode"`
		Difficulty string `json:"difficulty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	game.mu.Lock()
	game.GameMode = req.GameMode
	game.Difficulty = req.Difficulty
	game.InMenu = false
	game.mu.Unlock()
	game.resetGame()
	w.WriteHeader(http.StatusOK)
}

func handleBackToMenu(w http.ResponseWriter, r *http.Request) {
	game.mu.Lock()
	game.InMenu = true
	game.mu.Unlock()
	game.resetGame()
	w.WriteHeader(http.StatusOK)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Ping Pong Game</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            min-height: 100vh;
            color: white;
            overflow: hidden;
        }
        #mainMenu {
            background: rgba(0,0,0,0.8);
            padding: 50px;
            border-radius: 20px;
            text-align: center;
            box-shadow: 0 20px 60px rgba(0,0,0,0.5);
            max-width: 600px;
        }
        #mainMenu h1 {
            font-size: 56px;
            margin-bottom: 40px;
            color: #FFD700;
            text-shadow: 3px 3px 6px rgba(0,0,0,0.5);
        }
        .menu-section {
            margin: 30px 0;
        }
        .menu-section h2 {
            font-size: 24px;
            margin-bottom: 15px;
            color: #ADD8E6;
        }
        .button-group {
            display: flex;
            gap: 15px;
            justify-content: center;
            flex-wrap: wrap;
        }
        .menu-btn {
            padding: 15px 30px;
            font-size: 18px;
            cursor: pointer;
            border: 3px solid transparent;
            border-radius: 10px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            transition: all 0.3s;
            font-weight: bold;
            min-width: 150px;
        }
        .menu-btn:hover {
            transform: scale(1.1);
            box-shadow: 0 5px 25px rgba(255,255,255,0.3);
        }
        .menu-btn.selected {
            border-color: #FFD700;
            background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
        }
        #startBtn {
            margin-top: 40px;
            padding: 20px 60px;
            font-size: 24px;
            background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
        }
        #gameArea {
            display: none;
        }
        #controlPanel {
            background: rgba(0,0,0,0.7);
            padding: 20px 30px;
            border-radius: 15px;
            margin-bottom: 20px;
            display: flex;
            gap: 20px;
            align-items: center;
            flex-wrap: wrap;
            justify-content: center;
        }
        button {
            padding: 12px 25px;
            font-size: 16px;
            cursor: pointer;
            border: none;
            border-radius: 8px;
            background: #4CAF50;
            color: white;
            transition: all 0.3s;
            font-weight: bold;
        }
        button:hover { background: #45a049; transform: scale(1.05); }
        button:active { transform: scale(0.95); }
        #pauseBtn { background: #ff9800; }
        #pauseBtn:hover { background: #e68900; }
        #menuBtn { background: #9C27B0; }
        #menuBtn:hover { background: #7B1FA2; }
        #score {
            font-size: 36px;
            font-weight: bold;
            text-shadow: 2px 2px 4px rgba(0,0,0,0.5);
            min-width: 120px;
        }
        #gameContainer {
            position: relative;
            box-shadow: 0 15px 50px rgba(0,0,0,0.6);
            border-radius: 15px;
            overflow: hidden;
            border: 5px solid rgba(255,255,255,0.2);
        }
        canvas {
            display: block;
            background: #0a4d2e;
        }
        #gameOver {
            position: absolute;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -50%);
            background: rgba(0,0,0,0.95);
            padding: 50px;
            border-radius: 20px;
            text-align: center;
            display: none;
            border: 3px solid #FFD700;
        }
        #gameOver h1 {
            font-size: 52px;
            margin-bottom: 30px;
            color: #FFD700;
            animation: pulse 2s infinite;
        }
        @keyframes pulse {
            0%, 100% { transform: scale(1); }
            50% { transform: scale(1.05); }
        }
        #controls {
            margin-top: 20px;
            background: rgba(0,0,0,0.7);
            padding: 20px;
            border-radius: 15px;
            text-align: center;
        }
        #controls p { margin: 8px 0; font-size: 16px; }
        .control-key {
            background: rgba(255,255,255,0.2);
            padding: 5px 10px;
            border-radius: 5px;
            font-weight: bold;
        }
    </style>
</head>
<body>
    <div id="mainMenu">
        <h1>🏓 PING PONG</h1>
        <div class="menu-section">
            <h2>Game Mode</h2>
            <div class="button-group">
                <button class="menu-btn selected" onclick="selectMode('ai')">
                    🤖 vs Computer
                </button>
                <button class="menu-btn" onclick="selectMode('2player')">
                    👥 2 Players
                </button>
            </div>
        </div>
        <div class="menu-section" id="difficultySection">
            <h2>Difficulty</h2>
            <div class="button-group">
                <button class="menu-btn" onclick="selectDifficulty('easy')">
                    😊 Easy
                </button>
                <button class="menu-btn selected" onclick="selectDifficulty('medium')">
                    😐 Medium
                </button>
                <button class="menu-btn" onclick="selectDifficulty('hard')">
                    😈 Hard
                </button>
            </div>
        </div>
        <button id="startBtn" class="menu-btn" onclick="startGame()">
            ▶️ Start Game
        </button>
    </div>
    <div id="gameArea">
        <div id="controlPanel">
            <div id="score">0 : 0</div>
            <button id="pauseBtn" onclick="togglePause()">⏸️ Pause</button>
            <button id="menuBtn" onclick="backToMenu()">🏠 Menu</button>
        </div>
        <div id="gameContainer">
            <canvas id="canvas" width="1200" height="600"></canvas>
            <div id="gameOver">
                <h1 id="winnerText"></h1>
                <button onclick="playAgain()" style="font-size: 22px; padding: 15px 40px;">
                    🔄 Play Again
                </button>
                <button onclick="backToMenu()" style="font-size: 22px; padding: 15px 40px; margin-left: 15px;">
                    🏠 Menu
                </button>
            </div>
        </div>
        <div id="controls">
            <p id="controlsText"></p>
        </div>
    </div>
    <script>
        const canvas = document.getElementById('canvas');
        const ctx = canvas.getContext('2d');
        const keys = {};
        let selectedMode = 'ai';
        let selectedDifficulty = 'medium';
        window.addEventListener('keydown', e => keys[e.key.toLowerCase()] = true);
        window.addEventListener('keyup', e => keys[e.key.toLowerCase()] = false);
        function selectMode(mode) {
            selectedMode = mode;
            document.querySelectorAll('.menu-section')[0].querySelectorAll('.menu-btn').forEach(btn => {
                btn.classList.remove('selected');
            });
            event.target.classList.add('selected');
            document.getElementById('difficultySection').style.display =
                mode === 'ai' ? 'block' : 'none';
        }
        function selectDifficulty(difficulty) {
            selectedDifficulty = difficulty;
            document.querySelectorAll('.menu-section')[1].querySelectorAll('.menu-btn').forEach(btn => {
                btn.classList.remove('selected');
            });
            event.target.classList.add('selected');
        }
        async function startGame() {
            await fetch('/start', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({
                    gameMode: selectedMode,
                    difficulty: selectedDifficulty
                })
            });
            document.getElementById('mainMenu').style.display = 'none';
            document.getElementById('gameArea').style.display = 'block';
            updateControlsText();
        }
        function updateControlsText() {
            const text = selectedMode === 'ai'
                ? '<p><strong>Controls:</strong> <span class="control-key">W</span> (Up) / <span class="control-key">S</span> (Down)</p><p>The first player to reach 11 points wins!</p>'
                : '<p><strong>Left Player:</strong> <span class="control-key">W</span> (Up) / <span class="control-key">S</span> (Down)</p><p><strong>Right Player:</strong> <span class="control-key">↑</span> (Up) / <span class="control-key">↓</span> (Down)</p><p>The first player to reach 11 points wins!</p>';
            document.getElementById('controlsText').innerHTML = text;
        }
        async function backToMenu() {
            await fetch('/menu', {method: 'POST'});
            document.getElementById('mainMenu').style.display = 'block';
            document.getElementById('gameArea').style.display = 'none';
            document.getElementById('gameOver').style.display = 'none';
        }
        async function playAgain() {
            await fetch('/reset', {method: 'POST'});
            document.getElementById('gameOver').style.display = 'none';
        }
        async function movePaddle(paddle, direction) {
            await fetch('/move', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({paddle, direction})
            });
        }
        async function togglePause() {
            await fetch('/pause', {method: 'POST'});
        }
        function drawTable() {
            ctx.fillStyle = '#0a4d2e';
            ctx.fillRect(0, 0, canvas.width, canvas.height);
            ctx.strokeStyle = 'rgba(255,255,255,0.3)';
            ctx.lineWidth = 4;
            ctx.strokeRect(0, 0, canvas.width, canvas.height);
            ctx.setLineDash([15, 15]);
            ctx.lineWidth = 3;
            ctx.strokeStyle = 'rgba(255,255,255,0.5)';
            ctx.beginPath();
            ctx.moveTo(canvas.width/2, 0);
            ctx.lineTo(canvas.width/2, canvas.height);
            ctx.stroke();
            ctx.setLineDash([]);
        }
        function draw(state) {
            drawTable();
            ctx.fillStyle = '#2196F3';
            ctx.shadowBlur = 20;
            ctx.shadowColor = '#2196F3';
            ctx.fillRect(0, state.leftPaddle.y, state.leftPaddle.width, state.leftPaddle.height);
            ctx.fillStyle = '#F44336';
            ctx.shadowColor = '#F44336';
            ctx.fillRect(canvas.width - state.rightPaddle.width, state.rightPaddle.y,
                        state.rightPaddle.width, state.rightPaddle.height);
            ctx.shadowBlur = 25;
            ctx.shadowColor = '#FFFF00';
            ctx.fillStyle = 'white';
            ctx.beginPath();
            ctx.arc(state.ball.pos.x, state.ball.pos.y, state.ball.radius, 0, Math.PI * 2);
            ctx.fill();
            ctx.shadowBlur = 0;
            document.getElementById('score').textContent =
                state.leftScore + ' : ' + state.rightScore;
            document.getElementById('pauseBtn').innerHTML =
                state.paused ? '▶️ Resume' : '⏸️ Pause';
            if (state.gameOver) {
                document.getElementById('winnerText').textContent = state.winner;
                document.getElementById('gameOver').style.display = 'block';
            }
        }
        async function gameLoop() {
            const res = await fetch('/state');
            const state = await res.json();
            if (!state.inMenu) {
                if (keys['w']) await movePaddle('left', 'up');
                if (keys['s']) await movePaddle('left', 'down');
                if (state.gameMode === '2player') {
                    if (keys['arrowup']) await movePaddle('right', 'up');
                    if (keys['arrowdown']) await movePaddle('right', 'down');
                }
                draw(state);
            }
            requestAnimationFrame(gameLoop);
        }
        gameLoop();
    </script>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, html)
}

func gameLoop() {
	ticker := time.NewTicker(16 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		game.updateAI()
		game.update()
	}
}

func readInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func main() {
	game = newGame()
	go gameLoop()

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/state", handleState)
	http.HandleFunc("/move", handleMove)
	http.HandleFunc("/pause", handlePause)
	http.HandleFunc("/reset", handleReset)
	http.HandleFunc("/start", handleStartGame)
	http.HandleFunc("/menu", handleBackToMenu)

	fmt.Println("🏓 Ping Pong Game Server")
	fmt.Println("========================")

	useHTTPS := readInput("Use HTTPS? (yes/no): ")

	if strings.ToLower(useHTTPS) == "yes" || strings.ToLower(useHTTPS) == "y" {
		certFile := readInput("Certificate file path (cert.pem): ")
		keyFile := readInput("Private key file path (key.pem): ")

		if certFile == "" || keyFile == "" {
			log.Fatal("Error: Certificate and key files are required")
		}

		portInput := readInput("Port (press Enter for 443): ")
		port := ":443"
		if portInput != "" {
			port = ":" + portInput
		}

		tlsConfig := &tls.Config{
			MinVersion: tls.VersionTLS12,
		}

		server := &http.Server{
			Addr:      port,
			TLSConfig: tlsConfig,
		}

		fmt.Printf("✅ Server running with HTTPS on https://localhost%s\n", port)
		fmt.Println("Open your browser and start playing!")

		log.Fatal(server.ListenAndServeTLS(certFile, keyFile))
	} else {
		portInput := readInput("Port (press Enter for 80): ")
		port := ":80"
		if portInput != "" {
			port = ":" + portInput
		}

		fmt.Printf("✅ Server running with HTTP on http://localhost%s\n", port)
		fmt.Println("Open your browser and start playing!")

		log.Fatal(http.ListenAndServe(port, nil))
	}
}
