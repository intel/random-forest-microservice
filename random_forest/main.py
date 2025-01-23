# Copyright (c) 2025 Intel Corporation.

# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:

# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.

# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

"""
Code that houses the class that creates and uses the random forest classifier 
"""
import argparse
import tomllib
#TODO: Look into INtel Distribution of Modin
import pandas as pd 
#TODO: Look into intel extension for scikit-learn 
from sklearn.metrics import precision_score, recall_score
from sklearn.model_selection import train_test_split
import joblib

import odd, randomforestclassifier

# argparse for command-line arguments like config file location
parser = argparse.ArgumentParser(
    prog="Intel Optimized Data Discretion Random Forest Classifier (Intel OptForest)",
    description="Create a random forest classifier model with data automatically optimzed for data discretion. \
        Optimzied Data Discretion bins your input data to create a smaller, faster, and highly accurate model tailored to your use case.",
    epilog= "Copyright 2023 Intel Corporation. Distributed under (license)")
parser.add_argument('config_filepath')

# Parse provided TOML config file
def parse_config_file(filepath):
    try:
        try:
            with open(filepath, 'rb') as f:
                config = tomllib.load(f)
                return config
        except tomllib.TOMLDecodeError:
            #TODO log tomllib decode error in global logger
            print("Error decoding TOML file - please check provided file.")
            exit()
    except OSError:
        print("Problem opening provided file: " + filepath)
        exit()

# Handle reading input data from .csv
def read_input_data(filepath):
    try:
        f = pd.read_csv(filepath)
        return f
    except Exception as e:
        #TODO Log error
        print("Error with provided filepath + ", e)
        exit()

def inference(config, rf=None, splitdata=None):
    try:
        # load in model
        if rf == None:
            try:
                rf = joblib.load(config["path"])
            except FileNotFoundError:
                print("model not found, exiting...")
                exit
        if rf == "empty":
            print(0)
            print(0)
            return

        if not isinstance(splitdata, pd.DataFrame):
            d = read_input_data(config["input_data"])
        else:
            d = splitdata

        features = config["features"]
        yhat = rf.predict(d[features])
        d['yhat'] = yhat 

        # Measurring accuracy
        print(f"{precision_score(d[config['y_axis']], d['yhat'])}")
        print(f"{recall_score(d[config['y_axis']], d['yhat'])}")

    except Exception as e:
        print(e)

def show_trees(config, rf=None):
    if rf == None:
        try:
            rf = joblib.load(config["path"])
            rf.print_trees()
        except FileNotFoundError:
            print("model not found, exiting...")
            exit
    

if __name__ == '__main__':
    # Reading data for classification
    
    args = parser.parse_args()
    config = parse_config_file(args.config_filepath)
    print(config)
    try:
        if config["task"] == "train": 
            d = read_input_data(config["input_data"])
                # Setting the features used
            features = config["features"]
            print(features)
            # Get our train/test split
            d_test, d_train = train_test_split(d, test_size=config["data_split"])
            try:
                opt_array={i: odd.automated_optimal_binning(d_train[i].values)[2] for i in features}
            except Exception as e:
                    print(e)
            print("data binned")
            # Create the random forest without optimized data
            if config["show_unoptimized"] == True:
                try: 
                    rf_unopt = randomforestclassifier.RandomForestClassifier(
                        Y=d[config['y_axis']], 
                        X=d[features],
                        min_samples_split=config["min_samples_split"],
                        max_depth=config["max_depth"],
                        n_trees=config["n_trees"],
                        X_features_fraction=config["x_features_fraction"],
                        opti_array=opt_array,
                    )
                except Exception as e:
                    print(e)
                    exit()
                print("rf classifier created")
                print("rf forest growing")
                # Growing the random forest 
                rf_unopt.grow_random_forest()
            else:
                rf_unopt = "empty"
            # Create the random forest for optimized data
            try: 
                rf = randomforestclassifier.RandomForestClassifier(
                    Y=d[config['y_axis']], 
                    X=d[features],
                    min_samples_split=config["min_samples_split"],
                    max_depth=config["max_depth"],
                    n_trees=config["n_trees"],
                    X_features_fraction=config["x_features_fraction"],
                    opti_array=opt_array,
                )
            except Exception as e:
                print(e)
                exit()
            print("rf classifier created")
            print("rf forest growing")
            # Growing the random forest 
            rf.grow_random_forest()

            # Save random forest to disk
            try:
                print("saving classifier to disk")
                try:
                    joblib.dump(rf, config["path"], 3,5)
                except Exception as e:
                    print(e)
            except FileNotFoundError as e:
                try:
                    joblib.dump(rf, "test.model", 3,5)
                    print("error", e)
                except Exception as e:
                    print(e)
            # Printing out the trees 
            #rf.print_trees()
            rf.print_trees
            inference(config, rf_unopt, d_train)
            inference(config, rf_unopt, d_test)
            inference(config, rf, d_train)
            inference(config, rf, d_test)

        # Making predictions
        if config["task"] == "infer":
            inference(config)

        if config["task"] == "show_trees":
            show_trees(config)

        if config["task"] != "infer" and config["task"] != "train" and config["tast"] != "show_trees":
            print("Incorrect task type selected. Please choose from infer or train")
        

    except KeyError as e:
        print("Issue with provided config - please check your provided TOML configuration. " + e)
        exit()