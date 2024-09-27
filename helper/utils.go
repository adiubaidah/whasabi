package helper

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

func PanicIfError(message string, err error) {
	if err != nil {
		if message != "" {
			fmt.Println(message, err)
		}
		panic(err)
	}
}

func ReadFromRequestBody(request *http.Request, result any) {
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(result)
	PanicIfError("", err)
}

func WriteToResponseBody(writer http.ResponseWriter, response any) {
	writer.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(response)
	PanicIfError("", err)
}
func GetEnv(key string) string {
	return viper.GetString(key)
}

//get user from
