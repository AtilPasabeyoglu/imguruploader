package main

//This struct used for Image Details
type imgDetail struct {
	ImgurURL string
	Status string
}

//Creating an array with imgDetail struct. This way it can keep multiple details for multiple images
type detailArray []imgDetail
var imageDetails detailArray

//This struct is used for Job Status.
type jobDetails struct {
	jobID string
	imgDetail []imgDetail
	CreatedTimeStamp string
	FinishedTimeStamp string

}
//Variable for job status. It is the main output
var job jobDetails

//This struct is used for imgur response
type imgurResponse struct{
	Success bool `json:"success"`
	Status int `json:"status"`
	Data struct {
		Title         string   `json:"title"`
		Description   string   `json:"description"`
		Datetime      int      `json:"datetime"`
		Type          string   `json:"type"`
		Animated      bool     `json:"animated"`
		Width         int      `json:"width"`
		Height        int      `json:"height"`
		Size          int      `json:"size"`
		Views         int      `json:"views"`
		Bandwidth     int      `json:"bandwidth"`
		Vote          string   `json:"vote"`
		Favorite      bool     `json:"favorite"`
		Nsfw          bool     `json:"nsfw"`
		Section       string   `json:"section"`
		AccountURL    string   `json:"account_url"`
		AccountId     int      `json:"account_id"`
		IsAd          bool     `json:"is_ad"`
		In_most_viral bool     `json:"in_most_viral"`
		Has_sound     bool     `json:"has_sound"`
		Tags          []string `json:"tags"`
		Ad_type       int      `json:"ad_type"`
		Ad_url        string   `json:"ad_url"`
		In_gallery    bool     `json:"in_gallery"`
		Deletehash    string   `json:"deletehash"`
		Name          string   `json:"name"`
		Link          string   `json:"link"`
	}

}

//This struct is used for keeping the track of the given links
type uploadSlice struct {
	pending []string
	completed []string
	failure []string
}
//Variable created from uploadSlice struct
var jobStatus uploadSlice