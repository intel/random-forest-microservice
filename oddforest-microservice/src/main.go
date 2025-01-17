package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	session "intel.com/oddforest-microservice/session"
)

// Debug Flags

// Session object
var current_session session.Session

// setupRouter: Sets up the Gin-based http router with our options and our routes.
func setupRouter() *gin.Engine {
	router := gin.Default()

	//GET Methods
	router.GET("/status", getStatus)
	router.GET("/datasets", getDataset)
	router.GET("/models", getModel)
	router.GET("/models/tree", getModelTree)
	router.GET("/results", getResults)
	//POST Methods
	router.POST("/train", startTraining)
	router.POST("/data/upload", uploadData)
	router.POST("/model/upload", uploadModel)
	router.POST("/infer", infer)
	return router
}

// getStatus: Returns a list of: running + finished tasks; uploaded datasets; built models
func getStatus(c *gin.Context) {
	// Without ID, return everything. Begin building our return: start by querying the available task list and their status
	c.JSON(http.StatusOK, current_session)
	// Query the uploaded datasets

	// Query the available models

	// Return our list as a HTTP StatusOK

	// With ID, only return status of training job
}

// getDataset: Returns a list of available datasets in the microservice, or if an ID is provided returns information about a specific dataset
func getDataset(c *gin.Context) {
	// With no ID, returns a list of datasets available.
	c.JSON(http.StatusOK, current_session.Datasets)
	// If ID present, returns info about that dataset
}

// getModel: Returns a specific model based on the provided model ID. Without a model ID, it returns a list of available models
func getModel(c *gin.Context) {
	// If ID not present in request, query list of available models and return
	var downloadConfig session.DownloadConfig
	if err := c.BindTOML(&downloadConfig); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var model session.Model
	for _, mod := range current_session.Models {
		if downloadConfig.ModelID == mod.ID {
			model = mod
			log.Printf("Downloading model %s...", model.ID)
			file, err := os.ReadFile(model.Path)
			if err != nil {
				c.AbortWithError(http.StatusBadGateway, err)
			}
			c.Data(200, "application/octet-stream", file)
			return
		}
	}
	if model.ID == "" {
		log.Printf("Model not found: %s\nReturning models list", downloadConfig.ModelID)
	}
	c.JSON(http.StatusOK, current_session.Models)
	// If ID present, check if model exists. If yes, return model details (time of creation, dataset used, size)

	// If ID present and not exists, return error to user (http.StatusNotFound))

}

func getModelTree(c *gin.Context) {
	var downloadConfig session.DownloadConfig
	if err := c.BindTOML(&downloadConfig); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var model session.Model
	for _, mod := range current_session.Models {
		if downloadConfig.ModelID == mod.ID {
			model = mod
			treesTOMLPath := generateTreesTOML(model.Path, "show_trees")
			channel_status := make(chan []byte)
			go func() {
				cmd := exec.Command("python", "../..//random_forest/main.py", treesTOMLPath)
				log.Println("Showing model tree...")
				out, err := cmd.CombinedOutput()
				if err != nil {
					log.Println(err)
				}
				current_session.Tasks[len(current_session.Tasks)-1].Status = "Complete"
				channel_status <- out
			}()
			response := strings.Split(strings.TrimSuffix(string(<-channel_status), "\n"), "\n")
			c.JSON(http.StatusOK, response)
			return

		}
	}
	if model.ID == "" {
		log.Printf("Model not found: %s\nReturning models list", downloadConfig.ModelID)
	}
	c.JSON(http.StatusOK, current_session.Models)
}

// getResults: Returns a specfiic inference job's results based on the provided run ID. Without a run ID, it returns a list of available result runs
func getResults(c *gin.Context) {
	// If job ID not present, return list of available results
	c.JSON(http.StatusOK, current_session.Results)
	// If ID is present, check if run exists, and then return the run results (model used, dataset trained with, dataset ran, results from run)

	// If ID present and not exists, return error ot user (http.StatusNotFound)
}

// startTraining: Based on an input TOML file, builds a new model and assigns it an ID.
func startTraining(c *gin.Context) {
	// If TOML provided, check for valid dataset ID. If everything's ready, start a training task as a trackable async goroutine. Add said goroutine to the list of tasks
	// While job runs in background, return list of features from dataset and number of datapoints being trained.
	log.Println("Building new model...")
	// Check TOML:
	var training_body session.TrainingConfig
	if err := c.BindTOML(&training_body); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// Ready? Get model information ready, set new ID for this job, and add it to the list. Set status to "Getting Ready"
	// New Model
	var existing_model_ids []int
	max_id := 0
	for _, model := range current_session.Models {
		existing_model_ids = append(existing_model_ids, model.ID_num)
	}
	if len(existing_model_ids) != 0 {
		max_id = slices.Max(existing_model_ids)
	}
	log.Println(training_body)
	var new_model session.Model
	new_model.ID_num = max_id + 1
	new_model.ID = "m" + fmt.Sprint(new_model.ID_num)
	new_model.Name = training_body.Name
	new_model.TrainedDataset = training_body.DatasetID
	new_model.Path = "/storage/models/" + new_model.Name + ".model"
	new_model.Features = training_body.Features
	new_model.InferName = training_body.InferName
	// New Task
	var existing_task_ids []int
	for _, task := range current_session.Tasks {
		existing_task_ids = append(existing_task_ids, task.ID_num)
	}
	if len(existing_task_ids) != 0 {
		max_id = slices.Max(existing_task_ids)
	}
	var new_task session.Task
	new_task.ModelID = new_model.ID
	new_task.ID_num = max_id + 1
	new_task.ID = "t" + fmt.Sprint(new_task.ID_num)
	new_task.Status = "Running"
	current_session.Tasks = append(current_session.Tasks, new_task)

	// Get our train config TOML ready - get the dataset path, get the features, get the data, get the name to set the path
	// Prep the environment and send the config to the training tool. Start a goroutine to handle running and set status to "running"
	dataset_path := ""
	for _, dataset := range current_session.Datasets {
		if dataset.ID == new_model.TrainedDataset {
			dataset_path = dataset.Path
			break
		}
		response := "dataset not found, id: " + new_model.TrainedDataset
		c.JSON(http.StatusBadRequest, response)
	}

	// Read in our CSV file
	features := training_body.Features
	trainingtomlpath := generateTrainingTOML(dataset_path, new_model.Path, training_body.InferName, "train", features, training_body.MaxDepth, training_body.NTrees, training_body.SampleSplit, training_body.FeaturesFraction, training_body.DataSplit, training_body.ShowUnoptimzied)
	channel_status := make(chan []byte)
	go func() {
		cmd := exec.Command("python", "../..//random_forest/main.py", trainingtomlpath)
		log.Println("Starting training...")
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Println(err)
		}
		current_session.Tasks[len(current_session.Tasks)-1].Status = "Complete"
		channel_status <- out
	}()
	response := strings.Split(strings.TrimSuffix(string(<-channel_status), "\n"), "\n")
	log.Println(response)
	new_model.UnoptValAccuracy, _ = strconv.ParseFloat(response[len(response)-8], 64)
	new_model.UnoptValRecall, _ = strconv.ParseFloat(response[len(response)-7], 64)
	new_model.UnoptTestAccuracy, _ = strconv.ParseFloat(response[len(response)-6], 64)
	new_model.UnoptTestRecall, _ = strconv.ParseFloat(response[len(response)-5], 64)
	new_model.ValAccuracy, _ = strconv.ParseFloat(response[len(response)-4], 64)
	new_model.ValRecall, _ = strconv.ParseFloat(response[len(response)-3], 64)
	new_model.TestAccuracy, _ = strconv.ParseFloat(response[len(response)-2], 64)
	new_model.TestRecall, _ = strconv.ParseFloat(response[len(response)-1], 64)

	current_session.Models = append(current_session.Models, new_model)

	// Return a good status to the user.
	c.JSON(http.StatusOK, new_model)
}

// uploadData: Uploads a provided dataset (.csv) to the microservice datastore and assigns it an ID.
func uploadData(c *gin.Context) {
	// Check if dataset is a .csv. If valid (trust the user to add a valid dataset), extract the feature list and then assign it an ID for future use.

	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "Form error %s", err.Error())
	}
	filename := filepath.Base(file.Filename)
	fmt.Println(filename)
	path := "/storage/datasets/"
	if err := c.SaveUploadedFile(file, "/storage/datasets/"+filename); err != nil {
		if err := c.SaveUploadedFile(file, "./"+filename); err != nil {
			c.String(http.StatusBadRequest, "Error uploading file: %s", err.Error())
			return
		}
		c.String(http.StatusOK, "storage mount not available, saving locally")
		path = "./"
	}
	var existing_dataset_ids []int
	max_id := 0
	for _, model := range current_session.Models {
		existing_dataset_ids = append(existing_dataset_ids, model.ID_num)
	}
	if len(existing_dataset_ids) != 0 {
		max_id = slices.Max(existing_dataset_ids)
	}
	fmt.Println("assinging dataset id: " + fmt.Sprint(max_id+1))
	// Give our session the new dataset
	var new_dataset session.Dataset
	new_dataset.ID_num = max_id + 1
	new_dataset.ID = "d" + fmt.Sprint(new_dataset.ID_num)
	new_dataset.Name = strings.TrimSuffix(filename, ".csv")
	new_dataset.Path = path + filename
	current_session.Datasets = append(current_session.Datasets, new_dataset)
	//Return good status
	c.JSON(http.StatusOK, new_dataset)
}

// uploadModel: Uploads a previously downloaded model to the microservice datastore, verifies it, and assigns it an ID
func uploadModel(c *gin.Context) {
	log.Println("Uploading Model...")
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "Form error %s", err.Error())
	}
	filename := filepath.Base(file.Filename)
	fmt.Println(filename)
	if err := c.SaveUploadedFile(file, "/storage/models/"+filename); err != nil {
		if err := c.SaveUploadedFile(file, "./"+filename); err != nil {
			c.String(http.StatusBadRequest, "Error uploading file: %s", err.Error())
			return
		}
		c.String(http.StatusOK, "storage mount not available, saving locally")
	}
	// New Model
	var existing_model_ids []int
	max_id := 0
	for _, model := range current_session.Models {
		existing_model_ids = append(existing_model_ids, model.ID_num)
	}
	if len(existing_model_ids) != 0 {
		max_id = slices.Max(existing_model_ids)
	}
	var new_model session.Model
	new_model.ID_num = max_id + 1
	new_model.ID = "m" + fmt.Sprint(new_model.ID_num)
	new_model.Name = filename
	new_model.TrainedDataset = "unknown"
	new_model.Path = "/storage/models/" + new_model.Name + ".model"
	new_model.Features = []string{"features"}
	new_model.InferName = "unknown"
	current_session.Models = append(current_session.Models, new_model)
	// Return good status coode
	c.JSON(http.StatusOK, new_model)
}

// infer: infers on a defined dataset with a defined model, and returns the results.
func infer(c *gin.Context) {
	// take TOML with info on model and dataset. With only model, use the same dataset. Return results from inference
	// Load up our chosen model (make sure it exists, get its full path on disk),
	log.Println("Inferring with model...")
	var infer_body session.InferConfig
	if err := c.BindTOML(&infer_body); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	//load up our dataset (if available - otherwise use the dataset used to train it. If that's not available in the model information, return an error)
	dataset_path := ""
	for _, dataset := range current_session.Datasets {
		if dataset.ID == infer_body.DatasetID {
			dataset_path = dataset.Path
			break
		}
		response := "dataset not found, id: " + infer_body.DatasetID
		c.JSON(http.StatusBadRequest, response)
	}

	model_path := ""
	model_infer_name := ""
	var model_features []string
	for _, model := range current_session.Models {
		if model.ID == infer_body.ModelID {
			model_path = model.Path
			model_infer_name = model.InferName
			model_features = model.Features
			break
		}
		response := "dataset not found, id: " + infer_body.ModelID
		c.JSON(http.StatusBadRequest, response)
	}
	infertomlpath := generateInferenceTOML(dataset_path, model_path, model_infer_name, "infer", model_features)

	// Infer, store the results
	channel_status := make(chan []byte)
	go func() {
		cmd := exec.Command("python", "/app/random_forest/main.py", infertomlpath)
		log.Println("Starting inference...")
		output, err := cmd.Output()
		if err != nil {
			log.Println(err)
		}
		current_session.Tasks[len(current_session.Tasks)-1].Status = "Complete"
		channel_status <- output
	}()
	response := strings.Split(strings.TrimSuffix(string(<-channel_status), "\n"), "\n")
	log.Println(response)
	trainedprecision, _ := strconv.ParseFloat(response[0], 64)
	trainedrecall, _ := strconv.ParseFloat(response[1], 64)
	// Return a good status to the user.
	c.JSON(http.StatusOK, session.InferenceResponse{TrainedPrecision: trainedprecision, TrainedRecall: trainedrecall})
	// Return the inference results and a good status code
}

func generateTrainingTOML(filepath string, modelpath string, infername string, tasktype string, config []string, depth int, trees int, samplesplit int, fraction float64, datasplit float64, showunoptmizied bool) string {
	path := "/storage/train.toml"
	trainingToml := session.RandomForestTrainingConfig{TaskType: tasktype, FilePath: filepath, Features: config, InferenceName: infername, ModelPath: modelpath, NTrees: trees, SampleSplit: samplesplit, MaxDepth: depth, FeaturesFraction: fraction, DataSplit: datasplit, ShowUnoptimzied: showunoptmizied}
	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(trainingToml)
	if err != nil {
		log.Fatal(err)
	}
	// Write to file
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		path = "./train.toml"
		f, err = os.Create(path)
		log.Println("storage volume not available, writing to local directory")
		if err != nil {
			log.Fatal(err)
		}

	}
	_, err = f.Write(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	return path
}

func generateInferenceTOML(filepath string, modelpath string, infername string, tasktype string, config []string) string {
	path := "/storage/infer.toml"
	trainingToml := session.RandomForestTrainingConfig{TaskType: tasktype, FilePath: filepath, Features: config, InferenceName: infername, ModelPath: modelpath}
	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(trainingToml)
	if err != nil {
		log.Fatal(err)
	}
	// Write to file
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		path = "./infer.toml"
		f, err = os.Create(path)
		log.Println("storage volume not available, writing to local directory")
		if err != nil {
			log.Fatal(err)
		}

	}
	_, err = f.Write(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	return path
}

func generateTreesTOML(modelpath string, tasktype string) string {
	path := "/storage/trees.toml"
	treesTOML := session.RandomForestTrainingConfig{TaskType: tasktype, ModelPath: modelpath}
	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(treesTOML)
	if err != nil {
		log.Fatal(err)
	}
	// Write to file
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		path = "./infer.toml"
		f, err = os.Create(path)
		log.Println("storage volume not available, writing to local directory")
		if err != nil {
			log.Fatal(err)
		}

	}
	_, err = f.Write(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	return path
}

// main: our main function
func main() {
	gin.SetMode(gin.ReleaseMode)
	// Startup Tasks: Check for existing models, datasets in the mounted volume, add then to our model and dataset list
	volumePath := os.Getenv("VOLUMEPATH")
	if volumePath == "" {
		volumePath = "/storage"
	}
	//Create Router
	router := setupRouter()
	// Set up session variables
	current_session.Setup(volumePath)
	os.Setenv("PATH", os.Getenv("PATH")+":/home/oddforest/.pyenv/shims/")
	//Router Run
	router.Run(":9001")
}
