package cwssql

import (
	"context"
	"log"
	"testing"

	_ "ariga.io/atlas-provider-gorm/gormschema"
)

type Account struct {
	Account  string `gorm:"type:text;primaryKey"`
	Password string `gorm:"type:text;not null"`
	Enable   bool   `gorm:"default:true"`
	BaseTimeModel
}

type AccountRepository struct {
	Repository[Account]
}

func (r *AccountRepository) GetAccount(account string) (*Account, error) {
	return r.Get(map[string]any{"Account": account})
}

func NewAccountRepository(context context.Context, session *DBSession) AccountRepository {
	return AccountRepository{
		Repository: NewRepository[Account](context, session),
	}
}

func TestSqlRepository(t *testing.T) {
	db, err := NewSQLiteDB("data.db")
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&Account{})
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	repo := NewAccountRepository(ctx, NewSession(db))
	err = repo.Begin()
	if err != nil {
		panic(err)
	}

	account := &Account{
		Password: "XXXX",
		Account:  "test@xxx.com",
	}

	err = repo.Upsert(account)
	if err != nil {
		panic(err)
	}
	log.Println("Upsert account", account.Account)

	account.Account = "XXXXXXXXXXXX"
	repo.Refresh(account)

	account, err = repo.GetAccount("test@xxx.com")
	if err != nil {
		panic(err)
	}
	log.Printf("Retrieved account: %+v", account)

	accounts, err := repo.GetAll(map[string]any{"Account": "test@xxx.com"})
	if err != nil {
		log.Printf("Error getting all accounts: %v", err)
		return
	}
	log.Printf("All account: %+v", accounts)

	err = repo.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		return
	}

	err = repo.Delete(account)
	if err != nil {
		log.Printf("Error deleting account: %v", err)
		return
	}

	accounts, err = repo.DeleteAll(map[string]any{"Account": "test@xxx.com"})
	if err != nil {
		log.Printf("Error deleting all accounts: %v", err)
		return
	}
	log.Printf("All account: %+v", accounts)
}
