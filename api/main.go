package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Define the request body structure
type OllamaRequest struct {
	Prompt string `json:"prompt"`
	Model  string `json:"model"`
	Stream bool   `json:"stream"`
}

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		panic("Failed to load .env file")
	}

	// Get environment variables
	ollamaModel := os.Getenv("OLLAMA_MODEL")
	ollamaAPIURL := os.Getenv("OLLAMA_API_URL")

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "App running",
			"data": gin.H{
				"OLLAMA_MODEL":   ollamaModel,
				"OLLAMA_API_URL": ollamaAPIURL,
			},
		})
	})

	r.POST("/ask", func(c *gin.Context) {
		prompt := c.Query("prompt")
		if prompt == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Prompt is required",
			})
			return
		}

		// Create the request body using the struct
		requestBody, err := json.Marshal(OllamaRequest{
			Prompt: prompt,
			Model:  ollamaModel, // Use the model from .env
			Stream: false,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create request body",
			})
			return
		}

		// Call Ollama API
		response, err := http.Post("http://ollamadeepseek:11436/api/generate", "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to call Ollama API: " + err.Error(),
			})
			return
		}
		defer response.Body.Close()

		// Read the response body
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to read Ollama API response: " + err.Error(),
			})
			return
		}

		if response.StatusCode != http.StatusOK {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":        "Ollama API returned an error",
				"status_code":  response.StatusCode,
				"raw_response": string(responseBody),
			})
			return
		}

		// Parse the response body (assuming it's JSON)
		var ollamaResponse map[string]interface{}
		if err := json.Unmarshal(responseBody, &ollamaResponse); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":        "Failed to parse Ollama API response: " + err.Error(),
				"raw_response": string(responseBody), // Include raw response for debugging
			})
			return
		}

		// Return the Ollama API response
		c.JSON(http.StatusOK, gin.H{
			"response": ollamaResponse,
		})
	})

	r.Run(":8080") // Listen on port 8080
}
