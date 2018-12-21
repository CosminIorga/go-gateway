package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Hosts []struct {
		Host   string `json:"host"`
		Routes []struct {
			Route string `json:"route"`
		}
	} `json:"hosts"`
}

//Define a new structure that represents out API response (response status and body)
type HTTPResponse struct {
	status string
	body   interface{}
}

type ForwardPath struct {
	mainPath   string
	secondPath string
}

func main() {
	router := gin.Default()

	router.Any("*forwardPath", func(c *gin.Context) {
		forwardPath := splitForwardPath(c)

		ok := checkRouteAvailable(forwardPath.mainPath)

		if !ok {
			c.JSON(422, gin.H{
				"message": "Unauthorized",
			})

			return
		}

		c.JSON(200, forwardPath)
	})

	/**
	router.GET("/", func(c *gin.Context) {
		//Define a new channel
		var ch chan HTTPResponse = make(chan HTTPResponse)
		//List of APIs to call
		urls := [2]string{"https://jsonplaceholder.typicode.com/posts/1", "https://jsonplaceholder.typicode.com/posts/1/comments"}
		for _, url := range urls {
			go DoHTTPGet(url, ch)
		}

		response := gin.H{}

		// Get the response
		for _, url := range urls {
			response[url] = (<-ch).body
		}

		c.JSON(200, response)
	})
	*/
	router.Run() // listen and serve on 0.0.0.0:8080
}

func checkRouteAvailable(route string) bool {

	config := loadConfiguration("./config/config.json")
	for _, host := range config.Hosts {
		fmt.Println(host.Host)
		for _, routes := range host.Routes {
			fmt.Println(route)
			fmt.Println(host.Host+routes.Route)
			
			if route == host.Host+routes.Route {
				return true
			}
		}

	}
	return false
}

func splitForwardPath(c *gin.Context) ForwardPath {
	forwardPath, _ := c.Params.Get("forwardPath")

	splitPath := strings.SplitN(forwardPath, "/", 3)

	return ForwardPath{
		splitPath[1],
		splitPath[2],
	}
}

func DoHTTPGet(url string, ch chan<- HTTPResponse) {
	//Execute the HTTP get
	httpResponse, _ := http.Get(url)
	var httpBody interface{}
	err := json.NewDecoder(httpResponse.Body).Decode(&httpBody)

	if err != nil {
		panic(err)
	}
	// httpBody, _ := ioutil.ReadAll(httpResponse.Body)
	//Send an HTTPResponse back to the channel
	ch <- HTTPResponse{httpResponse.Status, httpBody}
}

func loadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}
