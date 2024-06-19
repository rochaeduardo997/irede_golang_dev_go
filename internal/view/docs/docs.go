package view_docs

import (
	http_adapter "github.com/rochaeduardo997/irede_golang_dev/pkg/http"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type DocsView struct{ HTTPAdapter http_adapter.IHTTP }

func NewDocsView(dv *DocsView) (result *DocsView, err error) {
	result = dv
	result.HTTPAdapter.Route().PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	return
}
