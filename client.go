package ansapi

import (
	"net/http"
	"io"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

//----------------------------------------------------------------------------------
type Client struct (
	BaseURL string
	APIKey  string
)
//----------------------------------------------------------------------------------
func (c *Client) SetClientConfig(baseURL, apiKey string) {
	BaseURL = baseURL
	APIKey = apiKey
}
//----------------------------------------------------------------------------------
func (c *Client) GetUsers() Users {
	req, _ := http.NewRequest("GET", BaseURL + "/getUsers", nil)
	req.Header.Add( "X-API-Key", APIKey)

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var learners Users
	json.Unmarshal([]byte(body), &learners)

	return learners
}
//----------------------------------------------------------------------------------
func (c *Client) GetUser(UID string) User {
	endpoint := fmt.Sprintf("/getUser?userUID=%s", UID)
	req, _ := http.NewRequest("GET", BaseURL + endpoint, nil)
	req.Header.Add( "X-API-Key", APIKey)

	res, err := http.DefaultClient.Do(req)
	fmt.Println(err)
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var learner User
	json.Unmarshal([]byte(body), &learner)

	return learner
}
//----------------------------------------------------------------------------------
func (c *Client) CreateUser(newUser User) {
	data := url.Values{}
	data.Set("firstName", newUser.FirstName)
	data.Set("lastName", newUser.LastName)
	data.Set("id", newUser.Id)
	data.Set("email", newUser.Email)
	data.Set("login", newUser.Email)
	data.Set("notify", "true")

	req, reqerr := http.NewRequest("POST", BaseURL + "/addUser", strings.NewReader(data.Encode()))
	if reqerr != nil {
		fmt.Println(reqerr)
	}
	req.Header.Add( "X-API-Key", APIKey)
	req.Header.Add( "Content-Type", "application/x-www-form-urlencoded")

	res, reserr := http.DefaultClient.Do(req)
	if reserr != nil {
		fmt.Println(reserr)
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	fmt.Println(body)
}
//----------------------------------------------------------------------------------
func (c *Client) GetTemplates() Templates {
	req, _ := http.NewRequest("GET", BaseURL + "/getTemplates", nil)
	req.Header.Add( "X-API-Key", APIKey)

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var templates Templates
	json.Unmarshal([]byte(body), &templates)

	return templates
}
//----------------------------------------------------------------------------------
func (c *Client) GetCourses() Courses {
	templates := GetTemplates()
	var courses Courses	
	for _, template := range templates.Templates {
		endpoint := fmt.Sprintf("/getCourses/%s", template.Id)
		req, _ := http.NewRequest("GET", BaseURL + endpoint, nil)
		req.Header.Add( "X-API-Key", APIKey)

		res, _ := http.DefaultClient.Do(req)
		defer res.Body.Close()

		body, _ := io.ReadAll(res.Body)
		var tempCourses Courses
		json.Unmarshal([]byte(body), &tempCourses)
		for _, retCourse := range tempCourses.Courses {
			courses.Courses = append(courses.Courses, retCourse)
		}
	}
	return courses
}
//----------------------------------------------------------------------------------
func (c *Client) GetSubscriptions(learnerEmail string) Courses {
	var courses Courses	
	var learnerUID string
	learners := GetUsers()

	for _, learner := range learners.Users {
		if learner.Email == learnerEmail {
			learnerUID = learner.UID
		}	
	}
	if learnerUID != "" {
		endpoint := fmt.Sprintf("/getSubscriptions?userUID=%s",learnerUID)
		req, _ := http.NewRequest("GET", BaseURL + endpoint, nil)
		req.Header.Add( "X-API-Key", APIKey)

		res, _ := http.DefaultClient.Do(req)
		defer res.Body.Close()

		body, _ := io.ReadAll(res.Body)

		json.Unmarshal([]byte(body), &courses)
	} else {
		fmt.Println("No learner found")
	}
	return courses
}
