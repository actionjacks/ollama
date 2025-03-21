package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "App running",
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

		// Prepare the request body for Ollama API
		requestBody, err := json.Marshal(map[string]interface{}{
			"prompt": prompt,
			"model":  "llama3:8b", // Replace with your model name
			"stream": false,       // Disable streaming for simplicity
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
