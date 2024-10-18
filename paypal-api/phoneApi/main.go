package phoneapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	url := "http://<your-server-ip>:5000/send-sms"
	apiKey := "your_api_key"
	payload := map[string]string{
		"number":  "+491234567890",
		"message": "Your verification code is: 123456",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("API-Key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("SMS successfully sent")
	} else {
		fmt.Printf("Failed to send SMS. Status code: %d\n", resp.StatusCode)
	}
}

/*
# *SMS Doğrulama API Dokümantasyonu*

Bu API, kullanıcılara doğrulama SMS mesajları göndermenizi sağlar. Yetkili istemciler, kimlik doğrulama için benzersiz bir API anahtarı kullanarak API'ye istekler gönderebilir.

## *Uç Nokta*:
⁠ POST /send-sms ⁠

### *Base URL*:
⁠ http://<your-server-ip>:5000/send-sms ⁠

## *Kimlik Doğrulama*:
İsteği doğrulamak için istek başlıklarında ⁠ API-Key ⁠'i eklemeniz gerekmektedir.

### *Headers*:
•⁠  ⁠⁠ Content-Type ⁠: ⁠ application/json ⁠
•⁠  ⁠⁠ API-Key ⁠: Sunucu tarafından sağlanan benzersiz API anahtarınız. Örnek:

⁠   API-Key: your_api_key
   ⁠

## *Request Format*:
Send a JSON payload with the recipient’s phone number and the message content.

### *Request Body (JSON)*:
⁠ json
{
  "number": "+491234567890",  // Recipient's phone number in international format
  "message": "Your verification code is: 123456"  // The SMS message content
}
 ⁠

### *Example Request (cURL)*:
⁠ bash
curl -X POST http://<your-server-ip>:5000/send-sms \
-H "API-Key: your_api_key" \
-H "Content-Type: application/json" \
-d '{"number": "+491234567890", "message": "Your verification code is: 123456"}'
 ⁠

## *Responses*:

•⁠  ⁠*200 OK*: The SMS was successfully sent.
  ⁠ json
  {
    "success": true,
    "message": "SMS successfully sent"
  }
   ⁠

•⁠  ⁠*400 Bad Request*: Missing required fields (⁠ number ⁠ or ⁠ message ⁠).
  ⁠ json
  {
    "error": "Number and message are required"
  }
   ⁠

•⁠  ⁠*401 Unauthorized*: Invalid or missing API key.
  ⁠ json
  {
    "error": "Unauthorized"
  }
   ⁠

## *Error Handling*:
Make sure to provide both the phone number and message in the correct format and include the valid API key in the headers. Incorrect or missing data will result in an error response.
*/
