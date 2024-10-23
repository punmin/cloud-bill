package cmd

import (
	"fmt"
)

type CloudBillExecutor[T any] struct {
	handler CloudBillHandler[T]
}

func (executor *CloudBillExecutor[T]) SyncBill(billMonth string, account CloudAccount) {
	if executor.handler.HasBill(billMonth, account) {
		fmt.Printf("%s bill for %s has been synced\n", account.AccountAliasName, billMonth)
		return
	}

	resourceSummarySet, err := executor.handler.GetBill(billMonth, account)
	if err != nil {
		panic(err)
	}

	if len(resourceSummarySet) > 0 {
		executor.handler.SaveBill(billMonth, account, resourceSummarySet)
	}
}
