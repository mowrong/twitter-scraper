package twitterscraper

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestMyGetTrends(t *testing.T) {
	scraper := New()
	err := scraper.SetProxy("http://127.0.0.1:7890")
	if err != nil {
		panic(err)
	}
	err = scraper.LoginOpenAccount()
	if err != nil {
		panic(err)
	}
	trends, err := scraper.GetTrends()
	if err != nil {
		panic(err)
	}
	fmt.Println(trends)
}

func TestGetSingleTweet(t *testing.T) {
	scraper := New()
	err := scraper.SetProxy("http://127.0.0.1:7890")
	if err != nil {
		panic(err)
	}

	//if f, err := os.Open("cookies_x.json"); os.IsNotExist(err) {
	//	panic("cookies not exist")
	//} else {
	//	defer f.Close()
	//	// deserialize from JSON
	//	var cookies []*http.Cookie
	//	err = json.NewDecoder(f).Decode(&cookies)
	//	if err != nil {
	//		panic(err)
	//	}
	//	// load cookies
	//	scraper.SetCookies(cookies)
	//	// check login status
	//	if !scraper.FakeIsLoggedIn() {
	//		panic(err)
	//	}
	//}

	err = scraper.LoginOpenAccount()
	if err != nil {
		panic(err)
	}
	//var collection = []string{
	//	"1857789201703645287",
	//}

	//for _, id := range collection {
	//	tweet, err := scraper.GetTweet(id)
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Println(tweet.Text)
	//
	//	profile, err := scraper.GetProfile(tweet.Username)
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Printf("%+v", profile)
	//	time.Sleep(time.Second * 2)
	//}

	profile, err := scraper.GetProfile("BeosinAlert")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", profile)
}

func TestSearch(t *testing.T) {
	scraper := New()
	err := scraper.SetProxy("http://127.0.0.1:7890")
	if err != nil {
		panic(err)
	}

	if f, err := os.Open("cookies_x.json"); os.IsNotExist(err) {
		panic("cookies not exist")
	} else {
		defer f.Close()
		// deserialize from JSON
		var cookies []*http.Cookie
		err = json.NewDecoder(f).Decode(&cookies)
		if err != nil {
			panic(err)
		}
		// load cookies
		scraper.SetCookies(cookies)
		// check login status
		if !scraper.FakeIsLoggedIn() {
			panic(err)
		}
	}

	scraper.SetSearchMode(SearchLatest)
	scraper.WithDelay(3)
	scraper.WithReplies(true)

	for tweet := range scraper.SearchTweets(context.Background(),
		"from:realScamSniffer", 10) {
		if tweet.Error != nil {
			panic(tweet.Error)
		}
		fmt.Println(tweet.Name, tweet.ID, time.Unix(tweet.Timestamp, 0).In(time.FixedZone("UTF+8", 8*60*60)).Format("2006-01-02 15:04:05"), tweet.PermanentURL, tweet.Text)
	}
}

func TestFollowing(t *testing.T) {
	scraper := New()
	err := scraper.SetProxy("http://127.0.0.1:7890")
	if err != nil {
		panic(err)
	}

	if f, err := os.Open("cookies_x.json"); os.IsNotExist(err) {
		panic("cookies not exist")
	} else {
		defer f.Close()
		// deserialize from JSON
		var cookies []*http.Cookie
		err = json.NewDecoder(f).Decode(&cookies)
		if err != nil {
			panic(err)
		}
		// load cookies
		scraper.SetCookies(cookies)
		// check login status
		if !scraper.FakeIsLoggedIn() {
			panic(err)
		}
	}

	following, err := scraper.GetFollowing("1497042007432515585")
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("%+v", following))
}
