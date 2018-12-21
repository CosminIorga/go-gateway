package main

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

//Define a new structure that represents out API response (response status and body)
type HTTPResponse struct {
	status string
	body   interface{}
}

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		//Define a new channel
		var ch chan HTTPResponse = make(chan HTTPResponse)
		//List of APIs to call
		urls := [2]string{"https://jsonplaceholder.typicode.com/posts/1", "https://jsonplaceholder.typicode.com/posts/1/comments"}
		for _, url := range urls {
			//For each URL call the DOHTTPGet function (notice the go keyword)
			go DoHTTPGet(url, ch)
		}

		response := gin.H{}

		// Get the response
		for _, url := range urls {
			response[url] = (<-ch).body
		}

		c.JSON(200, response)

		// postUrls := "https://jsonplaceholder.typicode.com/posts"

		// var myPostParam []map[string]string
		// value1 := map[string]string{"title": "test1", "body": "body1", "userId": "1"}
		// value2 := map[string]string{"title": "test2", "body": "body2", "userId": "2"}

		// myPostParam = append(myPostParam, value1, value2)

		// for _, value := range myPostParam {
		// 	//For each URL call the DOHTTPPost function (notice the go keyword)
		// 	go DoHTTPPost(postUrls, value, ch)
		// }

		// for range myPostParam {
		// 	// Use the response (<-ch).body
		// 	fmt.Println((<-ch).status)
		// }
	})

	// router.GET("/blocks/*path", func(c *gin.Context) {
	// 	forwardPath := c.Param("path")
	// 	forwardQuery := c.Request.URL.Query()

	// 	c.JSON(200, gin.H{
	// 		"path":  forwardPath,
	// 		"query": forwardQuery,
	// 	})
	// })

	// router.Any("/*name/*path", funct(c *gin.Context) {
	//   name := c.Param("name")
	//   path := c.Param("path")
	//   query := c.Request.URL.Query()

	//   reverseproxy.NewReverseProxy()

	// })

	router.Run() // listen and serve on 0.0.0.0:8080
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
