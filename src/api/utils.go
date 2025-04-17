package api

import (
	"bytes"
	"docryte/src/types"
	"fmt"
	"log"
	"net/http"
)

func notifyNewContact(cr *types.ContactRequest, token string, userId string) {
	addr := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	json := fmt.Appendf(nil, `{"chat_id":"%s","text":"%s"}`, userId, cr.String())
	req, _ := http.NewRequest("POST", addr, bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err.Error())
	}
	defer resp.Body.Close()

	fmt.Println("Notifying in Telegram Response Code:", resp.Status)
}
