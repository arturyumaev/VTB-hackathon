package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"payment/config"
	"payment/domain"
	"time"
)

type MongoDB struct {
	clientDB *mongo.Client
	db       *mongo.Database
	logger   *zerolog.Logger
}

func NewMongoDB(cnf *config.ConfigMongo, logger *zerolog.Logger) (*MongoDB, error) {
	clientOpts := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s/%s", cnf.Host, cnf.Port, cnf.Database))
	mongoClient, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		return nil, errors.New("cant connect to mongo")
	}
	ctx, _ := context.WithTimeout(context.TODO(), 3*time.Second)
	if err := mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		return nil, errors.New("mongo ping error")
	}
	db := mongoClient.Database(cnf.Database)

	return &MongoDB{
		clientDB: mongoClient,
		db:       db,
		logger:   logger,
	}, nil
}

func (m *MongoDB) GetBalance(login string) (*domain.UserBalance, error) {
	collection := m.db.Collection("vallets")
	filter := bson.M{"login": login}
	exists, _ := collection.CountDocuments(context.TODO(), filter)
	if exists == 0 {
		_, _ = collection.InsertOne(context.TODO(), bson.M{"login": login, "balance": 0})
	}
	var balance domain.UserBalance
	if err := collection.FindOne(context.Background(), filter).Decode(&balance); err != nil {
		return nil, err
	}
	return &balance, nil
}

func (m *MongoDB) Pay(login string, addresse string, amount int) error {
	if amount < 0 {
		return errors.New("amount less than 0 not allowed")
	}
	collection := m.db.Collection("vallets")
	filter := bson.M{"login": login}
	var userBalance domain.UserBalance
	if err := collection.FindOne(context.Background(), filter).Decode(&userBalance); err != nil {
		return errors.New("login not exists")
	}
	if userBalance.Balance < amount {
		return errors.New("no such money")
	}
	filter = bson.M{"login": addresse}
	addresseCount, err := collection.CountDocuments(context.Background(), filter)
	if addresseCount != 1 || err != nil {
		return errors.New("addressee not exists")
	}

	_, err = collection.UpdateOne(context.TODO(), bson.M{"login": login}, bson.M{"$inc": bson.M{"balance": amount*(-1)}})
	if err != nil {
		return err
	}
	_, err = collection.UpdateOne(context.TODO(), bson.M{"login": addresse}, bson.M{"$inc": bson.M{"balance": amount}})
	if err != nil {
		return err
	}
	collection = m.db.Collection("transactions")
	transaction := domain.Transaction{
		Login:    login,
		Addressee: addresse,
		Anount:    amount,
	}
	_, err = collection.InsertOne(context.TODO(), transaction)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoDB) AddMoney(login string, amount int) error {
	if amount < 0 {
		return errors.New("amount less than 0 not allowed")
	}
	collection := m.db.Collection("vallets")
	filter := bson.M{"login": login}
	var userBalance domain.UserBalance
	if err := collection.FindOne(context.Background(), filter).Decode(&userBalance); err != nil {
		return errors.New("login not exists")
	}
	_, err := collection.UpdateOne(context.TODO(), bson.M{"login": login}, bson.M{"$inc": bson.M{"balance": amount}})
	if err != nil {
		return err
	}
	return nil
}
