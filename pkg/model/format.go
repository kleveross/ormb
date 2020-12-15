package model

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// Format is the definition of model format.
type Format string

const (
	FormatSavedModel  Format = "SavedModel"
	FormatONNX        Format = "ONNX"
	FormatH5          Format = "H5"
	FormatPMML        Format = "PMML"
	FormatCaffeModel  Format = "CaffeModel"
	FormatNetDef      Format = "NetDef"
	FormatMXNetParams Format = "MXNetParams"
	FormatTorchScript Format = "TorchScript"
	FormatGraphDef    Format = "GraphDef"
	FormatTensorRT    Format = "TensorRT"
	FormatSKLearn     Format = "SKLearn"
	FormatXGBoost     Format = "XGBoost"
	FormatMLflow      Format = "MLflow"
	FormatOthers      Format = "Others"
)

type Interface interface {
	ValidateDirectory(rootPath string) error
}

func (f Format) ValidateDirectory(rootPath string) error {
	modelFilePath := path.Join(rootPath, "model")
	fileList, err := ioutil.ReadDir(modelFilePath)
	if err != nil {
		return err
	}

	switch f {
	case FormatSavedModel:
		err = f.validateForSavedModel(modelFilePath, fileList)
	case FormatONNX:
		err = f.validateForONNX(modelFilePath, fileList)
	case FormatH5:
		err = f.validateForH5(modelFilePath, fileList)
	case FormatPMML:
		err = f.validateForPMML(modelFilePath, fileList)
	case FormatCaffeModel:
		err = f.validateForCaffeModel(modelFilePath, fileList)
	case FormatNetDef:
		err = f.validateForNetDef(modelFilePath, fileList)
	case FormatMXNetParams:
		err = f.validateForMXNetParams(modelFilePath, fileList)
	case FormatTorchScript:
		err = f.validateForTorchScript(modelFilePath, fileList)
	case FormatGraphDef:
		err = f.validateForGraphDef(modelFilePath, fileList)
	case FormatTensorRT:
		err = f.validateForTensorRT(modelFilePath, fileList)
	case FormatSKLearn:
		err = f.validateForSKLearn(modelFilePath, fileList)
	case FormatXGBoost:
		err = f.validateForXGBoost(modelFilePath, fileList)
	case FormatMLflow:
		err = f.validateForMLflow(modelFilePath, fileList)
	case FormatOthers:
		return nil
	default:
		err = errors.New("unrecognized model format, please check the ormbfile.yaml")
	}

	if err != nil {
		return err
	}

	return nil
}

func ValidateError(modelPath string, modelName string, modelNum int32) error {
	if modelNum != 1 {
		return fmt.Errorf("Expected one %v file in %v directory, but found %v .", modelName, modelPath, modelNum)
	}
	return nil
}

func (f Format) validateForSavedModel(modelPath string, files []os.FileInfo) error {
	var pbFileNum int32
	var variablesDirNum int32
	for _, file := range files {
		if file.Name() == "saved_model.pb" {
			pbFileNum++
		}
		if file.IsDir() && file.Name() == "variables" {
			variablesDirNum++
		}
	}
	if e := ValidateError(modelPath, "saved_model.pb", pbFileNum); e != nil {
		return e
	}
	if e := ValidateError(modelPath, "variables", variablesDirNum); e != nil {
		return e
	}
	return nil
}

func (f Format) validateForONNX(modelPath string, files []os.FileInfo) error {
	var onnxFileNum int32
	for _, file := range files {
		if path.Ext(file.Name()) == ".onnx" {
			onnxFileNum++
		}
	}
	if e := ValidateError(modelPath, "*.onnx", onnxFileNum); e != nil {
		return e
	}
	return nil
}

func (f Format) validateForH5(modelPath string, files []os.FileInfo) error {
	var h5FileNum int32
	for _, file := range files {
		if path.Ext(file.Name()) == ".h5" {
			h5FileNum++
		}
	}
	if e := ValidateError(modelPath, "*.h5", h5FileNum); e != nil {
		return e
	}
	return nil
}

func (f Format) validateForPMML(modelPath string, files []os.FileInfo) error {
	if len(files) != 1 {
		return fmt.Errorf("there are too many pmml file")
	}

	if path.Ext(files[0].Name()) != ".pmml" {
		return fmt.Errorf("there are no *.pmml file in %v directory", modelPath)
	}
	return nil
}

func (f Format) validateForCaffeModel(modelPath string, files []os.FileInfo) error {
	var caffeModelFileNum int32
	var prototxtFileNum int32
	for _, file := range files {
		if path.Ext(file.Name()) == ".caffemodel" {
			caffeModelFileNum++
		}
		if path.Ext(file.Name()) == ".prototxt" {
			prototxtFileNum++
		}
	}
	if e := ValidateError(modelPath, "*.caffemodel", caffeModelFileNum); e != nil {
		return e
	}
	if e := ValidateError(modelPath, "*.prototxt", prototxtFileNum); e != nil {
		return e
	}
	return nil
}

func (f Format) validateForNetDef(modelPath string, files []os.FileInfo) error {
	var initFileNum int32
	var predictFileNum int32
	for _, file := range files {
		if file.Name() == "init_net.pb" {
			initFileNum++
		}
		if file.Name() == "predict_net.pb" {
			predictFileNum++
		}
	}
	if e := ValidateError(modelPath, "init_net.pb", initFileNum); e != nil {
		return e
	}
	if e := ValidateError(modelPath, "predict_net.pb", predictFileNum); e != nil {
		return e
	}
	return nil
}

func (f Format) validateForMXNetParams(modelPath string, files []os.FileInfo) error {
	var jsonFileNum int32
	var paramsFileNum int32
	for _, file := range files {
		if strings.HasSuffix(file.Name(), "symbol.json") {
			jsonFileNum++
		}
		if path.Ext(file.Name()) == ".params" {
			paramsFileNum++
		}
	}
	if e := ValidateError(modelPath, "*symbol.json", jsonFileNum); e != nil {
		return e
	}
	if e := ValidateError(modelPath, "*.params", paramsFileNum); e != nil {
		return e
	}
	return nil
}

func (f Format) validateForTorchScript(modelPath string, files []os.FileInfo) error {
	var ptFileNum int32
	for _, file := range files {
		if path.Ext(file.Name()) == ".pt" {
			ptFileNum++
		}
	}
	if e := ValidateError(modelPath, "*.pt", ptFileNum); e != nil {
		return e
	}
	return nil
}

func (f Format) validateForGraphDef(modelPath string, files []os.FileInfo) error {
	var graphdefFileNum int32
	for _, file := range files {
		if path.Ext(file.Name()) == ".graphdef" {
			graphdefFileNum++
		}
	}
	if e := ValidateError(modelPath, "*.graphdef", graphdefFileNum); e != nil {
		return e
	}
	return nil
}

func (f Format) validateForTensorRT(modelPath string, files []os.FileInfo) error {
	var tensorrtFileNum int32
	for _, file := range files {
		if path.Ext(file.Name()) == ".plan" || path.Ext(file.Name()) == ".engine" {
			tensorrtFileNum++
		}
	}
	if e := ValidateError(modelPath, "*.plan or *.engine", tensorrtFileNum); e != nil {
		return e
	}
	return nil
}

func (f Format) validateForSKLearn(modelPath string, files []os.FileInfo) error {
	var sklearnFileNum int32
	for _, file := range files {
		if path.Ext(file.Name()) == ".joblib" {
			sklearnFileNum++
		}
	}
	if e := ValidateError(modelPath, "*.joblib", sklearnFileNum); e != nil {
		return e
	}
	return nil
}

func (f Format) validateForXGBoost(modelPath string, files []os.FileInfo) error {
	var xgboostFileNum int32
	for _, file := range files {
		if path.Ext(file.Name()) == ".xgboost" {
			xgboostFileNum++
		}
	}
	if e := ValidateError(modelPath, "*.xgboost", xgboostFileNum); e != nil {
		return e
	}
	return nil
}

func (f Format) validateForMLflow(modelPath string, files []os.FileInfo) error {
	var MLflowFileNum int32
	for _, file := range files {
		if file.Name() == "MLmodel" {
			// assuming that user would not fool the tool
			MLflowFileNum++
		}
	}
	if e := ValidateError(modelPath, "MLmodel", MLflowFileNum); e != nil {
		return e
	}
	return nil
}
