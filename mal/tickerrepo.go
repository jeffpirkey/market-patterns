package mal

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"market-patterns/model"
	"market-patterns/model/report"
	"sort"
)

type TickerRepo struct {
	c          *mongo.Collection
	updateOpt  *options.FindOneAndUpdateOptions
	replaceOpt *options.FindOneAndReplaceOptions
}

var (
	TRUE  = true
	FALSE = false
)

func (repo *TickerRepo) Init() {

	repo.updateOpt = options.FindOneAndUpdate().SetUpsert(TRUE).SetReturnDocument(options.After)
	repo.replaceOpt = options.FindOneAndReplace().SetUpsert(TRUE).SetReturnDocument(options.After)

	idxModel := mongo.IndexModel{}
	idxModel.Keys = bsonx.Doc{{Key: "symbol", Value: bsonx.Int32(1)}}

	name := "idx_symbol"
	idxModel.Options = &options.IndexOptions{Background: &TRUE, Name: &name, Unique: &TRUE}

	tmp, err := repo.c.Indexes().CreateOne(context.TODO(), idxModel)
	if err != nil {
		log.Errorf("problem creating %v due to %v", tmp, err)
	}
}

func (repo *TickerRepo) InsertMany(data []model.Ticker) error {

	dataAry := make([]interface{}, len(data))
	for i, v := range data {
		dataAry[i] = v
	}
	_, err := repo.c.InsertMany(context.TODO(), dataAry)
	if err != nil {
		return errors.Wrap(err, "problem inserting many tickers")
	}
	return nil
}

func (repo *TickerRepo) DeleteAll() error {
	return repo.c.Drop(context.TODO())
}

func (repo *TickerRepo) FindOneAndReplace(ticker *model.Ticker) *model.Ticker {

	filter := bson.D{{"symbol", ticker.Symbol}}

	var update model.Ticker
	err := repo.c.FindOneAndReplace(context.TODO(), filter, ticker, repo.replaceOpt).Decode(&update)
	if err != nil {
		log.Warnf("problem replacing ticker due to %v", err)
	}

	return &update
}

func (repo *TickerRepo) FindAndReplace(ticker *model.Ticker) *model.Ticker {

	filter := bson.D{{"symbol", ticker.Symbol}}

	var update model.Ticker
	err := repo.c.FindOneAndReplace(context.TODO(), filter, ticker, repo.replaceOpt).Decode(&update)
	if err != nil {
		log.Warnf("problem replacing ticker due to %v", err)
	}

	return &update
}

func (repo *TickerRepo) FindOne(symbol string) (*model.Ticker, error) {

	filter := bson.D{{"symbol", symbol}}
	var ticker model.Ticker
	err := repo.c.FindOne(context.TODO(), filter).Decode(&ticker)
	return &ticker, err
}

func (repo *TickerRepo) FindOneAndUpdateCompanyName(symbol, company string) *model.Ticker {
	filter := bson.D{{"symbol", symbol}}
	update := bson.D{{"$set", bson.D{{"company", company}}}}

	var result model.Ticker
	err := repo.c.FindOneAndUpdate(context.TODO(), filter, update, repo.updateOpt).Decode(&result)
	if err != nil {
		log.Warnf("unable to update company of ticker with symbol %v due to %v", symbol, err)
		return nil
	}

	return &result
}

func (repo *TickerRepo) FindSymbols() *[]string {

	var symbols []string

	ary, err := repo.c.Distinct(context.TODO(), "symbol", bson.D{})
	if err != nil {
		log.Warnf("unable to load ticker symbols due to %v", err)
		return &symbols
	}

	if len(ary) > 0 {
		symbols = make([]string, len(ary))
		for i, v := range ary {
			symbols[i] = fmt.Sprint(v)
		}
	}

	return &symbols
}

func (repo *TickerRepo) FindSymbolsAndCompany() *report.TickerSymbolCompanySlice {

	var symbols report.TickerSymbolCompanySlice

	opts := options.Find()
	opts.Projection = bson.D{{"symbol", 1}, {"company", 1}, {"_id", 0}}
	cur, err := repo.c.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Warnf("unable to load ticker symbols and companies due to %v", err)
		return &symbols
	}

	for cur.Next(context.TODO()) {
		var doc *report.TickerSymbolCompany
		err := cur.Decode(&doc)
		if err != nil {
			log.Errorf("unable to unmarshal due to %v", err)
			continue
		}
		symbols = append(symbols, doc)
	}

	sort.Sort(symbols)

	return &symbols
}