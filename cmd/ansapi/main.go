package main

import(
	"fmt"
	"github.com/kroutled/ansapi"
)

func main() {
	fmt.Println("It's working...")
	client := ansapi.NewClient("https://worthacademy-discovery.anewspring.com/api","23bae3f2-e7bd-49b3-b5fc-ec975e95a790")
	fmt.Println(client.BaseURL)
	fmt.Println(client.APIKey)
	var users = client.GetCourses()
	for _, course := range(users.Courses) {
		fmt.Println(course)
	}
}

