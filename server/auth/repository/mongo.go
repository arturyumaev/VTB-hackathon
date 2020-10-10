package repository

import (
	"auth/config"
	"auth/domain"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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

func (m *MongoDB) GetConfig(name string) (*domain.OAuthConfig, error) {
	collection := m.db.Collection("config")
	filter := bson.M{"name": name}
	var authConfig domain.OAuthConfig
	if err := collection.FindOne(context.Background(), filter).Decode(&authConfig); err != nil {
		return nil, err
	}
	return &authConfig, nil
}

func (m *MongoDB) GetOAuthClient(name string) (*domain.OAuthClient, error) {
	collection := m.db.Collection("clients")
	filter := bson.M{"name": name}
	var client domain.OAuthClient
	if err := collection.FindOne(context.Background(), filter).Decode(&client); err != nil {
		return nil, err
	}
	return &client, nil
}

func (m *MongoDB) SaveYandexData(info *domain.YandexUser) (domain.UserData, error) {
	collection := m.db.Collection("users")
	user := domain.UserData{
		Login: info.Login,
		Name:  info.FirstName + " " + info.LastName,
		Email: info.DefaultEmail,
		Type:  "yandex",
	}
	exists, err := collection.CountDocuments(context.TODO(), bson.M{"type": "yandex", "login": info.Login})
	if err == nil && exists > 0 {
		if _, err := collection.UpdateOne(context.TODO(), bson.M{"type": "yandex", "login": info.Login},
			bson.M{"$set": user}); err != nil {
			return domain.UserData{}, err
		}
	} else {
		if _, err := collection.InsertOne(context.TODO(), user); err != nil {
			return domain.UserData{}, err
		}
	}
	return user, nil
}

func (m *MongoDB) GetUser(login string, password string) (*domain.UserData, error) {
	collection := m.db.Collection("users")
	filter := bson.M{"login": login, "type": "internal"}
	var data domain.UserData
	if err := collection.FindOne(context.Background(), filter).Decode(&data); err != nil {
		return nil, errors.New("user not exists")
	}
	if data.Password != getHash(password) {
		return nil, errors.New("password incorrect")
	}
	return &data, nil
}

func (m *MongoDB) SaveUser(user domain.UserRegData) error {
	collection := m.db.Collection("users")
	filter := bson.M{"login": user.Login, "type": "internal"}
	if count, _ := collection.CountDocuments(context.Background(), filter); count > 0 {
		return errors.New("user with this login already exists")
	}
	user.Password = getHash(user.Password)
	user.Type = "internal"
	if _, err := collection.InsertOne(context.TODO(), user); err != nil {
		return err
	}
	return nil
}

func (m *MongoDB) SaveUserSession(login string, session string) error {
	collection := m.db.Collection("sessions")
	sess := domain.Session{
		Login:     login,
		SessionId: session,
	}
	if _, err := collection.InsertOne(context.TODO(), sess); err != nil {
		return err
	}
	return nil
}

func (m *MongoDB) CheckAndInvalidateUserSession(login string, session string) (bool, error) {
	collection := m.db.Collection("sessions")
	exists, err := collection.CountDocuments(context.TODO(), bson.M{"login": login, "sessId": session})
	if err != nil || exists == 0 {
		return false, err
	}
	if _, err := collection.DeleteMany(context.TODO(), bson.M{"login": login, "sessId": session}); err != nil {
		return false, errors.New("cannot invalidate: " + err.Error())
	}
	return true, nil
}


func (m *MongoDB) CheckUserSession(login string, session string) (bool, error) {
	collection := m.db.Collection("sessions")
	exists, err := collection.CountDocuments(context.TODO(), bson.M{"login": login, "sessId": session})
	if err != nil || exists == 0 {
		return false, err
	}
	return true, nil
}

func (m *MongoDB) InvalidateAllUserSessions(login string) error {
	collection := m.db.Collection("sessions")
	if _, err := collection.DeleteMany(context.TODO(), bson.M{"login": login}); err != nil {
		return errors.New("cannot invalidate: " + err.Error())
	}
	return nil
}

func getHash(txt string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(txt)))
}
