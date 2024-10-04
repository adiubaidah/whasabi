package helper

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adiubaidah/wasabi/model"
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

func WriteToResponseBody(writer http.ResponseWriter, response *model.WebResponse) {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(response.Code)

	encoder := json.NewEncoder(writer)
	err := encoder.Encode(response)
	PanicIfError("Error encoding response", err)
}

//get user from
