package tmdb

type Movie struct {
	Title       string  `json:"title"`
	Overview    string  `json:"overview"`
	Popularity  float64 `json:"popularity"`
	ReleaseDate string  `json:"release_date"`
	VoteCount   int     `json:"vote_count"`
}
