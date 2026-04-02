package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Checklist struct {
	ID string `json:"id"`
	CustomerName string `json:"customer_name"`
	Template string `json:"template"`
	Steps string `json:"steps"`
	Progress int `json:"progress"`
	Assignee string `json:"assignee"`
	Status string `json:"status"`
	DueDate string `json:"due_date"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"concierge.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS checklists(id TEXT PRIMARY KEY,customer_name TEXT NOT NULL,template TEXT DEFAULT '',steps TEXT DEFAULT '[]',progress INTEGER DEFAULT 0,assignee TEXT DEFAULT '',status TEXT DEFAULT 'active',due_date TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Checklist)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO checklists(id,customer_name,template,steps,progress,assignee,status,due_date,created_at)VALUES(?,?,?,?,?,?,?,?,?)`,e.ID,e.CustomerName,e.Template,e.Steps,e.Progress,e.Assignee,e.Status,e.DueDate,e.CreatedAt);return err}
func(d *DB)Get(id string)*Checklist{var e Checklist;if d.db.QueryRow(`SELECT id,customer_name,template,steps,progress,assignee,status,due_date,created_at FROM checklists WHERE id=?`,id).Scan(&e.ID,&e.CustomerName,&e.Template,&e.Steps,&e.Progress,&e.Assignee,&e.Status,&e.DueDate,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Checklist{rows,_:=d.db.Query(`SELECT id,customer_name,template,steps,progress,assignee,status,due_date,created_at FROM checklists ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Checklist;for rows.Next(){var e Checklist;rows.Scan(&e.ID,&e.CustomerName,&e.Template,&e.Steps,&e.Progress,&e.Assignee,&e.Status,&e.DueDate,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM checklists WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM checklists`).Scan(&n);return n}
