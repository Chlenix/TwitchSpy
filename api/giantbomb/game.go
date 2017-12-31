package giantbomb

const (
	Ok                  = 1
	InvalidAPIKey       = 100
	ObjectNotFound      = 101
	BadUrlFormat        = 102
	FilterError         = 104
	SubscriberOnlyVideo = 105
)

type Genre struct {
	Name string `json:"name"`
}

type Results struct {
	Aliases string  `json:"aliases"`
	Deck    string  `json:"deck"`
	Genres  []Genre `json:"genres"`
}

type GBGameInfo struct {
	Error      string `json:"error"`
	StatusCode int    `json:"status_code"`
	Results           `json:"results"`
}
