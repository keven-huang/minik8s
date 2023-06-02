package main

// func main() {

// 	containers, err := dockerClient.GetAllContainers()
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	for _, con := range containers {
// 		fmt.Printf("%v %s\n", con.Names, con.ID)
// 	}
// 	id := "50e35238c50ee4b15401695b0ee558f8955620794013faa71a18d236b4e1dfeb"
// 	//err = dockerClient.RestartContainer(id)
// 	//if err != nil {
// 	//	panic(err.Error())
// 	//}
// 	TestKill(id)
// }

// func TestKill(id string) {
// 	err := dockerClient.KillContainer(id, "SIGKILL")
// 	if err != nil {
// 		panic(err.Error())
// 	}
// }
