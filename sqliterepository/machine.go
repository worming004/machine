package sqliterepository

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
	"github.com/worming004/machine/bubblemachine"
)

var machineRepositoryImpl bubblemachine.MachineRepository = MachineRepository{}

type MachineRepository struct {
	db *sql.DB
}

type NewRepositoryRequest struct {
	DataSourceName string
	Init           bool
}

func NewMachineRepository(r NewRepositoryRequest) (*MachineRepository, error) {
	db, err := sql.Open("sqlite3", r.DataSourceName)
	if err != nil {
		return nil, err
	}
	m := &MachineRepository{db: db}

	if r.Init {
		if err = initDb(m); err != nil {
			return nil, err
		}
	}

	return m, nil
}

func initDb(m *MachineRepository) error {
	_, err := m.db.Exec(`CREATE TABLE IF NOT EXISTS machines (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    state TEXT NOT NULL,
    count_of_ignored_transition INTEGER NOT NULL
    )`)
	if err != nil {
		return err
	}

	_, err = m.db.Exec(`CREATE TABLE IF NOT EXISTS bubbles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    machine_id INTEGER NOT NULL,
    FOREIGN KEY(machine_id) REFERENCES machines(id)
    )`)
	if err != nil {
		return err
	}

	_, err = m.db.Exec(`CREATE TABLE IF NOT EXISTS pieces (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    value INTEGER NOT NULL,
    machine_id INTEGER NOT NULL,
    FOREIGN KEY(machine_id) REFERENCES machines(id)
    )`)
	if err != nil {
		return err
	}

	return nil
}

// Get implements bubblemachine.MachineRepository.
func (m MachineRepository) Get(ctx context.Context, id int) (*bubblemachine.Machine, error) {
	panic("unimplemented")
}

// Save implements bubblemachine.MachineRepository.
func (mr MachineRepository) Save(ctx context.Context, m *bubblemachine.Machine) error {
	tx, err := mr.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	row := tx.QueryRowContext(ctx, "SELECT COUNT() FROM machines WHERE id = ?", m.GetId())
	var count int
	err = row.Scan(&count)
	if err != nil {
		tx.Rollback()
		return err
	}

	switch count {
	case 0:
		res, err := tx.ExecContext(ctx, "INSERT INTO machines (state, count_of_ignored_transition) VALUES (?, ?)", m.GetStateName(), m.CountOfIgnoredTransition())
		if err != nil {
			tx.Rollback()
			return err
		}
		id, err := res.LastInsertId()
		if err != nil {
			tx.Rollback()
			return err
		}
		m.SetId(id)
		err = mr.saveBubbles(ctx, tx, m)
		if err != nil {
			tx.Rollback()
			return err
		}
		err = mr.savePieces(ctx, tx, m)
		if err != nil {
			tx.Rollback()
			return err
		}

	case 1:
		_, err := mr.db.Exec("UPDATE machines (state, count_of_ignored_transition) VALUES (?, ?) WHERE id = ?", m.GetStateName(), m.CountOfIgnoredTransition(), m.GetId())
		if err != nil {
			return err
		}
		err = mr.saveBubbles(ctx, tx, m)
		if err != nil {
			tx.Rollback()
			return err
		}
		err = mr.savePieces(ctx, tx, m)
		if err != nil {
			tx.Rollback()
			return err
		}
		break
	case 2:
		tx.Rollback()
		return errors.New("multiple machines with same id exists")
	}

	return tx.Commit()
}

func (mr MachineRepository) saveBubbles(ctx context.Context, tx *sql.Tx, m *bubblemachine.Machine) error {
	_, err := tx.ExecContext(ctx, "DELETE FROM bubbles WHERE machine_id = ?", m.GetId())
	if err != nil {
		return err
	}
	for _, b := range m.GetBubbles() {
		_, err = tx.ExecContext(ctx, "INSERT INTO bubbles (name, machine_id) VALUES (?, ?)", b.String(), m.GetId())
		if err != nil {
			return err
		}
	}
	return nil
}
func (mr MachineRepository) savePieces(ctx context.Context, tx *sql.Tx, m *bubblemachine.Machine) error {
	_, err := tx.ExecContext(ctx, "DELETE FROM pieces WHERE machine_id = ?", m.GetId())
	if err != nil {
		return err
	}
	for _, b := range m.GetBubbles() {
		_, err = tx.ExecContext(ctx, "INSERT INTO pieces (value, machine_id) VALUES (?, ?)", b.String(), m.GetId())
		if err != nil {
			return err
		}
	}
	return nil
}

func (mr MachineRepository) Close() error {
	return mr.db.Close()
}
