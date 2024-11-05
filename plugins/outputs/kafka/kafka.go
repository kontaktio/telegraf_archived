//go:generate ../../../tools/readme_config_includer/generator
package kafka

import (
	_ "embed"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/gofrs/uuid"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/common/kafka"
	"github.com/influxdata/telegraf/plugins/common/proxy"
	"github.com/influxdata/telegraf/plugins/outputs"
	"github.com/influxdata/telegraf/plugins/serializers"
)

//go:embed sample.conf
var sampleConfig string

var ValidTopicSuffixMethods = []string{
	"",
	"measurement",
	"tags",
}

var zeroTime = time.Unix(0, 0)

var recordsSent = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "telegraf_output_kafka_records_sent",
	Help: "Number of records sent to kafka",
})

type Kafka struct {
	Brokers         []string    `toml:"brokers"`
	Topic           string      `toml:"topic"`
	TopicTag        string      `toml:"topic_tag"`
	ExcludeTopicTag bool        `toml:"exclude_topic_tag"`
	TopicSuffix     TopicSuffix `toml:"topic_suffix"`
	RoutingTag      string      `toml:"routing_tag"`
	RoutingKey      string      `toml:"routing_key"`

	proxy.Socks5ProxyConfig

	// Legacy TLS config options
	// TLS client certificate
	Certificate string
	// TLS client key
	Key string
	// TLS certificate authority
	CA string

	kafka.WriteConfig

	kafka.Logger

	Log telegraf.Logger `toml:"-"`

	saramaConfig *sarama.Config
	producerFunc func(addrs []string, config *sarama.Config) (sarama.SyncProducer, error)
	producer     sarama.SyncProducer

	serializer serializers.Serializer
}

type TopicSuffix struct {
	Method    string   `toml:"method"`
	Keys      []string `toml:"keys"`
	Separator string   `toml:"separator"`
}

type MurMurPartitioner struct {
}

// murmur2 implements hashing algorithm used by JVM clients for Kafka.
// See the original implementation: https://github.com/apache/kafka/blob/1.0.0/clients/src/main/java/org/apache/kafka/common/utils/Utils.java#L353
func murmur2(data []byte) int32 {
	length := int32(len(data))
	seed := uint32(0x9747b28c)
	m := int32(0x5bd1e995)
	r := uint32(24)

	h := int32(seed ^ uint32(length))
	length4 := length / 4

	for i := int32(0); i < length4; i++ {
		i4 := i * 4
		k := int32(data[i4+0]&0xff) + (int32(data[i4+1]&0xff) << 8) + (int32(data[i4+2]&0xff) << 16) + (int32(data[i4+3]&0xff) << 24)
		k *= m
		k ^= int32(uint32(k) >> r)
		k *= m
		h *= m
		h ^= k
	}

	switch length % 4 {
	case 3:
		h ^= int32(data[(length & ^3)+2]&0xff) << 16
		fallthrough
	case 2:
		h ^= int32(data[(length & ^3)+1]&0xff) << 8
		fallthrough
	case 1:
		h ^= int32(data[length & ^3] & 0xff)
		h *= m
	}

	h ^= int32(uint32(h) >> 13)
	h *= m
	h ^= int32(uint32(h) >> 15)

	return h
}

// toPositive converts i to positive number as per the original implementation in the JVM clients for Kafka.
// See the original implementation: https://github.com/apache/kafka/blob/1.0.0/clients/src/main/java/org/apache/kafka/common/utils/Utils.java#L741
func toPositive(i int32) int32 {
	return i & 0x7fffffff
}

func (m *MurMurPartitioner) Partition(message *sarama.ProducerMessage, numPartitions int32) (int32, error) {
	keyBytes, err := message.Key.Encode()
	if err != nil {
		return 0, err
	}
	hash := murmur2(keyBytes)
	return toPositive(hash) % numPartitions, nil
}

func (m *MurMurPartitioner) RequiresConsistency() bool {
	return true
}

func ValidateTopicSuffixMethod(method string) error {
	for _, validMethod := range ValidTopicSuffixMethods {
		if method == validMethod {
			return nil
		}
	}
	return fmt.Errorf("unknown topic suffix method provided: %s", method)
}

func (*Kafka) SampleConfig() string {
	return sampleConfig
}

func (k *Kafka) GetTopicName(metric telegraf.Metric) (telegraf.Metric, string) {
	topic := k.Topic
	if k.TopicTag != "" {
		if t, ok := metric.GetTag(k.TopicTag); ok {
			topic = t

			// If excluding the topic tag, a copy is required to avoid modifying
			// the metric buffer.
			if k.ExcludeTopicTag {
				metric = metric.Copy()
				metric.Accept()
				metric.RemoveTag(k.TopicTag)
			}
		}
	}

	var topicName string
	switch k.TopicSuffix.Method {
	case "measurement":
		topicName = topic + k.TopicSuffix.Separator + metric.Name()
	case "tags":
		var topicNameComponents []string
		topicNameComponents = append(topicNameComponents, topic)
		for _, tag := range k.TopicSuffix.Keys {
			tagValue := metric.Tags()[tag]
			if tagValue != "" {
				topicNameComponents = append(topicNameComponents, tagValue)
			}
		}
		topicName = strings.Join(topicNameComponents, k.TopicSuffix.Separator)
	default:
		topicName = topic
	}
	return metric, topicName
}

func (k *Kafka) SetSerializer(serializer serializers.Serializer) {
	k.serializer = serializer
}

func (k *Kafka) Init() error {
	k.SetLogger()

	err := ValidateTopicSuffixMethod(k.TopicSuffix.Method)
	if err != nil {
		return err
	}
	config := sarama.NewConfig()
	config.Producer.Flush.Bytes = 20000000
	config.Producer.Flush.Frequency = 100 * time.Millisecond
	config.Producer.Partitioner = func(topic string) sarama.Partitioner {
		return &MurMurPartitioner{}
	}

	if err := k.SetConfig(config, k.Log); err != nil {
		return err
	}

	k.saramaConfig = config

	// Legacy support ssl config
	if k.Certificate != "" {
		k.TLSCert = k.Certificate
		k.TLSCA = k.CA
		k.TLSKey = k.Key
	}

	if k.Socks5ProxyEnabled {
		config.Net.Proxy.Enable = true

		dialer, err := k.Socks5ProxyConfig.GetDialer()
		if err != nil {
			return fmt.Errorf("connecting to proxy server failed: %s", err)
		}
		config.Net.Proxy.Dialer = dialer
	}

	return nil
}

func (k *Kafka) Connect() error {
	producer, err := k.producerFunc(k.Brokers, k.saramaConfig)
	if err != nil {
		return err
	}
	k.producer = producer
	return nil
}

func (k *Kafka) Close() error {
	return k.producer.Close()
}

func (k *Kafka) routingKey(metric telegraf.Metric) (string, error) {
	if k.RoutingTag != "" {
		key, ok := metric.GetTag(k.RoutingTag)
		if ok {
			return key, nil
		}
	}

	if k.RoutingKey == "random" {
		u, err := uuid.NewV4()
		if err != nil {
			return "", err
		}
		return u.String(), nil
	}

	return k.RoutingKey, nil
}

func (k *Kafka) Write(metrics []telegraf.Metric) error {
	msgs := make([]*sarama.ProducerMessage, 0, len(metrics))
	for _, metric := range metrics {
		metric, topic := k.GetTopicName(metric)

		buf, err := k.serializer.Serialize(metric)
		if err != nil {
			k.Log.Debugf("Could not serialize metric: %v", err)
			continue
		}

		m := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.ByteEncoder(buf),
		}

		// Negative timestamps are not allowed by the Kafka protocol.
		if !metric.Time().Before(zeroTime) {
			m.Timestamp = metric.Time()
		}

		key, err := k.routingKey(metric)
		if err != nil {
			return fmt.Errorf("could not generate routing key: %v", err)
		}

		if key != "" {
			m.Key = sarama.StringEncoder(key)
		}
		msgs = append(msgs, m)
	}

	err := k.producer.SendMessages(msgs)
	if err != nil {
		// We could have many errors, return only the first encountered.
		if errs, ok := err.(sarama.ProducerErrors); ok {
			for _, prodErr := range errs {
				if prodErr.Err == sarama.ErrMessageSizeTooLarge {
					k.Log.Error("Message too large, consider increasing `max_message_bytes`; dropping batch")
					return nil
				}
				if prodErr.Err == sarama.ErrInvalidTimestamp {
					k.Log.Error(
						"The timestamp of the message is out of acceptable range, consider increasing broker `message.timestamp.difference.max.ms`; " +
							"dropping batch",
					)
					return nil
				}
				return prodErr //nolint:staticcheck // Return first error encountered
			}
		}
		return err
	}
	recordsSent.Add(float64(len(msgs)))
	return nil
}

func init() {
	outputs.Add("kafka", func() telegraf.Output {
		return &Kafka{
			WriteConfig: kafka.WriteConfig{
				MaxRetry:     3,
				RequiredAcks: -1,
			},
			producerFunc: sarama.NewSyncProducer,
		}
	})
	prometheus.MustRegister(recordsSent)
}
