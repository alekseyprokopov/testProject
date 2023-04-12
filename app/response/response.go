package response

type Response struct {
	Results []UserData `json:"results"`
}
type UserData struct {
	Gender string `json:"gender"`
	Name   struct {
		Title string `json:"title"`
		First string `json:"first"`
		Last  string `json:"last"`
	} `json:"name"`
	Location struct {
		Street struct {
			Number int    `json:"number"`
			Name   string `json:"name"`
		} `json:"street"`
		City        string `json:"city"`
		State       string `json:"state"`
		Country     string `json:"country"`
		Postcode    int    `json:"postcode"`
		Coordinates struct {
			Latitude  string `json:"latitude"`
			Longitude string `json:"longitude"`
		} `json:"coordinates"`
	} `json:"location"`
}
