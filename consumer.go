package kong

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
)

type Consumer struct {
	CreatedAt int    `json:"created_at"`
	ID        string `json:"id"`
	CustomID  string `json:"custom_id"`
	Username  string `json:"username"`
}
type ConsumerList struct {
	Data  []Consumer `json:"data"`
	Next  string     `json:"next"`
	Total int        `json:"total"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ConsumerService struct {
	sling *sling.Sling
}

type JWTAuth struct {
	Secret 		string `json:"secret"`
	ID          string `json:"id"`
	CreatedAt   int    `json:"created_at"`
	Key  		string `json:"key"`
	Algorithm   string `json:"algorithm"`
	ConsumerID  string `json:"consumer_id"`	
}

type JWTAuths struct {
	Data  []JWTAuth `json:"data"`
	Total int       `json:"total"`
}

func NewConsumerService(httpClient *http.Client, baseUrl string) *ConsumerService {
	return &ConsumerService{
		sling: sling.New().Client(httpClient).Base(baseUrl),
	}
}

func (s *ConsumerService) Get(name string) (*Consumer, *http.Response, error) {
	consumer := new(Consumer)
	
	resp, err := s.sling.New().Path("/consumers").Path(name).ReceiveSuccess(consumer)
	
	return consumer, resp, err
}

func (s *ConsumerService) GetJWT(name, auth string) (*JWTAuths, *http.Response, error) {
	jwtauth := new(JWTAuths)

	resp, err := s.sling.New().Path("/consumers"+name+auth).ReceiveSuccess(jwtauth)
	
	return jwtauth, resp, err
}

func (s *ConsumerService) List() (*ConsumerList, *http.Response, error) {
	consumers := new(ConsumerList)
	resp, err := s.sling.New().Path("/consumers/").ReceiveSuccess(consumers)
	return consumers, resp, err
}

func (s *ConsumerService) Create(customID, userName string) (*Consumer, *http.Response, error) {
	type Options struct {
		CustomID  string `url:"custom_id"`
		Username  string `url:"username"`
	}
	opt := Options{ customID, userName }
	consumers := new(Consumer)
	resp, err := s.sling.New().Post("/consumers/").BodyForm(opt).ReceiveSuccess(consumers)
	return consumers, resp, err
}

func (s *ConsumerService) CreateJWTCredential(userName string) (*JWTAuth, *http.Response, error) {
	type Options struct {
		algorithm  string `url:"algorithm"`
	}
	opt := Options{"HS256"}
	jwtAuth := new(JWTAuth)
	resp, err := s.sling.New().Post("/consumers/"+userName+"/jwt").BodyForm(opt).ReceiveSuccess(jwtAuth)
	return jwtAuth, resp, err
}

func (s *ConsumerService) BasicAuth(cons, username, password string) (*Consumer, *http.Response, error) {
	creds := new(Credentials)
	creds.Username = username
	creds.Password = password
	consumer := new(Consumer)
	path := fmt.Sprintf("%s/basic-auth", cons)
	resp, err := s.sling.New().Post("/consumers/").Path(path).BodyJSON(creds).ReceiveSuccess(consumer)
	return consumer, resp, err
}

