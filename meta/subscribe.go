package meta

import (
	"errors"
	"log"
	"net/http"

	"github.com/midorigreen/gosns/logging"
	"github.com/midorigreen/gosns/topic"
)

type subscribe struct {
	TopicData *topic.TopicData
}

type subscribeReq struct {
	Channel       string    `json:"channel"`
	ClientID      string    `json:"client_id"`
	Subscriptions []string  `json:"subscriptions"`
	Method        strMethod `json:"method"`
}

type strMethod struct {
	Format     string `json:"format"`
	WebHookURL string `json:"webhook_url"`
}

type subscribeRes struct {
	Channel       string   `json:"channel"`
	Successful    bool     `json:"successful"`
	ClientID      string   `json:"clientId"`
	Subscriptions []string `json:"subscriptions"`
	Error         string   `json:",omitempty"`
}

func (s *subscribe) handler(w http.ResponseWriter, r *http.Request) {
	logging.Logger.Info(r.URL.String())
	req := subscribeReq{}
	err := decodeBody(r, &req)
	if err != nil {
		log.Println(err)
		unsuccessed("parse request to json error", "", []string{}, w, http.StatusBadRequest)
		return
	}

	registered, err := register(req, s.TopicData)
	if err != nil {
		log.Println(err)
		unsuccessed(err.Error(), req.ClientID, []string{}, w, http.StatusBadRequest)
		return
	}

	successed(req.ClientID, registered, w)
}

func register(sReq subscribeReq, td *topic.TopicData) ([]string, error) {
	format := topic.FormatValue(sReq.Method.Format)
	if format == topic.Error {
		return nil, errors.New("subscribed format error: " + format.String())
	}

	var registered = []string{}
	var topics []topic.Topic
	for _, v := range sReq.Subscriptions {
		topic := topic.Topic{
			Channel: v,
			Subscribers: []topic.Subscriber{
				{
					ClientID: sReq.ClientID,
					Method: topic.Method{
						Format:     format,
						WebFookURL: sReq.Method.WebHookURL,
					},
				},
			},
		}
		topics = append(topics, topic)
		registered = append(registered, topic.Channel)
	}

	if len(topics) == 0 {
		return nil, errors.New("not found topic")
	}

	if err := td.Update(topics); err != nil {
		return nil, err
	}

	return registered, nil
}

func unsuccessed(errMes, clientID string, sub []string, w http.ResponseWriter, statusCode int) {
	write(false, errMes, clientID, sub, w, statusCode)
}

func successed(clientID string, sub []string, w http.ResponseWriter) {
	write(true, "", clientID, sub, w, http.StatusOK)
}

func write(success bool, mes, clientID string, sub []string, w http.ResponseWriter, statusCode int) {
	res := subscribeRes{
		Channel:       subscribePattarn,
		Successful:    success,
		ClientID:      clientID,
		Subscriptions: sub,
		Error:         mes,
	}
	writeRes(res, w, statusCode)
}
