# Introducing `ORMB`

The container virtualization technology like Docker has become the mainstay of cloud computing industry. Software engineers from around the world have flocked to support it. Based on the `Open Container Initiative` technology, the ecology of containers is evolving rapidly. There are also a lot of fast-growing projects such as `Docker Compose` and `Kata Containers`, in which `Kubernetes` has become the standard in the field of cluster scheduling.

When we look back, the reason why container virtualization technology represented by Docker can take a place in the world is because of its image distribution capability. The traditional application deployment is denounced by its complexity of deploying in different environments, which is solved by the image distribution. Docker brings the capability of `Build Once, Deploy Anywhere` to the traditional application scenario. Languages like Java, NodeJS, Rust, etc as well as various dependency libraries, can be built locally once into the OCI compliant container images, which can then be stored and distributed using image registry.

Docker solves the problem of traditional application distribution well, but when it comes to the machine learning, do we need the ability to distribute model?

## Why do we need to distribute the model

Compared with traditional application scenarios, machine learning scenarios have some differences, and it is these differences that lead to some deviations in the requirements of distribution capability. In the traditional application scenario, what needs to be distributed is usually binary executable files, script code with a language interpreter or bytecode with a language virtual machine. Moreover, the development team iterates their code in the traditional application scenario. The same version of code will produce the same output artifacts via multiple compilations.

However, in the machine learning scenario, things are a little different. When we need to deploy machine learning applications, we typically need a model inference server (also known as a model server) and a model that need to be deployed. As for machine learning applications, the model server is like a Tomcat for Java Web applications. The machine learning model itself is only a collection of information such as the weight and structure of the model, and cannot provide services directly to the outside world. It needs to provide RESTful or gRPC-based services externally through a server, and the server is so called a model server. The caller's request arrives at the model server first, which then leverages the model to perform forward calculations and returns the results to the caller as a response.

<p align="center">
<img src="./images/intro/workflow.png" height="150">
</p>

As the picture shown in the below, when we want to deploy model servings in a production environment using cloud-native technologies, we usually start by packaging the model server and the model file into an image for the sake of simplicity and convenience, and reuse the capability of image distribution to achieve the distribution of the model.

<p align="center">
<img src="./images/intro/state-of-art.png" height="150">
</p>

However, there are some problems with this approach:

- The model server itself is an immutable infrastructure, it is often maintained by the infrastructure team or the open source model server team. Besides, the model server is usually very heavy. Take [Nvidia Triton Inference Server](https://github.com/NVIDIA/triton-inference-server) 20.03 as an example, its mirror has 43 layers and the size is 6.3G after unzipping. Therefore, it's very easy to have large images when the model server and the model file are packaged together, which is not easy to maintain and distribution since the P2P distribution of image has not been widely implemented.
- In addition, model updates are very frequent. If the packaging process cannot be automated and requires the participation of algorithm engineers, the cost will be difficult to control.
- Finally, since the image itself does not contain the metadata of the model, we can not scale it horizontally. As the number of models and versions increases, managing multiple models and versions can be quite cumbersome. This problem is exacerbated by the absence of metadata such as hyperparameters, training metrics, storage formats and so on.

With the development of scale, model servers (such as [Nvidia Triton Inference Server](https://github.com/NVIDIA/triton-inference-server), [ONNX Runtime](https://github.com/microsoft/onnxruntime), etc.) can be regarded as low-level "runtime", which is a infrastructure that keeps unchanging. Since the model server can be distributed as a Docker image, the model, on the other hand, can be likened to the script code that runs at runtime. The model can be mounted by the model server container to the running container, which allows us to decouple the model server from the model as well as the infrastructure team and the algorithm team.

<p align="center">
<img src="./images/intro/new.png" height="150">
</p>

## State of the Art

Since model distribution is what we concern about in machine learning scenarios, it must be the feature that other model registry projects want to provide. As for now, the model distribution problem mainly has two kinds of solution ideas in the industry.

The first solution is the implementation of maintaining own-store storage represented by the [Caicloud Clever](https://caicloud.io/products/clever)'s first-generation model registry. The user needs to upload the model to the model registry via SDK or UI. After the model is uploaded, the model registry stores the model and its metadata in its own maintained storage backend. When users need to use the model, they need to use the SDK or interface provided by the model registry to download the model for inference serving.

The other solution is the implementation represented by [ModelDB](https://github.com/VertaAI/modeldb). In this solution, the model registry only stores the metadata of the model. As for the model itself, it will be stored in the third-party object storage such as MinIO, S3, etc. The path of the model storage will also be one of the metadata which will be stored in the model registry. The model can be downloaded using the SDK or interface provided by the object store of a third party.

These two methods both have their own advantages and disadvantages: the former can control the permissions of the model better, but the private SDK and interface has increased the learning costs for users.

The latter one uses mature third-party storage to save model files, and the learning cost is lower. But the separation of metadata and the model makes the control of permission become difficult since users can download the model files directly by bypassing the model registry through the SDK or interface. In addition, both methods require you to build your own wheels to implement the processing logic for your model's metadata and model files.

So, is there a solution that combines the advantages of both solutions while avoiding their disadvantages?

## Distribute machine learning models using image registry

We consider image registry as a way out, which have become a key part of infrastructure in the cloud native era. It provides standardized distribution capability for container images. If the capability of image registry can be reused to distribute models, it will not only avoid repeated wheel building, but also provide similar experience to `docker pull` and `docker push` to avoid excessive learning costs for users.

And so born the `ORMB`.`ORMB`(oh, RMB) is named from the OCI-based Registry for ML/DL Model Bundle, which is capable of distributing models and their metadata using existing image registry.

## An end-to-end example

In this section, we'll use image recognition as an example to show how `ORMB` can be used to distribute machine learning models.

We will train a simple CNN image recognition model locally using [Fashion MNIST](https://github.com/zalandoresearch/fashion-mnist) and push it to a remote image registry using `ORMB`. Later, we'll also use `ORMB` to pull the model from the registry in our remote server, and use the third-party model server to provide external serving. Then we can call the service with a RESTful interface and see the results.

The code of model training is shown below, but the specific training process will not be repeated in this article. Finally, the trained model is saved in SavedModel format. You can see the model in the [examples of ormb](https://github.com/caicloud/ormb/tree/master/examples/SavedModel-fashion).

```python
# Modeling, set optimizer, and train
model = keras.Sequential([
  keras.layers.Conv2D(input_shape=(28,28,1),
                      filters=8,
                      kernel_size=3,
                      strides=2,
                      activation='relu',
                      name='Conv1'),
  keras.layers.Flatten(),
  keras.layers.Dense(10,
                     activation=tf.nn.softmax,
                     name='Softmax')
])
model.compile(optimizer='adam',
              loss='sparse_categorical_crossentropy',
              metrics=['accuracy'])
model.fit(train_images, train_labels, epochs=epochs)

test_loss, test_acc = model.evaluate(test_images,
                                     test_labels)
import tempfile

# save the model to the sub model directory
MODEL_DIR = "./model"
version = 1
export_path = os.path.join(MODEL_DIR, str(version))
tf.keras.models.save_model(
    model,
    export_path,
    overwrite=True,
    include_optimizer=True,
    save_format=None,
    signatures=None,
    options=None
)
```

Next, we can push the trained model from local to remote image registry.

```bash
# Save the model in local cache first
$ ormb save gaocegege/fashion_model:v1
ref:       gaocegege/fashion_model:v1
digest:    6b08cd25d01f71a09c1eb852b3a696ee2806abc749628de28a71b507f9eab996
size:      162.1 KiB
format:    SavedModel
v1: saved

# Push the model from local cache to remote registry
$ ormb push gaocegege/fashion_model:v1
The push refers to repository [gaocegege/fashion_model]
ref:       gaocegege/fashion_model:v1
digest:    6b08cd25d01f71a09c1eb852b3a696ee2806abc749628de28a71b507f9eab996
size:      162.1 KiB
format:    SavedModel
v1: pushed to remote (1 layer, 162.1 KiB total)
```

Taking [Harbor]() as an example, we can see the model's metadata in Harbor registry.

<p align="center">
<img src="/docs/images/intro/harbor.png" height="350">
</p>

Then, we can download the model from the registry. The download process is similar to the push.

```bash
# Pull the model from remote registry to local cache
$ ormb pull gaocegege/fashion_model:v1
v1: Pulling from gaocegege/fashion_model
ref:     gaocegege/fashion_model:v1
digest:  6b08cd25d01f71a09c1eb852b3a696ee2806abc749628de28a71b507f9eab996
size:    162.1 KiB
Status: Downloaded newer model for gaocegege/fashion_model:v1

# Export the model from local cache to current directory
$ ormb export gaocegege/fashion_model:v1
ref:     localhost/gaocegege/fashion_model:v1
digest:  6b08cd25d01f71a09c1eb852b3a696ee2806abc749628de28a71b507f9eab996
size:    162.1 KiB

# View the local file directory
$ tree examples/SavedModel-fashion
examples/SavedModel-fashion
├── model
│   ├── saved_model.pb
│   └── variables
│       ├── variables.data-00000-of-00001
│       └── variables.index
├── ormbfile.yaml
└── training-serving.ipynb

2 directories, 5 files
```

Next, we can use `TFServing` to deploy the model as a RESTful service, and use the Fashion MNIST dataset for inference.

```bash
$ tensorflow_model_server --model_base_path=$(pwd)/model --model_name=fashion_model --rest_api_port=8501
2020-05-27 17:01:57.499303: I tensorflow_serving/model_servers/server.cc:358] Running gRPC ModelServer at 0.0.0.0:8500 ...
[evhttp_server.cc : 238] NET_LOG: Entering the event loop ...
2020-05-27 17:01:57.501354: I tensorflow_serving/model_servers/server.cc:378] Exporting HTTP/REST API at:localhost:8501 ...
```

<p align="center">
<img src="./images/intro/infer.png" height="350">
</p>

We could also use `Seldon Core` to deploy the model serving directly on the Kubernetes cluster, see [our documentation](https://github.com/caicloud/ormb/blob/master/docs/tutorial-serving-seldon.md) for more information.

```yaml
apiVersion: machinelearning.seldon.io/v1alpha2
kind: SeldonDeployment
metadata:
  name: tfserving
spec:
  name: mnist
  protocol: tensorflow
  predictors:
    - graph:
        children: []
        implementation: TENSORFLOW_SERVER
        modelUri: demo.goharbor.io/tensorflow/fashion_model:v1
        serviceAccountName: ormb
        name: mnist-model
        parameters:
          - name: signature_name
            type: STRING
            value: predict_images
          - name: model_name
            type: STRING
            value: mnist-model
      name: default
      replicas: 1
```

As the algorithm engineer iterates over the new version of the model and package it, the new image can be pulled with `ORMB` and redeployed. `ORMB` can be used with any image registry that meets the [OCI Distribution Specification](https://github.com/opencontainers/distribution-spec), which means that `ORMB` can both support image registry on the public cloud and other open source image registry projects like Harbor.

We can also use Harbor's Webhook to deploy model serving continuously. By registering a Webhook in the Harbor UI, all the request events of Harbor push will be forwarded to our defined HTTP endpoint. We can implement the corresponding deployment logic in Webhook, such as updating the version of Seldon deployment model serving according to the new model, completing the continuous deployment of model serving and so on.

<p align="center">
<img src="./images/intro/webhook.png" height="350">
</p>

## System design

In a word, `ORMB` is designed to be `Docker` in machine learning scenarios, aiming to solve the distribution problem of machine learning models. `ORMB` also stands on the shoulders of giants such as [OCI Artifacts](https://github.com/opencontainers/artifacts), [OCI Distribution Specification](https://github.com/opencontainers/distribution-spec) and [Harbor](https://github.com/goharbor/harbor).

### What happened with `docker pull`

Before we get to the design, let's take a look at what happens when we download a container image. For image registries that conform to the OCI specification, they all follow the same rule, which is the [OCI Distribution Specification](https://github.com/opencontainers/distribution-spec/blob/master/spec.md#pulling-an-image).

<p align="center">
<img src="./images/intro/docker-pull.png" height="550">
</p>

First, Docker will request the image manifest from the image registry. Manifest is a JSON file whose definition includes two parts, config and layers. Config is a JSON object, and layers is an array composed of JSON objects. As you can see, config has the same structure as each object in layers, including three fields, digest, mediaType and size. Digest is the ID of the object. The mediaType indicates the type of this content. And size is the size of this content.

The config of container image has a fixed mediaType: `application/vnd.oci.image.config.v1+json`. Here is an example of config, which contains the configuration of the container image (AKA the metadata of the image). The config(metadata) is often used by image registry to display information in the UI, distinguish between different operating system builds, and so on.

```json
{
  "created": "2015-10-31T22:22:56.015925234Z",
  "architecture": "amd64",
  "os": "linux",
  "config": {
    "Entrypoint": ["/bin/my-app-binary"]
  },
  "rootfs": {
    "diff_ids": [
      "sha256:c6f988f4874bb0add23a778f753c65efe992244e148a1d2ec2a8b664fb66bbd1",
      "sha256:5f70bf18a086007016e948b04aed3b82103a36bea41755b6cddfaf10ace3c6ef"
    ],
    "type": "layers"
  },
  "history": [
    {
      "created": "2015-10-31T22:22:54.690851953Z",
      "created_by": "/bin/sh -c #(nop) ADD file:a3bc1e842b69636f9df5256c49c5374fb4eef1e281fe3f282c65fb853ee171c5 in /"
    },
    {
      "created": "2015-10-31T22:22:55.613815829Z",
      "created_by": "/bin/sh -c #(nop) CMD [\"sh\"]",
      "empty_layer": true
    }
  ]
}
```

But the layers of container image is consisted by multiple `application/vnd.oci.image.layer.v1.*` layers(in which the most common type is application/vnd.oci.image.layer.v1.tar+gzip`). As we all know, container images are built in layers, with each layer in the image corresponding to an object in layers.

The config in image as well as each layer in layers are stored in the image registry as `Blob`, and their digest exists as `Key`. Therefore, after requesting the image manifest, Docker downloads all `Blobs` in parallel with digest, including config and all layers.

For more explanation, see the article [How Docker container images are made](https://mp.weixin.qq.com/s?__biz=Mzg3ODAzMTMyNQ==&mid=2247486285&idx=1&sn=42a8ec6caef52606b52e2f52bcb4c4a4&chksm=cf18b3fff86f3ae91364c2550a2388a12753b487ece57062f5c8a0e0e0128bc7b04b7c744612&mpshare=1&scene=1&srcid=&sharer_sharetime=1592534888462&sharer_shareid=398b8cdf40a05ba26d3dd04ea2871de0&exportkey=ATl%2F7vtmgPmdXxZu9m8Keoo%3D&pass_ticket=ueO%2Blu33I8%2FugN7BWl3s2elMGGjKsUSHTOptrPOvggpFPexr9IOFFwVErwhEniVu#rd).

### How can a image registry support model distribution

In the introduction above, we can see that by defining the manifest, image config of `application/vnd.oci.image.config.v1+json` mediaType and layers of `application/vnd.oci.image.layer.v1.*` mediaType, the image registry can support image well.

With the expansion of cloud native, there are many other types of artifacts except for container image, such as [Helm Chart](https://helm.sh/), [CNAB](https://cnab.io/), etc. We also hope to reuse the versioning management, distribution of the artifacts, as well as the ability of layered storage. So [OCI Artifacts](https://github.com/opencontainers/artifacts) were created to meet this need.

The image registry supports the storage and distribution of images through manifest, config, and layers. If we could define the types and structures of config and layers ourselves, we would be able to extend the ability to store and distribute other types of artifacts. [OCI Artifacts](https://github.com/opencontainers/artifacts) are the guidance documents that exactly provided for this purpose. Since the image registry already has this capability, [OCI Artifacts](https://github.com/opencontainers/artifacts) are not just a specification but a guidance document that guide developers on how to use the image registry's extend capabilities to support other artifact types.

When it comes to the specific supporting for model, we define our own [config](https://github.com/caicloud/ormb/blob/master/docs/spec-v1alpha1.md) for the model structure, which is called `ormbfile`, similar to `Dockerfile`. The tentative mediaType is `application/vnd.caicloud.model.config.v1alpha1+json`, and its sample configuration is shown below.

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
   "metrics": {
      "training": [
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
                   1001
               ],
               "dType": "float64",
           }
       ],
       "layers": {
               "conv": 1
        }
   },
   "training": {
       "git": {
           "repository": "git@github.com:caicloud/ormb.git",
           "revision": "22f1d8406d464b0c0874075539c1f2e96c253775"
       }
   },
   "dataset": {
       "git": {
           "repository": "git@github.com:caicloud/ormb.git",
           "revision": "22f1d8406d464b0c0874075539c1f2e96c253775"
       }
   }
}
```

For Layers, because the model files are difficult to store hierarchically, the model files are compressed and archived with `application/tar+gzip` mediaType and then uploaded to the image registry in the current design.

Therefore, the process of downloading a model from the image registry is as follows:

<p align="center">
<img src="./images/intro/ormb.png" height="550">
</p>

We work closely with the Harbor community to use Harbor as the default registry right now. We also contribute to the extensibility of Harbor in OCI Artifact support and so `ORMB` can reuse many of Harbor's capabilities.

For example, `ORMB` takes advantage of Harbor's Replication feature to enable model synchronization between image registry in multiple environments, such as training environments, model production environments, and so on. We can also configure webhooks for the model to add custom processing logic to the model lifecycle.

[ormb]: https://github.com/kleveross/ormb
