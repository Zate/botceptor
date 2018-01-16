package hector

import (
	"flag"
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"time"

	"github.com/xlvector/hector/algo"
	"github.com/xlvector/hector/ann"
	"github.com/xlvector/hector/dt"
	"github.com/xlvector/hector/fm"
	"github.com/xlvector/hector/gp"
	"github.com/xlvector/hector/lr"
	"github.com/xlvector/hector/sa"
	"github.com/xlvector/hector/svm"
)

func GetMutliClassClassifier(method string) algo.MultiClassClassifier {
	rand.Seed(time.Now().UTC().UnixNano())
	var classifier algo.MultiClassClassifier

	if method == "rf" {
		classifier = &(dt.RandomForest{})
	} else if method == "cart" {
		classifier = &(dt.CART{})
	} else if method == "rdt" {
		classifier = &(dt.RandomDecisionTree{})
	} else if method == "knn" {
		classifier = &(svm.KNN{})
	} else if method == "ann" {
		classifier = &(ann.NeuralNetwork{})
	}
	return classifier
}

func GetClassifier(method string) algo.Classifier {
	rand.Seed(time.Now().UTC().UnixNano())
	var classifier algo.Classifier

	if method == "lr" {
		classifier = &(lr.LogisticRegression{})
	} else if method == "ftrl" {
		classifier = &(lr.FTRLLogisticRegression{})
	} else if method == "ep" {
		classifier = &(lr.EPLogisticRegression{})
	} else if method == "rdt" {
		classifier = &(dt.RandomDecisionTree{})
	} else if method == "cart" {
		classifier = &(dt.CART{})
	} else if method == "cart-regression" {
		classifier = &(dt.RegressionTree{})
	} else if method == "rf" {
		classifier = &(dt.RandomForest{})
	} else if method == "fm" {
		classifier = &(fm.FactorizeMachine{})
	} else if method == "sa" {
		classifier = &(sa.SAOptAUC{})
	} else if method == "gbdt" {
		classifier = &(dt.GBDT{})
	} else if method == "svm" {
		classifier = &(svm.SVM{})
	} else if method == "linear_svm" {
		classifier = &(svm.LinearSVM{})
	} else if method == "l1vm" {
		classifier = &(svm.L1VM{})
	} else if method == "knn" {
		classifier = &(svm.KNN{})
	} else if method == "ann" {
		classifier = &(ann.NeuralNetwork{})
	} else if method == "lr_owlqn" {
		classifier = &(lr.LROWLQN{})
	} else {
		classifier = &(lr.LogisticRegression{})
	}
	return classifier
}

func GetRegressor(method string) algo.Regressor {
	rand.Seed(time.Now().UTC().UnixNano())

	var regressor algo.Regressor

	if method == "gp" {
		regressor = &(gp.GaussianProcess{})
	}
	return regressor
}

func PrepareParams() (string, string, string, string, map[string]string) {
	params := make(map[string]string)
	train_path := flag.String("train", "train.tsv", "path of training file")
	test_path := flag.String("test", "test.tsv", "path of testing file")
	pred_path := flag.String("pred", "", "path of pred file")
	output := flag.String("output", "", "output file path")
	verbose := flag.Int("v", 0, "verbose output if 1")
	learning_rate := flag.String("learning-rate", "0.01", "learning rate")
	learning_rate_discount := flag.String("learning-rate-discount", "1.0", "discount rate of learning rate per training step")
	regularization := flag.String("regularization", "0.01", "regularization")
	alpha := flag.String("alpha", "0.1", "alpha of ftrl")
	beta := flag.String("beta", "1", "beta of ftrl")
	c := flag.String("c", "1", "C in svm")
	e := flag.String("e", "0.01", "stop threshold")
	lambda1 := flag.String("lambda1", "0.1", "lambda1 of ftrl")
	lambda2 := flag.String("lambda2", "0.1", "lambda2 of ftrl")
	tree_count := flag.String("tree-count", "10", "tree count in rdt/rf")
	feature_count := flag.String("feature-count", "1.0", "feature count in rdt/rf")
	gini := flag.String("gini", "1.0", "gini threshold, between (0, 0.5]")
	min_leaf_size := flag.String("min-leaf-size", "10", "min leaf size in dt")
	max_depth := flag.String("max-depth", "10", "max depth of dt")
	factors := flag.String("factors", "10", "factor number in factorized machine")
	steps := flag.Int("steps", 1, "steps before convergent")
	global := flag.Int64("global", -1, "feature id of global bias")
	method := flag.String("method", "lr", "algorithm name")
	cv := flag.Int("cv", 7, "cross validation folder count")
	k := flag.String("k", "3", "neighborhood size of knn")
	radius := flag.String("radius", "1.0", "radius of RBF kernel")
	sv := flag.String("sv", "8", "support vector count for l1vm")
	hidden := flag.Int64("hidden", 1, "hidden neuron number")
	profile := flag.String("profile", "", "profile file name")
	model := flag.String("model", "", "model file name")
	action := flag.String("action", "", "train or test, do both if action is empty string")
	core := flag.Int("core", 1, "core number when run program")
	dt_sample_ratio := flag.String("dt-sample-ratio", "1.0", "sampling ratio when split feature in decision tree")
	dim := flag.String("dim", "1", "input space dimension")
	port := flag.String("port", "8080", "port")

	flag.Parse()
	runtime.GOMAXPROCS(*core)
	fmt.Println(*train_path)
	fmt.Println(*test_path)
	fmt.Println(*method)
	params["port"] = *port
	params["verbose"] = strconv.FormatInt(int64(*verbose), 10)
	params["learning-rate"] = *learning_rate
	params["learning-rate-discount"] = *learning_rate_discount
	params["regularization"] = *regularization
	params["alpha"] = *alpha
	params["beta"] = *beta
	params["lambda1"] = *lambda1
	params["lambda2"] = *lambda2
	params["tree-count"] = *tree_count
	params["feature-count"] = *feature_count
	params["max-depth"] = *max_depth
	params["min-leaf-size"] = *min_leaf_size
	params["steps"] = strconv.FormatInt(int64(*steps), 10)
	params["global"] = strconv.FormatInt(*global, 10)
	params["gini"] = *gini
	params["factors"] = *factors
	params["output"] = *output
	params["c"] = *c
	params["e"] = *e
	params["k"] = *k
	params["cv"] = strconv.FormatInt(int64(*cv), 10)
	params["radius"] = *radius
	params["sv"] = *sv
	params["hidden"] = strconv.FormatInt(int64(*hidden), 10)
	params["profile"] = *profile
	params["action"] = *action
	params["model"] = *model
	params["method"] = *method
	params["dt-sample-ratio"] = *dt_sample_ratio
	params["dim"] = *dim

	fmt.Println(params)
	return *train_path, *test_path, *pred_path, *method, params
}
