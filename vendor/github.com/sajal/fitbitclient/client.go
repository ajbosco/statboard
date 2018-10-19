package fitbitclient

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/fitbit"
	"io/ioutil"
	"net/http"
)

//Minimal config for new client. Do not modify once initialized
type Config struct {
	CredFile     string   //File where the tokens will be permanently cached in
	ClientID     string   //Fitbit client ID
	ClientSecret string   //Fitbit client secret
	Scopes       []string //Fitbit oauth scopes
}

//Adapted from https://github.com/golang/oauth2/issues/84#issuecomment-175834679
type fitBitTransport struct {
	Base   *oauth2.Transport
	Config *Config
}

func (c *fitBitTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	if err != nil {
		return nil, err
	}
	resp, err = c.Base.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	newTok, err := c.Base.Source.Token()
	if err != nil {
		return nil, err
	}
	//Might be a bit race-ey... YOLO
	err = savetokenfile(c.Config.CredFile, newTok)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func loadtokenfile(fname string) (*oauth2.Token, error) {
	//First try to get tok from file
	dat, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.Unmarshal(dat, tok)
	return tok, err
}

func savetokenfile(fname string, tok *oauth2.Token) error {
	//First try to get tok from file
	j, err := json.Marshal(tok)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fname, j, 0600)
	return err
}

//Create a http.Client which wraps a Transport which caches updated creds in cfg.CredFile
//If cfg.CredFile does not exist(or is unreadable) interactive oauth authentication process will initiate.
func NewFitBitClient(cfg *Config) (*http.Client, error) {
	//Initialise FitbitAPI
	conf := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Scopes:       cfg.Scopes,
		Endpoint:     fitbit.Endpoint,
	}
	tok, err := loadtokenfile(cfg.CredFile)
	if err != nil {
		//Means we need to generate token.
		url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
		fmt.Printf("Visit the URL for the auth dialog: %v\nEnter token here: ", url)
		var code string
		if _, err := fmt.Scan(&code); err != nil {
			return nil, err
		}
		tok, err = conf.Exchange(oauth2.NoContext, code)
		if err != nil {
			return nil, err
		}
		//Save the token for future.
		err = savetokenfile(cfg.CredFile, tok)
		if err != nil {
			return nil, err
		}
	}
	ts := conf.TokenSource(oauth2.NoContext, tok)
	tr := &oauth2.Transport{Source: ts}
	//Wrap an http.Client with our FitBitTransport
	client := &http.Client{
		Transport: &fitBitTransport{Base: tr, Config: cfg},
	}
	return client, nil
}
