package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/liqiongfan/leopards"
)

type BillOperation[T any] interface {
	GetBill(billMonth string, account CloudAccount) ([]T, error)
	HasBill(billMonth string, account CloudAccount) bool
	SaveBill(billMonth string, account CloudAccount, resourceSummarySet []T)
}

type CommonBillOperation[T any] struct {
	BillOperation BillOperation[T]
}

func (operation *CommonBillOperation[T]) SyncBill(billMonth string, account CloudAccount) {
	if operation.BillOperation.HasBill(billMonth, account) {
		fmt.Printf("%s bill for %s has been synced\n", account.AccountAliasName, billMonth)
		return
	}

	resourceSummarySet, err := operation.BillOperation.GetBill(billMonth, account)
	if err != nil {
		panic(err)
	}

	if len(resourceSummarySet) > 0 {
		operation.BillOperation.SaveBill(billMonth, account, resourceSummarySet)
	}
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
