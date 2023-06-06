package agent

import (
	"net"
	"time"

	"github.com/devldavydov/promytheus/internal/agent/publisher"
	"github.com/devldavydov/promytheus/internal/common/metric"
	"github.com/sirupsen/logrus"
)

type PublisherFactory func(
	threadID int,
	encrSettings publisher.EncryptionSettings,
) Publisher

func CreatePublisherFactory(
	settings ServiceSettings,
	shutdownTimeout time.Duration,
	ch chan metric.Metrics,
	hostIP net.IP,
	logger *logrus.Logger,
) PublisherFactory {
	fn := func(threadID int, encrSettings publisher.EncryptionSettings) Publisher {
		extraSettings := publisher.PublisherExtraSettings{
			HmacKey:         settings.HmacKey,
			EncrSettings:    encrSettings,
			ShutdownTimeout: &shutdownTimeout,
			HostIP:          hostIP,
		}

		switch settings.UseGRPC {
		case true:
			return publisher.NewGRPCPublisher(
				settings.ServerAddress,
				ch,
				threadID,
				logger,
				extraSettings)
		default:
			return publisher.NewHTTPPublisher(
				settings.ServerAddress,
				ch,
				threadID,
				logger,
				extraSettings)
		}
	}
	return fn
}
