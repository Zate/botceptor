package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	// "github.com/sirupsen/logrus"
	// prefixed "github.com/x-cray/logrus-prefixed-formatter"

	"gopkg.in/yaml.v2"

	"github.com/dghubble/oauth1"
	"github.com/xlvector/hector"
	"github.com/xlvector/hector/algo"
	"github.com/xlvector/hector/core"
	"github.com/zate/go-twitter/twitter"
)

func newTrue() *bool {
	b := true
	return &b
}

//var log = logrus.New()

// CheckErr to handle errors
func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type apikeys struct {
	AccessKey      string
	AccessSecret   string
	ConsumerKey    string
	ConsumerSecret string
	BotomKey       string
}

func (a *apikeys) getAPIKeys(filename string) *apikeys {
	yamlFile, err := ioutil.ReadFile(filename)
	CheckErr(err)
	err = yaml.Unmarshal(yamlFile, a)
	CheckErr(err)
	return a
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

func last200() {

	var k apikeys
	var consumerkey string
	var accesskey string
	var consumersecret string
	var accesssecret string
	k.getAPIKeys(".secrets.yaml")
	consumerkey = k.ConsumerKey
	accesskey = k.AccessKey
	consumersecret = k.ConsumerSecret
	accesssecret = k.AccessSecret

	// Pass in your consumer key (API Key) and your Consumer Secret (API Secret)
	config := oauth1.NewConfig(consumerkey, consumersecret)
	// Pass in your Access Token and your Access Token Secret
	token := oauth1.NewToken(accesskey, accesssecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	if FollowerCount > 200 {
		FollowerCount = 200
	}
	var cursor int64
	cursor = -1
	var allfollowers []twitter.User
	for cursor != 0 {

		fParams := &twitter.FollowerListParams{}
		fParams.Cursor = cursor
		fParams.Count = FollowerCount
		fParams.ScreenName = ScreenName
		followers, _, err := client.Followers.List(fParams)
		CheckErr(err)
		log.Println(followers)
		//cursor = 0 // if cursor is 0, then it will only run through 1 lot of followers from twitter of the size you specify above (default 200)
		cursor = followers.NextCursor // If you comment out cursor = 0 and uncomment this, it will iterate through ALL followers in batches as sized above (default 200)
		log.Println(ScreenName)
		for k := range followers.Users {

			allfollowers = append(allfollowers, followers.Users[k])
		}

		// Need this to not hit twitter rate limits.  If you are using NextCursor above, them have this uncommented also so that you will only make one request for a batch of followers per minute.  This will be right on with the API rate limit of 15 in 15 mins.
		//log.Println(cursor)
		if cursor != 0 && FollowerCount > 199 {
			log.Println("got here")
			time.Sleep(time.Duration(60) * time.Second)
		}
	}
	if GetTrainingData == false {
		do200(allfollowers, client)
	} else {
		log200(allfollowers, client)
	}
}

func do200(allfollowers []twitter.User, client *twitter.Client) {
	log.Println("do200")
	var fc []int
	var listed []int
	var tweets []int
	var favc []int

	// f, err := os.OpenFile("followers.t", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	// CheckErr(err)
	// defer f.Close()

	// w := bufio.NewWriter(f)

	for x := range allfollowers {
		fc = append(fc, allfollowers[x].FollowersCount)
		listed = append(listed, allfollowers[x].ListedCount)
		tweets = append(tweets, allfollowers[x].StatusesCount)
		favc = append(favc, allfollowers[x].FavouritesCount)
	}
	fcNumbers := sort.IntSlice(fc)
	sort.Sort(fcNumbers)
	fcMax := fcNumbers[len(fcNumbers)-1]
	listedNumbers := sort.IntSlice(listed)
	sort.Sort(listedNumbers)
	listedMax := listedNumbers[len(listedNumbers)-1]
	tweetsNumbers := sort.IntSlice(tweets)
	sort.Sort(tweetsNumbers)
	tweetsMax := tweetsNumbers[len(tweetsNumbers)-1]
	favcNumbers := sort.IntSlice(favc)
	sort.Sort(favcNumbers)
	favcMax := favcNumbers[len(favcNumbers)-1]

	for x := range allfollowers {
		idxfavs := float64(allfollowers[x].FavouritesCount) / float64(favcMax)
		idxfollowers := float64(allfollowers[x].FollowersCount) / float64(fcMax)
		idxlisted := float64(allfollowers[x].ListedCount) / float64(listedMax)
		idxtweets := float64(allfollowers[x].StatusesCount) / float64(tweetsMax)
		idxprofile := b2i(allfollowers[x].DefaultProfile)
		idxprofileimage := b2i(allfollowers[x].DefaultProfileImage)
		idxfollowing := b2i(allfollowers[x].Following)
		idxverified := b2i(allfollowers[x].Verified)
		fdata := fmt.Sprintf("0	1:%v 2:%v 3:%v 4:%v 5:%v 6:%v 7:%v 8:%v",
			idxfavs,
			idxfollowers,
			idxlisted,
			idxtweets,
			idxprofile,
			idxprofileimage,
			idxfollowing,
			idxverified)
		// fmt.Fprintf(w, "0	1:%v 2:%v 3:%v 4:%v 5:%v 6:%v 7:%v 8:%v\n",
		// 	idxfavs,
		// 	idxfollowers,
		// 	idxlisted,
		// 	idxtweets,
		// 	idxprofile,
		// 	idxprofileimage,
		// 	idxfollowing,
		// 	idxverified)
		_, _, _, method, params := doParams()
		model, _ := params["model"]
		c := &Classifier{
			classifier: hector.GetClassifier(method),
		}
		c.classifier.LoadModel(model)

		res := c.testFollower(fdata)
		if res > 0.02 && allfollowers[x].Verified != true && allfollowers[x].Following != true {
			log.Printf("%v is a bot : %v", allfollowers[x].ScreenName, res)
			if DeFang == false {
				user, resp, _ := client.Block.Create(&twitter.BlockUserParams{ScreenName: allfollowers[x].ScreenName})
				//log.Println(resp)
				if resp.StatusCode == 200 {
					log.Printf("%v was blocked", user.ScreenName)
				}
				user, resp, _ = client.Block.Destroy(&twitter.BlockUserParams{ScreenName: allfollowers[x].ScreenName})
				if resp.StatusCode == 200 {
					log.Printf("%v was unblocked", user.ScreenName)
				}
			}
		} else {
			log.Printf("%v not bot : %v Verified: %v Following: %v", allfollowers[x].ScreenName, res, allfollowers[x].Verified, allfollowers[x].Following)
		}
	}
	// w.Flush()
	log.Println("Run done")
	// 1 - ScreenName - unique - dont put these in training
	// 2 - user id - unique -  dont put these in training
	// 3 - number of accounts they follow ? need to index this ?
	// 4 - number of items they have favorited indexed to number of accounts they follow
	// 5 - number of accounts that follow them indexed to number of accounts they follow
	// 6 - number of times they are on a list indexed to number of accounts they follow
	// 7 - number of tweets indexed to number of accounts they follow
	// 8 - default profile? - binary
	// 9 - default profile image? - binary
	// 10 - am I following them? - binary
	// 11 - are they a verfied account? - binary

}

func log200(allfollowers []twitter.User, client *twitter.Client) {
	log.Println("log200")
	var fc []int
	var listed []int
	var tweets []int
	var favc []int
	var langc []int
	var desc []int

	f, err := os.OpenFile("followers_new_training.t", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	CheckErr(err)
	defer f.Close()

	w := bufio.NewWriter(f)

	for x := range allfollowers {
		fc = append(fc, allfollowers[x].FollowersCount)
		listed = append(listed, allfollowers[x].ListedCount)
		tweets = append(tweets, allfollowers[x].StatusesCount)
		favc = append(favc, allfollowers[x].FavouritesCount)
		langc = append(langc, int(new(big.Int).SetBytes([]byte(allfollowers[x].Lang)).Uint64()))
		desc = append(desc, int(new(big.Int).SetBytes([]byte(allfollowers[x].Description)).Uint64()))
	}
	fcNumbers := sort.IntSlice(fc)
	sort.Sort(fcNumbers)
	fcMax := fcNumbers[len(fcNumbers)-1]
	listedNumbers := sort.IntSlice(listed)
	sort.Sort(listedNumbers)
	listedMax := listedNumbers[len(listedNumbers)-1]
	tweetsNumbers := sort.IntSlice(tweets)
	sort.Sort(tweetsNumbers)
	tweetsMax := tweetsNumbers[len(tweetsNumbers)-1]
	favcNumbers := sort.IntSlice(favc)
	sort.Sort(favcNumbers)
	favcMax := favcNumbers[len(favcNumbers)-1]
	langcNumbers := sort.IntSlice(langc)
	sort.Sort(langcNumbers)
	langcMax := langcNumbers[len(langcNumbers)-1]
	descNumbers := sort.IntSlice(desc)
	sort.Sort(descNumbers)
	descMax := descNumbers[len(descNumbers)-1]

	for x := range allfollowers {
		idxfavs := math.Log10(float64(allfollowers[x].FavouritesCount)) / math.Log10(float64(favcMax))
		idxfollowers := math.Log10(float64(allfollowers[x].FollowersCount)) / math.Log10(float64(fcMax))
		idxlisted := math.Log10(float64(allfollowers[x].ListedCount)) / math.Log10(float64(listedMax))
		idxtweets := math.Log10(float64(allfollowers[x].StatusesCount)) / math.Log10(float64(tweetsMax))
		idxprofile := (b2i(allfollowers[x].DefaultProfile) + 1) / 2
		idxprofileimage := (b2i(allfollowers[x].DefaultProfileImage) + 1) / 2
		idxfollowing := (b2i(allfollowers[x].Following) + 1) / 2
		idxverified := (b2i(allfollowers[x].Verified) + 1) / 2
		lang := new(big.Int).SetBytes([]byte(allfollowers[x].Lang)).Uint64()
		idxlang := math.Log10(float64(lang)) / math.Log10(float64(langcMax))
		des := new(big.Int).SetBytes([]byte(allfollowers[x].Description)).Uint64()
		//log.Println(des)
		idxdesc := math.Log10(float64(des+1)) / math.Log10(float64(descMax+1))

		// fdata := fmt.Sprintf("0	1:%v 2:%v 3:%v 4:%v 5:%v 6:%v 7:%v 8:%v",
		// 	idxfavs,
		// 	idxfollowers,
		// 	idxlisted,
		// 	idxtweets,
		// 	idxprofile,
		// 	idxprofileimage,
		// 	idxfollowing,
		// 	idxverified)
		// fmt.Fprintf(w, "0	ScreenName:%v 1:%v 2:%v 3:%v 4:%v 5:%v 6:%v 7:%v 8:%v\n",
		// 	allfollowers[x].ScreenName,
		// 	idxfavs,
		// 	idxfollowers,
		// 	idxlisted,
		// 	idxtweets,
		// 	idxprofile,
		// 	idxprofileimage,
		// 	idxfollowing,
		// 	idxverified,
		// 	idxlang)
		fmt.Fprintf(w, "0	ScreenName:%v 1:%v 2:%v 3:%v 4:%v 5:%v 6:%v 7:%v 8:%v 9:%v 10:%v\n",
			allfollowers[x].ScreenName,
			idxfavs,
			idxfollowers,
			idxlisted,
			idxtweets,
			idxprofile,
			idxprofileimage,
			idxfollowing,
			idxverified,
			idxlang,
			idxdesc)
		// _, _, _, method, params := doParams()
		// model, _ := params["model"]
		// c := &Classifier{
		// 	classifier: hector.GetClassifier(method),
		// }
		// c.classifier.LoadModel(model)

		// res := c.testFollower(fdata)
		// if res > 0.02 && allfollowers[x].Verified != true && allfollowers[x].Following != true {
		// 	log.Printf("%v is a bot : %v", allfollowers[x].ScreenName, res)
		// 	user, resp, _ := client.Block.Create(&twitter.BlockUserParams{ScreenName: allfollowers[x].ScreenName})
		// 	//log.Println(resp)
		// 	if resp.StatusCode == 200 {
		// 		log.Printf("%v was blocked", user.ScreenName)
		// 	}
		// 	user, resp, _ = client.Block.Destroy(&twitter.BlockUserParams{ScreenName: allfollowers[x].ScreenName})
		// 	if resp.StatusCode == 200 {
		// 		log.Printf("%v was unblocked", user.ScreenName)
		// 	}
		// } else {
		// 	log.Printf("%v not bot : %v Verified: %v Following: %v", allfollowers[x].ScreenName, res, allfollowers[x].Verified, allfollowers[x].Following)
		// }
	}
	w.Flush()
	log.Println("Run done")
	// 1 - ScreenName - unique - dont put these in training
	// 2 - user id - unique -  dont put these in training
	// 3 - number of accounts they follow ? need to index this ?
	// 4 - number of items they have favorited indexed to number of accounts they follow
	// 5 - number of accounts that follow them indexed to number of accounts they follow
	// 6 - number of times they are on a list indexed to number of accounts they follow
	// 7 - number of tweets indexed to number of accounts they follow
	// 8 - default profile? - binary
	// 9 - default profile image? - binary
	// 10 - am I following them? - binary
	// 11 - are they a verfied account? - binary

}

// Classifier does classifying things.
type Classifier struct {
	classifier algo.Classifier
}

// tests a follower with the rf model.
func (c *Classifier) testFollower(data string) float64 {
	runtime.GOMAXPROCS(runtime.NumCPU())
	tks := strings.Split(data, "\t")
	sample := core.NewSample()
	for i, tk := range tks {
		if i == 0 {
			label, _ := strconv.Atoi(tk)
			sample.Label = label
		} else {
			kv := strings.Split(tk, " ")
			for _, v := range kv {
				meh := strings.Split(v, ":")
				featureID, _ := strconv.ParseInt(meh[0], 10, 64)
				featureValue, _ := strconv.ParseFloat(meh[1], 64)
				f := core.Feature{
					Id:    featureID,
					Value: featureValue,
				}
				sample.AddFeature(f)
			}
		}
	}

	// FUCK YEAH!

	prediction := c.classifier.Predict(sample)
	return prediction
}

// Created this to allow me to control the params as by default, hector likes to take them from the command line.  I hard code them here and use doParams() instead of the inbuilt hector function.
func doParams() (string, string, string, string, map[string]string) {
	params := make(map[string]string)
	trainPath := "followers_training.t"
	testPath := "followers_training.t"
	predPath := ""
	output := ""
	verbose := 0
	learningRate := "0.01"
	learningRateDiscount := "1.0"
	regularization := "0.01"
	alpha := "0.1"
	beta := "1"
	c := "1"
	e := "0.01"
	lambda1 := "0.1"
	lambda2 := "0.1"
	treeCount := "10"
	featureCount := "1.0"
	gini := "1.0"
	minLeafSize := "10"
	maxDepth := "10"
	factors := "10"
	steps := 1
	var global int64 = -1
	method := "rf" // I found this to be most accurate, feel free to change if you've run data through other models with hectorcv
	cv := 7
	k := "3"
	radius := "1.0"
	sv := "8"
	hidden := 1
	profile := ""
	model := "followers.mod"
	action := "test"
	dtSampleRatio := "1.0"
	dim := "1"
	port := "8080"

	params["port"] = port
	params["verbose"] = strconv.FormatInt(int64(verbose), 10)
	params["learning-rate"] = learningRate
	params["learning-rate-discount"] = learningRateDiscount
	params["regularization"] = regularization
	params["alpha"] = alpha
	params["beta"] = beta
	params["lambda1"] = lambda1
	params["lambda2"] = lambda2
	params["tree-count"] = treeCount
	params["feature-count"] = featureCount
	params["max-depth"] = maxDepth
	params["min-leaf-size"] = minLeafSize
	params["steps"] = strconv.FormatInt(int64(steps), 10)
	params["global"] = strconv.FormatInt(global, 10)
	params["gini"] = gini
	params["factors"] = factors
	params["output"] = output
	params["c"] = c
	params["e"] = e
	params["k"] = k
	params["cv"] = strconv.FormatInt(int64(cv), 10)
	params["radius"] = radius
	params["sv"] = sv
	params["hidden"] = strconv.FormatInt(int64(hidden), 10)
	params["profile"] = profile
	params["action"] = action
	params["model"] = model
	params["method"] = method
	params["dt-sample-ratio"] = dtSampleRatio
	params["dim"] = dim

	return trainPath, testPath, predPath, method, params
}

// GetTrainingData is a global command line flag to swap to just gathering training data.
var GetTrainingData bool

// DeFang stops the program from taking action and just reports on the users.
var DeFang bool

// FollowerCount is the number of followers to return and classify.
var FollowerCount int

// ScreenName is the ScreenName of the user to classify.
var ScreenName string

func init() {
	flag.BoolVar(&GetTrainingData, "gtd", false, "Dump followers to a file to be used to generate new training data.")
	flag.BoolVar(&DeFang, "df", true, "Process the followers list but do not block/unblock or take and action.")
	flag.IntVar(&FollowerCount, "fc", 1, "Number of followers to return and classify.")
	flag.StringVar(&ScreenName, "sn", "Zate", "ScreenName of the user to classify.")
	flag.Parse()
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	log.Println("Botceptor Coming Online ....")
	// if GetTrainingData == false {
	// 	for {
	// 		last200()
	// 		log.Println("Pausing for 60 seconds.")
	// 		time.Sleep(time.Duration(60) * time.Second)
	// 	}
	// } else {
	last200()
	// }
}
