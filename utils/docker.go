package utils

import (
	"bufio"
	"bytes"
	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"github.com/docker/docker/pkg/stdcopy"
	"golang.org/x/net/context"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

type ExecResult struct {
	StdOut string
	StdErr string
	ExitCode int
}

func Exec(ctx context.Context, containerID string, command []string) (types.IDResponse, error) {
	docker, err := docker.NewEnvClient()
	if err != nil {
		return types.IDResponse{}, err
	}
	//defer closer(docker)

	config :=  types.ExecConfig{
		AttachStderr: true,
		AttachStdout: true,
		Cmd: command,
		Privileged: true,
	}

	return docker.ContainerExecCreate(ctx, containerID, config)
}

func InspectExecResp(ctx context.Context, id string) (ExecResult, error) {
	var execResult ExecResult
	//docker, err := client.NewEnvClient()
	cli, err := docker.NewEnvClient()
	if err != nil {
		return execResult, err
	}
	// defer closer(docker)

	resp, err := cli.ContainerExecAttach(ctx, id, types.ExecConfig{})
	if err != nil {
		return execResult, err
	}
	defer resp.Close()

	// read the output
	var outBuf, errBuf bytes.Buffer
	outputDone := make(chan error)

	go func() {
		// StdCopy demultiplexes the stream into two buffers
		_, err = stdcopy.StdCopy(&outBuf, &errBuf, resp.Reader)
		outputDone <- err
	}()

	select {
	case err := <-outputDone:
		if err != nil {
			return execResult, err
		}
		break

	case <-ctx.Done():
		return execResult, ctx.Err()
	}

	stdout, err := ioutil.ReadAll(&outBuf)
	if err != nil {
		return execResult, err
	}
	stderr, err := ioutil.ReadAll(&errBuf)
	if err != nil {
		return execResult, err
	}

	res, err := cli.ContainerExecInspect(ctx, id)
	if err != nil {
		return execResult, err
	}

	execResult.ExitCode = res.ExitCode
	execResult.StdOut = string(stdout)
	execResult.StdErr = string(stderr)
	return execResult, nil
}

func RunCommand(containerID string, cmd []string) (ExecResult, error){
	resp, err := Exec(context.Background(), containerID, cmd)
	response, err := InspectExecResp(context.Background(), resp.ID)

	return response, err
}

func GetRunningContainers() []types.Container {
	cli, err := docker.NewEnvClient()
	if err != nil {
		panic(err)
	}
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	return containers
}

func CopyToContainer(containerID string, localPath string) error {
	cli, err := docker.NewEnvClient()
	if err != nil {
		return err
	}

	// dest:= "/home/matteo/git/saltshaker_states.tar.gz"
	dest := CreateTar(localPath)
	time.Sleep(3 * time.Second)
	reader, err := os.Open(dest)
	defer reader.Close()
	if err != nil {
		return err
	}

	err = cli.CopyToContainer(context.Background(), containerID, "/srv/salt/", bufio.NewReader(reader), types.CopyToContainerOptions{})
	if err != nil {
		return err
	}

	return nil
}

func ApplyState(containerID string, state string) (ExecResult, error) {
	cmd := strings.Split("salt salt state.apply " + state, " ")
	resp, err := RunCommand(containerID, cmd)

	return  resp, err
}

func BuildSaltshakerImage()  error {
	rootDir := RootDir()

	cmd := exec.Command("docker", "build", "-t", "saltshaker", rootDir + "/archive/")
	_, err := cmd.Output()
	if err != nil {
		return err
	}

	err = RunSaltshakerContainer()
	if err != nil {
		return err
	}

	return nil
}

func RunSaltshakerContainer() error {
	command := strings.Split("run --privileged --name saltshaker --hostname salt -d saltshaker:latest", " ")

	cmd := exec.Command("docker", command...)

	_, err := cmd.Output()
	if err != nil {
		return err
	}

	time.Sleep(10 * time.Second)

	return nil
}