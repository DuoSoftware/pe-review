package Twitter_Tweet

import "processengine/context"
import "github.com/ChimeraCoder/anaconda"
import "encoding/base64"
import "net/url"
import "fmt"
import "processengine/logger"
import "strconv"

func Invoke(FlowData map[string]interface{}) (flowResult map[string]interface{}, activityResult *context.ActivityContext) {

	//creating new instance of context.ActivityContext
	var activityContext = new(context.ActivityContext)

	//creating new instance of context.ActivityError
	var activityError context.ActivityError

	//setting activityError proprty values
	activityError.Encrypt = false
	activityError.ErrorString = "exception"
	activityError.Forward = false
	activityError.SeverityLevel = context.Info

	// getting the details from the input arguments
	consumerKey := FlowData["ConsumerKey"].(string)
	consumerSecret := FlowData["ConsumerSecret"].(string)
	accessToken := FlowData["AccessToken"].(string)
	accessTokenSecret := FlowData["AccessTokenSecret"].(string)
	tweetText := FlowData["TweetText"].(string)

	var tweetImage []byte

	var err error
	if FlowData["TweetImage"] != nil {
		tweetImage = FlowData["TweetImage"].([]byte)
		err = SendTweetWithMedia(consumerKey, consumerSecret, accessToken, accessTokenSecret, tweetText, tweetImage)
	} else {
		tweetImage = nil
		err = SendTweet(consumerKey, consumerSecret, accessToken, accessTokenSecret, tweetText)
	}

	if err != nil {
		//setting activityContext property values
		activityContext.ActivityStatus = false
		activityContext.Message = "Tweet Sending Failed : " + err.Error()
		activityContext.ErrorState = activityError
		FlowData["OutData"] = "Tweet Sending Failed : " + err.Error()
		FlowData["custMsg"] = "Tweet Sending Failed : " + err.Error()
		logger.Log_ACT(activityContext.Message, logger.Debug, FlowData["InSessionID"].(string))
	} else {
		//setting activityContext property values
		activityContext.ActivityStatus = true
		activityContext.Message = "Tweet Posted Successfully!"
		FlowData["OutData"] = "Tweet Posted Successfully!"
		FlowData["custMsg"] = "Tweet Posted Successfully!"
		logger.Log_ACT(activityContext.Message, logger.Debug, FlowData["InSessionID"].(string))
	}

	return FlowData, activityContext
}

func SendTweet(consumerKey, consumerSecret, accessToken, accessTokenSecret, tweet string) (err error) {
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)

	_, err = api.PostTweet(tweet, nil)
	return
}

func SendTweetWithMedia(consumerKey, consumerSecret, accessToken, accessTokenSecret, tweet string, data []byte) (err error) {

	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)

	mediaResponse, err := api.UploadMedia(base64.StdEncoding.EncodeToString(data))
	if err != nil {
		return
	}

	v := url.Values{}
	v.Set("media_ids", strconv.FormatInt(mediaResponse.MediaID, 10))

	_, err = api.PostTweet(tweet, v)
	return
}
