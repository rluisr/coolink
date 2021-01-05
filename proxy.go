package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func getImage(c *gin.Context) {
	targetURL := c.Query("q")
	if targetURL == "" {
		c.String(http.StatusBadRequest, "?q=<target url> is required")
		return
	}

	resp, err := http.Get(targetURL)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.String(http.StatusBadRequest, fmt.Sprintf("status code is %d", resp.StatusCode))
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	contentType := http.DetectContentType(body)

	c.Data(http.StatusOK, contentType, body)
}
