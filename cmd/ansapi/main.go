package main

import(
	"fmt"
	"github.com/kroutled/ansapi"
)

func main() {
	client, err := ansapi.NewClient("https://allangray.anewspring.com/api","b622d692-0289-4fbd-8eb2-7c0492d01ea2")
	if err != nil {
		panic(err)
	}
	//client.GenerateCourseExtIDs()
	var courses = client.GetAllCourses()
	for _, course := range courses {
		fmt.Println(course)
	}
}