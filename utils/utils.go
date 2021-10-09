package utils


import (

    "net/http"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"encoding/hex"
	"strings"


)



func GetHash(pwd string) string{        
	sha256Bytes := sha256.Sum256([]byte(pwd))
	return hex.EncodeToString(sha256Bytes[:])

}





func IdFromUrl(r *http.Request) (string, error) {
	parts := strings.Split(r.URL.String(), "/")
	// fmt.Println("hello",parts)
	if len(parts) > 3{

		if strings.Contains(parts[3], "?"){
			id := strings.Split(parts[3], "?")
			fmt.Println("hi",id[0])
			return id[0], nil
		}
		id := parts[3]
		return id, nil
	}
	id := parts[2]
	return id, nil
}


func RespondWithError(w http.ResponseWriter, code int, msg string) {
	RespondWithJSON(w, code, map[string]string{"error": msg})
}

func RespondWithJSON(w http.ResponseWriter, code int, data interface{}) {
	response, _ := json.Marshal(data)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}


