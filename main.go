package main

import (
	"log"
	"time"

	"github.com/caicloud/fornax/pkg/osutil"
	"github.com/caicloud/fornax/pkg/wait"
	"github.com/satori/go.uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	defaultDBName         string = "test"
	defaultCollectionName string = "rwCollection"
)

const (
	// Grace period waiting for mongodb to be running.
	MongoGracePeriod = 30 * time.Second
)

var (
	session *mgo.Session
)

type DataStore struct {
	s *mgo.Session
}

func InitStore(s *mgo.Session) {
	session = s
}

func NewStore() *DataStore {
	return &DataStore{
		s: session.Copy(),
	}
}

func (d *DataStore) Close() {
	d.s.Close()
}

type Doc struct {
	ID   string `bson:"_id,omitempty" json:"_id,omitempty"`
	Data string `bson:"data,omitempty" json:"data,omitempty"`
}

func (d *DataStore) NewDocument(doc *Doc) (string, error) {
	doc.ID = uuid.NewV4().String()
	col := d.s.DB(defaultDBName).C(defaultCollectionName)
	_, err := col.Upsert(bson.M{"_id": doc.ID}, doc)
	return doc.ID, err
}

func (d *DataStore) GetDocumentByID(id string) (*Doc, error) {
	doc := &Doc{}
	col := d.s.DB(defaultDBName).C(defaultCollectionName)
	err := col.Find(bson.M{"_id": id}).One(doc)
	return doc, err
}

func testDB() {
	ds := NewStore()
	defer ds.Close()

	write := &Doc{
		Data: time.Now().String(),
	}
	id, err := ds.NewDocument(write)
	if nil != err {
		log.Printf("write err: %v", err)
		return
	}

	read, err := ds.GetDocumentByID(id)
	if nil != err {
		log.Printf("read err: %v", err)
		return
	}
	log.Print(read)
}

func main() {
	log.Println("a mongo client for testing")

	// Get the IP of mongodb
	mongoIP := osutil.GetStringEnvWithDefault("MONGO_DB_IP", "180.101.191.213:31599")

	// Create mongo session with grace period.
	var err error
	var session *mgo.Session
	err = wait.Poll(time.Second, MongoGracePeriod, func() (done bool, err error) {
		session, err = mgo.Dial(mongoIP)
		return err == nil, nil
	})
	if err != nil {
		log.Fatal("Error dailing mongodb")
	}
	defer session.Close()
	session.SetMode(mgo.Strong, true)

	// Init data store
	InitStore(session)

	log.Println("client init finish")

	for {
		log.Println("alive")
		testDB()
		time.Sleep(time.Second)
	}
}
