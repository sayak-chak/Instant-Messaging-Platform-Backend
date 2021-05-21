package mocks

import "encoding/json"

func GetLoginJson(username string, password string) ([]byte, error){
	jsonString := map[string]interface{}{
		"username": username,
		"password": password,
	}
	return json.Marshal(jsonString)
}