package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bwmarrin/discordgo"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler function Using AWS Lambda Proxy Request
func Handler(request events.APIGatewayProxyRequest) (Response, error) {
	pubkey_b, err := hex.DecodeString("26df446c897e947190c7d6d22d4208dfebf09c5f0feed0496bdbfd730626248e")
	if err != nil {
		return Response{}, errors.New("Couldn't decode the public key")
	}
	if request.Body == "" {
		log.Print("400 No data")
		return Response{
			StatusCode: 400,
			Body:       `{"error":"No body data"}`,
		}, nil
	}

	var body []byte

	if request.IsBase64Encoded {
		body_b, err := base64.StdEncoding.DecodeString(request.Body)
		if err != nil {
			return Response{}, errors.New(fmt.Sprintf("Couldn't decode request body [%s]: %s", body, err))
		}
		body = body_b
	} else {
		body = []byte(request.Body)
	}
	pubkey := ed25519.PublicKey(pubkey_b)

	XSig, ok := request.Headers["x-signature-ed25519"]

	if !ok {
		log.Print("400 No Signature header")
		return Response{
			StatusCode: 400,
			Body:       `{"error": "Missing 'X-Signature-Ed25519' header"}`,
		}, nil
	}

	XSigTime, ok := request.Headers["x-signature-timestamp"]

	if !ok {
		log.Print("400 Missing Timestamp header")
		return Response{
			StatusCode: 400,
			Body:       `{"error": "Missing 'X-Signature-Timestamp' header"}`,
		}, nil
	}

	XSigB, err := hex.DecodeString(XSig)

	if err != nil {
		return Response{}, errors.New("Couldn't decode signature")
	}

	SignedData := []byte(XSigTime + string(body))

	if !ed25519.Verify(pubkey, SignedData, XSigB) {
		log.Print("401 Unauthorized")
		return Response{
			StatusCode: 401,
		}, nil
	} else {
		//authorized
		var inter discordgo.Interaction
		err := json.Unmarshal(body, &inter)

		if err != nil {
			log.Printf("Error decoding interaction: %s", err)
			return Response{
				StatusCode: 400,
			}, nil
		}

		log.Print("200 Type 1 Ping")
		return Response{
			StatusCode: 200,
			Body:       `{"type":1}`}, nil
	}

}

func main() {
	lambda.Start(Handler)
}
