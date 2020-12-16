<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [OCI Model Configuration Specification](#oci-model-configuration-specification)
  - [Terminology](#terminology)
    - [Tensor](#tensor)
    - [GitRepo](#gitrepo)
  - [Properties](#properties)
  - [Example](#example)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# OCI Model Configuration Specification

An OCI Model is a Machine Learning/Deep Learning model. This specification outlines the JSON format describing models for use and execution tool and its relationship to filesystem changesets.

This section defines the `application/vnd.caicloud.model.config.v1alpha1+json` media type.

## Terminology

This specification uses the following terms:

### [Tensor](https://www.tensorflow.org/guide/tensor)

Tensors are multi-dimensional arrays with a uniform type.

- **name** string, OPTIONAL

    Name of the tensor.

- **description** string, OPTIONAL

    Description of the tensor.

- **dType** string, OPTIONAL

    DType of the tensor.

- **size** object, OPTIONAL

    Size of the tensor.

- **opType** string, OPTIONAL

    OpType of the tensor,  It is only used for PMML.

- **values** object, OPTIONAL

    Values of the tensor,  It is only used for PMML.

### [GitRepo](https://kubernetes.io/docs/concepts/storage/volumes/#gitrepo)

GitRepo is the Git repository.

- **repository** string, OPTIONAL

    Git repository URL

- **revision** string, OPTIONAL

    Revision of the repository

- **description** string, OPTIONAL

    Description of the training script

## Properties

Note: Any OPTIONAL field MAY also be set to null, which is equivalent to being absent.

- **created** string, OPTIONAL

    A combined date and time at which the image was created, formatted as defined by RFC 3339, section 5.6.

- **author** string, OPTIONAL

    Gives the name and/or email address of the person or entity which created and is responsible for maintaining the image.

- **description** string, OPTIONAL

    Description of the model

- **tags** object, OPTIONAL

    Tags of the model

- **labels** object, OPTIONAL

    The field contains arbitrary metadata for the container. This property MUST use the [annotation rules](https://github.com/opencontainers/image-spec/blob/master/annotations.md#rules).

- **framework** string, OPTIONAL

    The framework of the model (e.g. TensorFlow)

- **format** string, REQUIRED

    The format of the model (SavedModel, ONNX and so on)

- **size** int, REQUIRED

    The size of the model

- **metrics** object, OPTIONAL

    The metrics of the model

    - **name** string, REQUIRED

        Name of the metric

    - **value** string,  REQUIRED
        
        Value of the metric

- **hyperParameters** object, OPTIONAL

    HyperParameters of the model

    - **name** string, REQUIRED

        Name of the hyperparameter

    - **value** string,  REQUIRED

        Value of the hyperparameter

- **signature** object, OPTIONAL

    Inputs and outputs of the model

    - **inputs** *Tensor*, OPTIONAL

        Inputs of the model
    
    - **outputs** *Tensor*, OPTIONAL

        Outputs of the model

    - **layers** object, OPTIONAL

        Layers of the model

- **training** object, OPTIONAL

    Training information of the model

    - **git** *GitRepo*, OPTIONAL

        Git repository of the training script

- **dataset** object, OPTIONAL

    - **git** *GitRepo*, OPTIONAL

        Git repository of the dataset

## Example

Here is an example model configuration JSON document:

```json
{
   "created": "2015-10-31T22:22:56.015925234Z",
   "author": "Ce Gao <gaoce@caicloud.io>",
   "description": "CNN Model",
   "tags": [
       "cv"
   ],
   "labels": {
       "tensorflow.version": "2.0.0"
   },
   "framework": "TensorFlow",
   "format": "SavedModel",
   "size": 9223372036854775807,
   "metrics":{ 
      "training":[
           {
               "name": "acc",
               "value": "0.928"
           }
      ],
   }
   "hyperParameters": [
       {
           "name": "batch_size",
           "value": "32"
       }
   ],
   "signature": {
       "inputs": [
           {
               "name": "input_1",
               "size": [
                   224,
                   224,
                   3
               ],
               "dType": "float64",
           }
       ],
       "outputs": [
           {
               "name": "output_1",
               "size": [
                   1,
                   1000
               ],
               "dType": "float64",
           }
       ],
       "layers": {
               "conv" : 1
       }
   },
   "training": {
       "git": {
           "repository": "git@github.com:kleveross/ormb.git",
           "revision": "22f1d8406d464b0c0874075539c1f2e96c253775"
       }
   },
   "dataset": {
       "git": {
           "repository": "git@github.com:kleveross/ormb.git",
           "revision": "22f1d8406d464b0c0874075539c1f2e96c253775"
       }
   }
}
```
