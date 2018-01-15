# botceptor
#### Tool to review and remove bot accounts from your twitter account.

Insert witty stuff about the tool here.





#### Creating your Machine Learning Model

`go get github.com/xlvector/hector`

`go install github.com/xlvector/hector/hectorcv`

`hectorcv --method [Method] --train [Data Path] --cv 10`


Method can be :

* `lr` : logistic regression with SGD and L2 regularization.
* `ftrl` : FTRL-proximal logistic regreesion with L1 regularization. Please review this paper for more details "Ad Click Prediction: a View from the Trenches".
* `ep` : bayesian logistic regression with expectation propagation. Please review this paper for more details "Web-Scale Bayesian Click-Through Rate Prediction for Sponsored Search Advertising in Microsoft’s Bing Search Engine"
* `fm` : factorization machine
* `cart` : classifiaction tree
* `cart-regression` : regression tree
* `rf` : random forest
* `rdt` : random decision trees
* `gbdt` : gradient boosting decisio tree
* `linear-svm` : linear svm with L1 regularization
* `svm` : svm optimizaed by SMO (current, its linear svm)
* `l1vm` : vector machine with L1 regularization by RBF kernel
* `knn` : k-nearest neighbor classification

hector-run.go will help you train one algorithm on train dataset and test it on test dataset, you can run it by following steps:

`cd src`

`go build hector-run.go`

`./hector-run --method [Method] --train [Data Path] --test [Data Path]`

Above methods will direct train algorithm on train dataset and then test on test dataset. If you want to train algorithm and get the model file, you can run it by following steps:

`./hector-run --method [Method] --action train --train [Data Path] --model [Model Path]`

Then, you can use model file to test any test dataset:

`./hector-run --method [Method] --action test --test [Data Path] --model [Model Path]`
