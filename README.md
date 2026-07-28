# Optimized Random Forest Classifier Training and Inference Microservice
The Optimized Random Forest Classifier: Training and Inference Microservice is a random forest classifier creation and inference tool designed to use new Intel-desigend optimization techniques to create smaller, faster, and still as accurate random forest models for classification and regression tasks.

The service is delivered as a flexible microservice, capable of being used either on its own through its RESTful HTTP API or integrated into a wider microservice-based system.

## Security Considerations

> **This is a reference implementation. The following known-unsafe patterns are intentional placeholders that a downstream integrator MUST address before any production or network-exposed deployment.**

### 1. No Authentication or Authorization
All HTTP routes (`/train`, `/infer`, `/model/upload`, `/data/upload`, etc.) are completely unauthenticated. Any network peer can invoke them without credentials. Integrators must add an authentication and authorization layer at the API gateway or within the service before exposure to any untrusted network.

### 2. Unsafe Model Deserialization (Pickle RCE)
`joblib.load()` — used to load `.model` files — is `pickle.load()` under the hood. A crafted pickle file can execute **arbitrary OS commands at load time** via `__reduce__`. Because `/model/upload` accepts arbitrary files with no validation and `/infer` loads whatever was uploaded, a network peer can achieve **unauthenticated remote code execution** by uploading a malicious model.

Adding authentication alone does **not** fix this: even an authenticated user could upload a malicious file. Before production use, integrators must replace `joblib.load` with a safe alternative — for example:
- Verify a cryptographic signature on every model file before loading it.
- Use a non-pickle serialization format (e.g., ONNX, PMML) and a corresponding safe loader.
- Validate the model file structure before deserialization.

### 3. Container Runs as Root
The container previously ran all processes as `root`. This has been corrected in this repository (a non-root `appuser` is now created and used), but integrators should verify their own derived images maintain this practice and further restrict container capabilities (e.g., `--cap-drop ALL`, read-only root filesystem) appropriate to their environment.

## Licensing Information
This repository is licensed under the MIT [LICENSE](License.txt). This project also utilizes the Random Forest algorithm provided by SciKit Learn, which is licensed under the BSD 3-Clause License. Information on that algorithms license can be found at their repository [here](https://github.com/scikit-learn/scikit-learn/tree/main).

## Running as a Docker Container
The application is built to run as a Docker container. This application is built and tested with Docker Engine 24.0.4 and Docker Compose v2.19.1.

```
docker compose build
docker compose up
```

## Sample API Commands
### Status
```
curl --location 'localhost:9001/status'
```
### Upload Dataset
```
curl --location 'localhost:9001/data/upload' \
--form 'file=@"/<full path to>/telecom_churn.csv"'
```
### Train Model
```
curl --location 'localhost:9001/train' \
--header 'Content-Type: text/plain' \
--data 'title = "OOD Random Forest Classifier Configuration"


dataset_id = "d1"
model_id = "m1"
infer_name = "Churn"
name="test1"
features = [
    '\''AccountWeeks'\'',
    '\''DataUsage'\'',
    '\''DayMins'\'',
    '\''DayCalls'\'',
    '\''MonthlyCharge'\'',
    '\''OverageFee'\'',
    '\''RoamMins'\''
]
n_trees = 10
min_samples_split = 5
max_depth = 10
x_features_fraction= 0.5
data_split= 0.7
show_unoptimized = true'
```
### Infer with Model
```
curl --location 'localhost:9001/infer' \
--header 'Content-Type: text/plain' \
--data 'datasetid="d1"
modelid="m1"'
```
### Download Model
```
curl --location --request GET 'localhost:9001/models' \
--header 'Content-Type: text/plain' \
--data 'modelid="m1"'
```

### Get Model Tree
```
curl --location --request GET 'localhost:9001/models/tree' \
--header 'Content-Type: text/plain' \
--data 'modelid="m1"'
```