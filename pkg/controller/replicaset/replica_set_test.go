package replicaset

import (
	"fmt"
	"minik8s/pkg/cmd"
	"minik8s/pkg/kubelet/dockerClient"
	"regexp"
	"testing"
	"time"
)

func Count(t *testing.T, regex *regexp.Regexp, number int, prefix string) {
	num := 0
	containers, err := dockerClient.GetAllContainers()

	if err != nil {
		t.Fatal(err)
	}

	for _, con := range containers {
		flag := false
		for _, name := range con.Names {
			fmt.Println(name)
			if regex.MatchString(name) {
				flag = true
				break
			}
		}
		if flag {
			num++
		}
	}
	if num != number {
		t.Fatal(prefix, num, " != ", number)
	}
}

func TestReplicaSetController(t *testing.T) {
	c := cmd.NewKubectlCommand()
	//创建replicaset
	c.SetArgs([]string{"create", "-f", "./test_file/replicaset-example.yaml"})
	err := c.Execute()
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(12 * time.Second)

	regex := regexp.MustCompile("^/my-replicaset-")
	Count(t, regex, 6, "创建replicaset ^/my-replicaset- ")

	// 创建pod, 会自动加入到replicaset中,删除一个之前的pod
	c.SetArgs([]string{"create", "-f", "./test_file/pod-example.yaml"})
	err = c.Execute()
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(8 * time.Second)

	regex = regexp.MustCompile("^/my-replicaset-")
	Count(t, regex, 4, "创建pod ^/my-replicaset- ")
	regex = regexp.MustCompile("^/test5")
	Count(t, regex, 2, "创建pod ^/test5 ")

	// 删除pod, 会自动创建新的pod
	c.SetArgs([]string{"delete", "pod", "test5"})
	err = c.Execute()
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(8 * time.Second)

	regex = regexp.MustCompile("^/my-replicaset-")
	Count(t, regex, 6, "删除pod  ^/my-replicaset- ")
	regex = regexp.MustCompile("^/test5")
	Count(t, regex, 0, "删除pod  ^/test5 ")

	// 删除replicaset, 会删除所有的pod
	c.SetArgs([]string{"delete", "replicaset", "my-replicaset"})
	err = c.Execute()
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(8 * time.Second)

	regex = regexp.MustCompile("^/my-replicaset-")
	Count(t, regex, 0, "删除replicaset ^/my-replicaset- ")
	regex = regexp.MustCompile("^/test5")
	Count(t, regex, 0, "删除replicaset  ^/test5 ")
}
