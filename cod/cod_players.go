// Package players allows other plugins to reference players within the database.
package players

import (
	"database/sql"
	"github.com/adabei/goldenbot/events"
	"github.com/adabei/goldenbot/events/cod"
	"github.com/adabei/goldenbot/rcon"
	"log"
)

type Players struct {
	requests chan rcon.RCONQuery
	events   chan interface{}
	db       *sql.DB
}

func NewPlayers(requests chan rcon.RCONQuery, ea events.Aggregator, db *sql.DB) *Players {
	p := new(Players)
	p.requests = requests
	p.events = ea.Subscribe(p)
	p.db = db
	return p
}

const schema = `
create table players (
  id text primary key
);`

func (p *Players) Setup() error {
	_, err := p.db.Exec(schema)
	return err
}

func (p *Players) Start() {
	p.db.Exec("insert into players(id) values('kaiser');")
	log.Println("yes inserting")

	for {
		ev := <-p.events

		switch ev := ev.(type) {
		case cod.Join:
			if !exists(p.db, ev.GUID) {
				log.Println("inserting ", ev.GUID, "into players")
				_, err := p.db.Exec("insert into players(id) values(?);", ev.GUID)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

func exists(db *sql.DB, id string) bool {
	var guid string
	err := db.QueryRow("select id from players where id = ?", id).Scan(&guid)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
