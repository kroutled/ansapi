package main

import(
	"fmt"
	"github.com/kroutled/ansapi"
)

func main() {
	client, err := ansapi.NewClient("","")
	if err != nil {
		panic(err)
	}
	//client.GenerateCourseExtIDs()
	var courses = client.GetAllCourses()
	for _, course := range courses {
		fmt.Println(course)
	}
}