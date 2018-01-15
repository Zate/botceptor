# botceptor
### Tool to review and remove bot accounts from your twitter account.

## Getting Started

### Clone the repo
```
git clone git@github.com:Zate/botceptor.git
```
Rename skel.secrets.yaml to .secrets.yaml
```
cd botceptor
mv skel.secrets.yaml .secrets.yaml
```
### Set up your twitter app for auth

Make sure you are logged into twitter.com with the account you want to work with.

Go to [the twitter apps page](https://apps.twitter.com/app/new)

Fill in Name, Description, Website and leave Callback URL blank.  Accept the agreement and create your twitter application.

Go to the Keys and Access Tokens tab.

Get the Consumer Key (API Key) and Consumer Secret (API Secret) and put them in .secrets.yaml (create this from skel.secrets.yaml) as `consumerkey:` and `consumersecret:`

Down under Your Access Token, click Create my access token which will give the app you have created access to your twitter account.

It will generate an Access Token and an Access Token Secret.  Place those in .secrets.yaml as `accesskey:` and  `accesssecret:`

#### **SECURE THIS FILE**.  *It will allow someone to access your twitter account via the API.*

### Running the app
By default the app is setup to run over the last 20 followers and output if they are a bot or not.  You will need to adjust some comments to make it work on all your followers and to make it take action.

Look for `followers, _, err := client.Followers.List(&twitter.FollowerListParams{Cursor: cursor, Count: 20})` and adjust the 20 to be up to 200.  200 is the maximum the twitter API allows.

If you want to process all followers, find `cursor = 0` and comment it out and uncomment `//cursor = followers.NextCursor` below it.  You will also need to find `//time.Sleep(time.Duration(60) * time.Second)` and uncomment that so that you do not hit the twitter API limit (15 requests in 15mins)

> Running over all followers can take a while, 1 minute per 200 followers to be exact.

Find `log.Printf("%v is a bot : %v", allfollowers[x].ScreenName, res)` and uncomment the lines below it down to `} else {` to turn on the block/unblock action to remove a follower.  If you do not do that, it will only tell you whether it thinks the user is a bot or not.

In it's default form, it classifies anyone with a score over 0.02 as a bot.  You can change this to suite yourself by modifying `if res > 0.02 {`

If you are ready, simple run the program by:

`go build && ./botceptor` on linux/mac and `go build && ./botceptor.exe` on Windows.

> This is still very much a WIP, so much of this might not work.  YMMV.  Not responsible if you nuke all your twitter followers.


## Creating your Machine Learning Model
```
go get github.com/xlvector/hector
go install github.com/xlvector/hector/hectorcv
hectorcv --method [Method] --train [Data Path] --cv 10
```

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
```
cd src
go build hector-run.go
./hector-run --method [Method] --train [Data Path] --test [Data Path]
```
Above methods will direct train algorithm on train dataset and then test on test dataset. If you want to train algorithm and get the model file, you can run it by following steps:
```
./hector-run --method [Method] --action train --train [Data Path] --model [Model Path]
```
Then, you can use model file to test any test dataset:
```
./hector-run --method [Method] --action test --test [Data Path] --model [Model Path]
```
## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

## Acknowledgments

* Thankyou to the creator of [hector](https://github.com/xlvector/hector) which I have used for the machine learning code
* Thankyou to the create or [go-twitter](https://github.com/dghubble/go-twitter) which allows me to connect to twitter.
