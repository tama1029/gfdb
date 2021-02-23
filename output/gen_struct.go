package output

import (
	"fmt"
	"os"
)

type GenStructOutput struct {
	// output path
	outputDirectory string
}

func NewGenStructOutput(outputd string) *GenStructOutput {
	return &GenStructOutput{
		outputDirectory: outputd,
	}
}

func (o *GenStructOutput) ToFile(renders map[string][]byte) error {
	if f, err := os.Stat(o.outputDirectory); os.IsNotExist(err) || !f.IsDir() {
		return fmt.Errorf("output directory not found, please create directory")
	}
	for key, render := range renders {
		w, err := os.Create(o.outputDirectory + "/" + key + ".go")
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

func (o *GenStructOutput) ToFmt(renders map[string][]byte) {
	for key, render := range renders {
		fmt.Println("---------------------------------")
		fmt.Println(key)
		fmt.Println("")
		fmt.Println(string(render))
		fmt.Println("---------------------------------")
	}
}
