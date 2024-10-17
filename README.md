# raindrop-sdk-go

Unofficial [Raindrop.io](https://raindrop.io/) Go SDK.

## Create an Access Token
Follow the [offiical raindrop.io authentication documentation](https://developer.raindrop.io/v1/authentication/token) to create an access token, which will be used by the SDK.

## Installation
```
go get github.com/jiachengxu/raindrop-sdk-go@v0.1.0
```
## Examples
```golang
package main

import (
	"fmt"

	"github.com/jiachengxu/raindrop-sdk-go"
)

func main() {
	client := raindrop.NewClient("your_test_token_here")

	// Example: Get collections
	collections, err := client.GetRootCollections()
	if err != nil {
		fmt.Println("Error fetching collections:", err)
		return
	}

	fmt.Println(collections)
}
```