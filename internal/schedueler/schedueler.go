package schedueler

import (
	"time"

	"github.com/madflojo/tasks"
	"github.com/thejixer/shop-api/internal/mailer"
	"github.com/thejixer/shop-api/internal/models"
)

type ScheduelerService struct {
	mailer    *mailer.MailerService
	Scheduler *tasks.Scheduler
}

func NewTaskSchedueler(m *mailer.MailerService) *ScheduelerService {

	scheduler := tasks.New()

	return &ScheduelerService{
		Scheduler: scheduler,
		mailer:    m,
	}

}

func (s *ScheduelerService) SchedueleShipmentNotification(order *models.OrderDto) error {

	_, err := s.Scheduler.Add(&tasks.Task{
		Interval: time.Duration(10 * time.Minute),
		RunOnce:  true,
		TaskFunc: func() error {
			err := s.mailer.NotifyShipmentEmail(order)
			if err != nil {
				return err
			}
			return nil
		},
	})
	if err != nil {
		return err
	}
	return nil
}
