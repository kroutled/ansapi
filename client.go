package ansapi

import (
	"net/http"
	"io"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"
)
//----------------------------------------------------------------------------------
type Client struct {
	BaseURL string
	APIKey  string
}

//----------------------------------------------------------------------------------
func NewClient(baseURL, apiKey string) (*Client, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("baseURL cannot be empty")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("apiKey cannot be empty")
	}
	return &Client{
		BaseURL: baseURL,
		APIKey: apiKey,
	}, nil
}
//----------------------------------------------------------------------------------
func (c *Client) SetClientConfig(baseURL, apiKey string) error {
	if baseURL == "" {
		return fmt.Errorf("baseURL cannot be empty")
	}
	if apiKey == "" {
		return fmt.Errorf("apiKey cannot be empty")
	}
	c.BaseURL = baseURL
	c.APIKey = apiKey
	return nil
}
//----------------------------------------------------------------------------------
func (c *Client) GetUsers() Users {
	req, _ := http.NewRequest("GET", c.BaseURL + "/getUsers", nil)
	req.Header.Add( "X-API-Key", c.APIKey)

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
	req, _ := http.NewRequest("GET", c.BaseURL + endpoint, nil)
	req.Header.Add( "X-API-Key", c.APIKey)

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

	req, reqerr := http.NewRequest("POST", c.BaseURL + "/addUser", strings.NewReader(data.Encode()))
	if reqerr != nil {
		fmt.Println(reqerr)
	}
	req.Header.Add( "X-API-Key", c.APIKey)
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
	req, _ := http.NewRequest("GET", c.BaseURL + "/getTemplates", nil)
	req.Header.Add( "X-API-Key", c.APIKey)

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var templates Templates
	json.Unmarshal([]byte(body), &templates)

	return templates
}
//----------------------------------------------------------------------------------
func (c *Client) GetCourses() Courses {
	templates := c.GetTemplates()
	var courses Courses	
	var mu sync.Mutex
	var wg sync.WaitGroup

	concurrency := 50
	sem := make(chan struct{}, concurrency)

	for _, template := range templates.Templates {
		wg.Add(1)
		sem <- struct{}{}

		go func(t Template) {
			defer wg.Done()
			defer func(){ <-sem}()

			endpoint := fmt.Sprintf("/getCourses/%s?includeWithoutId=true", template.Id)
			req, _ := http.NewRequest("GET", c.BaseURL + endpoint, nil)
			req.Header.Add( "X-API-Key", c.APIKey)

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()

			body, _ := io.ReadAll(res.Body)
			var tempCourses Courses
			json.Unmarshal([]byte(body), &tempCourses)

			mu.Lock()
			courses.Courses = append(courses.Courses, tempCourses.Courses...)
			mu.Unlock()
		}(template)
	}
	
	wg.Wait()
	return courses
}
//----------------------------------------------------------------------------------
func (c *Client)InitCourseExtID(crs *Course) {
	endpoint := fmt.Sprintf("/initializeCourseId/%s/%s", crs.UID, crs.UID)
	req, err := http.NewRequest("POST", c.BaseURL + endpoint, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add( "X-API-Key", c.APIKey)
	req.Header.Add( "Content-Type", "application/x-www-form-urlencoded")

	res, reserr := http.DefaultClient.Do(req)
	fmt.Println(res.StatusCode)
	if reserr != nil {
		fmt.Println(reserr)
		fmt.Println(res.StatusCode)
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	fmt.Println(body)
}
//----------------------------------------------------------------------------------
func (c *Client) GenerateCourseExtIDs() {
	rn := time.Now()
	var courses = c.GetCourses()
	
	i := 0
	for _, course := range(courses.Courses) {
		if course.Id == "" {
			fmt.Println("course with no extid: ", course.Name)
			c.InitCourseExtID(&course)
			i++
		}
	}
	fmt.Println(i)
	fmt.Println("Took: ", time.Since(rn))
}
//----------------------------------------------------------------------------------
func (c *Client) GetSubscriptions(learnerEmail string) Courses {
	var courses Courses	
	var learnerUID string
	learners := c.GetUsers()

	for _, learner := range learners.Users {
		if learner.Email == learnerEmail {
			learnerUID = learner.UID
		}	
	}
	if learnerUID != "" {
		endpoint := fmt.Sprintf("/getSubscriptions?userUID=%s",learnerUID)
		req, _ := http.NewRequest("GET", c.BaseURL + endpoint, nil)
		req.Header.Add( "X-API-Key", c.APIKey)

		res, _ := http.DefaultClient.Do(req)
		defer res.Body.Close()

		body, _ := io.ReadAll(res.Body)

		json.Unmarshal([]byte(body), &courses)
	} else {
		fmt.Println("No learner found")
	}
	return courses
}