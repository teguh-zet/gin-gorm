package websocket

import (
	"fmt"
	"gin-gonic/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Mengizinkan koneksi dari domain mana saja (CORS)
	CheckOrigin: func(r *http.Request) bool {
		return true 
	},
}

func ServeWS(manager *Manager, c *gin.Context) {
	// 1. Ambil token dari Query Param: ws://localhost:8080/ws?token=...
	tokenString := c.Query("token")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
		return
	}

	// 2. Validasi Token menggunakan fungsi utils Anda
	// Fungsi ini sudah sekaligus mengecek validitas token dan mengambil User ID
	userIDUint, err := utils.GetUserIDFromToken(tokenString)
	if err != nil {
		// Jika error, berarti token tidak valid atau expired
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// 3. Upgrade koneksi HTTP ke WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	// 4. Daftarkan Client ke Manager
	// Kita ubah userID dari uint ke string agar seragam di struct Client
	client := &Client{
		Manager: manager, 
		Conn:    conn, 
		Send:    make(chan []byte, 256),
		UserID:  fmt.Sprintf("%d", userIDUint), // Konversi uint ke string
	}

	client.Manager.Register <- client

	// Jalankan routine baca & tulis
	go client.WritePump()
	go client.ReadPump()
}