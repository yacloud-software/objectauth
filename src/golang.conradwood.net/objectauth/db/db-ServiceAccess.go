package db

/*
 This file was created by mkdb-client.
 The intention is not to modify thils file, but you may extend the struct DBServiceAccess
 in a seperate file (so that you can regenerate this one from time to time)
*/

/*
 PRIMARY KEY: ID
*/

/*
 postgres:
 create sequence serviceaccess_seq;

Main Table:

 CREATE TABLE serviceaccess (id integer primary key default nextval('serviceaccess_seq'),callingservice text not null  ,subjectservice text not null  );

Alter statements:
ALTER TABLE serviceaccess ADD COLUMN IF NOT EXISTS callingservice text not null default '';
ALTER TABLE serviceaccess ADD COLUMN IF NOT EXISTS subjectservice text not null default '';


Archive Table: (structs can be moved from main to archive using Archive() function)

 CREATE TABLE serviceaccess_archive (id integer unique not null,callingservice text not null,subjectservice text not null);
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
	default_def_DBServiceAccess *DBServiceAccess
)

type DBServiceAccess struct {
	DB                  *sql.DB
	SQLTablename        string
	SQLArchivetablename string
}

func DefaultDBServiceAccess() *DBServiceAccess {
	if default_def_DBServiceAccess != nil {
		return default_def_DBServiceAccess
	}
	psql, err := sql.Open()
	if err != nil {
		fmt.Printf("Failed to open database: %s\n", err)
		os.Exit(10)
	}
	res := NewDBServiceAccess(psql)
	ctx := context.Background()
	err = res.CreateTable(ctx)
	if err != nil {
		fmt.Printf("Failed to create table: %s\n", err)
		os.Exit(10)
	}
	default_def_DBServiceAccess = res
	return res
}
func NewDBServiceAccess(db *sql.DB) *DBServiceAccess {
	foo := DBServiceAccess{DB: db}
	foo.SQLTablename = "serviceaccess"
	foo.SQLArchivetablename = "serviceaccess_archive"
	return &foo
}

// archive. It is NOT transactionally save.
func (a *DBServiceAccess) Archive(ctx context.Context, id uint64) error {

	// load it
	p, err := a.ByID(ctx, id)
	if err != nil {
		return err
	}

	// now save it to archive:
	_, e := a.DB.ExecContext(ctx, "archive_DBServiceAccess", "insert into "+a.SQLArchivetablename+" (id,callingservice, subjectservice) values ($1,$2, $3) ", p.ID, p.CallingService, p.SubjectService)
	if e != nil {
		return e
	}

	// now delete it.
	a.DeleteByID(ctx, id)
	return nil
}

// Save (and use database default ID generation)
func (a *DBServiceAccess) Save(ctx context.Context, p *savepb.ServiceAccess) (uint64, error) {
	qn := "DBServiceAccess_Save"
	rows, e := a.DB.QueryContext(ctx, qn, "insert into "+a.SQLTablename+" (callingservice, subjectservice) values ($1, $2) returning id", p.CallingService, p.SubjectService)
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
func (a *DBServiceAccess) SaveWithID(ctx context.Context, p *savepb.ServiceAccess) error {
	qn := "insert_DBServiceAccess"
	_, e := a.DB.ExecContext(ctx, qn, "insert into "+a.SQLTablename+" (id,callingservice, subjectservice) values ($1,$2, $3) ", p.ID, p.CallingService, p.SubjectService)
	return a.Error(ctx, qn, e)
}

func (a *DBServiceAccess) Update(ctx context.Context, p *savepb.ServiceAccess) error {
	qn := "DBServiceAccess_Update"
	_, e := a.DB.ExecContext(ctx, qn, "update "+a.SQLTablename+" set callingservice=$1, subjectservice=$2 where id = $3", p.CallingService, p.SubjectService, p.ID)

	return a.Error(ctx, qn, e)
}

// delete by id field
func (a *DBServiceAccess) DeleteByID(ctx context.Context, p uint64) error {
	qn := "deleteDBServiceAccess_ByID"
	_, e := a.DB.ExecContext(ctx, qn, "delete from "+a.SQLTablename+" where id = $1", p)
	return a.Error(ctx, qn, e)
}

// get it by primary id
func (a *DBServiceAccess) ByID(ctx context.Context, p uint64) (*savepb.ServiceAccess, error) {
	qn := "DBServiceAccess_ByID"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,callingservice, subjectservice from "+a.SQLTablename+" where id = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByID: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByID: error scanning (%s)", e))
	}
	if len(l) == 0 {
		return nil, a.Error(ctx, qn, fmt.Errorf("No ServiceAccess with id %v", p))
	}
	if len(l) != 1 {
		return nil, a.Error(ctx, qn, fmt.Errorf("Multiple (%d) ServiceAccess with id %v", len(l), p))
	}
	return l[0], nil
}

// get it by primary id (nil if no such ID row, but no error either)
func (a *DBServiceAccess) TryByID(ctx context.Context, p uint64) (*savepb.ServiceAccess, error) {
	qn := "DBServiceAccess_TryByID"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,callingservice, subjectservice from "+a.SQLTablename+" where id = $1", p)
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
		return nil, a.Error(ctx, qn, fmt.Errorf("Multiple (%d) ServiceAccess with id %v", len(l), p))
	}
	return l[0], nil
}

// get all rows
func (a *DBServiceAccess) All(ctx context.Context) ([]*savepb.ServiceAccess, error) {
	qn := "DBServiceAccess_all"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,callingservice, subjectservice from "+a.SQLTablename+" order by id")
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

// get all "DBServiceAccess" rows with matching CallingService
func (a *DBServiceAccess) ByCallingService(ctx context.Context, p string) ([]*savepb.ServiceAccess, error) {
	qn := "DBServiceAccess_ByCallingService"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,callingservice, subjectservice from "+a.SQLTablename+" where callingservice = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByCallingService: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByCallingService: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBServiceAccess) ByLikeCallingService(ctx context.Context, p string) ([]*savepb.ServiceAccess, error) {
	qn := "DBServiceAccess_ByLikeCallingService"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,callingservice, subjectservice from "+a.SQLTablename+" where callingservice ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByCallingService: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByCallingService: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBServiceAccess" rows with matching SubjectService
func (a *DBServiceAccess) BySubjectService(ctx context.Context, p string) ([]*savepb.ServiceAccess, error) {
	qn := "DBServiceAccess_BySubjectService"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,callingservice, subjectservice from "+a.SQLTablename+" where subjectservice = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("BySubjectService: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("BySubjectService: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBServiceAccess) ByLikeSubjectService(ctx context.Context, p string) ([]*savepb.ServiceAccess, error) {
	qn := "DBServiceAccess_ByLikeSubjectService"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,callingservice, subjectservice from "+a.SQLTablename+" where subjectservice ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("BySubjectService: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("BySubjectService: error scanning (%s)", e))
	}
	return l, nil
}

/**********************************************************************
* Helper to convert from an SQL Query
**********************************************************************/

// from a query snippet (the part after WHERE)
func (a *DBServiceAccess) FromQuery(ctx context.Context, query_where string, args ...interface{}) ([]*savepb.ServiceAccess, error) {
	rows, err := a.DB.QueryContext(ctx, "custom_query_"+a.Tablename(), "select "+a.SelectCols()+" from "+a.Tablename()+" where "+query_where, args...)
	if err != nil {
		return nil, err
	}
	return a.FromRows(ctx, rows)
}

/**********************************************************************
* Helper to convert from an SQL Row to struct
**********************************************************************/
func (a *DBServiceAccess) Tablename() string {
	return a.SQLTablename
}

func (a *DBServiceAccess) SelectCols() string {
	return "id,callingservice, subjectservice"
}
func (a *DBServiceAccess) SelectColsQualified() string {
	return "" + a.SQLTablename + ".id," + a.SQLTablename + ".callingservice, " + a.SQLTablename + ".subjectservice"
}

func (a *DBServiceAccess) FromRows(ctx context.Context, rows *gosql.Rows) ([]*savepb.ServiceAccess, error) {
	var res []*savepb.ServiceAccess
	for rows.Next() {
		foo := savepb.ServiceAccess{}
		err := rows.Scan(&foo.ID, &foo.CallingService, &foo.SubjectService)
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
func (a *DBServiceAccess) CreateTable(ctx context.Context) error {
	csql := []string{
		`create sequence if not exists ` + a.SQLTablename + `_seq;`,
		`CREATE TABLE if not exists ` + a.SQLTablename + ` (id integer primary key default nextval('` + a.SQLTablename + `_seq'),callingservice text not null ,subjectservice text not null );`,
		`CREATE TABLE if not exists ` + a.SQLTablename + `_archive (id integer primary key default nextval('` + a.SQLTablename + `_seq'),callingservice text not null ,subjectservice text not null );`,
		`ALTER TABLE serviceaccess ADD COLUMN IF NOT EXISTS callingservice text not null default '';`,
		`ALTER TABLE serviceaccess ADD COLUMN IF NOT EXISTS subjectservice text not null default '';`,

		`ALTER TABLE serviceaccess_archive ADD COLUMN IF NOT EXISTS callingservice text not null default '';`,
		`ALTER TABLE serviceaccess_archive ADD COLUMN IF NOT EXISTS subjectservice text not null default '';`,
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
func (a *DBServiceAccess) Error(ctx context.Context, q string, e error) error {
	if e == nil {
		return nil
	}
	return fmt.Errorf("[table="+a.SQLTablename+", query=%s] Error: %s", q, e)
}
