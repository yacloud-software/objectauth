package db

/*
 This file was created by mkdb-client.
 The intention is not to modify thils file, but you may extend the struct DBUserToObject
 in a seperate file (so that you can regenerate this one from time to time)
*/

/*
 PRIMARY KEY: ID
*/

/*
 postgres:
 create sequence usertoobject_seq;

Main Table:

 CREATE TABLE usertoobject (id integer primary key default nextval('usertoobject_seq'),objecttype integer not null  ,objectid bigint not null  ,userid text not null  ,active boolean not null  ,read boolean not null  ,write boolean not null  ,execute boolean not null  ,view boolean not null  );

Alter statements:
ALTER TABLE usertoobject ADD COLUMN objecttype integer not null default 0;
ALTER TABLE usertoobject ADD COLUMN objectid bigint not null default 0;
ALTER TABLE usertoobject ADD COLUMN userid text not null default '';
ALTER TABLE usertoobject ADD COLUMN active boolean not null default false;
ALTER TABLE usertoobject ADD COLUMN read boolean not null default false;
ALTER TABLE usertoobject ADD COLUMN write boolean not null default false;
ALTER TABLE usertoobject ADD COLUMN execute boolean not null default false;
ALTER TABLE usertoobject ADD COLUMN view boolean not null default false;


Archive Table: (structs can be moved from main to archive using Archive() function)

 CREATE TABLE usertoobject_archive (id integer unique not null,objecttype integer not null,objectid bigint not null,userid text not null,active boolean not null,read boolean not null,write boolean not null,execute boolean not null,view boolean not null);
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
	default_def_DBUserToObject *DBUserToObject
)

type DBUserToObject struct {
	DB                  *sql.DB
	SQLTablename        string
	SQLArchivetablename string
}

func DefaultDBUserToObject() *DBUserToObject {
	if default_def_DBUserToObject != nil {
		return default_def_DBUserToObject
	}
	psql, err := sql.Open()
	if err != nil {
		fmt.Printf("Failed to open database: %s\n", err)
		os.Exit(10)
	}
	res := NewDBUserToObject(psql)
	ctx := context.Background()
	err = res.CreateTable(ctx)
	if err != nil {
		fmt.Printf("Failed to create table: %s\n", err)
		os.Exit(10)
	}
	default_def_DBUserToObject = res
	return res
}
func NewDBUserToObject(db *sql.DB) *DBUserToObject {
	foo := DBUserToObject{DB: db}
	foo.SQLTablename = "usertoobject"
	foo.SQLArchivetablename = "usertoobject_archive"
	return &foo
}

// archive. It is NOT transactionally save.
func (a *DBUserToObject) Archive(ctx context.Context, id uint64) error {

	// load it
	p, err := a.ByID(ctx, id)
	if err != nil {
		return err
	}

	// now save it to archive:
	_, e := a.DB.ExecContext(ctx, "archive_DBUserToObject", "insert into "+a.SQLArchivetablename+"+ (id,objecttype, objectid, userid, active, read, write, execute, view) values ($1,$2, $3, $4, $5, $6, $7, $8, $9) ", p.ID, p.ObjectType, p.ObjectID, p.UserID, p.Active, p.Read, p.Write, p.Execute, p.View)
	if e != nil {
		return e
	}

	// now delete it.
	a.DeleteByID(ctx, id)
	return nil
}

// Save (and use database default ID generation)
func (a *DBUserToObject) Save(ctx context.Context, p *savepb.UserToObject) (uint64, error) {
	qn := "DBUserToObject_Save"
	rows, e := a.DB.QueryContext(ctx, qn, "insert into "+a.SQLTablename+" (objecttype, objectid, userid, active, read, write, execute, view) values ($1, $2, $3, $4, $5, $6, $7, $8) returning id", p.ObjectType, p.ObjectID, p.UserID, p.Active, p.Read, p.Write, p.Execute, p.View)
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
func (a *DBUserToObject) SaveWithID(ctx context.Context, p *savepb.UserToObject) error {
	qn := "insert_DBUserToObject"
	_, e := a.DB.ExecContext(ctx, qn, "insert into "+a.SQLTablename+" (id,objecttype, objectid, userid, active, read, write, execute, view) values ($1,$2, $3, $4, $5, $6, $7, $8, $9) ", p.ID, p.ObjectType, p.ObjectID, p.UserID, p.Active, p.Read, p.Write, p.Execute, p.View)
	return a.Error(ctx, qn, e)
}

func (a *DBUserToObject) Update(ctx context.Context, p *savepb.UserToObject) error {
	qn := "DBUserToObject_Update"
	_, e := a.DB.ExecContext(ctx, qn, "update "+a.SQLTablename+" set objecttype=$1, objectid=$2, userid=$3, active=$4, read=$5, write=$6, execute=$7, view=$8 where id = $9", p.ObjectType, p.ObjectID, p.UserID, p.Active, p.Read, p.Write, p.Execute, p.View, p.ID)

	return a.Error(ctx, qn, e)
}

// delete by id field
func (a *DBUserToObject) DeleteByID(ctx context.Context, p uint64) error {
	qn := "deleteDBUserToObject_ByID"
	_, e := a.DB.ExecContext(ctx, qn, "delete from "+a.SQLTablename+" where id = $1", p)
	return a.Error(ctx, qn, e)
}

// get it by primary id
func (a *DBUserToObject) ByID(ctx context.Context, p uint64) (*savepb.UserToObject, error) {
	qn := "DBUserToObject_ByID"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,objecttype, objectid, userid, active, read, write, execute, view from "+a.SQLTablename+" where id = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByID: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByID: error scanning (%s)", e))
	}
	if len(l) == 0 {
		return nil, a.Error(ctx, qn, fmt.Errorf("No UserToObject with id %d", p))
	}
	if len(l) != 1 {
		return nil, a.Error(ctx, qn, fmt.Errorf("Multiple (%d) UserToObject with id %d", len(l), p))
	}
	return l[0], nil
}

// get all rows
func (a *DBUserToObject) All(ctx context.Context) ([]*savepb.UserToObject, error) {
	qn := "DBUserToObject_all"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,objecttype, objectid, userid, active, read, write, execute, view from "+a.SQLTablename+" order by id")
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

// get all "DBUserToObject" rows with matching ObjectType
func (a *DBUserToObject) ByObjectType(ctx context.Context, p uint32) ([]*savepb.UserToObject, error) {
	qn := "DBUserToObject_ByObjectType"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,objecttype, objectid, userid, active, read, write, execute, view from "+a.SQLTablename+" where objecttype = $1", p)
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
func (a *DBUserToObject) ByLikeObjectType(ctx context.Context, p uint32) ([]*savepb.UserToObject, error) {
	qn := "DBUserToObject_ByLikeObjectType"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,objecttype, objectid, userid, active, read, write, execute, view from "+a.SQLTablename+" where objecttype ilike $1", p)
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

// get all "DBUserToObject" rows with matching ObjectID
func (a *DBUserToObject) ByObjectID(ctx context.Context, p uint64) ([]*savepb.UserToObject, error) {
	qn := "DBUserToObject_ByObjectID"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,objecttype, objectid, userid, active, read, write, execute, view from "+a.SQLTablename+" where objectid = $1", p)
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
func (a *DBUserToObject) ByLikeObjectID(ctx context.Context, p uint64) ([]*savepb.UserToObject, error) {
	qn := "DBUserToObject_ByLikeObjectID"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,objecttype, objectid, userid, active, read, write, execute, view from "+a.SQLTablename+" where objectid ilike $1", p)
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

// get all "DBUserToObject" rows with matching UserID
func (a *DBUserToObject) ByUserID(ctx context.Context, p string) ([]*savepb.UserToObject, error) {
	qn := "DBUserToObject_ByUserID"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,objecttype, objectid, userid, active, read, write, execute, view from "+a.SQLTablename+" where userid = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByUserID: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByUserID: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBUserToObject) ByLikeUserID(ctx context.Context, p string) ([]*savepb.UserToObject, error) {
	qn := "DBUserToObject_ByLikeUserID"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,objecttype, objectid, userid, active, read, write, execute, view from "+a.SQLTablename+" where userid ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByUserID: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByUserID: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBUserToObject" rows with matching Active
func (a *DBUserToObject) ByActive(ctx context.Context, p bool) ([]*savepb.UserToObject, error) {
	qn := "DBUserToObject_ByActive"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,objecttype, objectid, userid, active, read, write, execute, view from "+a.SQLTablename+" where active = $1", p)
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
func (a *DBUserToObject) ByLikeActive(ctx context.Context, p bool) ([]*savepb.UserToObject, error) {
	qn := "DBUserToObject_ByLikeActive"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,objecttype, objectid, userid, active, read, write, execute, view from "+a.SQLTablename+" where active ilike $1", p)
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

// get all "DBUserToObject" rows with matching Read
func (a *DBUserToObject) ByRead(ctx context.Context, p bool) ([]*savepb.UserToObject, error) {
	qn := "DBUserToObject_ByRead"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,objecttype, objectid, userid, active, read, write, execute, view from "+a.SQLTablename+" where read = $1", p)
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
func (a *DBUserToObject) ByLikeRead(ctx context.Context, p bool) ([]*savepb.UserToObject, error) {
	qn := "DBUserToObject_ByLikeRead"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,objecttype, objectid, userid, active, read, write, execute, view from "+a.SQLTablename+" where read ilike $1", p)
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

// get all "DBUserToObject" rows with matching Write
func (a *DBUserToObject) ByWrite(ctx context.Context, p bool) ([]*savepb.UserToObject, error) {
	qn := "DBUserToObject_ByWrite"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,objecttype, objectid, userid, active, read, write, execute, view from "+a.SQLTablename+" where write = $1", p)
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
func (a *DBUserToObject) ByLikeWrite(ctx context.Context, p bool) ([]*savepb.UserToObject, error) {
	qn := "DBUserToObject_ByLikeWrite"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,objecttype, objectid, userid, active, read, write, execute, view from "+a.SQLTablename+" where write ilike $1", p)
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

// get all "DBUserToObject" rows with matching Execute
func (a *DBUserToObject) ByExecute(ctx context.Context, p bool) ([]*savepb.UserToObject, error) {
	qn := "DBUserToObject_ByExecute"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,objecttype, objectid, userid, active, read, write, execute, view from "+a.SQLTablename+" where execute = $1", p)
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
func (a *DBUserToObject) ByLikeExecute(ctx context.Context, p bool) ([]*savepb.UserToObject, error) {
	qn := "DBUserToObject_ByLikeExecute"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,objecttype, objectid, userid, active, read, write, execute, view from "+a.SQLTablename+" where execute ilike $1", p)
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

// get all "DBUserToObject" rows with matching View
func (a *DBUserToObject) ByView(ctx context.Context, p bool) ([]*savepb.UserToObject, error) {
	qn := "DBUserToObject_ByView"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,objecttype, objectid, userid, active, read, write, execute, view from "+a.SQLTablename+" where view = $1", p)
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
func (a *DBUserToObject) ByLikeView(ctx context.Context, p bool) ([]*savepb.UserToObject, error) {
	qn := "DBUserToObject_ByLikeView"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,objecttype, objectid, userid, active, read, write, execute, view from "+a.SQLTablename+" where view ilike $1", p)
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
func (a *DBUserToObject) FromQuery(ctx context.Context, query_where string, args ...interface{}) ([]*savepb.UserToObject, error) {
	rows, err := a.DB.QueryContext(ctx, "custom_query_"+a.Tablename(), "select "+a.SelectCols()+" from "+a.Tablename()+" where "+query_where, args...)
	if err != nil {
		return nil, err
	}
	return a.FromRows(ctx, rows)
}

/**********************************************************************
* Helper to convert from an SQL Row to struct
**********************************************************************/
func (a *DBUserToObject) Tablename() string {
	return a.SQLTablename
}

func (a *DBUserToObject) SelectCols() string {
	return "id,objecttype, objectid, userid, active, read, write, execute, view"
}
func (a *DBUserToObject) SelectColsQualified() string {
	return "" + a.SQLTablename + ".id," + a.SQLTablename + ".objecttype, " + a.SQLTablename + ".objectid, " + a.SQLTablename + ".userid, " + a.SQLTablename + ".active, " + a.SQLTablename + ".read, " + a.SQLTablename + ".write, " + a.SQLTablename + ".execute, " + a.SQLTablename + ".view"
}

func (a *DBUserToObject) FromRows(ctx context.Context, rows *gosql.Rows) ([]*savepb.UserToObject, error) {
	var res []*savepb.UserToObject
	for rows.Next() {
		foo := savepb.UserToObject{}
		err := rows.Scan(&foo.ID, &foo.ObjectType, &foo.ObjectID, &foo.UserID, &foo.Active, &foo.Read, &foo.Write, &foo.Execute, &foo.View)
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
func (a *DBUserToObject) CreateTable(ctx context.Context) error {
	csql := []string{
		`create sequence if not exists ` + a.SQLTablename + `_seq;`,
		`CREATE TABLE if not exists ` + a.SQLTablename + ` (id integer primary key default nextval('` + a.SQLTablename + `_seq'),objecttype integer not null  ,objectid bigint not null  ,userid text not null  ,active boolean not null  ,read boolean not null  ,write boolean not null  ,execute boolean not null  ,view boolean not null  );`,
		`CREATE TABLE if not exists ` + a.SQLTablename + `_archive (id integer primary key default nextval('` + a.SQLTablename + `_seq'),objecttype integer not null  ,objectid bigint not null  ,userid text not null  ,active boolean not null  ,read boolean not null  ,write boolean not null  ,execute boolean not null  ,view boolean not null  );`,
	}
	for i, c := range csql {
		_, e := a.DB.ExecContext(ctx, fmt.Sprintf("create_"+a.SQLTablename+"_%d", i), c)
		if e != nil {
			return e
		}
	}
	return nil
}

/**********************************************************************
* Helper to meaningful errors
**********************************************************************/
func (a *DBUserToObject) Error(ctx context.Context, q string, e error) error {
	if e == nil {
		return nil
	}
	return fmt.Errorf("[table="+a.SQLTablename+", query=%s] Error: %s", q, e)
}
