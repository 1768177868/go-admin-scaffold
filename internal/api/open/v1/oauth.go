package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GithubOAuth handles the GitHub OAuth login request
func GithubOAuth(c *gin.Context) {
	// TODO: Implement GitHub OAuth login
	c.JSON(http.StatusOK, gin.H{"message": "GitHub OAuth login"})
}

// GithubOAuthCallback handles the GitHub OAuth callback
func GithubOAuthCallback(c *gin.Context) {
	// TODO: Implement GitHub OAuth callback
	c.JSON(http.StatusOK, gin.H{"message": "GitHub OAuth callback"})
}
