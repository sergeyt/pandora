package dgraph

import (
	"context"
	"net/http"

	dgo "github.com/dgraph-io/dgo/v2"
	"github.com/sergeyt/pandora/modules/send"
)

func TransactionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		dg, close, err := NewClient(ctx)
		if err != nil {
			send.Error(w, err)
			return
		}
		defer close()

		tx := dg.NewTxn()
		defer Discard(ctx, tx)

		ctx = context.WithValue(ctx, "tx", tx)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func RequestTransaction(r *http.Request) *dgo.Txn {
	return r.Context().Value("tx").(*dgo.Txn)
}
