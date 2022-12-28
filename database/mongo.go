package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrDuplicateMongoDatabase error = errors.New("duplicate mongo database")
var ErrMongoDBInfoNotExist error = errors.New("database info not exist")
var ErrMongoConnNotExist error = errors.New("mongo connection not exist")

func NewMongoConnection(appName string, _mongoConnSettings ...MongoConnSetting) (conn *MongoConnection, err error) {
	conn = &MongoConnection{
		tempDBInfo:      make(map[string]MongoConnSetting),
		tempMongoClient: make(map[string]*mongo.Client),
	}

	for _, _connSetting := range _mongoConnSettings {
		err = conn.addMongoDB(appName, _connSetting)
		if err != nil {
			return conn, err
		}
	}

	return conn, nil
}

type MongoConnSetting struct {
	Address                string `default:"127.0.0.1"`
	Port                   int    `default:"27017"`
	UserName               string `default:"-"`
	Password               string `default:"-"`
	DatabaseName           string `default:"MongoDatabaseName"`
	MaxPoolSize            uint64 `default:"100"`
	Timeout                int    `default:"1"`
	ConnectTimeout         int    `default:"30"`
	SocketTimeout          int    `default:"1"`
	ServerSelectionTimeout int    `default:"30"`
	Comment                string `default:"Comment for db"`
}

type MongoConnection struct {
	tempDBInfo      map[string]MongoConnSetting
	tempMongoClient map[string]*mongo.Client
}

func (conn *MongoConnection) addMongoDB(appName string, _config MongoConnSetting) (err error) {
	var isExist bool
	// 檢查是否有相同資料庫名稱
	if _, isExist = conn.tempDBInfo[_config.DatabaseName]; isExist {
		return ErrDuplicateMongoDatabase
	}

	conn.tempDBInfo[_config.DatabaseName] = _config

	// 檢查是否已有對應位址的連線，避免進行多餘連線
	var _client *mongo.Client
	var _addr string = fmt.Sprintf("%v:%v", _config.Address, _config.Port)
	if _, isExist = conn.tempMongoClient[_addr]; isExist {
		return nil
	}

	_opt := options.Client()
	// 設置應用程式名稱
	_opt.SetAppName(appName)
	// 設定位址和連接埠
	_opt.ApplyURI(fmt.Sprintf("mongodb://%v:%v", _config.Address, _config.Port))

	// 檢查是否需要設定憑證
	if _config.UserName != "empty" {
		_cred := options.Credential{
			Username:    _config.UserName,
			PasswordSet: false,
		}
		if _config.Password != "empty" {
			_cred.PasswordSet = true
			_cred.Password = _config.Password
		}

		_opt.SetAuth(_cred)
	}

	// 設置連線超時
	_opt.SetConnectTimeout(time.Duration(_config.ConnectTimeout) * time.Second)
	// 設置連線池上限
	_opt.SetMaxPoolSize(_config.MaxPoolSize)
	// 設置指令超時
	_opt.SetTimeout(time.Duration(_config.Timeout) * time.Second)
	// 設置Socket連接超時
	_opt.SetSocketTimeout(time.Duration(_config.SocketTimeout) * time.Second)
	// 設置Mongo Server切換超時
	_opt.SetServerSelectionTimeout(time.Duration(_config.ServerSelectionTimeout) * time.Second)

	_client, err = mongo.Connect(context.TODO(), _opt)
	if err != nil {
		return err
	}

	_ctx, _cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer _cancel()
	err = _client.Ping(_ctx, nil)
	if err != nil {
		return err
	}

	conn.tempMongoClient[_addr] = _client
	return nil
}

func (conn *MongoConnection) GetDB(dbName string) (db *mongo.Database, err error) {
	var isExist bool
	var _config MongoConnSetting
	if _config, isExist = conn.tempDBInfo[dbName]; !isExist {
		return nil, ErrMongoDBInfoNotExist
	}

	var _client *mongo.Client
	if _client, isExist = conn.tempMongoClient[fmt.Sprintf("%v:%v", _config.Address, _config.Port)]; !isExist {
		return nil, ErrMongoConnNotExist
	}

	return _client.Database(dbName), nil
}

func (conn *MongoConnection) DisconnectAll() (err error) {
	for addr := range conn.tempMongoClient {
		err = conn.tempMongoClient[addr].Disconnect(context.Background())
		if err != nil {
			return fmt.Errorf("mongo '%v' disconnect fail. msg => %v", addr, err.Error())
		}
	}

	conn.tempDBInfo = nil
	conn.tempMongoClient = nil

	return nil
}
