package topic

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"time"

	"errors"

	"github.com/emluque/dscache"
)

var file = "./subscribed.json"
var ds, dsErr = dscache.New(2 * dscache.MB)

func init() {
	if dsErr != nil {
		log.Fatalf("err: %s", dsErr)
	}
	topics, err := LoadFile(file)
	if err != nil {
		log.Fatalf("err: %s", err)
	}
	for _, v := range topics {
		subByte, err := json.Marshal(v.Subscribers)
		if err != nil {
			log.Fatalln(err)
		}
		if err = ds.Set(v.Channel, string(subByte), 24*time.Hour); err != nil {
			log.Fatalln(err)
		}
	}
}

// TopicDataService is interface dealing with topic data
type TopicDataService interface {
	Add(topic Topic) error
	Update(topic []Topic) error
	Fetch(channel string) (Topic, error)
}

// TopicData is topic file path and cache
type TopicData struct {
	Path string
	Ds   *dscache.Dscache
}

// CreateTopicData is generating topic data
// TopicData struct should be only created by this func
func CreateTopicData() *TopicData {
	return &TopicData{
		Path: file,
		Ds:   ds,
	}
}

// Add is add topic data to file and cache
func (d *TopicData) Add(topic Topic) error {
	topics, err := LoadFile(d.Path)
	if err != nil {
		return err
	}
	topics = append(topics, topic)

	// add file
	if err = writeFile(topics, d.Path); err != nil {
		return err
	}

	return addCache(topic, d.Ds)
}

// Update is update topic to file and cache
func (d *TopicData) Update(topics []Topic) error {
	nowTopics, err := LoadFile(d.Path)
	if err != nil {
		return err
	}
	var upTopics []Topic
	for _, topic := range topics {
		for i, nowTopic := range nowTopics {
			if topic.Channel == nowTopic.Channel {
				// TODO: duplication check
				nowTopics[i].Subscribers = append(nowTopic.Subscribers, topic.Subscribers...)
				upTopics = append(upTopics, topic)
				break
			}
		}
	}
	if len(upTopics) == 0 {
		return errors.New("not existing topic which try to update")
	}

	// update file
	if err = writeFile(nowTopics, d.Path); err != nil {
		return err
	}

	// update cache
	for _, v := range upTopics {
		// cache clear
		ok := d.Ds.Purge(v.Channel)
		if ok != true {
			return errors.New("failed update cache")
		}
		if err = addCache(v, d.Ds); err != nil {
			return err
		}
	}
	return nil
}

// Fetch is fetching topic from cache (and file)
func (d *TopicData) Fetch(channel string) (Topic, error) {
	str, ok := d.Ds.Get(channel)
	if ok != true {
		log.Println("Not Found Cache")
		// fail safe for cache
		topics, err := LoadFile(d.Path)
		if err != nil {
			return Topic{}, errors.New("not found channel from cache and " + err.Error())
		}
		for _, v := range topics {
			if channel == v.Channel {
				addCache(v, d.Ds)
				return v, nil
			}
		}
		return Topic{}, errors.New("not found channel from cache and file")
	}
	var subscribers []Subscriber
	err := json.Unmarshal([]byte(str), &subscribers)
	if err != nil {
		return Topic{}, err
	}
	return Topic{
		Channel:     channel,
		Subscribers: subscribers,
	}, nil
}

// LoadFile is loading topic file and reformat topic struct
func LoadFile(path string) ([]Topic, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	v, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if len(v) == 0 {
		return nil, nil
	}

	s := []Topic{}
	err = json.Unmarshal(v, &s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func writeFile(topics []Topic, path string) error {
	byte, err := json.MarshalIndent(topics, "", "\t")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(path, byte, 0666); err != nil {
		return err
	}
	return nil
}

func addCache(topic Topic, ds *dscache.Dscache) error {
	subByte, err := json.Marshal(topic.Subscribers)
	if err != nil {
		return err
	}
	return ds.Set(topic.Channel, string(subByte), 24*time.Hour)
}
