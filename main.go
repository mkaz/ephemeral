package main

/*
	ephemeral - a program to delete tweets

	Originally from Vicky Lai whose version was written to run on AWS Lambda
	See her blog post: https://vickylai.com/verbose/delete-old-tweets-ephemeral/
	   or Github code: https://github.com/hivickylai/ephemeral

	Modified to add twepoch, which is a date to not delete prior to. I like
	keeping around my first tweet. This version will not delete anything prior
	to the TWEPOCH date. Format: YYYY-DD-MM

	Usage:
		* Set environment variables
		* Test run: ./ephemeral --test
*/

import (
	"flag"
	"net/url"
	"os"
	"time"

	"github.com/ChimeraCoder/anaconda"
	log "github.com/sirupsen/logrus"
)

var (
	consumerKey       = getenv("TWITTER_CONSUMER_KEY")
	consumerSecret    = getenv("TWITTER_CONSUMER_SECRET")
	accessToken       = getenv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret = getenv("TWITTER_ACCESS_TOKEN_SECRET")
	maxTweetAge       = getenv("MAX_TWEET_AGE")
	twepochDate       = getenv("TWEPOCH")
)

var testRun bool
var verbose bool

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		log.Fatal("missing required environment variable " + name)
	}
	return v
}

func getTimeline(api *anaconda.TwitterApi) ([]anaconda.Tweet, error) {
	args := url.Values{}
	args.Add("count", "200")        // Twitter only returns most recent 20 tweets by default, so override
	args.Add("include_rts", "true") // When using count argument, RTs are excluded, so include them as recommended
	timeline, err := api.GetUserTimeline(args)
	if err != nil {
		return make([]anaconda.Tweet, 0), err
	}
	return timeline, nil
}

func deleteFromTimeline(api *anaconda.TwitterApi, ageLimit time.Duration) {
	timeline, err := getTimeline(api)

	if err != nil {
		log.Fatal("Could not get timeline")
	}

	for _, t := range timeline {
		createdTime, err := t.CreatedAtTime()
		if err != nil {
			log.Fatal("Couldn't parse time ", err)
		}

		// dont delete the first tweet
		twepoch, err := time.Parse("2006-01-02", twepochDate)
		if err != nil {
			log.Fatal("Error parsing twepochDate")
		}

		if createdTime.Before(twepoch) {
			log.WithFields(log.Fields{
				"tweet":   t.Text,
				"twepoch": twepoch,
			}).Info("Tweet before the twepoch")
			continue
		}

		if time.Since(createdTime) > ageLimit {
			if !testRun {
				_, err := api.DeleteTweet(t.Id, true)
				if err != nil {
					log.Warn("Failed to delete! ", err)
				}
			}
			log.Info("DELETED: Age - ", time.Since(createdTime).Round(1*time.Minute), " - ", t.Text)
		} else {
			log.WithFields(log.Fields{
				"tweet":    t.Text[:45],
				"ageLimit": ageLimit,
			}).Info("Tweet within save window.")
		}
	}
	log.Info("No more tweets to delete.")

}

func main() {

	flag.BoolVar(&testRun, "test", false, "Just test what would happen")
	flag.BoolVar(&verbose, "verbose", false, "Verbose output")
	flag.Parse()

	if testRun {
		verbose = true
	}

	if verbose {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	log.Info("Running ephemeral...")

	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)
	api.SetLogger(anaconda.BasicLogger)

	h, _ := time.ParseDuration(maxTweetAge)

	deleteFromTimeline(api, h)
}
