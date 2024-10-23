package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/liqiongfan/leopards"
)

type CloudBillHandler[T any] interface {
	GetBill(billMonth string, account CloudAccount) ([]T, error)
	HasBill(billMonth string, account CloudAccount) bool
	SaveBill(billMonth string, account CloudAccount, resourceSummarySet []T)
}

func HasBill(billMonth string, account CloudAccount, tableName string, mainAccountIDFieldName string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	var result = struct {
		Count int `json:"count"`
	}{}

	err := db.Query().From(tableName).Select(leopards.As(leopards.Count(`id`), `count`)).Where(
		leopards.And(
			leopards.EQ("bill_month", fmt.Sprintf("%s-01 00:00:00", billMonth)),
			leopards.EQ(mainAccountIDFieldName, account.MainAccountID),
		),
	).Scan(ctx, &result)

	if err != nil {
		panic(err)
	}

	return result.Count > 0

}
