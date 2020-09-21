package main

import (
	"encoding/json"
	"fmt"
	liquid "github.com/osteele/liquid"
	_ "github.com/spf13/cobra"
	cli "github.com/urfave/cli"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

func main() {
	app := cli.NewApp()
	app.Name = "goplate"
	app.Usage = " -d '{name: \"value\"}' {src glob pattern} {output dir}"
	app.Flags = []cli.Flag {
		&cli.StringFlag{
			Name: "data",
			Aliases: []string{"d"},
			Value: "{}",
			Usage: "template binding values",
		},
	}
	app.Action = func(c *cli.Context) error {
		srcPattern := c.Args().Get(0)
		dstDir := c.Args().Get(1)
		data := c.String("data")
		files, err := filepath.Glob(srcPattern)
		if err != nil {
			panic(err)
		}
		engine := liquid.NewEngine()
		for _, file := range files {
			base := path.Base(file)
			dstFile := base[0:len(base) - len(path.Ext(file))]
			file, err := os.Open(file)
			if err != nil {
				continue
			}
			defer file.Close()

			// ディレクトリがなければを作る
			_, err = os.Stat(dstDir)
			if err != nil {
				fmt.Printf("make dir %s\n", dstDir)
				if err := os.MkdirAll(dstDir, 0777); err != nil {
					// error
					fmt.Printf("error make dir %s\n", dstDir)
					continue
				}
			}

			b, err := ioutil.ReadAll(file)
			if err != nil {
				continue
			}

			var values map[string]interface{}
			err = json.Unmarshal([]byte(data), &values)
			fmt.Println(values)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			output, err := engine.ParseAndRenderString(string(b), values)
			if err != nil {
				log.Fatalln(err)
			}

			outFile := path.Join(dstDir, dstFile)
			fp, err := os.Create(outFile)
			if err != nil {
				return err
			}
			defer fp.Close()

			_, err = fp.WriteString(output)
			if err != nil {
				return err
			}
			fmt.Printf("make: %s\n", outFile)
		}

		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
