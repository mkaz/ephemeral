# ephemeral: delete tweets

**ephemeral** is a program that deletes tweets. It is a modified version of
Vicky Lai's program which was written to run on AWS Lambda. See her [blog
post](https://vickylai.com/verbose/delete-old-tweets-ephemeral/) or [Github
code](https://github.com/hivickylai/ephemeral).

You can use ephemeral to automatically delete all tweets from your timeline that
are older than a certain number of hours that you can choose. For instance, you
can ensure that your tweets are deleted after one week (168h), or one day (24h).

You can also set a date to not delete prior to, which I named the twepoch. This
is used in case you want to keep early tweets around.


### Twitter API

This program uses the Go client library [Anaconda](https://github.com/ChimeraCoder/anaconda) to access the Twitter API. You will need to [create a new Twitter application and generate API keys](https://apps.twitter.com/).


### Environment Variables

The program assumes the following environment variables are set:

```
TWITTER_CONSUMER_KEY
TWITTER_CONSUMER_SECRET
TWITTER_ACCESS_TOKEN
TWITTER_ACCESS_TOKEN_SECRET
MAX_TWEET_AGE
TWEPOCH
```

`MAX_TWEET_AGE` expects a value of hours, such as: `MAX_TWEET_AGE=72h`

`TWEPOCH` expects a date value in `YYYY-MM-DD` format


### Usage

1. See releases for binary downloads for Linux and Raspberry Pi.

2. Create `env.sh` file, see `env.sh.sample`

3. Test run

```
$ source env.sh
$ ephemeral --test
```

4. Remove `--test` flag when you're ready to delete away

### Automate

I run this program using cron on a schedule time running on a Raspberry Pi.
Here is the cron entry which shows soucing env.sh prior to running.

```
0 8 * * * . path/to/dir/env.sh; path/to/dir/ephemeral
```

### Build

The program is a standard golang program, so the build process is quite
straight-forward.

```
$ git clone https://github.com/mkaz/ephemeral
$ cd ephemeral
$ go get
$ go build
```

##### Run on a Raspberry Pi

The great thing about golang is its ability to target multiple architectures.
You can do a cross-compile build for Raspberry Pi using:

`env GOOS=linux GOARCH=arm GOARM=5 go build`


### License

This code is modified from Vicky Lai, who modified her code from Adam Drake,
both of which released under the MIT license, so let's keep that license
train going.


