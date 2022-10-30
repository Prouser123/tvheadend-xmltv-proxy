package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
	_ "github.com/joho/godotenv/autoload" // Import .env file if one exists
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var tvhServer string
var tvhServerAuthUrl string

func ensureEnv() {
	// Check for TVH_SERVER env variable.
	var ok bool
	tvhServer, ok = os.LookupEnv("TVH_SERVER")
	if !ok {
		log.Fatal("TVH_SERVER is not present") // log.Fatal will terminate the program.
	} else {
		fmt.Printf("Found tvh server at: '%s'\n", tvhServer)
	}

	// Check for TVH_AUTH env variable
	creds, ok2 := os.LookupEnv("TVH_SERVER_BASIC_AUTH")
	if !ok2 {
		fmt.Println("TVH_SERVER_BASIC_AUTH is not present, assuming auth not needed")
		tvhServerAuthUrl = tvhServer
	} else {
		// Split tvhServer at the http://
		protoSplit := strings.Split(tvhServer, "://")
		fmt.Printf("Found server with protocol '%s' and url '%s'\n", protoSplit[0], protoSplit[1])
		tvhServerAuthUrl = fmt.Sprintf("%s://%s@%s", protoSplit[0], creds, protoSplit[1])
	}

	fmt.Printf("Authenticated URL: %s\n", tvhServerAuthUrl)

}

func main() {
	ensureEnv()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/epg.xml", epg)

	e.Logger.Fatal(e.Start(":8080"))
}

func epg(c echo.Context) error {
	// Step 1: Request the xmltv page from tvheadend
	client := resty.New()

	resp, err := client.R().Get(fmt.Sprintf("%s/xmltv/channels", tvhServerAuthUrl))

	if err != nil {
		return c.NoContent(500)
	} else {
		xmlData := resp.String()
		// Now we need to do string replacement on imagecache stuff.
		oldStr := fmt.Sprintf("%s/imagecache/", tvhServer)
		newStr := fmt.Sprintf("%s/imagecache/", c.Request().Host)
		fmt.Printf("Replacing '%s' with '%s'", oldStr, newStr)
		new := strings.ReplaceAll(xmlData, oldStr, newStr)

		return c.String(http.StatusOK, new)
	}
}
