package db

/*
 This file was created by mkdb-client.
 The intention is not to modify thils file, but you may extend the struct DBGroupToObject
 in a seperate file (so that you can regenerate this one from time to time)
*/

/*
 PRIMARY KEY: ID
*/

/*
 postgres:
 create sequence grouptoobject_seq;

Main Table:

 CREATE TABLE grouptoobject (id integer primary key default nextval('grouptoobject_seq'),objecttype integer not null  ,objectid bigint not null  ,groupid text not null  ,active boolean not null  ,read boolean not null  ,write boolean not null  ,execute boolean not null  ,view boolean not null  );

Alter statements:
ALTER TABLE grouptoobject ADD COLUMN IF NOT EXISTS objecttype integer not null default 0;
ALTER TABLE grouptoobject ADD COLUMN IF NOT EXISTS objectid bigint not null default 0;
ALTER TABLE grouptoobject ADD COLUMN IF NOT EXISTS groupid text not null default '';
ALTER TABLE grouptoobject ADD COLUMN IF NOT EXISTS active boolean not null default false;
ALTER TABLE grouptoobject ADD COLUMN IF NOT EXISTS read boolean not null default false;
ALTER TABLE grouptoobject ADD COLUMN IF NOT EXISTS write boolean not null default false;
ALTER TABLE grouptoobject ADD COLUMN IF NOT EXISTS execute boolean not null default false;
ALTER TABLE grouptoobject ADD COLUMN IF NOT EXISTS view boolean not null default false;


Archive Table: (structs can be moved from main to archive using Archive() function)

 CREATE TABLE grouptoobject_archive (id integer unique not null,objecttype integer not null,objectid bigint not null,groupid text not null,active boolean not null,read boolean not null,write boolean not null,execute boolean not null,view boolean not null);
*/

import (
	"context"
	gosql "database/sql"
	"fmt"
	savepb "golang.conradwood.net/apis/objectauth"
	"golang.conradwood.net/go-easyops/sql"
	"os"
)

var (
	default_def_DBGroupToObject *DBGroupToObject
)

type DBGroupToObject struct {
	DB                  *sql.DB
	SQLTablename        string
	SQLArchivetablename string
}

func DefaultDBGroupToObject() *DBGroupToObject {
	if default_def_DBGroupToObject != nil {
		return default_def_DBGroupToObject
	}
	psql, err := sql.Open()
	if err != nil {
		fmt.Printf("Failed to open database: %s\n", err)
		os.Exit(10)
	}
	res := NewDBGroupToObject(psql)
	ctx := context.Background()
	err = res.CreateTable(ctx)
	if err != nil {
		fmt.Printf("Failed to create table: %s\n", err)
		os.Exit(10)
	}
	default_def_DBGroupToObject = res
	return res
}
func NewDBGroupToObject(db *sql.DB) *DBGroupToObject {
	foo := DBGroupToObject{DB: db}
	foo.SQLTablename = "grouptoobject"
	foo.SQLArchivetablename = "grouptoobject_archive"
	return &foo
}

// archive. It is NOT transactionally save.
func (a *DBGroupToObject) Archive(ctx context.Context, id uint64) error {

	// load it
	p, err := a.ByID(ctx, id)
	if err != nil {
		return err
	}

	// now save it to archive:
	_, e := a.DB.ExecContext(ctx, "archive_DBGroupToObject", "insert into "+a.SQLArchivetablename+" (id,objecttype, objectid, groupid, active, read, write, execute, view) values ($1,$2, $3, $4, $5, $6, $7, $8, $9) ", p.ID, p.ObjectType, p.ObjectID, p.GroupID, p.Active, p.Read, p.Write, p.Execute, p.View)
	if e != nil {
		return e
	}

	// now delete it.
	a.DeleteByID(ctx, id)
	return nil
}

// Save (and use database default ID generation)
func (a *DBGroupToObject) Save(ctx context.Context, p *savepb.GroupToObject) (uint64, error) {
	qn := "DBGroupToObject_Save"
	rows, e := a.DB.QueryContext(ctx, qn, "insert into "+a.SQLTablename+" (objecttype, objectid, groupid, active, read, write, execute, view) values ($1, $2, $3, $4, $5, $6, $7, $8) returning id", p.ObjectType, p.ObjectID, p.GroupID, p.Active, p.Read, p.Write, p.Execute, p.View)
	if e != nil {
		return 0, a.Error(ctx, qn, e)
	}
	defer rows.Close()
	if !rows.Next() {
		return 0, a.Error(ctx, qn, fmt.Errorf("No rows after insert"))
	}
	var id uint64
	e = rows.Scan(&id)
	if e != nil {
		return 0, a.Error(ctx, qn, fmt.Errorf("failed to scan id after insert: %s", e))
	}
	p.ID = id
	return id, nil
}

// Save using the ID specified
func (a *DBGroupToObject) SaveWithID(ctx context.Context, p *savepb.GroupToObject) error {
	qn := "insert_DBGroupToObject"
	_, e := a.DB.ExecContext(ctx, qn, "insert into "+a.SQLTablename+" (id,objecttype, objectid, groupid, active, read, write, execute, view) values ($1,$2, $3, $4, $5, $6, $7, $8, $9) ", p.ID, p.ObjectType, p.ObjectID, p.GroupID, p.Active, p.Read, p.Write, p.Execute, p.View)
	return a.Error(ctx, qn, e)
}

func (a *DBGroupToObject) Update(ctx context.Context, p *savepb.GroupToObject) error {
	qn := "DBGroupToObject_Update"
	_, e := a.DB.ExecContext(ctx, qn, "update "+a.SQLTablename+" set objecttype=$1, objectid=$2, groupid=$3, active=$4, read=$5, write=$6, execute=$7, view=$8 where id = $9", p.ObjectType, p.ObjectID, p.GroupID, p.Active, p.Read, p.Write, p.Execute, p.View, p.ID)

	return a.Error(ctx, qn, e)
}

// delete by id field
func (a *DBGroupToObject) DeleteByID(ctx context.Context, p uint64) error {
	qn := "deleteDBGroupToObject_ByID"
	_, e := a.DB.ExecContext(ctx, qn, "delete from "+a.SQLTablename+" where id = $1", p)
	return a.Error(ctx, qn, e)
}

// get it by primary id
func (a *DBGroupToObject) ByID(ctx context.Context, p uint64) (*savepb.GroupToObject, error) {
	qn := "DBGroupToObject_ByID"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,objecttype, objectid, groupid, active, read, write, execute, view from "+a.SQLTablename+" where id = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByID: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByID: error scanning (%s)", e))
	}
	if len(l) == 0 {
		return nil, a.Error(ctx, qn, fmt.Errorf("No GroupToObject with id %v", p))
	}
	if len(l) != 1 {
		return nil, a.Error(ctx, qn, fmt.Errorf("Multiple (%d) GroupToObject with id %v", len(l), p))
	}
	return l[0], nil
}

// get it by primary id (nil if no such ID row, but no error either)
func (a *DBGroupToObject) TryByID(ctx context.Context, p uint64) (*savepb.GroupToObject, error) {
	qn := "DBGroupToObject_TryByID"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,objecttype, objectid, groupid, active, read, write, execute, view from "+a.SQLTablename+" where id = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("TryByID: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("TryByID: error scanning (%s)", e))
	}
	if len(l) == 0 {
		return nil, nil
	}
	if len(l) != 1 {
		return nil, a.Error(ctx, qn, fmt.Errorf("Multiple (%d) GroupToObject with id %v", len(l), p))
	}
	return l[0], nil
}

// get all rows
func (a *DBGroupToObject) All(ctx context.Context) ([]*savepb.GroupToObject, error) {
	qn := "DBGroupToObject_all"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,objecttype, objectid, groupid, active, read, write, execute, view from "+a.SQLTablename+" order by id")
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("All: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, fmt.Errorf("All: error scanning (%s)", e)
	}
	return l, nil
}

/**********************************************************************
* GetBy[FIELD] functions
**********************************************************************/

// get all "DBGroupToObject" rows with matching ObjectType
func (a *DBGroupToObject) ByObjectType(ctx context.Context, p uint32) ([]*savepb.GroupToObject, error) {
	qn := "DBGroupToObject_ByObjectType"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,objecttype, objectid, groupid, active, read, write, execute, view from "+a.SQLTablename+" where objecttype = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByObjectType: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByObjectType: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBGroupToObject) ByLikeObjectType(ctx context.Context, p uint32) ([]*savepb.GroupToObject, error) {
	qn := "DBGroupToObject_ByLikeObjectType"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,objecttype, objectid, groupid, active, read, write, execute, view from "+a.SQLTablename+" where objecttype ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByObjectType: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByObjectType: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBGroupToObject" rows with matching ObjectID
func (a *DBGroupToObject) ByObjectID(ctx context.Context, p uint64) ([]*savepb.GroupToObject, error) {
	qn := "DBGroupToObject_ByObjectID"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,objecttype, objectid, groupid, active, read, write, execute, view from "+a.SQLTablename+" where objectid = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByObjectID: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByObjectID: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBGroupToObject) ByLikeObjectID(ctx context.Context, p uint64) ([]*savepb.GroupToObject, error) {
	qn := "DBGroupToObject_ByLikeObjectID"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,objecttype, objectid, groupid, active, read, write, execute, view from "+a.SQLTablename+" where objectid ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByObjectID: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByObjectID: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBGroupToObject" rows with matching GroupID
func (a *DBGroupToObject) ByGroupID(ctx context.Context, p string) ([]*savepb.GroupToObject, error) {
	qn := "DBGroupToObject_ByGroupID"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,objecttype, objectid, groupid, active, read, write, execute, view from "+a.SQLTablename+" where groupid = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByGroupID: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByGroupID: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBGroupToObject) ByLikeGroupID(ctx context.Context, p string) ([]*savepb.GroupToObject, error) {
	qn := "DBGroupToObject_ByLikeGroupID"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,objecttype, objectid, groupid, active, read, write, execute, view from "+a.SQLTablename+" where groupid ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByGroupID: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByGroupID: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBGroupToObject" rows with matching Active
func (a *DBGroupToObject) ByActive(ctx context.Context, p bool) ([]*savepb.GroupToObject, error) {
	qn := "DBGroupToObject_ByActive"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,objecttype, objectid, groupid, active, read, write, execute, view from "+a.SQLTablename+" where active = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByActive: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByActive: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBGroupToObject) ByLikeActive(ctx context.Context, p bool) ([]*savepb.GroupToObject, error) {
	qn := "DBGroupToObject_ByLikeActive"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,objecttype, objectid, groupid, active, read, write, execute, view from "+a.SQLTablename+" where active ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByActive: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByActive: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBGroupToObject" rows with matching Read
func (a *DBGroupToObject) ByRead(ctx context.Context, p bool) ([]*savepb.GroupToObject, error) {
	qn := "DBGroupToObject_ByRead"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,objecttype, objectid, groupid, active, read, write, execute, view from "+a.SQLTablename+" where read = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByRead: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByRead: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBGroupToObject) ByLikeRead(ctx context.Context, p bool) ([]*savepb.GroupToObject, error) {
	qn := "DBGroupToObject_ByLikeRead"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,objecttype, objectid, groupid, active, read, write, execute, view from "+a.SQLTablename+" where read ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByRead: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByRead: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBGroupToObject" rows with matching Write
func (a *DBGroupToObject) ByWrite(ctx context.Context, p bool) ([]*savepb.GroupToObject, error) {
	qn := "DBGroupToObject_ByWrite"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,objecttype, objectid, groupid, active, read, write, execute, view from "+a.SQLTablename+" where write = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByWrite: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByWrite: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBGroupToObject) ByLikeWrite(ctx context.Context, p bool) ([]*savepb.GroupToObject, error) {
	qn := "DBGroupToObject_ByLikeWrite"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,objecttype, objectid, groupid, active, read, write, execute, view from "+a.SQLTablename+" where write ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByWrite: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByWrite: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBGroupToObject" rows with matching Execute
func (a *DBGroupToObject) ByExecute(ctx context.Context, p bool) ([]*savepb.GroupToObject, error) {
	qn := "DBGroupToObject_ByExecute"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,objecttype, objectid, groupid, active, read, write, execute, view from "+a.SQLTablename+" where execute = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByExecute: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByExecute: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBGroupToObject) ByLikeExecute(ctx context.Context, p bool) ([]*savepb.GroupToObject, error) {
	qn := "DBGroupToObject_ByLikeExecute"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,objecttype, objectid, groupid, active, read, write, execute, view from "+a.SQLTablename+" where execute ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByExecute: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByExecute: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBGroupToObject" rows with matching View
func (a *DBGroupToObject) ByView(ctx context.Context, p bool) ([]*savepb.GroupToObject, error) {
	qn := "DBGroupToObject_ByView"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,objecttype, objectid, groupid, active, read, write, execute, view from "+a.SQLTablename+" where view = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByView: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByView: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBGroupToObject) ByLikeView(ctx context.Context, p bool) ([]*savepb.GroupToObject, error) {
	qn := "DBGroupToObject_ByLikeView"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,objecttype, objectid, groupid, active, read, write, execute, view from "+a.SQLTablename+" where view ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByView: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByView: error scanning (%s)", e))
	}
	return l, nil
}

/**********************************************************************
* Helper to convert from an SQL Query
**********************************************************************/

// from a query snippet (the part after WHERE)
func (a *DBGroupToObject) FromQuery(ctx context.Context, query_where string, args ...interface{}) ([]*savepb.GroupToObject, error) {
	rows, err := a.DB.QueryContext(ctx, "custom_query_"+a.Tablename(), "select "+a.SelectCols()+" from "+a.Tablename()+" where "+query_where, args...)
	if err != nil {
		return nil, err
	}
	return a.FromRows(ctx, rows)
}

/**********************************************************************
* Helper to convert from an SQL Row to struct
**********************************************************************/
func (a *DBGroupToObject) Tablename() string {
	return a.SQLTablename
}

func (a *DBGroupToObject) SelectCols() string {
	return "id,objecttype, objectid, groupid, active, read, write, execute, view"
}
func (a *DBGroupToObject) SelectColsQualified() string {
	return "" + a.SQLTablename + ".id," + a.SQLTablename + ".objecttype, " + a.SQLTablename + ".objectid, " + a.SQLTablename + ".groupid, " + a.SQLTablename + ".active, " + a.SQLTablename + ".read, " + a.SQLTablename + ".write, " + a.SQLTablename + ".execute, " + a.SQLTablename + ".view"
}

func (a *DBGroupToObject) FromRows(ctx context.Context, rows *gosql.Rows) ([]*savepb.GroupToObject, error) {
	var res []*savepb.GroupToObject
	for rows.Next() {
		foo := savepb.GroupToObject{}
		err := rows.Scan(&foo.ID, &foo.ObjectType, &foo.ObjectID, &foo.GroupID, &foo.Active, &foo.Read, &foo.Write, &foo.Execute, &foo.View)
		if err != nil {
			return nil, a.Error(ctx, "fromrow-scan", err)
		}
		res = append(res, &foo)
	}
	return res, nil
}

/**********************************************************************
* Helper to create table and columns
**********************************************************************/
func (a *DBGroupToObject) CreateTable(ctx context.Context) error {
	csql := []string{
		`create sequence if not exists ` + a.SQLTablename + `_seq;`,
		`CREATE TABLE if not exists ` + a.SQLTablename + ` (id integer primary key default nextval('` + a.SQLTablename + `_seq'),objecttype integer not null ,objectid bigint not null ,groupid text not null ,active boolean not null ,read boolean not null ,write boolean not null ,execute boolean not null ,view boolean not null );`,
		`CREATE TABLE if not exists ` + a.SQLTablename + `_archive (id integer primary key default nextval('` + a.SQLTablename + `_seq'),objecttype integer not null ,objectid bigint not null ,groupid text not null ,active boolean not null ,read boolean not null ,write boolean not null ,execute boolean not null ,view boolean not null );`,
		`ALTER TABLE grouptoobject ADD COLUMN IF NOT EXISTS objecttype integer not null default 0;`,
		`ALTER TABLE grouptoobject ADD COLUMN IF NOT EXISTS objectid bigint not null default 0;`,
		`ALTER TABLE grouptoobject ADD COLUMN IF NOT EXISTS groupid text not null default '';`,
		`ALTER TABLE grouptoobject ADD COLUMN IF NOT EXISTS active boolean not null default false;`,
		`ALTER TABLE grouptoobject ADD COLUMN IF NOT EXISTS read boolean not null default false;`,
		`ALTER TABLE grouptoobject ADD COLUMN IF NOT EXISTS write boolean not null default false;`,
		`ALTER TABLE grouptoobject ADD COLUMN IF NOT EXISTS execute boolean not null default false;`,
		`ALTER TABLE grouptoobject ADD COLUMN IF NOT EXISTS view boolean not null default false;`,

		`ALTER TABLE grouptoobject_archive ADD COLUMN IF NOT EXISTS objecttype integer not null default 0;`,
		`ALTER TABLE grouptoobject_archive ADD COLUMN IF NOT EXISTS objectid bigint not null default 0;`,
		`ALTER TABLE grouptoobject_archive ADD COLUMN IF NOT EXISTS groupid text not null default '';`,
		`ALTER TABLE grouptoobject_archive ADD COLUMN IF NOT EXISTS active boolean not null default false;`,
		`ALTER TABLE grouptoobject_archive ADD COLUMN IF NOT EXISTS read boolean not null default false;`,
		`ALTER TABLE grouptoobject_archive ADD COLUMN IF NOT EXISTS write boolean not null default false;`,
		`ALTER TABLE grouptoobject_archive ADD COLUMN IF NOT EXISTS execute boolean not null default false;`,
		`ALTER TABLE grouptoobject_archive ADD COLUMN IF NOT EXISTS view boolean not null default false;`,
	}
	for i, c := range csql {
		_, e := a.DB.ExecContext(ctx, fmt.Sprintf("create_"+a.SQLTablename+"_%d", i), c)
		if e != nil {
			return e
		}
	}

	// these are optional, expected to fail
	csql = []string{
		// Indices:

		// Foreign keys:

	}
	for i, c := range csql {
		a.DB.ExecContextQuiet(ctx, fmt.Sprintf("create_"+a.SQLTablename+"_%d", i), c)
	}
	return nil
}

/**********************************************************************
* Helper to meaningful errors
**********************************************************************/
func (a *DBGroupToObject) Error(ctx context.Context, q string, e error) error {
	if e == nil {
		return nil
	}
	return fmt.Errorf("[table="+a.SQLTablename+", query=%s] Error: %s", q, e)
}


