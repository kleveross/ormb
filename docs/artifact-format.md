# Preparing Artifact format

 When saving an artifact with ORMB you have to declare its format in ormbfile.yaml, which is case sensitive. Here is the list of the accepted format values and extensions.

| Format | Required format value | Required artifact extension |
| ------------- | ------------- | ------------- |
| SavedModel  |  "SavedModel"  | "saved_model.pb"  |
| ONNX  |  "ONNX"  | ".onnx"  |
| H5  | "H5"  | ".h5"  |
| PMML  | "PMML"  | ".pmml"  |
| CaffeModel  | "CaffeModel"  | ".caffemodel"  |
| NetDef  | "NetDef"  | "init_net.pb"  |
| MXNetParams  | "MXNetParams"  | ".params"  |
| TorchScript  | "TorchScript"  | ".pt"  |
| GraphDef  | "GraphDef"  | ".graphdef"  |
| TensorRT  | "TensorRT"  | ".engine"   |
| SKLearn  | "SKLearn"  | ".joblib"  |
| XGBoost  | "XGBoost"  | ".xgboost"  |
| MLflow  | "MLflow"  | -  |
| Pickle  | "Pickle"  | ".pickle"  |
| Others  | "Others"  | -  |
