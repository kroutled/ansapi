package ansapi

import (
	"net/http"
	"io"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"sync"
)
//----------------------------------------------------------------------------------
type Client struct {
	BaseURL string
	APIKey  string
}
//----------------------------------------------------------------------------------
//-------------------------------------CLIENT---------------------------------------
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
func (c Client) SetClientConfig(baseURL, apiKey string) error {
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
//-------------------------------------USERS----------------------------------------
//----------------------------------------------------------------------------------
func (c Client) GetUsers() Users {
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
func (c Client) UserExistsExtID(userExtID string) bool {
	endpoint := fmt.Sprintf("/userExists/%s", userExtID)
	req, _ := http.NewRequest("GET", c.BaseURL + endpoint, nil)
	req.Header.Add( "X-API-Key", c.APIKey)

	res, err := http.DefaultClient.Do(req)
	fmt.Println(err)
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var result BaseResponse
	json.Unmarshal(body, &result)

	if result.Result == true {
		return true
	}
	return false
}
//----------------------------------------------------------------------------------
func (c Client) UserExistsUID(userUID string) bool {
	endpoint := fmt.Sprintf("/userExists?userUID=%s", userUID)
	req, _ := http.NewRequest("GET", c.BaseURL + endpoint, nil)
	req.Header.Add( "X-API-Key", c.APIKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	
	var result BaseResponse
	json.Unmarshal(body, &result)

	if result.Result == true {
		return true
	}
	return false
}
//----------------------------------------------------------------------------------
func (c Client) GetUserByUID(UID string) User {
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
func (c Client) CreateUser(newUser User) {
	data := url.Values{}
	data.Set("firstName", newUser.FirstName)
	data.Set("lastName", newUser.LastName)
	data.Set("id", newUser.ID)
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
func (c Client) UpdateUser(user User) {
	data := url.Values{}
	data.Set("firstName", user.FirstName)
	data.Set("lastName", user.LastName)
	data.Set("id", user.ID)
	data.Set("email", user.Email)
	data.Set("login", user.Email)

	endpoint := fmt.Sprintf("/updateUser?UID=%s", user.UID)
	req, reqerr := http.NewRequest("POST", c.BaseURL + endpoint, strings.NewReader(data.Encode()))
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
//-------------------------------------COURSES--------------------------------------
//----------------------------------------------------------------------------------
func (c Client) GetTemplates() []Template {
	req, _ := http.NewRequest("GET", c.BaseURL + "/getTemplates", nil)
	req.Header.Add( "X-API-Key", c.APIKey)

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var templates Templates
	json.Unmarshal([]byte(body), &templates)

	return templates.Templates
}
//----------------------------------------------------------------------------------
func (c Client) GetTemplateCourses(tuid string, coursesresultsch chan<-Course) {
	endpoint := fmt.Sprintf("/getCourses?templateUID=%s", tuid)
	req,err := http.NewRequest("GET", c.BaseURL + endpoint, nil)
	req.Header.Add( "X-API-Key", c.APIKey)

	if err != nil {
		fmt.Println(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	var courses Courses 
	json.Unmarshal(body, &courses)

	for _, course := range courses.Courses {
		coursesresultsch <- course
	}
}
//----------------------------------------------------------------------------------
func (c Client) GetAllCourses() []Course {
	templates := c.GetTemplates()
	var wg = sync.WaitGroup{}

	workers := 80
	queuedjobsch := make(chan Template)
	resultsch := make(chan Course, 2000)

	for w:=0;w<workers;w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for template := range queuedjobsch {
				c.GetTemplateCourses(template.UID, resultsch)
			}
		}() 
	}

	go func(){
		for _, template := range templates {
			queuedjobsch <- template
		}
		close(queuedjobsch)
	}()

	go func(){
		wg.Wait()
		close(resultsch)
	}()

	var courseResp []Course
	for course := range resultsch {
		courseResp = append(courseResp, course)
	}

	return courseResp
}
//----------------------------------------------------------------------------------
func (c Client) GetCourseByUID(crsUID string) Course {
	endpoint := fmt.Sprintf("/getCourse?UID=%s", crsUID)
	req, _ := http.NewRequest("GET", c.BaseURL + endpoint, nil)
	req.Header.Add( "X-API-Key", c.APIKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	var course Course
	json.Unmarshal(body, &course)
	
	return course
}
//----------------------------------------------------------------------------------
func (c Client) GetCourseByExtID(crsExtID string) Course {
	endpoint := fmt.Sprintf("/getCourse/%s", crsExtID)
	req, _ := http.NewRequest("GET", c.BaseURL + endpoint, nil)
	req.Header.Add( "X-API-Key", c.APIKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	var course Course
	json.Unmarshal(body, &course)
	
	return course
}
//----------------------------------------------------------------------------------
func (c Client)InitCourseExtID(crs Course) {
	endpoint := fmt.Sprintf("/initializeCourseId/%s/%s", crs.UID, crs.UID)
	req, err := http.NewRequest("POST", c.BaseURL + endpoint, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add( "X-API-Key", c.APIKey)
	req.Header.Add( "Content-Type", "application/x-www-form-urlencoded")

	res, reserr := http.DefaultClient.Do(req)
	fmt.Printf("Course initialised: %d\n", res.StatusCode)
	if reserr != nil {
		fmt.Println(reserr)
		fmt.Println(res.StatusCode)
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	fmt.Println(body)
}
//----------------------------------------------------------------------------------
func (c Client) GenerateCourseExtIDs() {
	var courses = c.GetAllCourses()
	
	i := 0
	for _, course := range(courses) {
		if course.ID == "" {
			fmt.Println("Course with no ExtID: ", course.Name)
			fmt.Println("Initialising ExtID... ")
			c.InitCourseExtID(course)
			i++
		}
	}
}
///----------------------------------------------------------------------------------
//-------------------------------------SUBSCRIPTIONS---------------------------------
//-----------------------------------------------------------------------------------
func (c *Client) GetSubscriptionsByEmail(learnerEmail string) Courses {
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
//-----------------------------------------------------------------------------------
func (c *Client) GetLeanrerSubscriptionResults(learnerUID string) []Subscription {
	var subscriptionResults Subscriptions

	if learnerUID != "" {
		endpoint := fmt.Sprintf("/getResults?userUID=%s",learnerUID)
		req, _ := http.NewRequest("GET", c.BaseURL + endpoint, nil)
		req.Header.Add( "X-API-Key", c.APIKey)

		res, _ := http.DefaultClient.Do(req)
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
		}

		json.Unmarshal(body, &subscriptionResults)
	}

	return subscriptionResults.Courses	
}
//-----------------------------------------------------------------------------------
func (c *Client) GetSubscriptions(learner User, resultsch chan <- Subscription) {
	var subscriptionResults Subscriptions

	if learner.UID != "" {
		endpoint := fmt.Sprintf("/getResults?userUID=%s",learner.UID)
		req, _ := http.NewRequest("GET", c.BaseURL + endpoint, nil)
		req.Header.Add( "X-API-Key", c.APIKey)

		res, _ := http.DefaultClient.Do(req)
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
		}

		json.Unmarshal(body, &subscriptionResults)
	}

	for _, subres := range subscriptionResults.Courses {
		resultsch <- subres
	}
}
//-----------------------------------------------------------------------------------
func (c *Client) GetAllSubscriptions() []Subscription {
	learners := c.GetUsers()
	workers := 80
	jobsch := make(chan User)
	resultsch := make(chan Subscription, 2000)
	wg := sync.WaitGroup{}

	//Get workers ready
	for w:=0; w<workers; w++ {
		wg.Add(1)
		go func(){
			defer wg.Done()
			for learnerjob := range jobsch {
				c.GetSubscriptions(learnerjob, resultsch)
			}
		}()
	}

	//feed the jobs
	go func(){
		for _, learner := range learners.Users {
			jobsch <- learner
		}
		close(jobsch)
	}()
	
	//wait and listen for to all end
	go func(){
		wg.Wait()
		close(resultsch)
	}()

	var allSubscriptionResults []Subscription
	for result := range resultsch {
		allSubscriptionResults = append(allSubscriptionResults, result)
	}
	return allSubscriptionResults
}
//-----------------------------------------------------------------------------------
func (c *Client) ExtendSupscription() {
	subscriptions := c.GetAllSubscriptions()
	var expiredSubscriptions []Subscription
	//var nonCommencedSubs []Subscription
	//var lesson3Subs []Subscription
	//var otherSubs []Subscription

	for _, sub := range subscriptions {
		if sub.Expired == true && sub.Completed == false {
			expiredSubscriptions = append(expiredSubscriptions, sub)
		}
	}

    nonCommenced, stoppedAfter3, in4Or5 := GroupExpiredCourses(expiredSubscriptions)

    fmt.Println("Non-commenced & expired:", len(nonCommenced))
    fmt.Println("Stopped after lesson 3:", len(stoppedAfter3))
    fmt.Println("In lesson 4 or 5:", len(in4Or5))

	//for _, expiredSub := range expiredSubscriptions {
		//for i, part := range expiredSub.Parts {
			//if part.Name =="Introduction Video" && part.Completed == false {
				//nonCommencedSubs = append(nonCommencedSubs, expiredSub)
			//}

			//if part.BlockName == "Lesson 3" && part.Name == "Lesson 3 Badge"  && part.Completed == true {
				//if expiredSub.Parts[i+1].Completed == false {
					//lesson3Subs = append(lesson3Subs, expiredSub)
				//}
			//}

			//if part.BlockName == "Lesson 4" && part.Name == "The Power of Assets"  && part.Completed == false {
				//if expiredSub.Parts[i-1].BlockName == "Lesson 3" && expiredSub.Parts[i-1].Name == "Lesson 3 Badge"  && expiredSub.Parts[i-1].Completed == true {
					//otherSubs = append(otherSubs, expiredSub)
				//}
			//}
			
			//if part.BlockName == "Lesson 5" && part.Name == "The Plan"  && part.Completed == false {
				//if expiredSub.Parts[i-1].BlockName == "Lesson 4" && expiredSub.Parts[i-1].Name == "Lesson 4 Badge"  && expiredSub.Parts[i-1].Completed == true {
					//otherSubs = append(otherSubs, expiredSub)
				//}
			//}
		//}
	//}
	//fmt.Printf("Non Commenced: %d\n", len(nonCommencedSubs))
	//fmt.Printf("Up to Lesson 3: %d\n", len(lesson3Subs))
	//fmt.Printf("Either in leasson 4 or 5: %d\n", len(otherSubs))

	fmt.Println(len(expiredSubscriptions))
}

func GroupExpiredCourses(courses []Subscription) (nonCommenced, stoppedAfter3, in4Or5 []Subscription) {
    for _, course := range courses {
        if !course.Expired || course.Completed {
            // skip fully completed courses
            continue
        }

        lesson3Done := false
        lesson4Or5Started := false
        anyLessonStarted := false

        for _, p := range course.Parts {
            // Only consider actual lessons
            if p.Type != "Lesson" && p.Type != "Scorm" && p.Type != "Assessment" {
                continue
            }

            if p.Completed {
                anyLessonStarted = true
            }

            switch p.BlockName {
            case "Lesson 3":
                if p.Completed {
                    lesson3Done = true
                }
            case "Lesson 4", "Lesson 5":
                if p.Completed {
                    lesson4Or5Started = true
                }
            }
        }

        // Assign to group based on progress
        switch {
        case !anyLessonStarted:
            nonCommenced = append(nonCommenced, course)
        case lesson3Done && !lesson4Or5Started:
            stoppedAfter3 = append(stoppedAfter3, course)
        case lesson3Done && lesson4Or5Started:
            in4Or5 = append(in4Or5, course)
        }
    }

    return
}