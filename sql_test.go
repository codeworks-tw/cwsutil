package cwsutil

import (
	"context"
	"log"
	"testing"

	_ "ariga.io/atlas-provider-gorm/gormschema"
	"github.com/codeworks-tw/cwsutil/cwssql"
)

type Account struct {
	cwssql.BaseIdModel
	Account  string `gorm:"type:text;not null;uniqueIndex"`
	Password string `gorm:"type:text;not null"`
	Enable   bool   `gorm:"default:true"`
	cwssql.BaseTimeModel
}

type AccountRepository struct {
	cwssql.Repository[Account]
}

func (r *AccountRepository) GetAccount(account string) (*Account, error) {
	return r.Get(map[string]any{"Account": account})
}

func NewAccountRepository(context context.Context, session *cwssql.DBSession) AccountRepository {
	return AccountRepository{
		Repository: cwssql.NewRepository[Account](context, session),
	}
}

func TestSqlRepository(t *testing.T) {
	db, err := cwssql.NewSQLiteDB("database.db")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	repo := NewAccountRepository(ctx, cwssql.NewSession(db))
	err = repo.Begin()
	if err != nil {
		panic(err)
	}

	account, err := repo.GetAccount("test@xxx.com")
	if err != nil {
		panic(err)
	}
	log.Printf("Retrieved account: %+v", account)

	if account == nil {
		account = &Account{
			Password: "XXXX",
			Account:  "test@xxx.com",
		}
	}

	err = repo.Upsert(account)
	if err != nil {
		panic(err)
	}
	log.Println("Upsert account", account.Id)

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
