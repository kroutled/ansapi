package ansapi

type BaseResponse struct {
	Result bool `json:"result"`
}

type Users struct {
	Users		[]User		`json:"users"`
}

type User struct {
	ID			string		`json:"id"`
	UID			string		`json:"uid"`
	FirstName	string 		`json:"firstName"`
	Name	    string 		`json:"name"`
	LastName	string 		`json:"lastName"`
	Login		string 		`json:"login"`
	Email		string 		`json:"email"`
	CellNumber	string 		`json:"cellPhoneNumber"`
    TelephoneNumber string  `json:"telephoneNumber"`
} 

type Courses struct {
	Courses		[]Course		`json:"courses"`
}

type Course struct {
	ID						string		`json:"id"`
	UID						string		`json:"uid"`
	Name					string 		`json:"name"`
	title					string 		`json:"title"`
	description		        string 		`json:"description"`
	isTemplate		        bool 		`json:"isTemplate"`
	templateUid		        string 		`json:"templateUid"`
	active				    bool 		`json:"active"`
}

type Templates struct {
	Templates	[]Template		`json:"templates"`
}

type Template struct {
	ID			string		`json:"id"`
	UID			string		`json:"uid"`
	Name		string 		`json:"name"`
}

type Subscriptions struct {
	Courses []Subscription				`json:"courses"`
}

type Subscription struct {
    LearnerUID                string
    LearnerFirstName          string
    LearnerLastName           string
    LearnerEmail              string
    LearnerCellNumber         string
    Inidcator                 string
    ID                        string    `json:"id"`
    UID                       string    `json:"uid"`
    Name                      string    `json:"name"`
    Reseller                  string    `json:"reseller"`
    Active                    bool      `json:"active"`
    SubscribeDate             string    `json:"subscribeDate"`
    StartDate                 string    `json:"startDate"`
    Finished                  bool      `json:"finished"`
    Expired                   bool      `json:"expired"`
    Progress                  string    `json:"progress"` 
    KnowledgeIntake           float64   `json:"knowledgeIntake"`
    Efficiency                float64   `json:"efficiency"`
    ObjectivesProgress        float64   `json:"objectivesProgress"`
    LastTrainingCompletedDate string    `json:"lastTrainingCompletedDate"`
    LastTrainingCompletedDateTime string `json:"lastTrainingCompletedDateTime"`
    Completed                 bool      `json:"completed"`
    CompleteDate              string    `json:"completeDate"`
    Grade                     string    `json:"grade"`
    Passed                    bool      `json:"passed"`
    GradeDate                 string    `json:"gradeDate"`
    Parts	                  []SubscriptionPart    `json:"parts"`
}

type SubscriptionPart struct {
    ID              string     `json:"id"`
    UID             string     `json:"uid"`
    Name            string     `json:"name"`
    Type            string     `json:"type"`
    SubType         string     `json:"subType"`
    BlockName       string     `json:"blockName"`
    Completed       bool       `json:"completed"`
    CompleteDate    string	   `json:"completeDate"`
    CompleteDateTime string    `json:"completeDateTime"`
    Score           string     `json:"score"`
    Passed          bool       `json:"passed"`
    Attempts        []Attempt  `json:"attempts"`
}

type Attempt struct {
    AttemptNumber   int        `json:"attemptNumber"`
    StartDate       string		 `json:"startDate"`
    StartDateTime   string  	`json:"startDateTime"`
    Completed       bool       `json:"completed"`
    CompleteDate    string 		`json:"completeDate"`
    CompleteDateTime string 	`json:"completeDateTime"`
    ExpireDate      string		 `json:"expireDate"`
    Progress        string     `json:"progress"`
    Score           string     `json:"score"`  
    Passed          bool       `json:"passed"`
    Terms           []Term     `json:"terms"`
    Criteria        []Criterion `json:"criteria"`
}

type Term struct {
    Name      string `json:"name"`
    Score     int    `json:"score"`
    MaxScore  int    `json:"maxScore"`
    Threshold int    `json:"threshold"`
}

type Criterion struct {
    ID       string `json:"id"`
    Name     string `json:"name"`
    MaxScore int    `json:"maxScore"`
    Score    int    `json:"score"`
}