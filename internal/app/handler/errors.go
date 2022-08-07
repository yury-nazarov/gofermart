package handler

import (
	"encoding/json"
	"io"
	"net/http"
)

// JSONError400 проверяет кейсы для возврата HTTP 400 (http.StatusBadRequest)
// 			    Передаем HTTP Request и пустой интерфейс который
//			    по сути является любой структурой для анмаршала JSON из HTTP Request
func  (c *Controller) JSONError400(r *http.Request, anyData interface{})  error {
	bodyData , err := io.ReadAll(r.Body)
	if err != nil || len(bodyData) == 0 {
		c.logger.Printf("the HTTP Body parsing error: %s", err)
		return err
	}

	// Unmarshal JSON
	if err = json.Unmarshal(bodyData, &anyData); err != nil {
		c.logger.Printf("unmarshal json error: %s", err)
		return err
	}
	return nil
}
