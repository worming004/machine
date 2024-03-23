package sqliterepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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
	machinerows, err := m.db.Query(`SELECT state, count_of_ignored_transition FROM machines WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer machinerows.Close()

	var (
		stateName                string
		countOfIgnoredTransition int
	)
	machinerows.Next()
	if err = machinerows.Scan(&stateName, &countOfIgnoredTransition); err != nil {
		return nil, fmt.Errorf("error while scanning machines: %w", err)
	}

	bubblerows, err := m.db.Query(`SELECT name FROM bubbles where machine_id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer bubblerows.Close()
	var bubbles []bubblemachine.Bubble
	for bubblerows.Next() {
		var b string
		err = bubblerows.Scan(&b)
		if err != nil {
			return nil, fmt.Errorf("error while scanning bubbles: %w", err)
		}
		bubbles = append(bubbles, bubblemachine.NewBubble(b))
	}

	piecesrows, err := m.db.Query(`SELECT value FROM pieces where machine_id = ?;`, id)
	if err != nil {
		return nil, err
	}
	var pieces []bubblemachine.Piece
	for piecesrows.Next() {
		var p int
		err = piecesrows.Scan(&p)
		if err != nil {
			return nil, fmt.Errorf("error while scanning pieces: %w", err)
		}
		pieces = append(pieces, bubblemachine.NewPiece(p))
	}

	var machine bubblemachine.Machine
	bubblemachine.NewMachineSetterForDb(&machine).
		SetBubbles(bubbles).
		SetPieces(pieces).
		SetCount(countOfIgnoredTransition).
		SetStateByName(bubblemachine.StateName(stateName))

	return &machine, nil
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
	default:
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
	for _, p := range m.GetPieces() {
		_, err = tx.ExecContext(ctx, "INSERT INTO pieces (value, machine_id) VALUES (?, ?)", p.Value(), m.GetId())
		if err != nil {
			return err
		}
	}
	return nil
}

func (mr MachineRepository) Close() error {
	return mr.db.Close()
}
