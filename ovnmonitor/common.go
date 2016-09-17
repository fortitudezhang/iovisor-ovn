package ovnmonitor

import (
	"reflect"

	l "github.com/op/go-logging"

	"github.com/socketplane/libovsdb"
)

var log = l.MustGetLogger("politoctrl")

type MonitorHandler struct {
	Quit         chan bool
	Update       chan *libovsdb.TableUpdates
	Bufupdate    chan string
	BufupdateOvs chan string
	Cache        *map[string]map[string]libovsdb.Row
	Db           *libovsdb.OvsdbClient
}

func PrintRow(row libovsdb.Row) {
	for key, value := range row.Fields {
		log.Debugf("%20s : %s\n", key, value)
	}
}

//DUMP ALL THE DB
func PrintCache(h *MonitorHandler) {
	var cache = *h.Cache
	log.Noticef("print all tables in cache\n")
	for tableName, table := range cache {
		log.Noticef("%20s:%s\n", "TABLE", tableName)
		for uuid, row := range table {
			log.Noticef("%20s:%s\n", "UUID", uuid)
			PrintRow(row)
		}
	}
}

//Print a table
func PrintCacheTable(h *MonitorHandler, tab string) {
	var cache = *h.Cache
	log.Noticef("print table %s\n", tab)
	for tableName, table := range cache {
		if tableName == tab {
			log.Noticef("%20s:%s\n", "TABLE", tableName)
			for uuid, row := range table {
				log.Noticef("%20s:%s\n", "UUID", uuid)
				PrintRow(row)
			}
		}
	}
}

func PopulateCache(h *MonitorHandler, updates libovsdb.TableUpdates) {
	var cache = *h.Cache
	for table, tableUpdate := range updates.Updates {
		if _, ok := cache[table]; !ok {
			cache[table] = make(map[string]libovsdb.Row)
		}
		for uuid, row := range tableUpdate.Rows {
			empty := libovsdb.Row{}
			if !reflect.DeepEqual(row.New, empty) {
				cache[table][uuid] = row.New
			} else {
				delete(cache[table], uuid)
			}
		}
	}
}

type MyNotifier struct {
	handler *MonitorHandler
}

func (n MyNotifier) Update(context interface{}, tableUpdates libovsdb.TableUpdates /*, h *MonitorHandler*/) {
	PopulateCache(n.handler, tableUpdates)
	n.handler.Update <- &tableUpdates
}
func (n MyNotifier) Locked([]interface{}) {
}
func (n MyNotifier) Stolen([]interface{}) {
}
func (n MyNotifier) Echo([]interface{}) {
}
func (n MyNotifier) Disconnected(client *libovsdb.OvsdbClient) {
}
