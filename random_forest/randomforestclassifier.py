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
import pandas as pd 
from collections import Counter
from tqdm import tqdm

import randomforesttree

class RandomForestClassifier():
    """
    Class that creates a random forest for classification problems
    """
    def __init__(
        self,
        Y: list,
        X: pd.DataFrame,
        min_samples_split=None,
        max_depth=None,
        n_trees=None,
        X_features_fraction=None,
        X_obs_fraction=None,
        opti_array=None,
        
    ):  
        # Saving the data for the random forest
        self.Y = Y 
        self.X = X

        # Saving the hyper parameters
        self.min_samples_split = min_samples_split if min_samples_split else 20
        self.max_depth = max_depth if max_depth else 5

        # Saving the final feature list 
        self.features = list(X.columns)

        # Getting the number of features 
        self.n_features = len(self.features)

        # Saving the hyper parameters specific to the random forest 
        self.n_trees = n_trees if n_trees is not None else 30
        self.X_features_fraction = X_features_fraction if X_features_fraction is not None else 1.0
        self.X_obs_fraction = X_obs_fraction if X_obs_fraction is not None else 1.0
        self.opti_array=opti_array

    def bootstrap_sample(self):
        """
        Function that creates a bootstraped sample with the class instance parameters 
        """
        # Sampling the number of rows with repetition
        Xbootstrap = self.X.sample(frac=self.X_obs_fraction, replace=True) 

        # Getting the index of samples 
        indexes = Xbootstrap.index

        # Getting the corresponding Y variables
        Ybootstrap = [self.Y[x] for x in indexes]

        # Droping the index of X 
        Xbootstrap.reset_index(inplace=True, drop=True)

        # Returning the X, Y pair
        return Xbootstrap, Ybootstrap

    def grow_random_forest(self):
        """
        Main method of the class; Creates **n_trees** random trees
        """
        # List to hold trees in 
        random_forest = []

        # Iterating 
        for _ in tqdm(range(self.n_trees), disable=True):
            # Getting the bootstrapped sample
            X, Y = self.bootstrap_sample()
            
            # Initiating the random tree
            tree = randomforesttree.RandomForestTree(
                Y=Y, 
                X=X, 
                min_samples_split=self.min_samples_split,
                max_depth=self.max_depth,
                X_features_fraction=self.X_features_fraction,
                opti_array=self.opti_array
                )

            # Growing the tree
            tree.grow_tree()

            # Appending the tree to the list of trees (the forest)
            random_forest.append(tree)
        
        # Saving the random forest list to memory
        self.random_forest = random_forest

    def print_trees(self):
        """
        Method to print out all the grown trees in the classifier 
        """
        for i in range(self.n_trees):
            print("------ \n")
            print(f"Tree number: {i + 1} \n")
            self.random_forest[i].print_tree()
            print("------ \n")

    def tree_predictions(self, X: pd.DataFrame) -> list:
        """
        Method to get the predictions from all the trees 
        """
        predictions = []
        for i in range(self.n_trees):
            yhat = self.random_forest[i].predict(X)
            
            # Apending to prediction placeholder
            predictions.append(yhat)
        
        # Returning the prediction list 
        return predictions
    
    def predict(self, X: pd.DataFrame) -> list:
        """
        Method to get the final prediction of the whole random forest 
        """
        # Getting the individual tree predictions
        yhat = self.tree_predictions(X)

        # Saving the number of obs in X 
        n = X.shape[0]

        # Getting the majority vote of each coordinate of the prediction list 
        yhat_final = []

        for i in range(n):
            yhat_obs = [x[i] for x in yhat]

            # Getting the most frequent entry 
            counts = Counter(yhat_obs)
            most_common = counts.most_common(1)[0][0]

            # Appending the most common entry to final yhat list 
            yhat_final.append(most_common)
        
        # Returning the final predictions 
        return yhat_final
