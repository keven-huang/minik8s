/*
	this client is based on docker-go-client

*/

package dockerClient

import (
	"github.com/docker/docker/client"
)

func GetNewClient() (*client.Client, error) {
	return client.NewClientWithOpts()
}
