package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type budgetCategory struct {
	ID            int64     `json:"id"`
	UserCompanyId NullInt64 `json:"user_company_id"`
	Name          string    `json:"name"`
}

type lawFirmActions struct {
	ID            int64     `json:"id"`
	UserCompanyId NullInt64 `json:"user_company_id"`
	Name          string    `json:"name"`
}

type lawFirmCost struct {
	ID               int64     `json:"id"`
	UserId           NullInt64 `json:"user_id"`
	UserCompanyId    NullInt64 `json:"user_company_id"`
	CustomerNumberId NullInt64 `json:"customer_number_id"`
	ValidFrom        NullTime  `json:"valid_from"`
}

type lawFirmCostUsptoFees struct {
	ID            int64     `json:"id"`
	UserId        NullInt64 `json:"user_id"`
	UserCompanyId NullInt64 `json:"user_company_id"`
	ValidFrom     NullTime  `json:"valid_from"`
}

//------------------------------------------------------------------------------------------------
func (lfa *lawFirmActions) lawFirmAction_createNew(db *sql.DB) error {
	user_company_id, err := lfa.UserCompanyId.MarshalDB()
	statement := "INSERT INTO law_firm_action_list(id, user_company_id, name) VALUES (NULL, ?, ?)"
	if err != nil {
		return err
	}
	res, err := db.Exec(statement, user_company_id, lfa.Name)
	if err != nil {
		return err
	}
	lfa.ID, err = res.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}

func (lfa *lawFirmActions) lawFirmActions_getOne(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT user_company_id, name FROM law_firm_action_list WHERE id=%d", lfa.ID)
	return db.QueryRow(statement).Scan(&lfa.UserCompanyId, &lfa.Name)
}

func lawFirmActions_getAll(db *sql.DB, start, count int64) ([]lawFirmActions, error) {
	statement := fmt.Sprintf("SELECT id, user_company_id, name FROM law_firm_action_list LIMIT %d OFFSET %d", count, start)
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lfas := []lawFirmActions{}
	var lfa lawFirmActions
	for rows.Next() {
		if err := rows.Scan(&lfa.ID, &lfa.UserCompanyId, &lfa.Name); err != nil {
			return nil, err
		}
		lfas = append(lfas, lfa)
	}
	return lfas, nil
}

func (lfa *lawFirmActions) lawFirmActions_update(db *sql.DB) error {
	user_company_id, err := lfa.UserCompanyId.MarshalDB()
	statement := "UPDATE law_firm_action_list SET user_company_id=?, name=? WHERE id=?"
	_, err = db.Exec(statement, user_company_id, lfa.Name, lfa.ID)
	return err
}

func (lfa *lawFirmActions) lawFirmActions_delete(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM law_firm_action_list WHERE id=%d", lfa.ID)
	_, err := db.Exec(statement)
	return err
}

//--------------------------------------------------------------------------------------------------------------------------
func (b *budgetCategory) budgetCategory_createNew(db *sql.DB) error {
	user_company_id, err := b.UserCompanyId.MarshalDB()
	statement := "INSERT INTO budget_categories_list(id, user_company_id, name) VALUES (NULL, ?, ?)"
	if err != nil {
		return err
	}
	res, err := db.Exec(statement, user_company_id, b.Name)
	if err != nil {
		return err
	}
	b.ID, err = res.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}

func (b *budgetCategory) budgetCategory_getOne(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT user_company_id, name FROM budget_categories_list WHERE id=%d", b.ID)
	return db.QueryRow(statement).Scan(&b.UserCompanyId, &b.Name)
}

func budgetCategory_getAll(db *sql.DB, start, count int64) ([]budgetCategory, error) {
	statement := fmt.Sprintf("SELECT id, user_company_id, name FROM budget_categories_list LIMIT %d OFFSET %d", count, start)
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bdg := []budgetCategory{}
	for rows.Next() {
		var b budgetCategory
		if err := rows.Scan(&b.ID, &b.UserCompanyId, &b.Name); err != nil {
			return nil, err
		}
		bdg = append(bdg, b)
	}
	return bdg, nil
}

func (b *budgetCategory) budgetCategory_update(db *sql.DB) error {
	user_company_id, err := b.UserCompanyId.MarshalDB()
	statement := "UPDATE budget_categories_list SET user_company_id=?, name=? WHERE id=?"
	_, err = db.Exec(statement, user_company_id, b.Name, b.ID)
	return err
}

func (b *budgetCategory) budgetCategory_delete(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM budget_categories_list WHERE id=%d", b.ID)
	_, err := db.Exec(statement)
	return err
}

//------------------------------------------------------------------------------------
func (lfc *lawFirmCost) lawFirmCost_createNew(db *sql.DB) error {
	user_id, err := lfc.UserId.MarshalDB()
	user_company_id, err := lfc.UserCompanyId.MarshalDB()
	customer_number_id, err := lfc.CustomerNumberId.MarshalDB()
	valid_from, err := lfc.ValidFrom.MarshalDB()
	//fmt.Println(valid_from)
	statement := "INSERT INTO law_firm_cost(user_id, user_company_id, customer_number_id, valid_from) VALUES (?, ?, ?, ?)"
	if err != nil {
		return err
	}
	res, err := db.Exec(statement, user_id, user_company_id, customer_number_id, valid_from)
	if err != nil {
		return err
	}
	lfc.ID, err = res.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}

func (lfc *lawFirmCost) lawFirmCost_getOne(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT user_id, user_company_id, customer_number_id, valid_from FROM law_firm_cost WHERE id=%d", lfc.ID)
	return db.QueryRow(statement).Scan(&lfc.UserId, &lfc.UserCompanyId, &lfc.CustomerNumberId, &lfc.ValidFrom)
}

func lawFirmCost_getAll(db *sql.DB, start, count int64) ([]lawFirmCost, error) {
	statement := fmt.Sprintf("SELECT id, user_id, user_company_id, customer_number_id, valid_from FROM law_firm_cost LIMIT %d OFFSET %d", count, start)
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lfcs := []lawFirmCost{}
	var lfc lawFirmCost
	for rows.Next() {
		if err := rows.Scan(&lfc.ID, &lfc.UserId, &lfc.UserCompanyId, &lfc.CustomerNumberId, &lfc.ValidFrom); err != nil {
			return nil, err
		}
		lfcs = append(lfcs, lfc)
	}
	return lfcs, nil
}

func (lfc *lawFirmCost) lawFirmCost_update(db *sql.DB) error {
	user_id, err := lfc.UserId.MarshalDB()
	user_company_id, err := lfc.UserCompanyId.MarshalDB()
	customer_number_id, err := lfc.CustomerNumberId.MarshalDB()
	valid_from, err := lfc.ValidFrom.MarshalDB()

	statement := "UPDATE law_firm_cost SET user_id=?, user_company_id=?, customer_number_id=?, valid_from=? WHERE id=?"
	_, err = db.Exec(statement, user_id, user_company_id, customer_number_id, valid_from, lfc.ID)
	return err
}

func (lfc *lawFirmCost) lawFirmCost_delete(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM law_firm_cost WHERE id=%d", lfc.ID)
	_, err := db.Exec(statement)
	return err
}

//--------------------------------------------------------------------------------------------------------------------------

func (lfcuf *lawFirmCostUsptoFees) lawFirmCostUsptoFees_createNew(db *sql.DB) error {
	user_id, err := lfcuf.UserId.MarshalDB()
	user_company_id, err := lfcuf.UserCompanyId.MarshalDB()
	valid_from, err := lfcuf.ValidFrom.MarshalDB()

	statement := "INSERT INTO law_firm_cost_uspto_fees(id, user_id, user_company_id, valid_from) VALUES (NULL, ?, ?)"
	if err != nil {
		return err
	}
	res, err := db.Exec(statement, user_id, user_company_id, valid_from)
	if err != nil {
		return err
	}
	lfcuf.ID, err = res.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}

func (lfcuf *lawFirmCostUsptoFees) lawFirmCostUsptoFees_getOne(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT user_id, user_company_id, valid_from FROM law_firm_cost_uspto_fees WHERE id=%d", lfcuf.ID)
	return db.QueryRow(statement).Scan(&lfcuf.UserId, &lfcuf.UserCompanyId, &lfcuf.ValidFrom)
}

func lawFirmCostUsptoFees_getAll(db *sql.DB, start, count int64) ([]lawFirmCostUsptoFees, error) {
	statement := fmt.Sprintf("SELECT id, user_id, user_company_id, valid_from FROM law_firm_cost_uspto_fees LIMIT %d OFFSET %d", count, start)
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lfcufs := []lawFirmCostUsptoFees{}
	var lfcuf lawFirmCostUsptoFees
	for rows.Next() {
		if err := rows.Scan(&lfcuf.ID, &lfcuf.UserId, &lfcuf.UserCompanyId, &lfcuf.ValidFrom); err != nil {
			return nil, err
		}
		lfcufs = append(lfcufs, lfcuf)
	}
	return lfcufs, nil
}

func (lfcuf *lawFirmCostUsptoFees) lawFirmCostUsptoFees_update(db *sql.DB) error {
	user_id, err := lfcuf.UserId.MarshalDB()
	user_company_id, err := lfcuf.UserCompanyId.MarshalDB()
	valid_from, err := lfcuf.ValidFrom.MarshalDB()

	statement := "UPDATE law_firm_cost_uspto_fees SET user_id=?, user_company_id=?, valid_from=? WHERE id=?"
	_, err = db.Exec(statement, user_id, user_company_id, valid_from, lfcuf.ID)
	return err
}

func (lfcuf *lawFirmCostUsptoFees) lawFirmCostUsptoFees_delete(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM law_firm_cost_uspto_fees WHERE id=%d", lfcuf.ID)
	_, err := db.Exec(statement)
	return err
}
