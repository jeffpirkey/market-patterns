package main

import (
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"market-patterns/report"
	"net/http"
	_ "net/http/pprof"
)

func start() {

	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("./ui/build", true)))

	apiLatest := router.Group("/api/latest")
	apiLatest.GET("/predict", handlePredict)
	apiLatest.GET("/ticker-names", handleTickerNames)

	log.Info("market-pattern server listening...")

	log.Fatal(router.Run(":7666"))
}

func handlePredict(ctx *gin.Context) {

	tsym := "ibm"
	prediction, err := predict(tsym)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
	}
	ctx.JSON(http.StatusOK, prediction)
}

func handleTickerNames(ctx *gin.Context) {
	tickerNames := report.TickerNames{Names: Tickers.FindNames()}
	ctx.JSON(http.StatusOK, tickerNames)
}

func startProfile() {
	log.Info("Starting profile server...")
	log.Fatal(http.ListenAndServe("localhost:6060", nil))
}
