// Package main предоставляет multichecker для статического анализа Go-кода.
// Включает стандартные анализаторы, SA-анализаторы из staticcheck,
// анализатор ST1000, публичные анализаторы (errcheck, ineffassign)
// и собственный анализатор на запрет os.Exit в main.
//
// Запуск: staticlint [flags] [packages]
// Пример: staticlint ./...
//
// Анализаторы:
// - Стандартные анализаторы (asmdecl, assign и др.) — проверка стандартных проблем.
// - SA-анализаторы — проверка сложных семантических ошибок.
// - ST1000 — проверка формата комментариев к пакетам.
// - errcheck — проверка необработанных ошибок.
// - ineffassign — поиск неэффективных присваиваний.
// - osexit — запрет прямого вызова os.Exit в main.main.
package main

import (
	"github.com/VladSnap/shortener/cmd/staticlint/osexit"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/directive"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/slog"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

// Analyzer проверяет наличие прямых вызовов os.Exit в функции main пакета main.

func main() {
	// Список стандартных анализаторов
	stdAnalyzers := []*analysis.Analyzer{
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		atomicalign.Analyzer,
		bools.Analyzer,
		buildssa.Analyzer,
		buildtag.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		ctrlflow.Analyzer,
		deepequalerrors.Analyzer,
		directive.Analyzer,
		errorsas.Analyzer,
		findcall.Analyzer,
		framepointer.Analyzer,
		httpresponse.Analyzer,
		ifaceassert.Analyzer,
		inspect.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		pkgfact.Analyzer,
		printf.Analyzer,
		reflectvaluecompare.Analyzer,
		shift.Analyzer,
		sigchanyzer.Analyzer,
		slog.Analyzer,
		sortslice.Analyzer,
		stdmethods.Analyzer,
		stringintconv.Analyzer,
		structtag.Analyzer,
		testinggoroutine.Analyzer,
		tests.Analyzer,
		timeformat.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
		unusedwrite.Analyzer,
		usesgenerics.Analyzer,
	}

	// Анализаторы SA из staticcheck
	var saAnalyzers []*analysis.Analyzer
	for _, a := range staticcheck.Analyzers {
		if a.Analyzer.Name[:2] == "SA" {
			saAnalyzers = append(saAnalyzers, a.Analyzer)
		}
	}

	// Анализатор ST1000 из staticcheck
	var stAnalyzers []*analysis.Analyzer
	for _, a := range staticcheck.Analyzers {
		if a.Analyzer.Name == "ST1000" {
			stAnalyzers = append(stAnalyzers, a.Analyzer)
		}
	}

	// Собственный анализатор
	customAnalyzers := []*analysis.Analyzer{
		osexit.Analyzer,
	}

	allAnalyzers := []*analysis.Analyzer{}
	// Объединение всех анализаторов
	allAnalyzers = append(allAnalyzers, stdAnalyzers...)
	allAnalyzers = append(allAnalyzers, saAnalyzers...)
	allAnalyzers = append(allAnalyzers, stAnalyzers...)
	// Добавляем stylecheck и simple
	for _, v := range stylecheck.Analyzers {
		allAnalyzers = append(allAnalyzers, v.Analyzer)
	}
	for _, v := range simple.Analyzers {
		allAnalyzers = append(allAnalyzers, v.Analyzer)
	}
	allAnalyzers = append(allAnalyzers, customAnalyzers...)

	multichecker.Main(allAnalyzers...)
}
