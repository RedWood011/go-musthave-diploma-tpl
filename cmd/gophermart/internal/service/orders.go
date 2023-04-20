package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/RedWood011/cmd/gophermart/internal/apperrors"
	"github.com/RedWood011/cmd/gophermart/internal/entity"
)

const (
	register = "REGISTERED"
	process  = "PROCESSING"
	finish   = "PROCESSED"
)

type responseOrder struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual"`
}

func (s *Service) CreateOrder(ctx context.Context, order entity.Order) error {
	err := s.storage.SaveOrder(ctx, order)
	if err != nil && isUniqueViolationError(err) {
		existOrder, ierr := s.storage.GetOrder(ctx, order.Number)
		if err != nil {
			s.logger.Info(fmt.Sprintf("CreateOrder: Save order: %v", ierr))
			return ierr
		}
		if existOrder.UserID != order.UserID {
			return apperrors.ErrOrderOwnedByAnotherUser
		}
	}
	s.logger.Info(fmt.Sprintf("CreateOrder: Save order: %v", err))
	return err
}

func (s *Service) GetOrders(ctx context.Context, userID string) ([]entity.Order, error) {
	return s.storage.GetUserOrders(ctx, userID)

}

func (s *Service) UpdateOrders(ctx context.Context) {
	orders, err := s.storage.GetAllOrders(ctx)
	s.logger.Info(fmt.Sprintf("UpdateOrders: Get all orders: %v", orders))

	if err != nil {
		return
		s.logger.Info(fmt.Sprintf("Error get all orders: %s", err.Error()))
	}
	saveOrders := make([]entity.Order, 0, len(orders))
	usersBalanse := make([]entity.User, 0, len(orders))
loop:
	for _, order := range orders {
		var response *http.Response
		for i := 0; i <= s.cfg.CountRepetitionClient; i++ {
			url := fmt.Sprintf("%s/api/orders/%s", s.cfg.AccrualSystemAddress, order.Number)
			response, err = http.Get(url)
			if err != nil || response.StatusCode != http.StatusOK {
				switch {
				case i == s.cfg.CountRepetitionClient:
					var code int
					if response != nil {
						code = response.StatusCode
					}
					s.logger.Info(fmt.Sprintf("Error request: UserID = %s, order number =%s, statusCode = %d", order.UserID, order.Number, code))
					continue loop

				default:
					continue
				}
			}
		}
		body, err := io.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			s.logger.Info(fmt.Sprintf("Error read body response: UserID = %s, order number =%s", order.UserID, order.Number))
			continue
		}
		var res responseOrder
		err = json.Unmarshal(body, &res)
		if err != nil {
			s.logger.Info(fmt.Sprintf("Error parse json response: UserID = %s, order number =%s", order.UserID, order.Number))
			continue
		}

		if res.Status == register || res.Status == process {
			continue
		}

		order.Accrual = res.Accrual
		order.Status = res.Status
		saveOrders = append(saveOrders, order)
		if res.Status == finish {
			usersBalanse = append(usersBalanse, entity.User{
				ID: order.UserID,
				UserBalance: entity.UserBalance{
					Balance: res.Accrual,
				},
			})
		}
	}

	if len(saveOrders) > 0 {
		s.logger.Info(fmt.Sprintf("UpdateOrders: Save orders: %v", saveOrders))
		err = s.storage.UpdateOrders(ctx, saveOrders)
		if err != nil {
			s.logger.Info("Error save orders")
		}
	}
	if len(usersBalanse) > 0 {
		s.logger.Info(fmt.Sprintf("UpdateOrders: Save orders: %v", usersBalanse))
		err = s.storage.UpdateUsersBalance(ctx, usersBalanse)
		if err != nil {
			s.logger.Info("Error update user balance")
		}
	}
}
