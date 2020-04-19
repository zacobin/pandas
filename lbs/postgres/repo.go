//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use p file except in compliance with the License. You may obtain
//  a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//  WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//  License for the specific language governing permissions and limitations
//  under the License.
package proxy

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	sslSuffix = "?ssl=true"
)

const (
	dbnameRepository      = "lbs"
	collectionCollections = "collections"
	collectionGeofences   = "geofences"
	collectionEntity      = "entity"
)



type UserID struct {
	UserID string `bson:"user_id"`
}

type Repository interface {
	// Helper
	AddCollection(userID string, collectionID string) error
	RemoveCollection(userID string, collectionID string) error
	GetAllCollections() ([]*Collection, error)
	UpdateCollection(userID string, p *Collection) error

	// Geofences
	AddGeofence(userID string, collectionID string, fenceName string, fenceID string) error
	RemoveGeofence(userID string, collectionID string, fenceID string) error
	IsGeofenceExistWithName(userID string, collectionID string, fenceName string) bool
	IsGeofenceExistWithId(userID string, collectionID string, fenceID string) bool
	GetFences(userID, collectionID string) ([]*geofenceRecord, error)
	GetFenceUserID(fenceID string) (string, error)

	//Entity
	AddEntity(userID string, collectionID string, entityName string) error
	DeleteEntity(userID string, collectionID string, entityName string) error
	UpdateEntity(userID string, collectionID string, entityName string, entity EntityRecord) error
	IsEntityExistWithName(userID string, collectionID string, entityName string) bool
	GetEntities(userID string, collectionID string) ([]*EntityRecord, error)

	Close()
}

type lbsRepository struct {
	session *mgo.Session
}

// ConnectMongo connects to a mongo database collection, using the provided username, password, and certificate file
// It returns a pointer to the session and collection objects, or an error if the connection attempt fails.
func connectWithMongo(mongoURI string, username string, password string, cert string) (*mgo.Session, error) {
	uri := strings.TrimSuffix(mongoURI, sslSuffix)
	dialInfo, err := mgo.ParseURL(uri)
	if err != nil {
		logrus.WithError(err).Errorf("Cannot parse Mongo Connection URI")
		return nil, err
	}
	dialInfo.FailFast = true
	dialInfo.Timeout = 10 * time.Second

	// only do ssl if we have the suffix
	if strings.HasSuffix(mongoURI, sslSuffix) {
		logrus.Debugf("Using TLS for mongo ...")
		tlsConfig := &tls.Config{}
		roots := x509.NewCertPool()
		if ca, err := ioutil.ReadFile(cert); err == nil {
			roots.AppendCertsFromPEM(ca)
		}
		tlsConfig.RootCAs = roots
		tlsConfig.InsecureSkipVerify = false
		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
			return conn, err
		}
	}

	// in case the username/password are not part of the URL string
	if username != "" && password != "" {
		dialInfo.Username = username
		dialInfo.Password = password
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func NewRepository() Repository {
	/*
		mongoURI := conf.GetMongoAddr()
		username := viper.GetString("mongo.username")
		password := viper.GetString("mongo.password")
		cert := viper.GetString("mongo.cert")

		session, err := connectWithMongo(mongoURI, username, password, cert)
		if err != nil {
			logrus.WithError(err).Fatalf("Cannot create repository with %s %s %s", mongoURI, username, password)
		}
		ensureIndex(session, collectionCollections, "collection_id")

		return &lbsRepository{session: session}
	*/
	return &lbsRepository{}
}

// Helper
func (r *lbsRepository) Close() {
	r.session.Close()
}

func getCollection(repo *lbsRepository, collectionName string) (*mgo.Session, *mgo.Collection) {
	sess := repo.session.Clone()
	c := sess.DB(dbnameRepository).C(collectionName)
	return sess, c
}

func ensureIndex(s *mgo.Session, collection string, keys ...string) {
	session := s.Clone()
	defer session.Close()

	c := session.DB(dbnameRepository).C(collection)
	index := mgo.Index{
		Key:        keys,
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		logrus.Errorf("Create dev INDEX error! %s\n", err)
	}
}

func (r *lbsRepository) AddCollection(userID string, collectionID string) error {
	sess, c := getCollection(r, collectionCollections)
	defer sess.Close()
	return c.Insert(&Collection{
		UserID:        userID,
		CollectionID:  collectionID,
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
	})
}

func (r *lbsRepository) RemoveCollection(userID, collectionID string) error {
	sess, c := getCollection(r, collectionCollections)
	defer sess.Close()
	return c.Remove(bson.M{"user_id": userID, "collection_id": collectionID})
}

func (r *lbsRepository) GetFences(userID, collectionID string) ([]*geofenceRecord, error) {
	sess, c := getCollection(r, collectionGeofences)
	defer sess.Close()

	fences := []*geofenceRecord{}
	err := c.Find(bson.M{"user_id": userID, "collection_id": collectionID}).All(&fences)
	return fences, err
}

func (r *lbsRepository) GetFenceUserID(fenceID string) (string, error) {
	userID := &UserID{}
	sess, c := getCollection(r, collectionGeofences)
	defer sess.Close()

	err := c.Find(bson.M{"fence_id": fenceID}).One(&userID)
	return userID.UserID, err
}

func (r *lbsRepository) GetAllCollections() ([]*Collection, error) {
	sess, c := getCollection(r, collectionCollections)
	defer sess.Close()

	collections := []*Collection{}
	err := c.Find(bson.M{}).All(&collections)
	return collections, err
}

func (r *lbsRepository) UpdateCollection(userID string, p *Collection) error {
	sess, c := getCollection(r, collectionCollections)
	defer sess.Close()

	return c.Update(bson.M{"user_id": userID, "collection_id": p.CollectionID}, p)
}

func (r *lbsRepository) AddGeofence(userID string, collectionID string, fenceName string, fenceID string) error {
	sess, c := getCollection(r, collectionGeofences)
	defer sess.Close()

	return c.Insert(&geofenceRecord{
		UserID:        userID,
		CollectionID:  collectionID,
		FenceName:     fenceName,
		FenceID:       fenceID,
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
	})
}

func (r *lbsRepository) AddEntity(userID string, collectionID string, entityName string) error {
	sess, c := getCollection(r, collectionEntity)
	defer sess.Close()

	return c.Insert(&EntityRecord{
		UserID:        userID,
		CollectionID:  collectionID,
		EntityName:    entityName,
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
	})
}

func (r *lbsRepository) DeleteEntity(userID string, collectionID string, entityName string) error {
	sess, c := getCollection(r, collectionEntity)
	defer sess.Close()

	return c.Remove(bson.M{"user_id": userID, "collection_id": collectionID, "entity_name": entityName})
}

func (r *lbsRepository) UpdateEntity(userID string, collectionID string, entityName string, entity EntityRecord) error {
	sess, c := getCollection(r, collectionEntity)
	defer sess.Close()

	return c.Update(bson.M{"user_id": userID, "collection_id": collectionID, "entity_name": entityName}, entity)
}

func (r *lbsRepository) IsEntityExistWithName(userID string, collectionID string, entityName string) bool {
	sess, c := getCollection(r, collectionEntity)
	defer sess.Close()

	if err := c.Find(bson.M{"user_id": userID, "collection_id": collectionID, "entity_name": entityName}).One(nil); err == nil {
		return true
	}

	return false
}

func (r *lbsRepository) GetEntities(userID, collectionID string) ([]*EntityRecord, error) {
	sess, c := getCollection(r, collectionEntity)
	defer sess.Close()

	entities := []*EntityRecord{}
	err := c.Find(bson.M{"user_id": userID, "collection_id": collectionID}).All(&entities)
	return entities, err
}

func (r *lbsRepository) RemoveGeofence(userID string, collectionID string, fenceID string) error {
	sess, c := getCollection(r, collectionGeofences)
	defer sess.Close()
	return c.Remove(bson.M{"user_id": userID, "collection_id": collectionID, "fence_id": fenceID})
}

func (r *lbsRepository) IsGeofenceExistWithName(userID string, collectionID string, fenceName string) bool {
	sess, c := getCollection(r, collectionGeofences)
	defer sess.Close()

	if err := c.Find(bson.M{"user_id": userID, "collection_id": collectionID, "fence_name": fenceName}).One(nil); err == nil {
		return true
	}

	return false
}
func (r *lbsRepository) IsGeofenceExistWithId(userID string, collectionID string, fenceID string) bool {
	sess, c := getCollection(r, collectionGeofences)
	defer sess.Close()

	if err := c.Find(bson.M{"user_id": userID, "collection_id": collectionID, "fence_id": fenceID}).One(nil); err == nil {
		return true
	}

	return false

}
