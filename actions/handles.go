package actions

import (
	"net/http"
	"os"
	"os/exec"

	"github.com/gobuffalo/buffalo"
)

type AntenaRequest struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}

func AntenaHandler(c buffalo.Context) error {
	var req AntenaRequest
	c.Bind(&req)

	cmd_app := "/opt/shells/antenaShell"
	cmd := exec.Command(cmd_app, req.Url, req.Name)
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	return c.Render(http.StatusOK, r.JSON(req))
}

func CamHandler(c buffalo.Context) error {
	cmd_app := "/opt/shells/camShell"
	cmd := exec.Command(cmd_app, c.Param("name"), c.Param("time"))
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	return c.Render(http.StatusOK, r.String("OK"))
}

func StartStreamCamHandler(c buffalo.Context) error {
	cmd_app := "/opt/shells/streamStartCamShell"
	cmd := exec.Command(cmd_app, c.Param("feed"))
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	return c.Render(http.StatusOK, r.String(`http://localhost:6563/`+c.Param("feed")+`.jpeg`))
}

type RTSPRequest struct {
	Url  string `json:"url"`
	Feed string `json:"feed"`
}

func StartRTSPStreamHandler(c buffalo.Context) error {
	var req RTSPRequest
	c.Bind(&req)

	cmd_app := "/opt/shells/streamRTSPShell"
	cmd := exec.Command(cmd_app, req.Url, req.Feed)
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	return c.Render(http.StatusOK, r.String(`http://localhost:6563/`+req.Feed+`.jpeg`))
}

func StopStreamCamHandler(c buffalo.Context) error {
	cmd_app := "/opt/shells/streamStopShell"
	cmd := exec.Command(cmd_app, c.Param("feed"))
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	return c.Render(http.StatusOK, r.String("OK"))
}

type saveOutput struct {
	savedOutput []byte
}

func (so *saveOutput) Write(p []byte) (n int, err error) {
	so.savedOutput = append(so.savedOutput, p...)
	return os.Stdout.Write(p)
}

func ListStreamHandler(c buffalo.Context) error {
	cmd_app := "/opt/shells/streamListShell"
	cmd := exec.Command(cmd_app)

	var so saveOutput
	cmd.Stdout = &so
	_ = cmd.Run()

	var output string
	if len(string(so.savedOutput)) > 1 {
		output = string(so.savedOutput)
	} else {
		output = "Não existe nenhum stream em andamento!"
	}

	return c.Render(http.StatusOK, r.String(output))
}
