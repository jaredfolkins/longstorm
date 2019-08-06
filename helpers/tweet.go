package helpers

import (
	"fmt"
	"net/http"
	"strings"
	"unicode"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

const TweetSize = 270 // leave space for the meta data yo! [%d/%d]

const tweetTmpl = `%s

[%d/%d]
`

type TweetWorker struct {
	Client     *twitter.Client
	HttpClient *http.Client
	Token      string
	Config     string
}

func NewTweetWorker(consumerToken, consumerSecret, accessToken, accessSecret string) *TweetWorker {
	config := oauth1.NewConfig(consumerToken, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	tw := &TweetWorker{
		HttpClient: httpClient,
		Client:     twitter.NewClient(httpClient),
	}
	return tw
}

func (tw *TweetWorker) Breakup(story string, honorNewlines bool) []string {

	var words []string
	// break the story apart by space, keeping newlines and carriage returns
	if honorNewlines {
		words = strings.FieldsFunc(story, IsSrslySpace)
	} else { // break the story apart and strip all newlines and carriages returns
		words = strings.Fields(story)
	}
	// totalWords
	totalWords := len(words) - 1 // account for array starting at 0
	// current-tweet-size
	var cts int
	// target-tweet, all-tweets
	var tt, at []string

	for k, word := range words {
		// add white space to the end of the word to build the tweet storm's structure as the space must be included in the current-tweet-size
		word = word + " "
		// calc and add the current word length to the current-tweet-size
		cts = cts + len(word)

		if k == totalWords { // if this is the last word of the entire storm, perform a special "finishing" case

			tt = append(tt, word)                            // append the current word to the new target tweet
			tweet := strings.TrimSpace(strings.Join(tt, "")) // join the slice of words into a tweet seperated by space
			at = append(at, tweet)                           // append the tweet to the all-tweets slice

		} else if cts < TweetSize { // if the tweet still has space to appened words without going over tweetSize, append

			tt = append(tt, word) // append current word to target-tweet

		} else { // else we understand that the tweet is full and we need to append this current tweet to all-tweets and continute to iterate

			cts = len(word)                                  // reset size of current-tweet-size counter
			tweet := strings.TrimSpace(strings.Join(tt, "")) // join the slice of words into a tweet seperated by space
			at = append(at, tweet)                           // append the tweet to the all-tweets slice
			tt = make([]string, 0)                           // clear the target-tweet
			tt = append(tt, word)                            // append the current word to the new target tweet

		}
	}

	return at
}

func (tw *TweetWorker) Storm(story string, honorNewlines bool) []string {
	// get back all-tweets
	at := tw.Breakup(story, honorNewlines)
	// count the tweets and give us the total dawg
	totalTweets := len(at)

	// compiled-tweets
	var ct []string
	for k, v := range at {
		// account for the fact that in the breakup() function the array starts at zero, but humans do not. // example:[1/32]
		k = k + 1
		// use the template to compile the tweet storm
		ct = append(ct, fmt.Sprintf(tweetTmpl, v, k, totalTweets))
	}
	return ct
}

func (tw *TweetWorker) FirstTweet(longstorm []string) (*twitter.Tweet, error) {
	var t *twitter.Tweet
	var err error
	for k, v := range longstorm {
		if k == 0 {
			t, _, err = tw.Client.Statuses.Update(v, nil)
			if err != nil {
				return nil, err
			}
		}
	}
	return t, nil

}

func (tw *TweetWorker) AppendTweet(text string, params *twitter.StatusUpdateParams) (*twitter.Tweet, error) {
	t, _, err := tw.Client.Statuses.Update(text, params)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// ripped off from Golang source and removed the \n \t \r characters
func IsSrslySpace(r rune) bool {
	if uint32(r) <= unicode.MaxLatin1 {
		switch r {
		//case '\t', '\n', '\v', '\f', '\r', ' ', 0x85, 0xA0:
		case '\v', '\f', ' ', 0x85, 0xA0:
			return true
		}
		return false
	}
	return isExcludingLatin(unicode.White_Space, r)
}

// linearMax is the maximum size table for linear search for non-Latin1 rune.
// Derived by running 'go test -calibrate'.
const linearMax = 18

func isExcludingLatin(rangeTab *unicode.RangeTable, r rune) bool {
	r16 := rangeTab.R16
	if off := rangeTab.LatinOffset; len(r16) > off && r <= rune(r16[len(r16)-1].Hi) {
		return is16(r16[off:], uint16(r))
	}
	r32 := rangeTab.R32
	if len(r32) > 0 && r >= rune(r32[0].Lo) {
		return is32(r32, uint32(r))
	}
	return false
}

// is16 reports whether r is in the sorted slice of 16-bit ranges.
func is16(ranges []unicode.Range16, r uint16) bool {
	if len(ranges) <= linearMax || r <= unicode.MaxLatin1 {
		for i := range ranges {
			range_ := &ranges[i]
			if r < range_.Lo {
				return false
			}
			if r <= range_.Hi {
				return range_.Stride == 1 || (r-range_.Lo)%range_.Stride == 0
			}
		}
		return false
	}

	// binary search over ranges
	lo := 0
	hi := len(ranges)
	for lo < hi {
		m := lo + (hi-lo)/2
		range_ := &ranges[m]
		if range_.Lo <= r && r <= range_.Hi {
			return range_.Stride == 1 || (r-range_.Lo)%range_.Stride == 0
		}
		if r < range_.Lo {
			hi = m
		} else {
			lo = m + 1
		}
	}
	return false
}

// is32 reports whether r is in the sorted slice of 32-bit ranges.
func is32(ranges []unicode.Range32, r uint32) bool {
	if len(ranges) <= linearMax {
		for i := range ranges {
			range_ := &ranges[i]
			if r < range_.Lo {
				return false
			}
			if r <= range_.Hi {
				return range_.Stride == 1 || (r-range_.Lo)%range_.Stride == 0
			}
		}
		return false
	}

	// binary search over ranges
	lo := 0
	hi := len(ranges)
	for lo < hi {
		m := lo + (hi-lo)/2
		range_ := ranges[m]
		if range_.Lo <= r && r <= range_.Hi {
			return range_.Stride == 1 || (r-range_.Lo)%range_.Stride == 0
		}
		if r < range_.Lo {
			hi = m
		} else {
			lo = m + 1
		}
	}
	return false
}

/*
	tw := models.NewTweetWorker()
	tv, err := tw.TweetTitleAndMeta(st.Title, st.PubKey)
	c.Logger().Error(err)
	if err != nil {
		c.Logger().Error(err)
		return c.Render(500, r.JSON("{'response': '500'}"))
	}

	msgs := tw.Breakup(st.Text)
	param := &twitter.StatusUpdateParams{
		InReplyToStatusID: tv.ID,
	}

	for _, m := range msgs {
		t, err := tw.AppendTweet(m, param)
		if err != nil {
			c.Logger().Error(err)
			return c.Render(500, r.JSON("{'response': '500'}"))
		}
		param.InReplyToStatusID = t.ID
	}

*/
