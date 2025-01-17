# Tutorials

## Tutorial 1: Build and Download a Model

In this tutorial, you will learn how to upload a dataset, build a model, and download that model.

### Time to Complete
15 minutes

### Learning Objectives

-   By the end of this tutorial, you will be able to create a new model using the microservice

### Prerequisites

You should follow the Build and Run steps in the [Getting Started Guide](get-started-guide.md) before running any tutorial.

### Build a New Model

Using the microservice requires communicating with the RESTful HTTP API. The below commands use *cURL**, but can be adapted to your tool of choice. 

The microservice requires a dataset to train a model. The below commands use an example dataset called `telecom_churn.csv` 

1.  Upload the dataset to the `/data/upload` endpoint. Datasets should be in `.csv` format:
```
curl --location 'localhost:9001/data/upload' \
--form 'file=@"/<full path to>/telecom_churn.csv"'
```
You will receive a JSON response with information about the dataset that you've uploaded. Importantly, the response contains the internal reference ID that you'll use for other requests:
```
{
    "Name": "telecom_churn",
    "ID": "d1",
    "Datapoints": 0,
    "ID_num": 1,
    "Path": "/storage/datasets/telecom_churn.csv"
}
```

2.  Train the classifier model. The below request is specific to the dataset. It informs the microservice what keys to use in the model, such as the name of the column that the model should be filtering for and the features of the rest of the data:
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
The `dataset_id` key should match the response from when you've uploaded your dataset. The rqeuest also includes options for the model itself, most importantly the `max_depth` which defines the depth and complexity of the finalized tree. A depth of 10 usually takes around 8-10 minutes to finish training. The `show_unoptimized` key will allow you to generate a comparison model and will provide you with performance difference between an unoptimzied and optimized model.

>**Note:** It is not uncommon for the request to take a few minutes or more, depending on your `max_depth` key.

 You should see a response similar to below:
 ```
 {
    "Name": "test1",
    "ID": "m1",
    "TrainedDataset": "d1",
    "Features": [
        "AccountWeeks",
        "DataUsage",
        "DayMins",
        "DayCalls",
        "MonthlyCharge",
        "OverageFee",
        "RoamMins"
    ],
    "ID_num": 1,
    "Path": "/storage/models/test1.model",
    "InferName": "Churn",
    "UnoptValAccuracy": 0.9828571428571429,
    "UnoptValRecall": 0.5073746312684366,
    "UnoptTestAccuracy": 0.92,
    "UnoptTestRecall": 0.4791666666666667,
    "ValAccuracy": 0.9823529411764705,
    "ValRecall": 0.49262536873156343,
    "TestAccuracy": 0.9710144927536232,
    "TestRecall": 0.4652777777777778
}
```

With the request provided above, the response contains information about optimized and unoptimized versions of the model.

3. Download the Model

You can download the model with a simple request. Note that this will return the binary representation of the model, so you should pipe this output into a file if using *cURL* or save the response in your request tool.
```
curl --location --request GET 'localhost:9001/models' \
--header 'Content-Type: text/plain' \
--data 'modelid="m1"'
```

### Summary

In this tutorial, you learned how to:
 - Build the microservice
 - Start the microservice
 - Upload a dataset to the microservice
 - Train a new model with the uploaded dataset
 - Download the model

## Tutorial 2: Fetch the tree diagram of the model

In this tutorial, you will learn how to generate a tree diagram of an existing model.

### Time to Complete
5 minutes

### Learning Objectives

-   By the end of this tutorial, you will be able to view the tree diagram of a model you've already created to get more information about the model.

### Prerequisites

You should follow the Build and Run steps in the [Getting Started Guide](get-started-guide.md) before running any tutorial. You should also build a model; you can follow the above [tutorial](##tutorial-1-build-and-download-a-model)

### Step 1: Make the request

1.  Make the request to see the tree of the model. In the below request, we're choosing `modelid=m1`:
```
curl --location --request GET 'localhost:9001/models/tree' \
--header 'Content-Type: text/plain' \
--data 'modelid="m1"'
```
You should receive a response containing the model tree like the below, though below has been truncated as the response may be long depending on your choice of `max_depth` when building your model:
```
    "------ ",
    "",
    "Tree number: 1 ",
    "",
    "Root",
    "   | Tree Depth: 0",
    "   | GINI impurity of the node: 0.24",
    "   | Class distribution in the node: {0: 2864, 1: 469}",
    "   | Predicted class: 0",
    "|-------- Split rule: DataUsage <= 1.186",
    "           | Tree Depth: 1",
    "           | GINI impurity of the node: 0.27",
    "           | Class distribution in the node: {0: 1994, 1: 389}",
    "           | Predicted class: 0",
    "|---------------- Split rule: MonthlyCharge <= 61.52",
    "                   | Tree Depth: 2",
    "                   | GINI impurity of the node: 0.18",
    "                   | Class distribution in the node: {0: 1900, 1: 205}",
    "                   | Predicted class: 0",
    "|------------------------ Split rule: RoamMins <= 17.242",
    "                           | Tree Depth: 3",
    "                           | GINI impurity of the node: 0.17",
    "                           | Class distribution in the node: {0: 1893, 1: 196}",
    "                           | Predicted class: 0",
    "|-------------------------------- Split rule: DayMins <= 139.023",
    "                                   | Tree Depth: 4",
    "                                   | GINI impurity of the node: 0.27",
    "                                   | Class distribution in the node: {0: 438, 1: 85}",
    "                                   | Predicted class: 0",
    ...
    ...
    ...
```

### Summary

In this tutorial, you learned how to view the underlying trees for the model you generated.

## Learn More

-   Understand the architecture in
    the [Overview](overview.md).

## Troubleshooting

The microservice runs as a Docker Compose service. Data in containers should be considered ephemeral. If the service stops working, first you should try to bring down, dispose, rebuild, and bring back up the service:
```
docker compose down -v
docker compose build
docker compose up
```

### Error Logs

The microservice runs as a Docker Compose service. You can view the logs of the container using Docker:
```
docker logs oddforest_microservice
```

## Known Issues

-   Uploaded datasets may show an incorrect number of datapoints in the upload response.
-   Model tree responses may include error data. Model tree data is still accurate