package topic

import (
	"encoding/json"
	"net/http"

	"github.com/midorigreen/gosns/logging"
)

const (
	topicPath = "/topic"
)

type topic struct {
	TopicData *TopicData
}

type topicReq struct {
	Channel string `json:"channel"`
	Data    string `json:"data"`
}

// Handler is topic handler
func Handler() {
	t := topic{
		TopicData: CreateTopicData(),
	}
	http.HandleFunc(topicPath, t.handler)
}

func (t *topic) handler(w http.ResponseWriter, r *http.Request) {
	logging.Logger.Info(r.URL.String())
	var tReq topicReq
	decodeBody(r, &tReq)
	topic, err := t.TopicData.Fetch(tReq.Channel)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("not found channel"))
		return
	}
	go sends(topic, tReq)
	w.Write([]byte("ok"))
}

func decodeBody(req *http.Request, out interface{}) error {
	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	return decoder.Decode(out)
}

func findChannel(topics []Topic, channel string) Topic {
	for _, v := range topics {
		if channel == v.Channel {
			return v
		}
	}
	return Topic{}
}
