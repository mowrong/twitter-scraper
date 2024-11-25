package twitterscraper

import (
	"net/http"
	"net/url"
)

type FollowingResponse struct {
	Data struct {
		User struct {
			Result struct {
				Typename string `json:"__typename"`
				Timeline struct {
					Timeline Timeline `json:"timeline"`
				} `json:"timeline"`
			} `json:"result"`
		} `json:"user"`
	} `json:"data"`
}

type Timeline struct {
	Instructions []Instruction `json:"instructions"`
}

type Instruction struct {
	Type      string  `json:"type"`
	Direction string  `json:"direction,omitempty"`
	Entries   []Entry `json:"entries,omitempty"`
	Entry     Entry   `json:"entry,omitempty"`
}

type Following struct {
	Profiles []Profile
}

func parseFollowing(timeline Timeline) Following {
	var profiles []Profile
	for _, instruction := range timeline.Instructions {
		if instruction.Type == "TimelineAddEntries" {
			for _, entry := range instruction.Entries {
				if entry.Content.ItemContent.UserResults.Result.RestId == "" || entry.Content.ItemContent.UserResults.Result.Typename != "User" {
					continue
				}
				entry.Content.ItemContent.UserResults.Result.Legacy.IDStr = entry.Content.ItemContent.UserResults.Result.RestId
				profiles = append(profiles, parseProfile(entry.Content.ItemContent.UserResults.Result.Legacy))
			}
		}
	}

	return Following{profiles}
}

type Entry struct {
	EntryId   string `json:"entryId"`
	SortIndex string `json:"sortIndex"`
	Content   struct {
		EntryType   string `json:"entryType"`
		Typename    string `json:"__typename"`
		ItemContent struct {
			ItemType    string `json:"itemType"`
			Typename    string `json:"__typename"`
			UserResults struct {
				Result struct {
					Typename                   string `json:"__typename"`
					Id                         string `json:"id"`
					RestId                     string `json:"rest_id"`
					AffiliatesHighlightedLabel struct {
						Label struct {
							Url struct {
								Url     string `json:"url"`
								UrlType string `json:"urlType"`
							} `json:"url"`
							Badge struct {
								Url string `json:"url"`
							} `json:"badge"`
							Description          string `json:"description"`
							UserLabelType        string `json:"userLabelType"`
							UserLabelDisplayType string `json:"userLabelDisplayType"`
						} `json:"label,omitempty"`
					} `json:"affiliates_highlighted_label"`
					HasGraduatedAccess bool       `json:"has_graduated_access"`
					IsBlueVerified     bool       `json:"is_blue_verified"`
					ProfileImageShape  string     `json:"profile_image_shape"`
					Legacy             legacyUser `json:"legacy"`
					Professional       struct {
						RestId           string `json:"rest_id"`
						ProfessionalType string `json:"professional_type"`
						Category         []struct {
							Id       int    `json:"id"`
							Name     string `json:"name"`
							IconName string `json:"icon_name"`
						} `json:"category"`
					} `json:"professional,omitempty"`
				} `json:"result"`
			} `json:"user_results"`
			UserDisplayType string `json:"userDisplayType"`
		} `json:"itemContent,omitempty"`
		ClientEventInfo struct {
			Component string `json:"component"`
			Element   string `json:"element"`
		} `json:"clientEventInfo,omitempty"`
		Value      string `json:"value,omitempty"`
		CursorType string `json:"cursorType,omitempty"`
	} `json:"content"`
}

func (s *Scraper) GetFollowing(userId string) (Following, error) {
	var jsn FollowingResponse
	req, err := http.NewRequest("GET", "https://x.com/i/api/graphql/eWTmcJY3EMh-dxIR7CYTKw/Following", nil)
	if err != nil {
		return Following{}, err
	}

	variables := map[string]interface{}{
		"userId":                 userId,
		"count":                  20,
		"includePromotedContent": false,
	}

	features := map[string]interface{}{
		"communities_web_enable_tweet_community_results_fetch":                    true,
		"articles_preview_enabled":                                                true,
		"rweb_tipjar_consumption_enabled":                                         false,
		"responsive_web_graphql_exclude_directive_enabled":                        true,
		"verified_phone_label_enabled":                                            false,
		"creator_subscriptions_tweet_preview_api_enabled":                         true,
		"responsive_web_graphql_timeline_navigation_enabled":                      true,
		"responsive_web_graphql_skip_user_profile_image_extensions_enabled":       false,
		"c9s_tweet_anatomy_moderator_badge_enabled":                               true,
		"tweetypie_unmention_optimization_enabled":                                true,
		"responsive_web_edit_tweet_api_enabled":                                   true,
		"graphql_is_translatable_rweb_tweet_is_translatable_enabled":              true,
		"view_counts_everywhere_api_enabled":                                      true,
		"longform_notetweets_consumption_enabled":                                 true,
		"responsive_web_twitter_article_tweet_consumption_enabled":                false,
		"tweet_awards_web_tipping_enabled":                                        false,
		"creator_subscriptions_quote_tweet_preview_enabled":                       false,
		"freedom_of_speech_not_reach_fetch_enabled":                               true,
		"standardized_nudges_misinfo":                                             true,
		"tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled": false,
		"rweb_video_timestamps_enabled":                                           false,
		"longform_notetweets_rich_text_read_enabled":                              true,
		"longform_notetweets_inline_media_enabled":                                false,
		"responsive_web_enhance_cards_enabled":                                    false,
	}

	query := url.Values{}
	query.Set("variables", mapToJSONString(variables))
	query.Set("features", mapToJSONString(features))

	req.URL.RawQuery = query.Encode()

	err = s.RequestAPI(req, &jsn)
	if err != nil {
		return Following{}, err
	}

	return parseFollowing(jsn.Data.User.Result.Timeline.Timeline), nil
}
