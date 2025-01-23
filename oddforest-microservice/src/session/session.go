package session

import (
	"fmt"
	"log"
	"os"
)

// Internal data types to hold session, model, dataset, result, and task data during runtime.
type Session struct {
	Models   []Model
	Datasets []Dataset
	Results  []Result
	Tasks    []Task
}

type Model struct {
	Name              string
	ID                string
	TrainedDataset    string
	Features          []string
	ID_num            int
	Path              string
	InferName         string
	UnoptValAccuracy  float64
	UnoptValRecall    float64
	UnoptTestAccuracy float64
	UnoptTestRecall   float64
	ValAccuracy       float64
	ValRecall         float64
	TestAccuracy      float64
	TestRecall        float64
}

type Dataset struct {
	Name       string
	ID         string
	Datapoints int
	ID_num     int
	Path       string
}

type Result struct {
	ModelID   string
	DatasetID string
	Tree      string
	Precision string
	Recall    string
}

type Task struct {
	ID      string
	ModelID string
	Status  string
	ID_num  int
}

// Validation structures to capture our incoming TOML data
type TrainingConfig struct {
	Name             string   `form:"name" toml:"name" binding:"required"`
	DatasetID        string   `form:"datasetid" toml:"dataset_id" binding:"required"`
	InferName        string   `form:"inferencename" toml:"infer_name" binding:"required"`
	Features         []string `toml:"features" binding:"required"`
	NTrees           int      `toml:"n_trees"`
	SampleSplit      int      `toml:"min_samples_split"`
	MaxDepth         int      `toml:"max_depth"`
	FeaturesFraction float64  `toml:"x_features_fraction"`
	DataSplit        float64  `toml:"data_split"`
	ShowUnoptimzied  bool     `toml:"show_unoptimized"`
}

type UploadConfig struct {
	Name string `form:"name" toml:"name" binding:"required"`
}

type InferConfig struct {
	ModelID   string `form:"modelid" toml:"modelid" binding:"required"`
	DatasetID string `form:"datasetid" toml:"datasetid" binding:"required"`
}

type DownloadConfig struct {
	ModelID string `form:"modelid" toml:"modelid" binding:"required"`
}

// Structs for our structured responses to the client
type TrainingResponse struct {
	Response string
}

type InferenceResponse struct {
	TrainedPrecision float64
	TrainedRecall    float64
}

// Struct for holding toml info to submit to the training script
type RandomForestTrainingConfig struct {
	TaskType         string   `toml:"task"`
	FilePath         string   `toml:"input_data"`
	Features         []string `toml:"features"`
	InferenceName    string   `toml:"y_axis"`
	ModelPath        string   `toml:"path"`
	NTrees           int      `toml:"n_trees"`
	SampleSplit      int      `toml:"min_samples_split"`
	MaxDepth         int      `toml:"max_depth"`
	FeaturesFraction float64  `toml:"x_features_fraction"`
	DataSplit        float64  `toml:"data_split"`
	ShowUnoptimzied  bool     `toml:"show_unoptimized"`
}

func (self *Session) Setup(volumePath string) {
	// Check along the volumePath for any existing models, datasets
	models_exists := false
	datasets_exists := false
	default_feature := []string{"Unknown"}
	files, err := os.ReadDir(volumePath)
	if err != nil {
		log.Print(err)
	} else {
		for _, file := range files {
			if file.Name() == "models" {
				models_exists = true
				models, err := os.ReadDir("models")
				if err != nil {
					log.Print(err)
				}
				for i, model := range models {
					self.Models = append(self.Models, Model{model.Name(), "m" + fmt.Sprint(i), "Unknown", default_feature, i, volumePath + "models" + model.Name() + ".model", "Unknown", 0, 0, 0, 0, 0, 0, 0, 0})
				}
			}
			if file.Name() == "datasets" {
				datasets_exists = true
				datasets, err := os.ReadDir("datasets")
				if err != nil {
					log.Print(err)
				}
				for i, dataset := range datasets {
					self.Datasets = append(self.Datasets, Dataset{dataset.Name(), "d" + fmt.Sprint(i), 0, i, volumePath + "datasets" + dataset.Name() + ".csv"})
				}
			}
		}
	}

	if !models_exists {
		err := os.Mkdir(volumePath+"/models", 0777)
		if err != nil && !os.IsExist(err) {
			log.Print(err)
		}
	}

	if !datasets_exists {
		err := os.Mkdir(volumePath+"/datasets", 0777)
		if err != nil && !os.IsExist(err) {
			log.Print(err)
		}
	}
}
