package gpu_server

import (
	"fmt"
	"minik8s/cmd/kube-apiserver/app/apiconfig"
	"strings"
	"time"
)

type Server struct {
	ssh_cli *Cli
	jobName string
	jobID   string
	OutFile string
	ErrFile string
}

func NewServer(jobName string, outfile string, errfile string) *Server {
	return &Server{
		ssh_cli: NewSSHClient(User, Pwd, Host, Port),
		jobName: jobName,
		OutFile: outfile,
		ErrFile: errfile,
	}
}

const (
	User = "stu1652"
	Pwd  = "M2o$iqsz"
	Host = "pilogin.hpc.sjtu.edu.cn"
	Port = "22"
	Path = "/lustre/home/acct-stu/stu1652/"
)

func (s *Server) JobUpload() error {
	//scp upload program
	homePath := apiconfig.JOB_FILE_DIR_PATH + "/" + s.jobName + "/"
	// homePath := "./"
	cudafile := s.jobName + ".cu"
	n, err := s.ssh_cli.UploadFile(homePath+cudafile, Path+cudafile)
	if err != nil {
		fmt.Printf("upload failed: %v\n", err)
		return err
	}
	// scp upload slurm
	slurmfile := s.jobName + ".slurm"
	n, err = s.ssh_cli.UploadFile(homePath+slurmfile, Path+slurmfile)
	if err != nil {
		fmt.Printf("upload failed: %v\n", err)
		return err
	}
	fmt.Printf("upload file[%v] ok, status=[%d]\n", s.jobName, n)
	return nil
}

func (s *Server) SubmitJob() (string, error) {
	cmd := "sbatch " + Path + s.jobName + ".slurm"
	fmt.Printf("cmd=[%v]\n", cmd)
	backinfo, err := s.ssh_cli.Run(cmd)
	if err != nil {
		fmt.Printf("failed to run shell,err=[%v]\n", err)
		return "", err
	}
	var jobID string
	fmt.Printf("%v back info: \n[%v]\n", cmd, backinfo)
	n, err := fmt.Sscanf(backinfo, "Submitted batch job %s", &jobID)
	if err != nil || n != 1 {
		return "", err
	}
	return jobID, nil
}

func (s *Server) GetJobStatus() bool {
	cmd := "sacct -j " + s.jobID + "| tail -n +3 | awk '{print $6}'"
	fmt.Printf("cmd=[%v]\n", cmd)
	backinfo, err := s.ssh_cli.Run(cmd)
	if err != nil {
		fmt.Printf("failed to run shell,err=[%v]\n", err)
		return false
	}
	fmt.Println("back info: \n", backinfo)
	rows := strings.Split(backinfo, "\n")
	if len(rows) > 0 {
		row := rows[0]
		cols := strings.Split(row, " ")
		if len(cols) > 0 {
			status := cols[0]
			if status == "COMPLETED" {
				return true
			}
		}
	}
	return false
}

func (s *Server) JobDownload() error {
	homePath := apiconfig.JOB_FILE_DIR_PATH + "/" + s.jobName + "/"
	// homePath := "./" + s.jobName + "/"
	n, err := s.ssh_cli.DownloadFile(Path+s.OutFile+".out", homePath+s.OutFile+".out")
	if err != nil {
		return err
	}
	n, err = s.ssh_cli.DownloadFile(Path+s.ErrFile+".err", homePath+s.ErrFile+".err")
	if err != nil {
		return err
	}
	fmt.Printf("download file[%v] ok, status=[%d]\n", s.jobName, n)
	return nil
}

func (s *Server) Run() {
	err := s.JobUpload()
	if err != nil {
		fmt.Printf("failed to upload job,err=[%v]\n", err)
		return
	}
	s.jobID, err = s.SubmitJob()
	if err != nil {
		fmt.Printf("failed to submit job,err=[%v]\n", err)
		return
	}
	fmt.Printf("job[%v] submitted, jobID=[%v]\n", s.jobName, s.jobID)
	for {
		if s.GetJobStatus() {
			fmt.Printf("job[%v] is completed\n", s.jobName)
			break
		}
		time.Sleep(5 * time.Second)
	}
	err = s.JobDownload()
	if err != nil {
		fmt.Printf("failed to download job,err=[%v]\n", err)
		return
	}
}
