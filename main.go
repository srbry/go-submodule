package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	toml "github.com/pelletier/go-toml"
	"gopkg.in/urfave/cli.v1"
)

type rawManifest struct {
	Constraints  []rawProject    `toml:"constraint,omitempty"`
	Overrides    []rawProject    `toml:"override,omitempty"`
	Ignored      []string        `toml:"ignored,omitempty"`
	Required     []string        `toml:"required,omitempty"`
	PruneOptions rawPruneOptions `toml:"prune,omitempty"`
}

type rawProject struct {
	Name     string `toml:"name"`
	Branch   string `toml:"branch,omitempty"`
	Revision string `toml:"revision,omitempty"`
	Version  string `toml:"version,omitempty"`
	Source   string `toml:"source,omitempty"`
}

type rawPruneOptions struct {
	UnusedPackages bool `toml:"unused-packages,omitempty"`
	NonGoFiles     bool `toml:"non-go,omitempty"`
	GoTests        bool `toml:"go-tests,omitempty"`

	//Projects []map[string]interface{} `toml:"project,omitempty"`
	Projects []map[string]interface{}
}

func main() {
	app := cli.NewApp()
	app.Name = "go-submodule"
	app.Usage = "convert a golang dep toml to git submodule commands"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "source, s",
			Usage: "Path to toml file",
		},
	}
	app.Action = parseToml
	app.Run(os.Args)
}

func parseToml(c *cli.Context) error {
	tomlBytes, err := ioutil.ReadFile(c.String("source"))
	if err != nil {
		return err
	}
	var depManifest rawManifest
	if err := toml.Unmarshal(tomlBytes, &depManifest); err != nil {
		return err
	}
	mergedDeps := make(map[string]rawProject)
	fmt.Println("####### Constraints #######")
	for _, constraint := range depManifest.Constraints {
		mergedDeps[constraint.Name] = constraint
	}
	fmt.Println("####### Overrides #######")
	for _, override := range depManifest.Overrides {
		mergedDeps[override.Name] = override
	}
	for _, dep := range mergedDeps {
		filePath := strings.Split(dep.Name, "/")
		dirPath := strings.Join(filePath[:len(filePath)-1], "/")
		fmt.Printf("mkdir -p %s\n", dirPath)
		fmt.Printf("pushd %s\n", dirPath)
		gitSource := dep.Name
		if dep.Source != "" {
			gitSource = dep.Source
		}
		fmt.Printf("  git submodule add git@%s.git\n", gitSource)
		checkout := checkoutString(dep)
		if checkout != "" {
			fmt.Printf("  cd %s\n", filePath[len(filePath)-1])
			fmt.Printf("  git checkout %s\n", checkout)
		}
		fmt.Printf("popd\n\n")
	}
	return nil
}

func checkoutString(dep rawProject) string {
	var checkout string
	if dep.Branch != "" {
		checkout = dep.Branch
	}
	if dep.Version != "" {
		checkout = dep.Version
	}
	if dep.Revision != "" {
		checkout = dep.Revision
	}
	return checkout
}
