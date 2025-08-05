package ansapi

type Users struct {
	Users		[]User		`json:"users"`
}

type User struct {
	Id			string		`json:"id"`
	UID			string		`json:"uid"`
	FirstName	string 		`json:"firstName"`
	LastName	string 		`json:"lastName"`
	Login		string 		`json:"login"`
	Email		string 		`json:"email"`
} 

type Courses struct {
	Courses		[]Course		`json:"courses"`
}

type Course struct {
	Id						string		`json:"id"`
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
	Id			string		`json:"id"`
	UID			string		`json:"uid"`
	Name		string 		`json:"name"`
}

