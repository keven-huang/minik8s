package main

// It communicates with the deamon process in host machine.
// // equal to `docker ps`
// func main() {
// 	// the client version should correspond to the docker api version in host machine
// 	cli, err := client.NewClientWithOpts(client.WithVersion("1.41"))
// 	if err != nil {
// 		panic(err)
// 	}
// 	//print(cli.ClientVersion())
// 	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
// 	if err != nil {
// 		panic(err)
// 	}

// 	for _, container := range containers {
// 		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
// 	}
// }
