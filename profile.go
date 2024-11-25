package twitterscraper

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Global cache for user IDs
var cacheIDs sync.Map

// Profile of twitter user.
type Profile struct {
	Avatar         string
	Banner         string
	Biography      string
	Birthday       string
	FollowersCount int
	FollowingCount int
	FriendsCount   int
	IsPrivate      bool
	IsVerified     bool
	Joined         *time.Time
	LikesCount     int
	ListedCount    int
	Location       string
	Name           string
	PinnedTweetIDs []string
	TweetsCount    int
	URL            string
	UserID         string
	Username       string
	Website        string
}

type user struct {
	Data struct {
		User struct {
			Result struct {
				Typename string     `json:"__typename"`
				RestID   string     `json:"rest_id"`
				Legacy   legacyUser `json:"legacy"`
			} `json:"result"`
		} `json:"user"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// GetProfile return parsed user profile.
func (s *Scraper) GetProfile(username string) (Profile, error) {
	var jsn user
	req, err := http.NewRequest("GET", "https://x.com/i/api/graphql/laYnJPCAcVo0o6pzcnlVxQ/UserByScreenName", nil)
	if err != nil {
		return Profile{}, err
	}

	variables := map[string]interface{}{
		"screen_name":          username,
		"withHighlightedLabel": true,
	}

	features := map[string]interface{}{"hidden_profile_subscriptions_enabled": true, "rweb_tipjar_consumption_enabled": true, "responsive_web_graphql_exclude_directive_enabled": true, "verified_phone_label_enabled": false, "subscriptions_verification_info_is_identity_verified_enabled": true, "subscriptions_verification_info_verified_since_enabled": true, "highlights_tweets_tab_ui_enabled": true, "responsive_web_twitter_article_notes_tab_enabled": true, "subscriptions_feature_can_gift_premium": true, "creator_subscriptions_tweet_preview_api_enabled": true, "responsive_web_graphql_skip_user_profile_image_extensions_enabled": false, "responsive_web_graphql_timeline_navigation_enabled": true}
	fieldToggles := map[string]interface{}{"withAuxiliaryUserLabels": false}

	query := url.Values{}
	query.Set("variables", mapToJSONString(variables))
	query.Set("features", mapToJSONString(features))
	query.Set("fieldToggles", mapToJSONString(fieldToggles))
	req.URL.RawQuery = query.Encode()

	err = s.RequestAPI(req, &jsn)
	if err != nil {
		return Profile{}, err
	}

	//if len(jsn.Errors) > 0 {
	//	return Profile{}, fmt.Errorf("%s", jsn.Errors[0].Message)
	//}

	if jsn.Data.User.Result.RestID == "" {
		return Profile{}, fmt.Errorf("rest_id not found")
	}
	jsn.Data.User.Result.Legacy.IDStr = jsn.Data.User.Result.RestID

	if jsn.Data.User.Result.Legacy.ScreenName == "" {
		return Profile{}, fmt.Errorf("either @%s does not exist or is private", username)
	}

	return parseProfile(jsn.Data.User.Result.Legacy), nil
}

// GetUserIDByScreenName from API
func (s *Scraper) GetUserIDByScreenName(screenName string) (string, error) {
	id, ok := cacheIDs.Load(screenName)
	if ok {
		return id.(string), nil
	}

	profile, err := s.GetProfile(screenName)
	if err != nil {
		return "", err
	}

	cacheIDs.Store(screenName, profile.UserID)

	return profile.UserID, nil
}
