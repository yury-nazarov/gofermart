package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// JSONError400 проверяет кейсы для возврата HTTP 400 (http.StatusBadRequest)
// 			    Передаем HTTP Request и пустой интерфейс который
//			    по сути является любой структурой для анмаршала JSON из HTTP Request
func  JSONError400(r *http.Request, anyData interface{}, logger *log.Logger)  error {
	bodyData , err := io.ReadAll(r.Body)
	if err != nil || len(bodyData) == 0 {
		logger.Printf("the HTTP Body parsing error: %s", err)
		return err
	}

	// Unmarshal JSON
	if err = json.Unmarshal(bodyData, &anyData); err != nil {
		logger.Printf("unmarshal json error: %s", err)
		return err
	}
	return nil
}
