# Optimized Random Forest Classifier Training and Inference Microservice
The Optimized Random Forest Classifier: Training and Inference Microservice is a random forest classifier creation and inference tool designed to use new Intel-desigend optimization techniques to create smaller, faster, and still as accurate random forest models for classification and regression tasks.

The service is delivered as a flexible microservice, capable of being used either on its own through its RESTful HTTP API or integrated into a wider microservice-based system.

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