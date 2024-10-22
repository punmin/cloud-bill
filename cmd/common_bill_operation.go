package cmd

import "fmt"

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
