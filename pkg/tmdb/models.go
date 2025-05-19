package tmdb

type Movie struct {
	Title       string  `json:"title"`
	Overview    string  `json:"overview"`
	ReleaseDate string  `json:"release_date"`
	Popularity  float64 `json:"popularity"`
	VoteCount   int     `json:"vote_count"`
}
