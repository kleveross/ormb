package model

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
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
	FormatMXNETParams Format = "MXNETParams"
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
	case FormatMXNETParams:
		err = f.validateForMXNETParams(modelFilePath, fileList)
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
	default:
		err = errors.New("unrecognized model format, please check the ormbfile.yaml")
	}

	if err != nil {
		return err
	}

	return nil
}

func (f Format) validateForSavedModel(modelPath string, files []os.FileInfo) error {
	var pbFileFlag bool
	var variablesDirFlag bool
	for _, file := range files {
		if path.Ext(file.Name()) == ".pb" {
			pbFileFlag = true
		}
		if file.IsDir() && file.Name() == "variables" {
			variablesDirFlag = true
		}
	}
	if !pbFileFlag {
		return fmt.Errorf("there are no *.pb file in %v directory", modelPath)
	}
	if !variablesDirFlag {
		return fmt.Errorf("there are no variables dir in %v directory", modelPath)
	}
	return nil
}

func (f Format) validateForONNX(modelPath string, files []os.FileInfo) error {
	var onnxFileFlag bool
	for _, file := range files {
		if path.Ext(file.Name()) == ".onnx" {
			onnxFileFlag = true
		}
	}
	if !onnxFileFlag {
		return fmt.Errorf("there are no *.onnx file in %v directory", modelPath)
	}
	return nil
}

func (f Format) validateForH5(modelPath string, files []os.FileInfo) error {
	var h5FileFlag bool
	for _, file := range files {
		if path.Ext(file.Name()) == ".h5" {
			h5FileFlag = true
		}
	}
	if !h5FileFlag {
		return fmt.Errorf("there are no *.h5 file in %v directory", modelPath)
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
	var caffeModelFileFlag bool
	var prototxtFileFlag bool
	for _, file := range files {
		if path.Ext(file.Name()) == ".caffemodel" {
			caffeModelFileFlag = true
		}
		if path.Ext(file.Name()) == ".prototxt" {
			prototxtFileFlag = true
		}
	}
	if !caffeModelFileFlag {
		return fmt.Errorf("there are no *.caffemodel file in %v directory", modelPath)
	}
	if !prototxtFileFlag {
		return fmt.Errorf("there are no *.prototxt file in %v directory", modelPath)
	}
	return nil
}

func (f Format) validateForNetDef(modelPath string, files []os.FileInfo) error {
	var initFileFlag bool
	var predictFileFlag bool
	for _, file := range files {
		if file.Name() == "init_net.pb" {
			initFileFlag = true
		}
		if file.Name() == "predict_net.pb" {
			predictFileFlag = true
		}
	}
	if !initFileFlag {
		return fmt.Errorf("there are no init_net.pb file in %v directory", modelPath)
	}
	if !predictFileFlag {
		return fmt.Errorf("there are no predict_net.pb file in %v directory", modelPath)
	}
	return nil
}

func (f Format) validateForMXNETParams(modelPath string, files []os.FileInfo) error {
	var jsonFileFlag bool
	var paramsFileFlag bool
	for _, file := range files {
		if path.Ext(file.Name()) == ".json" {
			jsonFileFlag = true
		}
		if path.Ext(file.Name()) == ".params" {
			paramsFileFlag = true
		}
	}
	if !jsonFileFlag {
		return fmt.Errorf("there are no *.json file in %v directory", modelPath)
	}
	if !paramsFileFlag {
		return fmt.Errorf("there are no *.params file in %v directory", modelPath)
	}
	return nil
}

func (f Format) validateForTorchScript(modelPath string, files []os.FileInfo) error {
	var ptFileFlag bool
	for _, file := range files {
		if path.Ext(file.Name()) == ".pt" {
			ptFileFlag = true
		}
	}
	if !ptFileFlag {
		return fmt.Errorf("there are no *.pt file in %v directory", modelPath)
	}
	return nil
}

func (f Format) validateForGraphDef(modelPath string, files []os.FileInfo) error {
	var pbFileFlag bool
	for _, file := range files {
		if path.Ext(file.Name()) == ".pb" {
			pbFileFlag = true
			break
		}
	}
	if !pbFileFlag {
		return fmt.Errorf("there are no *.pb file in %v directory", modelPath)
	}
	return nil
}

func (f Format) validateForTensorRT(modelPath string, files []os.FileInfo) error {
	var tensorrtFileFlag bool
	for _, file := range files {
		if path.Ext(file.Name()) == ".plan" {
			tensorrtFileFlag = true
		}
	}
	if !tensorrtFileFlag {
		return fmt.Errorf("there are no *.plan file in %v directory", modelPath)
	}
	return nil
}

func (f Format) validateForSKLearn(modelPath string, files []os.FileInfo) error {
	var sklearnFileFlag bool
	for _, file := range files {
		if path.Ext(file.Name()) == ".joblib" {
			sklearnFileFlag = true
		}
	}
	if !sklearnFileFlag {
		return fmt.Errorf("there are no *.joblib file in %v directory", modelPath)
	}
	return nil
}

func (f Format) validateForXGBoost(modelPath string, files []os.FileInfo) error {
	var xgboostFileFlag bool
	for _, file := range files {
		if path.Ext(file.Name()) == ".xgboost" {
			xgboostFileFlag = true
		}
	}
	if !xgboostFileFlag {
		return fmt.Errorf("there are no *.xgboost file in %v directory", modelPath)
	}
	return nil
}

func (f Format) validateForMLflow(modelPath string, files []os.FileInfo) error {
	var isMLflowFile bool
	for _, file := range files {
		if file.Name() == "MLmodel" {
			// assuming that user would not fool the tool
			isMLflowFile = true
		}
	}
	if !isMLflowFile {
		return fmt.Errorf("there are no MLmodel file in %v, directory", modelPath)
	}
	return nil
}
