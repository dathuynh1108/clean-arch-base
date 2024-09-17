package tracer

import (
	"context"
	"runtime"
	"strings"

	"github.com/dathuynh1108/clean-arch-base/pkg/comerr"
	"github.com/dathuynh1108/clean-arch-base/pkg/comjson"

	"go.elastic.co/apm/v2"
)

func StartTransaction(ctx context.Context, transactionType string, metadata any) context.Context {
	funcName := "worker.job.transaction"
	pc, _, _, ok := runtime.Caller(1)
	if ok {
		fullFuncName := runtime.FuncForPC(pc).Name()
		funcIdx := strings.LastIndexByte(fullFuncName, '/') // Use / to get full name path with object (prevent same func name)
		funcName = fullFuncName[funcIdx+1:]
	}
	tx := apm.DefaultTracer().StartTransaction(funcName, transactionType)
	if metadata != nil {
		tx.Context.SetCustom("metadata", metadata)
	}
	return apm.ContextWithTransaction(ctx, tx)
}

func FinishTransaction(ctx context.Context, executeErr *error) {
	var (
		errLike = recover()
		err     = comerr.ToError(errLike)
	)

	tx := apm.TransactionFromContext(ctx)
	if tx != nil {
		jobResult := "success"
		if err != nil {
			jobResult = "failed"
			apm.CaptureError(ctx, err).Send()
		} else if executeErr != nil {
			if *executeErr != nil {
				jobResult = "failed"
				apm.CaptureError(ctx, *executeErr).Send()
			}
		}
		tx.Result = jobResult
		tx.End()
	}
}

func StartSpan(ctx context.Context, spanType string, metadata any) context.Context {
	funcName := "worker.job.span"
	pc, _, _, ok := runtime.Caller(1)
	if ok {
		fullFuncName := runtime.FuncForPC(pc).Name()
		funcIdx := strings.LastIndexByte(fullFuncName, '.')
		funcName = fullFuncName[funcIdx+1:]
	}
	span, ctx := apm.StartSpan(ctx, funcName, spanType)
	if metadata != nil {
		// Span not auto marshal to json
		jsonMetadata, _ := comjson.MarshalString(metadata)
		span.Context.SetLabel("metadata", jsonMetadata)
	}
	return ctx
}

func StartSpanWithName(ctx context.Context, name string, spanType string, metadata any) context.Context {
	span, ctx := apm.StartSpan(ctx, name, spanType)
	if metadata != nil {
		span.Context.SetLabel("metadata", metadata)
	}
	return ctx
}

func FinishSpan(ctx context.Context, executeErr error) {
	var (
		errLike = recover()
		err     = comerr.ToError(errLike)
	)

	tx := apm.SpanFromContext(ctx)
	if tx != nil {
		if err != nil {
			apm.CaptureError(ctx, err).Send()
		} else if executeErr != nil {
			apm.CaptureError(ctx, executeErr).Send()
		}
		tx.End()
	}
}
