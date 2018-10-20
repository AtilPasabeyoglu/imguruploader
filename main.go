package main

import (
	"bytes"
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"log"
	//Gorilla mux library for routing
	//For more info please visit "http://www.gorillatoolkit.org/pkg/mux"
	"github.com/gorilla/mux"
	"image"
	"image/jpeg"
	"image/png"
	"image/gif"
	"sync"
)

func main(){
	//Handled routing in the main function. Used Gorilla Mux package to create easier routes.
	//For more info please visit "http://www.gorillatoolkit.org/pkg/mux"
	var r = mux.NewRouter()

	//Routing for handlers
	r.HandleFunc("/v1/images/upload/{jobID}",GetTasks).Methods("GET")
	r.HandleFunc("/v1/images",GetLinks).Methods("GET")
	r.HandleFunc("/v1/images/upload",CreateJob).Methods("POST")

	http.ListenAndServe(":8080", r)
	http.Handle("/",r)

}


//This function handles the async upload tasks. It takes url(string) and c(channel) as parameters.
func uploadTask(url string,wg *sync.WaitGroup){

	//Download Status
	var status string
	status="Pending"

	//Variable for upload functions return
	var task string
	task=upload(url)

	//Imgur link for uploaded file
	if task!=""{
		//Download Successful!
		status = "Success!"
		jobStatus.completed = append(jobStatus.completed,url)
		fmt.Println("Upload Successful")
	} else {
		//Download Failed
		status = "Failed!"
		jobStatus.failure = append(jobStatus.failure,url)
	}

	imageDetails=append(imageDetails,imgDetail{ImgurURL:task,Status:status})

	//For loop to remove completed items from pending slice
	for i, v := range jobStatus.pending {
		if v == url {
			jobStatus.pending = append(jobStatus.pending[:i], jobStatus.pending[i+1:]...)
			break
		}
	}

	wg.Done()
}



/*
This function checks the image in the given URL for the extension and after that decodes that image.
After that encodes the image variable and uploads it to imgur with the given clientID
*/
func upload(url string) string {

	//Does the image encode and decode
	var imageBuf bytes.Buffer
	imageBuf=imageEncodeDecode(url)

	//POST request for imgur image upload
	req, err := http.NewRequest("POST", "https://api.imgur.com/3/image", &imageBuf)
	if isError(err) {
		//If the image upload fails, this is for loggin the failure
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Client-ID "+clientID)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")


	client := &http.Client{}
	res, err := client.Do(req)
	if isError(err) { return err.Error() }

	defer res.Body.Close()

	//Reading response body
	body, err := ioutil.ReadAll(res.Body)
	if isError(err) { return err.Error() }

	//Decoding the json to result variable
	var result imgurResponse
	error := json.Unmarshal([]byte(body), &result)
	if error != nil {
		fmt.Printf("err was %v", error)
	}

	return result.Data.Link

}

//This function determines Image Type (jpeg,png or gif)
func imageType(url string) string{
	response, e := http.Get(url)
	if e != nil {
		log.Fatal(e)
	}

	_, format, _ := image.DecodeConfig(response.Body)

	return format
}

// Decodes the image from given URL and after that encodes it into a variable
func imageEncodeDecode(url string) bytes.Buffer{

	response, e := http.Get(url)
	if e != nil {
		log.Fatal(e)
	}

	defer response.Body.Close()

	//Returns extension of the given image.
	imageExtension:=imageType(url)

	//Variables for image encoding and decoding
	var myImage image.Image
	var imageBuf bytes.Buffer
	_=myImage

	if imageExtension=="jpeg"{
		myImage, error := jpeg.Decode(response.Body)
		if error != nil {
		}

		error = jpeg.Encode(&imageBuf, myImage,nil)
	}else if imageExtension=="png"{
		myImage, error := png.Decode(response.Body)
		if error != nil {
		}

		error = png.Encode(&imageBuf, myImage)
	}else if imageExtension=="gif"{
		myImage, error := gif.Decode(response.Body)
		if error != nil {
		}
		error = gif.Encode(&imageBuf, myImage,nil)
	}

	return imageBuf

}

//Error Handling function
func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}