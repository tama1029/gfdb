package output

import (
	"fmt"
	"os"
)

type GenOutput struct {
	// output path
	outputDirectory string
}

func NewGenOutput(outputd string) *GenOutput {
	return &GenOutput{
		outputDirectory: outputd,
	}
}

func (o *GenOutput) ToFile(renders map[string][]byte, subDirestory string) error {
	if f, err := os.Stat(o.outputDirectory + "/" + subDirestory); os.IsNotExist(err) || !f.IsDir() {
		if err := os.MkdirAll(o.outputDirectory+"/"+subDirestory, 0755); err != nil {
			return fmt.Errorf("output directory not found, tried to create a directory but it failed")
		}
	}
	for key, render := range renders {
		w, err := os.Create(o.outputDirectory + "/" + subDirestory + "/" + key + ".go")
		if err != nil {
			return err
		}
		defer w.Close()

		_, err = w.Write(render)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *GenOutput) ToFmt(renders map[string][]byte) {
	for key, render := range renders {
		fmt.Println("---------------------------------")
		fmt.Println(key)
		fmt.Println("")
		fmt.Println(string(render))
		fmt.Println("---------------------------------")
	}
}
