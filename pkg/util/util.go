package util

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v2"

	ormbmodel "github.com/kleveross/ormb/pkg/model"
)

// InferModelFormat infer model format by file ext.
func InferModelFormat(dir string) (ormbmodel.Format, error) {
	fileList, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}

	netdefFileNum := 0
	mxnetFileNum := 0

	for _, file := range fileList {
		if file.IsDir() {
			continue
		}

		fileExtName := strings.Trim(path.Ext(file.Name()), ".")
		switch fileExtName {
		case "pb":
			if strings.HasSuffix(file.Name(), "saved_model.pb") {
				return ormbmodel.FormatSavedModel, nil
			} else if strings.HasSuffix(file.Name(), "init_net.pb") || strings.HasSuffix(file.Name(), "predict_net.pb") {
				netdefFileNum++
			}
		case "onnx":
			return ormbmodel.FormatONNX, nil
		case "graphdef":
			return ormbmodel.FormatGraphDef, nil
		case "caffemodel":
			return ormbmodel.FormatCaffeModel, nil
		case "pt":
			return ormbmodel.FormatTorchScript, nil
		case "plan", "engine":
			return ormbmodel.FormatTensorRT, nil
		case "pmml":
			return ormbmodel.FormatPMML, nil
		case "params":
			mxnetFileNum++
		case "json":
			if strings.HasSuffix(file.Name(), "symbol.json") {
				mxnetFileNum++
			}
		case "h5":
			return ormbmodel.FormatH5, nil
		case "xgboost":
			return ormbmodel.FormatXGBoost, nil
		case "joblib":
			return ormbmodel.FormatSKLearn, nil
		}
	}

	if netdefFileNum == 2 {
		return ormbmodel.FormatNetDef, nil
	}

	if mxnetFileNum == 2 {
		return ormbmodel.FormatMXNetParams, nil
	}

	return ormbmodel.FormatOthers, nil
}

// WriteORMBFile write ormbfile.yaml if file file is not exist.
func WriteORMBFile(filePath string, format ormbmodel.Format) error {
	metadata := &ormbmodel.Metadata{
		Format: string(format),
	}
	data, err := yaml.Marshal(metadata)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close() // nolint

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}
