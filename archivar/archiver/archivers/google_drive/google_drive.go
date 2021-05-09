package google_drive

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Eun/gdriver"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleDrive struct {
	OAuthToken      string
	ClientSecrets   string
	uploadDirectory string
	drive           *gdriver.GDriver
	client          *http.Client
	logger          *logrus.Logger
	token           *oauth2.Token
}

func New(config interface{}, logger *logrus.Logger) (gdrive *GoogleDrive) {
	jsonM, _ := json.Marshal(config)
	json.Unmarshal(jsonM, &gdrive)
	return gdrive
}

// Retrieve a token, saves the token, then returns the generated client.
func (g *GoogleDrive) getClient(config *oauth2.Config, ctx context.Context) (*http.Client, error) {
	if g.OAuthToken != "" {
		tok, err := tokenFromString(g.OAuthToken)
		if err != nil {
			return nil, err
		}

		return config.Client(ctx, tok), nil
	}

	tok := getTokenFromWeb(config)
	g.token = tok

	return config.Client(ctx, tok), nil
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	jt, _ := json.Marshal(tok)
	fmt.Printf("%s", jt)
	return tok
}

// Retrieves a token from a local file.
// func tokenFromFile(file string) (*oauth2.Token, error) {
// 	f, err := os.Open(file)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer f.Close()
// 	tok := &oauth2.Token{}
// 	err = json.NewDecoder(f).Decode(tok)
// 	return tok, err
// }

// Retrieves a token from a local file.
func tokenFromString(tokenString string) (*oauth2.Token, error) {
	tok := &oauth2.Token{}
	err := json.NewDecoder(strings.NewReader(tokenString)).Decode(tok)
	return tok, err
}

func (g *GoogleDrive) Connect() (newSession bool, err error) {
	if g.drive != nil {
		return false, nil
	}
	ctx := context.Background()

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON([]byte(g.ClientSecrets))
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	g.client, err = g.getClient(config, ctx)
	if err != nil {
		return false, err
	}

	g.drive, err = gdriver.New(g.client)
	if err != nil {
		return false, err
	}

	return
}
